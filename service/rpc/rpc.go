package rpc

import (
	"bytes"
	"compress/gzip"
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/mikeydub/go-gallery/service/logger"
	"github.com/mikeydub/go-gallery/service/tracing"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/everFinance/goar"
	"github.com/gorilla/websocket"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/mikeydub/go-gallery/contracts"
	"github.com/mikeydub/go-gallery/service/persist"
	"github.com/mikeydub/go-gallery/util"
	"github.com/spf13/viper"
)

var keepAliveTimeout = 600 * time.Second
var client = &http.Client{
	Timeout: time.Second * 30,
	Transport: tracing.NewTracingTransport(&http.Transport{
		Dial: (&net.Dialer{
			KeepAlive: keepAliveTimeout,
		}).Dial,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
	}, true),
}

// Transfer represents a Transfer from the RPC response
type Transfer struct {
	BlockNumber     persist.BlockNumber
	From            persist.EthereumAddress
	To              persist.EthereumAddress
	TokenID         persist.TokenID
	TokenType       persist.TokenType
	Amount          uint64
	ContractAddress persist.EthereumAddress
}

// TokenContractMetadata represents a token contract's metadata
type TokenContractMetadata struct {
	Name   string
	Symbol string
}

// ErrHTTP represents an error returned from an HTTP request
type ErrHTTP struct {
	URL    string
	Status int
}

// NewEthClient returns an ethclient.Client
func NewEthClient() *ethclient.Client {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dialer := *websocket.DefaultDialer
	dialer.ReadBufferSize = 1024 * 20
	rpcClient, err := rpc.DialWebsocketWithDialer(ctx, viper.GetString("RPC_URL"), "", dialer)
	if err != nil {
		panic(err)
	}

	return ethclient.NewClient(rpcClient)

}

// NewIPFSShell returns an IPFS shell
func NewIPFSShell() *shell.Shell {
	sh := shell.NewShell(viper.GetString("IPFS_URL"))
	sh.SetTimeout(time.Minute * 2)
	return sh
}

// NewArweaveClient returns an Arweave client
func NewArweaveClient() *goar.Client {
	return goar.NewClient("https://arweave.net")
}

// GetTokenContractMetadata returns the metadata for a given contract (without URI)
func GetTokenContractMetadata(address persist.EthereumAddress, ethClient *ethclient.Client) (*TokenContractMetadata, error) {
	contract := address.Address()
	instance, err := contracts.NewIERC721MetadataCaller(contract, ethClient)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()
	name, err := instance.Name(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		return nil, err
	}
	symbol, err := instance.Symbol(&bind.CallOpts{
		Context: ctx,
	})
	if err != nil {
		return nil, err
	}

	return &TokenContractMetadata{Name: name, Symbol: symbol}, nil
}

// GetMetadataFromURI parses and returns the NFT metadata for a given token URI
func GetMetadataFromURI(ctx context.Context, turi persist.TokenURI, ipfsClient *shell.Shell, arweaveClient *goar.Client) (persist.TokenMetadata, error) {

	ctx, cancel := context.WithTimeout(ctx, time.Minute*1)
	defer cancel()
	var meta persist.TokenMetadata
	err := DecodeMetadataFromURI(ctx, turi, &meta, ipfsClient, arweaveClient)
	if err != nil {
		return nil, err
	}

	return meta, nil

}

