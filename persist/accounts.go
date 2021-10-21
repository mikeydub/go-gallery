package persist

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"runtime/debug"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

const accountCollName = "accounts"

// Address represents an Ethereum address
type Address string

// BlockNumber represents an Ethereum block number
type BlockNumber uint64

// Account represents an ethereum account in the database
type Account struct {
	Version      int64              `bson:"version"              json:"version"` // schema version for this model
	ID           DBID               `bson:"_id"                  json:"id" binding:"required"`
	CreationTime primitive.DateTime `bson:"created_at"        json:"created_at"`
	Deleted      bool               `bson:"deleted" json:"-"`
	LastUpdated  primitive.DateTime `bson:"last_updated" json:"last_updated"`

	Address         Address     `bson:"address" json:"address"`
	LastSyncedBlock BlockNumber `bson:"last_synced_block" json:"last_synced_block"`
}

// AccountRepository is the interface for interacting with the account persistence layer
type AccountRepository interface {
	GetByAddress(context.Context, Address) (*Account, error)
	UpsertByAddress(context.Context, Address, *Account) error
}

// ErrAccountNotFoundByAddress is an error that occurs when an account is not found by an address
type ErrAccountNotFoundByAddress struct {
	Address Address
}

func (e ErrAccountNotFoundByAddress) Error() string {
	return fmt.Sprintf("account not found by address: %s", e.Address)
}

func (a Address) String() string {
	return normalizeAddress(string(a))
}

// Lower returns the Ethereum address number as a lowercase hex string
func (a Address) Lower() Address {
	return Address(strings.ToLower(a.String()))
}

// Address returns the ethereum address byte array
func (a Address) Address() common.Address {
	return common.HexToAddress(a.String())
}

// Uint64 returns the ethereum block number as a uint64
func (b BlockNumber) Uint64() uint64 {
	return uint64(b)
}

// BigInt returns the ethereum block number as a big.Int
func (b BlockNumber) BigInt() *big.Int {
	return new(big.Int).SetUint64(b.Uint64())
}

func (b BlockNumber) String() string {
	return strings.ToLower(b.BigInt().String())
}

// Hex returns the ethereum block number as a hex string
func (b BlockNumber) Hex() string {
	return strings.ToLower(b.BigInt().Text(16))
}

func normalizeAddress(address string) string {
	if len(address) != 40 && len(address) != 42 {
		log.Printf("invalid address: %s len: %d\n", address, len(address))
		debug.PrintStack()
		return ""
	}
	withoutPrefix := strings.TrimPrefix(address, "0x")
	return "0x" + withoutPrefix[len(withoutPrefix)-40:]
}
