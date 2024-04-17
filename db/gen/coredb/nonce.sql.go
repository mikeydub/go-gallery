// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.26.0
// source: nonce.sql

package coredb

import (
	"context"

	"github.com/mikeydub/go-gallery/service/persist"
)

const consumeNonce = `-- name: ConsumeNonce :one
update nonces
    set
        consumed = true
    where
        value = $1
        and not consumed
        and nonces.created_at > (now() - interval '1 hour')
    returning id, value, created_at, consumed
`

func (q *Queries) ConsumeNonce(ctx context.Context, value string) (Nonce, error) {
	row := q.db.QueryRow(ctx, consumeNonce, value)
	var i Nonce
	err := row.Scan(
		&i.ID,
		&i.Value,
		&i.CreatedAt,
		&i.Consumed,
	)
	return i, err
}

const insertNonce = `-- name: InsertNonce :one
insert into nonces (id, value) values ($1, $2)
    on conflict (value)
        do nothing
    returning id, value, created_at, consumed
`

type InsertNonceParams struct {
	ID    persist.DBID `db:"id" json:"id"`
	Value string       `db:"value" json:"value"`
}

func (q *Queries) InsertNonce(ctx context.Context, arg InsertNonceParams) (Nonce, error) {
	row := q.db.QueryRow(ctx, insertNonce, arg.ID, arg.Value)
	var i Nonce
	err := row.Scan(
		&i.ID,
		&i.Value,
		&i.CreatedAt,
		&i.Consumed,
	)
	return i, err
}
