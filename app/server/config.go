package server

import (
	"log"
	"medods-auth/persistance/postgres"
	"os"
	"strconv"
	"sync"
)

var once sync.Once
var conf ServerConfig

const defaultPort = "8080"

type ServerConfig struct {
	Port       string
	HashSecret []byte
	Postgres   *postgres.PostgresConfig
}

func Config() ServerConfig {
	once.Do(func() {
		port := os.Getenv("JWT_PORT")
		if port == "" {
			log.Printf("JWT_PORT env variable not specified, defaulting to %s", defaultPort)
			port = defaultPort
		}
		_, err := strconv.Atoi(port)
		if err != nil {
			log.Printf("failed to parse port: %s, defaulting to %s", port, defaultPort)
			port = defaultPort
		}
		conf = ServerConfig{
			Port:     port,
			Postgres: &postgres.PostgresConfig{},
		}
		conf.HashSecret = []byte(os.Getenv("HASH_SECRET"))

		// TODO: Validation
		conf.Postgres.Host = os.Getenv("POSTGRES_HOST")
		conf.Postgres.Port = os.Getenv("POSTGRES_PORT")
		conf.Postgres.User = os.Getenv("POSTGRES_USER")
		conf.Postgres.Password = os.Getenv("POSTGRES_USER")
		conf.Postgres.Name = os.Getenv("POSTGRES_DB")
		conf.Postgres.HashDatabase = true
		conf.Postgres.BlackListDatabase = true
		if os.Getenv("JWT_SERVER_MODE") == "test" {
			conf.Postgres.SkipSSL = true
		}
	})
	return conf
}
