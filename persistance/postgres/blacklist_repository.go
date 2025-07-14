package postgres

import (
	"context"
	"database/sql"
	"medods-auth/token"
	"time"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type BlacklistRepository struct {
	db *sqlx.DB
}

func NewBlackListRepository(db *sqlx.DB) *BlacklistRepository {
	return &BlacklistRepository{
		db,
	}
}

func (repo *BlacklistRepository) Add(ctx context.Context, jti token.JTI) error {
	_, err := repo.db.ExecContext(
		ctx,
		"INSERT INTO blacklist (jti, created_at) VALUES ($1, $2)",
		jti, time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (repo *BlacklistRepository) Contains(ctx context.Context, jti token.JTI) (bool, error) {
	row := repo.db.QueryRowxContext(ctx, "SELECT jti FROM blacklist WHERE jti = $1", jti)
	var gotId token.JTI
	err := row.Scan(&gotId)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	return true, nil
}