// GetDataFromURI calls URI and returns the data
func GetDataFromURI(ctx context.Context, turi persist.TokenURI, ipfsClient *shell.Shell, arweaveClient *goar.Client) ([]byte, error) {

	d, _ := ctx.Deadline()
	logger.For(ctx).Debugf("Getting data from URI: %s -timeout: %s", turi.String(), time.Until(d))
	asString := turi.String()

	logger.For(ctx).Debugf("Getting data from %s with type %s", asString, turi.Type())

	switch turi.Type() {
	case persist.URITypeBase64JSON, persist.URITypeBase64SVG:
		// decode the base64 encoded json
		b64data := asString[strings.IndexByte(asString, ',')+1:]
		decoded, err := base64.StdEncoding.DecodeString(string(b64data))
		if err != nil {
			return nil, fmt.Errorf("error decoding base64 data: %s \n\n%s", err, b64data)
		}

		return removeBOM(decoded), nil
	case persist.URITypeIPFS:
		path := strings.ReplaceAll(asString, "ipfs://", "")
		path = strings.ReplaceAll(path, "ipfs/", "")
		path = strings.Split(path, "?")[0]

		it, err := ipfsClient.Cat(path)
		if err != nil {
			if err == context.Canceled {
				c, cancel := context.WithTimeout(context.Background(), time.Second*30)
				defer cancel()
				ctx = c
			}
			bs, nextErr := GetIPFSData(ctx, path)
			if nextErr == nil {
				return removeBOM(bs), nil
			}

			return nil, fmt.Errorf("error getting data from ipfs: %s | %s - cat: %s", err, nextErr, path)
		}
		defer it.Close()

		buf := &bytes.Buffer{}
		err = util.CopyMax(buf, it, 1024*1024*1024)
		if err != nil {
			if err == context.Canceled {
				c, cancel := context.WithTimeout(context.Background(), time.Second*30)
				defer cancel()
				ctx = c
			}
			bs, nextErr := GetIPFSData(ctx, path)
			if nextErr == nil {
				return removeBOM(bs), nil
			}
			return nil, fmt.Errorf("error getting data from ipfs: %s - cat: %s", err, path)
		}

		return removeBOM(buf.Bytes()), nil
	case persist.URITypeArweave:
		path := strings.ReplaceAll(asString, "arweave://", "")
		path = strings.ReplaceAll(path, "ar://", "")
		return GetArweaveData(arweaveClient, path)
	case persist.URITypeHTTP:

		req, err := http.NewRequestWithContext(ctx, "GET", asString, nil)
		if err != nil {
			return nil, fmt.Errorf("error creating request: %s", err)
		}
		resp, err := client.Do(req)
		if err != nil {
			return nil, fmt.Errorf("error getting data from http: %s", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode > 399 || resp.StatusCode < 200 {
			return nil, ErrHTTP{Status: resp.StatusCode, URL: asString}
		}
		buf := &bytes.Buffer{}
		err = util.CopyMax(buf, resp.Body, 1024*1024*1024)
		if err != nil {
			return nil, fmt.Errorf("error getting data from http: %s - %s", err, asString)
		}

		return removeBOM(buf.Bytes()), nil
	case persist.URITypeIPFSAPI:
		parsedURL, err := url.Parse(asString)
		if err != nil {
			return nil, err
		}
		query := parsedURL.Query().Get("arg")
		it, err := ipfsClient.Cat(query)
		if err != nil {
			return nil, err
		}
		defer it.Close()
		buf := &bytes.Buffer{}
		err = util.CopyMax(buf, it, 1024*1024*1024)
		if err != nil {
			return nil, fmt.Errorf("error getting data from ipfs: %s - cat: %s", err, query)
		}

		return removeBOM(buf.Bytes()), nil
	case persist.URITypeJSON, persist.URITypeSVG:
		idx := strings.IndexByte(asString, '{')
		if idx == -1 {
			return []byte(asString), nil
		}
		return removeBOM([]byte(asString[idx:])), nil

	default:
		return nil, fmt.Errorf("unknown token URI type: %s - %s", turi.Type(), turi)
	}

}

// DecodeMetadataFromURI calls URI and decodes the data into a metadata map
func DecodeMetadataFromURI(ctx context.Context, turi persist.TokenURI, into *persist.TokenMetadata, ipfsClient *shell.Shell, arweaveClient *goar.Client) error {

	d, _ := ctx.Deadline()
	logger.For(ctx).Debugf("Getting data from URI: %s -timeout: %s", turi.String(), time.Until(d))
	asString := turi.String()

	logger.For(ctx).Debugf("Getting data from %s with type %s", asString, turi.Type())

	switch turi.Type() {
	case persist.URITypeBase64JSON:
		// decode the base64 encoded json
		b64data := asString[strings.IndexByte(asString, ',')+1:]
		decoded, err := base64.StdEncoding.DecodeString(string(b64data))
		if err != nil {
			return fmt.Errorf("error decoding base64 data: %s \n\n%s", err, b64data)
		}

		return json.Unmarshal(removeBOM(decoded), into)
	case persist.URITypeBase64SVG:
		b64data := asString[strings.IndexByte(asString, ',')+1:]
		decoded, err := base64.StdEncoding.DecodeString(string(b64data))
		if err != nil {
			return fmt.Errorf("error decoding base64 data: %s \n\n%s", err, b64data)
		}
		into = &persist.TokenMetadata{"image": string(decoded)}
		return nil
	case persist.URITypeIPFS:
		path := strings.ReplaceAll(asString, "ipfs://", "")
		path = strings.ReplaceAll(path, "ipfs/", "")
		path = strings.Split(path, "?")[0]

		it, err := ipfsClient.Cat(path)
		if err != nil {
			if err == context.Canceled {
				c, cancel := context.WithTimeout(context.Background(), time.Second*10)
				defer cancel()
				ctx = c
			}
			bs, nextErr := GetIPFSData(ctx, path)
			if nextErr == nil {
				return json.Unmarshal(bs, into)
			}

			return fmt.Errorf("error getting data from ipfs: %s | %s - cat: %s", err, nextErr, path)
		}
		defer it.Close()
		return json.NewDecoder(it).Decode(into)
	case persist.URITypeArweave:
		path := strings.ReplaceAll(asString, "arweave://", "")
		path = strings.ReplaceAll(path, "ar://", "")
		result, err := GetArweaveData(arweaveClient, path)
		if err != nil {
			return err
		}
		return json.Unmarshal(result, into)
	case persist.URITypeHTTP:

		req, err := http.NewRequestWithContext(ctx, "GET", asString, nil)
		if err != nil {
			return fmt.Errorf("error creating request: %s", err)
		}
		resp, err := client.Do(req)
		if err != nil {
			return fmt.Errorf("error getting data from http: %s", err)
		}
		defer resp.Body.Close()
		if resp.StatusCode > 399 || resp.StatusCode < 200 {
			return ErrHTTP{Status: resp.StatusCode, URL: asString}
		}
		return json.NewDecoder(resp.Body).Decode(into)
	case persist.URITypeIPFSAPI:
		parsedURL, err := url.Parse(asString)
		if err != nil {
			return err
		}
		query := parsedURL.Query().Get("arg")
		it, err := ipfsClient.Cat(query)
		if err != nil {
			return err
		}
		defer it.Close()
		return json.NewDecoder(it).Decode(into)
	case persist.URITypeJSON, persist.URITypeSVG:
		idx := strings.IndexByte(asString, '{')
		if idx == -1 {
			return json.Unmarshal(removeBOM([]byte(asString)), into)
		}
		return json.Unmarshal(removeBOM([]byte(asString[idx:])), into)

	default:
		return fmt.Errorf("unknown token URI type: %s", turi.Type())
	}

}

func removeBOM(bs []byte) []byte {
	if len(bs) > 3 && bs[0] == 0xEF && bs[1] == 0xBB && bs[2] == 0xBF {
		return bs[3:]
	}
	return bs
}

func GetIPFSData(pCtx context.Context, path string) ([]byte, error) {
	url := fmt.Sprintf("https://ipfs.io/ipfs/%s", path)

	req, err := http.NewRequestWithContext(pCtx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %s", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error getting data from http: %s", err)
	}
	if resp.StatusCode > 399 || resp.StatusCode < 200 {
		return nil, ErrHTTP{Status: resp.StatusCode, URL: url}
	}
	defer resp.Body.Close()

	buf := &bytes.Buffer{}
	err = util.CopyMax(buf, resp.Body, 1024*1024*1024)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// GetIPFSHeaders returns the headers for the given IPFS hash
func GetIPFSHeaders(pCtx context.Context, path string) (http.Header, error) {
	url := fmt.Sprintf("https://ipfs.io/ipfs/%s", path)

	req, err := http.NewRequestWithContext(pCtx, "GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %s", err)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error getting data from http: %s", err)
	}
	if resp.StatusCode > 399 || resp.StatusCode < 200 {
		return nil, ErrHTTP{Status: resp.StatusCode, URL: url}
	}
	defer resp.Body.Close()

	return resp.Header, nil
}

// GetTokenURI returns metadata URI for a given token address.
func GetTokenURI(ctx context.Context, pTokenType persist.TokenType, pContractAddress persist.EthereumAddress, pTokenID persist.TokenID, ethClient *ethclient.Client) (persist.TokenURI, error) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Minute)
	defer cancel()
	contract := common.HexToAddress(string(pContractAddress))
	switch pTokenType {
	case persist.TokenTypeERC721:

		instance, err := contracts.NewIERC721MetadataCaller(contract, ethClient)
		if err != nil {
			return "", err
		}

		logger.For(ctx).Debugf("Token ID: %s\tToken Address: %s", pTokenID.String(), contract.Hex())

		turi, err := instance.TokenURI(&bind.CallOpts{
			Context: ctx,
		}, pTokenID.BigInt())
		if err != nil {
			return "", err
		}

		return persist.TokenURI(strings.ReplaceAll(turi, "\x00", "")), nil
	case persist.TokenTypeERC1155:

		instance, err := contracts.NewIERC1155MetadataURI(contract, ethClient)
		if err != nil {
			return "", err
		}

		logger.For(ctx).Debugf("Token ID: %d\tToken Address: %s", pTokenID.BigInt().Uint64(), contract.Hex())

		turi, err := instance.Uri(&bind.CallOpts{
			Context: ctx,
		}, pTokenID.BigInt())
		if err != nil {
			return "", err
		}

		return persist.TokenURI(strings.ReplaceAll(turi, "\x00", "")), nil

	default:
		return "", fmt.Errorf("unknown token type: %s", pTokenType)
	}
}

