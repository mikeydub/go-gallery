package eth

import (
	"context"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/mikeydub/go-gallery/contracts"
	"github.com/mikeydub/go-gallery/persist"
)

// Client represents an abstraction over the geth ethereum client
type Client struct {
	EthClient       *ethclient.Client
	contractAddress common.Address
}

// NewEthClient initializes a new eth client
func NewEthClient(ethClient *ethclient.Client, contractAddress string) *Client {
	return &Client{
		EthClient:       ethClient,
		contractAddress: common.HexToAddress(contractAddress),
	}
}

const ensContractAddress = "0xFaC7BEA255a6990f749363002136aF6556b31e04"

// HasNFT checks if a wallet address has a given NFT
func (c *Client) HasNFT(pCtx context.Context, id persist.TokenID, userAddr persist.Address) (bool, error) {

	instance, err := contracts.NewIERC1155Caller(c.contractAddress, c.EthClient)
	if err != nil {
		return false, err
	}

	call, err := instance.BalanceOf(&bind.CallOpts{Context: pCtx}, userAddr.Address(), id.BigInt())
	if err != nil {
		return false, err
	}

	return call.Cmp(big.NewInt(0)) > 0, nil

}

// HasNFTs checks if a wallet address has a given set of NFTs
func (c *Client) HasNFTs(pCtx context.Context, ids []persist.TokenID, userAddr persist.Address) (bool, error) {

	instance, err := contracts.NewIERC1155Caller(c.contractAddress, c.EthClient)
	if err != nil {
		return false, err
	}

	bigIntIDs := make([]*big.Int, len(ids))
	addrs := make([]common.Address, len(ids))
	for i := 0; i < len(ids); i++ {

		bigIntIDs[i] = ids[i].BigInt()
		addrs[i] = userAddr.Address()
	}

	call, err := instance.BalanceOfBatch(&bind.CallOpts{Context: pCtx}, addrs, bigIntIDs)
	if err != nil {
		return false, err
	}
	for _, v := range call {
		if v.Cmp(big.NewInt(0)) > 0 {
			return true, nil
		}
	}

	return false, nil

}

// ResolvesENS checks if an ENS resolves to a given address
func (c *Client) ResolvesENS(pCtx context.Context, ens string, userAddr persist.Address) (bool, error) {

	instance, err := contracts.NewIENSCaller(c.contractAddress, c.EthClient)
	if err != nil {
		return false, err
	}

	nh := namehash(ens)
	asBytes32 := [32]byte{}
	for i := 0; i < len(nh); i++ {
		asBytes32[i] = nh[i]
	}

	call, err := instance.Resolver(&bind.CallOpts{Context: pCtx}, asBytes32)
	if err != nil {
		return false, err
	}

	return strings.EqualFold(userAddr.String(), call.String()), nil

}

// function that computes the namehash for a given ENS domain
func namehash(name string) common.Hash {
	node := common.Hash{}

	if len(name) > 0 {
		labels := strings.Split(name, ".")

		for i := len(labels) - 1; i >= 0; i-- {
			labelSha := crypto.Keccak256Hash([]byte(labels[i]))
			node = crypto.Keccak256Hash(node.Bytes(), labelSha.Bytes())
		}
	}

	return node
}
