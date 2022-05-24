package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/mikeydub/go-gallery/service/persist"
)

// WalletRepository is a repository for wallets
type WalletRepository struct {
	db *sql.DB

	insertStmt            *sql.Stmt
	getByIDStmt           *sql.Stmt
	getByChainAddressStmt *sql.Stmt
}

// NewWalletRepository creates a new postgres repository for interacting with wallets
func NewWalletRepository(db *sql.DB) *WalletRepository {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	insertStmt, err := db.PrepareContext(ctx, `INSERT INTO wallets (ID,VERSION,ADDRESS,CHAIN,WALLET_TYPE) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (ADDRESS,CHAIN) DO NOTHING;`)
	checkNoErr(err)

	getByIDStmt, err := db.PrepareContext(ctx, `SELECT ID,VERSION,CREATED_AT,LAST_UPDATED,ADDRESS,WALLET_TYPE,CHAIN FROM wallets WHERE ID = $1 AND DELETED = FALSE;`)
	checkNoErr(err)

	getByChainAddressStmt, err := db.PrepareContext(ctx, `SELECT ID,VERSION,CREATED_AT,LAST_UPDATED,ADDRESS,WALLET_TYPE,CHAIN FROM wallets WHERE ADDRESS = $1 AND CHAIN = $2 AND DELETED = FALSE;`)
	checkNoErr(err)

	return &WalletRepository{
		db:                    db,
		getByIDStmt:           getByIDStmt,
		getByChainAddressStmt: getByChainAddressStmt,
		insertStmt:            insertStmt,
	}
}

// GetByID returns a wallet by its ID
func (w *WalletRepository) GetByID(ctx context.Context, ID persist.DBID) (persist.Wallet, error) {
	var wallet persist.Wallet
	err := w.getByIDStmt.QueryRowContext(ctx, ID).Scan(&wallet.ID, &wallet.Version, &wallet.CreationTime, &wallet.LastUpdated, &wallet.Address, &wallet.WalletType, &wallet.Chain)
	if err != nil {
		if err == sql.ErrNoRows {
			return wallet, persist.ErrWalletNotFoundByID{WalletID: ID}
		}
		return wallet, err
	}
	return wallet, nil

}

// GetByChainAddress returns a wallet by address and chain
func (w *WalletRepository) GetByChainAddress(ctx context.Context, chainAddress persist.ChainAddress) (persist.Wallet, error) {
	var wallet persist.Wallet
	err := w.getByChainAddressStmt.QueryRowContext(ctx, chainAddress.Address, chainAddress.Chain).Scan(&wallet.ID, &wallet.Version, &wallet.CreationTime, &wallet.LastUpdated, &wallet.Address, &wallet.WalletType, &wallet.Chain)
	if err != nil {
		if err == sql.ErrNoRows {
			return wallet, persist.ErrWalletNotFoundByChainAddress{ChainAddress: chainAddress}
		}
		return wallet, err
	}
	return wallet, nil

}

// Insert inserts a wallet by its address and chain
func (w *WalletRepository) Insert(ctx context.Context, chainAddress persist.ChainAddress, walletType persist.WalletType) (persist.DBID, error) {

	_, err := w.insertStmt.ExecContext(ctx, persist.GenerateID(), 0, chainAddress.Address, chainAddress.Chain, walletType)
	if err != nil {
		return "", err
	}

	// rather than using the ID generated above, we must retrieve it because in the case of conflict the ID above would be inaccurate.
	wa, err := w.GetByChainAddress(ctx, chainAddress)
	if err != nil {
		return "", err
	}

	return wa.ID, nil
}