// GetBalanceOfERC1155Token returns the balance of an ERC1155 token
func GetBalanceOfERC1155Token(pOwnerAddress, pContractAddress persist.EthereumAddress, pTokenID persist.TokenID, ethClient *ethclient.Client) (*big.Int, error) {
	contract := common.HexToAddress(string(pContractAddress))
	owner := common.HexToAddress(string(pOwnerAddress))
	instance, err := contracts.NewIERC1155(contract, ethClient)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()
	bal, err := instance.BalanceOf(&bind.CallOpts{
		Context: ctx,
	}, owner, pTokenID.BigInt())
	if err != nil {
		return nil, err
	}

	return bal, nil
}

// GetContractCreator returns the address of the contract creator
func GetContractCreator(ctx context.Context, contractAddress persist.EthereumAddress, ethClient *ethclient.Client) (persist.EthereumAddress, error) {
	highestBlock, err := ethClient.BlockNumber(ctx)
	if err != nil {
		return "", fmt.Errorf("error getting highest block: %s", err.Error())
	}
	_, err = ethClient.CodeAt(ctx, contractAddress.Address(), big.NewInt(int64(highestBlock)))
	if err != nil {
		return "", fmt.Errorf("error getting code at: %s", err.Error())
	}
	lowestBlock := uint64(0)

	for lowestBlock <= highestBlock {
		midBlock := uint64((highestBlock + lowestBlock) / 2)
		codeAt, err := ethClient.CodeAt(ctx, contractAddress.Address(), big.NewInt(int64(midBlock)))
		if err != nil {
			return "", fmt.Errorf("error getting code at: %s", err.Error())
		}
		if len(codeAt) > 0 {
			highestBlock = midBlock
		} else {
			lowestBlock = midBlock
		}

		if lowestBlock+1 == highestBlock {
			break
		}
	}
	block, err := ethClient.BlockByNumber(ctx, big.NewInt(int64(highestBlock)))
	if err != nil {
		return "", fmt.Errorf("error getting block: %s", err.Error())
	}
	txs := block.Transactions()
	for _, tx := range txs {
		receipt, err := ethClient.TransactionReceipt(ctx, tx.Hash())
		if err != nil {
			return "", fmt.Errorf("error getting transaction receipt: %s", err.Error())
		}
		if receipt.ContractAddress == contractAddress.Address() {
			msg, err := tx.AsMessage(types.HomesteadSigner{}, nil)
			if err != nil {
				return "", fmt.Errorf("error getting message: %s", err.Error())
			}
			return persist.EthereumAddress(fmt.Sprintf("0x%s", strings.ToLower(msg.From().String()))), nil
		}
	}
	return "", fmt.Errorf("could not find contract creator")
}

