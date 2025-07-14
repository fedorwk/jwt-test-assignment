package postgres

import (
	"context"
	"medods-auth/service/auth"
	"medods-auth/token"
	"medods-auth/user"
	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type HashRepository struct {
	db *sqlx.DB
}

func NewHashRepository(db *sqlx.DB) *HashRepository {
	return &HashRepository{
		db,
	}
}

type TokenDBRecord struct {
	JTI       uuid.UUID `db:"jti"`
	UserID    uuid.UUID `db:"user_id"`
	UserAgent string    `db:"user_agent"`
	Hash      []byte    `db:"hash"`
	CreatedAt time.Time `db:"created_at"`
}

func dbRecordFromAuthRecord(in auth.RefreshTokenRecord) *TokenDBRecord {
	r := &TokenDBRecord{}
	r.JTI = in.JTI
	r.UserID = in.User.Id
	r.UserAgent = in.User.UserAgent
	r.Hash = in.Hash
	r.CreatedAt = in.CreatedAt
	return r
}

func (r *TokenDBRecord) toAuthRecord() *auth.RefreshTokenRecord {
	out := &auth.RefreshTokenRecord{
		JTI: r.JTI,
		User: user.User{
			Id:        r.UserID,
			UserAgent: r.UserAgent,
		},
		Hash:      r.Hash,
		CreatedAt: r.CreatedAt,
	}
	return out
}

func (r *HashRepository) Store(ctx context.Context, rec *auth.RefreshTokenRecord) error {
	_, err := r.db.NamedExecContext(ctx,
		"INSERT INTO token (jti, user_id, user_agent, hash, created_at) VALUES (:jti, :user_id, :user_agent, :hash, :created_at)",
		dbRecordFromAuthRecord(*rec),
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *HashRepository) Get(ctx context.Context, jti token.JTI) (*auth.RefreshTokenRecord, error) {
	row := r.db.QueryRowxContext(
		ctx,
		"SELECT (jti, user_id, user_agent, hash, created_at) FROM token WHERE jti = $1",
		jti,
	)

	var record TokenDBRecord
	err := row.StructScan(&record)
	if err != nil {
		return nil, err
	}

	return record.toAuthRecord(), nil
}

func (r *HashRepository) DeleteByUserId(ctx context.Context, userId uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM token WHERE user_id = $1", userId)
	if err != nil {
		return err
	}
	return nil
}
