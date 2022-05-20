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

	insertStmt              *sql.Stmt
	getByAddressDetailsStmt *sql.Stmt
}

// NewWalletRepository creates a new postgres repository for interacting with wallets
func NewWalletRepository(db *sql.DB) *WalletRepository {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	insertStmt, err := db.PrepareContext(ctx, `INSERT INTO wallets (ID,VERSION,ADDRESS,CHAIN,WALLET_TYPE) VALUES ($1,$2,$3,$4,$5) ON CONFLICT (ADDRESS,CHAIN) DO NOTHING;`)
	checkNoErr(err)

	getByAddressDetailsStmt, err := db.PrepareContext(ctx, `SELECT ID,VERSION,CREATED_AT,LAST_UPDATED,ADDRESS,WALLET_TYPE,CHAIN FROM wallets WHERE ADDRESS = $1 AND CHAIN = $2 AND DELTED = FALSE;`)
	checkNoErr(err)

	return &WalletRepository{
		db:                      db,
		getByAddressDetailsStmt: getByAddressDetailsStmt,
		insertStmt:              insertStmt,
	}
}

// GetByAddressDetails returns a wallet by address and chain
func (w *WalletRepository) GetByAddressDetails(ctx context.Context, addr persist.Address, chain persist.Chain) (persist.Wallet, error) {
	var wallet persist.Wallet
	err := w.getByAddressDetailsStmt.QueryRowContext(ctx, addr, chain).Scan(&wallet.ID, &wallet.Version, &wallet.CreationTime, &wallet.LastUpdated, &wallet.Address, &wallet.WalletType, &wallet.Chain)
	if err != nil {
		if err == sql.ErrNoRows {
			return wallet, persist.ErrWalletNotFoundByAddressDetails{Address: addr, Chain: chain}
		}
		return wallet, err
	}
	return wallet, nil

}

// Insert inserts a wallet by its address and chain
func (w *WalletRepository) Insert(ctx context.Context, addr persist.Address, chain persist.Chain, walletType persist.WalletType) (persist.DBID, error) {

	_, err := w.insertStmt.ExecContext(ctx, persist.GenerateID(), 0, addr, chain, walletType)
	if err != nil {
		return "", err
	}

	// rather than using the ID generated above, we must retrieve it because in the case of conflict the ID above would be innacurate.
	wa, err := w.GetByAddressDetails(ctx, addr, chain)
	if err != nil {
		return "", err
	}

	return wa.ID, nil
}