// GetArweaveData returns the data from an Arweave transaction
func GetArweaveData(client *goar.Client, id string) ([]byte, error) {
	tx, err := client.GetTransactionByID(id)
	if err != nil {
		return nil, err
	}
	data, err := client.GetTransactionData(id)
	if err != nil {
		return nil, err
	}

	decoded, err := base64.RawURLEncoding.DecodeString(string(data))
	if err == nil {
		data = decoded
	}

	for _, tag := range tx.Tags {
		decodedName, err := base64.RawURLEncoding.DecodeString(tag.Name)
		if err != nil {
			return nil, err
		}
		if strings.EqualFold(string(decodedName), "Content-Encoding") {
			decodedValue, err := base64.RawURLEncoding.DecodeString(tag.Value)
			if err != nil {
				return nil, err
			}
			if strings.EqualFold(string(decodedValue), "gzip") {
				zipped, err := gzip.NewReader(bytes.NewReader(data))
				if err != nil {
					return nil, err
				}
				buf := new(bytes.Buffer)
				_, err = io.Copy(buf, zipped)
				if err != nil {
					return nil, err
				}
				data = buf.Bytes()
			}
		}
	}
	return removeBOM(data), nil
}

// GetArweaveContentType returns the content-type from an Arweave transaction
func GetArweaveContentType(client *goar.Client, id string) (string, error) {
	data, err := client.GetTransactionTags(id)
	if err != nil {
		return "", err
	}

	for _, tag := range data {
		decodedName, err := base64.RawURLEncoding.DecodeString(tag.Name)
		if err != nil {
			return "", err
		}
		if strings.EqualFold(string(decodedName), "Content-Encoding") || strings.EqualFold(string(decodedName), "Content-Type") {
			decodedValue, err := base64.RawURLEncoding.DecodeString(tag.Value)
			if err != nil {
				return "", err
			}
			return string(decodedValue), nil
		}
	}
	return "", nil
}

func padHex(pHex string, pLength int) string {
	for len(pHex) < pLength {
		pHex = "0" + pHex
	}
	return pHex
}

func (h ErrHTTP) Error() string {
	return fmt.Sprintf("HTTP Error Status - %d | URL - %s", h.Status, h.URL)
}
