package testutil

import (
	"context"
	"database/sql"
	"fmt"
	"medods-auth/service/auth"
	"medods-auth/token"
	"medods-auth/user"

	"time"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type testRepo struct {
	db *sqlx.DB
}

func (tr *testRepo) Close() error {
	return tr.db.Close()
}

var schemaToken = `CREATE TABLE token (
    jti BLOB PRIMARY KEY,
    user_id BLOB NOT NULL,
    user_agent TEXT,
    hash BLOB NOT NULL,
    created_at TIMESTAMP NOT NULL
);`

var schemaBlacklist = `CREATE TABLE blacklist (
    jti BLOB PRIMARY KEY,
	created_at TIMESTAMP NOT NULL
);`

func (tr *testRepo) init() {
	var err error
	tr.db, err = sqlx.Connect("sqlite3", ":memory:")
	if err != nil {
		panic(err)
	}
	_, err = tr.db.Exec(schemaToken)
	if err != nil {
		panic(err)
	}
	_, err = tr.db.Exec(schemaBlacklist)
	if err != nil {
		panic(err)
	}
}

func NewTestInmemoryRepo() *testRepo {
	t := &testRepo{}
	t.init()
	return t
}

type TestDBRecord struct {
	JTI       uuid.UUID `db:"jti"`
	UserID    uuid.UUID `db:"user_id"`
	UserAgent string    `db:"user_agent"`
	Hash      []byte    `db:"hash"`
	CreatedAt time.Time `db:"created_at"`
}

func dbRecordFromAuthRecord(in auth.RefreshTokenRecord) *TestDBRecord {
	r := &TestDBRecord{}
	r.JTI = in.JTI
	r.UserID = in.User.Id
	r.UserAgent = in.User.UserAgent
	r.Hash = in.Hash
	r.CreatedAt = in.CreatedAt
	return r
}

func (r *TestDBRecord) toAuthRecord() *auth.RefreshTokenRecord {
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

func (r *testRepo) Store(ctx context.Context, rec *auth.RefreshTokenRecord) error {
	_, err := r.db.NamedExecContext(ctx,
		"INSERT INTO token (jti, user_id, user_agent, hash, created_at) VALUES (:jti, :user_id, :user_agent, :hash, :created_at)",
		dbRecordFromAuthRecord(*rec),
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *testRepo) Get(ctx context.Context, jti token.JTI) (*auth.RefreshTokenRecord, error) {
	row := r.db.QueryRowxContext(
		ctx,
		"SELECT (jti, user_id, user_agent, hash, created_at) FROM token WHERE jti = $1",
		jti,
	)

	var record TestDBRecord
	err := row.StructScan(&record)
	if err != nil {
		return nil, err
	}

	return record.toAuthRecord(), nil
}

func (r *testRepo) DeleteByUserId(ctx context.Context, userId uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, "DELETE FROM token WHERE user_id = $1", userId)
	if err != nil {
		return err
	}
	return nil
}

// type TokenBlackList interface {
// 	Add(context.Context, token.JTI) error
// 	Contains(context.Context, token.JTI) (bool, error)
// }

func (r *testRepo) Add(ctx context.Context, t token.JTI) error {
	_, err := r.db.ExecContext(
		ctx,
		"INSERT INTO blacklist (jti, created_at) VALUES ($1, $2)",
		t, time.Now(),
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *testRepo) Contains(ctx context.Context, jti token.JTI) (bool, error) {
	row := r.db.QueryRowxContext(ctx, "SELECT jti FROM blacklist WHERE jti = $1", jti)
	var gotId token.JTI
	err := row.Scan(&gotId)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil
		}
		return false, err
	}
	if jti.String() == gotId.String() {
		return true, nil
	}
	return false, nil
}

func (r *testRepo) DumpContents() error {
	rows, err := r.db.Queryx("SELECT * FROM token")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var row TestDBRecord
		rows.StructScan(&row)
		fmt.Printf("%+v\n", row)
	}

	rows, err = r.db.Queryx("SELECT * FROM blacklist")
	if err != nil {
		return err
	}
	defer rows.Close()
	for rows.Next() {
		var jti, created string
		rows.Scan(&jti, &created)
		fmt.Println(jti, created)
	}
	return nil
}
