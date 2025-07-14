package postgres

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

var schemaToken = `CREATE TABLE IF NOT EXISTS token (
    jti UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    user_agent TEXT,
    hash BYTEA NOT NULL,
    created_at TIMESTAMP NOT NULL
);`

var schemaBlacklist = `CREATE TABLE IF NOT EXISTS blacklist (
    jti UUID PRIMARY KEY,
    created_at TIMESTAMP NOT NULL
);`

type PostgresConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string

	MaxOpenConns    *int
	MaxIdleConns    *int
	ConnMaxLifetime *time.Duration

	HashDatabase      bool
	BlackListDatabase bool

	SkipSSL bool
}

func InitDatabase(conf *PostgresConfig) (*sqlx.DB, error) {
	db, err := connect(conf)
	if err != nil {
		return nil, err
	}

	if !conf.HashDatabase && !conf.BlackListDatabase {
		return db, nil
	}

	tx, err := db.Beginx()
	if err != nil {
		return nil, err
	}
	if conf.HashDatabase {
		_, err = tx.ExecContext(context.TODO(), schemaToken)
		if err != nil {
			db.Close()
			return nil, err
		}
	}
	if conf.BlackListDatabase {
		_, err = tx.ExecContext(context.TODO(), schemaBlacklist)
		if err != nil {
			db.Close()
			return nil, err
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func connect(conf *PostgresConfig) (*sqlx.DB, error) {
	err := conf.validate()
	if err != nil {
		return nil, err
	}
	connstr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		conf.Host, conf.Port, conf.User, conf.Password, conf.Name)

	if conf.SkipSSL {
		connstr += " sslmode=disable"
	}

	fmt.Println(connstr)

	db, err := sqlx.Connect("postgres", connstr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	if conf.MaxOpenConns != nil {
		db.SetMaxOpenConns(*conf.MaxOpenConns)
	} else {
		db.SetMaxOpenConns(25)
	}
	if conf.MaxIdleConns != nil {
		db.SetMaxIdleConns(*conf.MaxOpenConns)
	} else {
		db.SetMaxIdleConns(25)
	}
	if conf.MaxIdleConns != nil {
		db.SetConnMaxLifetime(*conf.ConnMaxLifetime)
	} else {
		db.SetConnMaxLifetime(5 * time.Minute)
	}

	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

func (c *PostgresConfig) validate() error {
	if c.Host == "" ||
		c.Port == "" ||
		c.User == "" ||
		c.Password == "" ||
		c.Name == "" {
		res := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
			c.Host, c.Port, c.User, c.Password, c.Name)
		return errors.New("invalid postgres connection string:\n" + res)
	}
	return nil
}
