package db

import (
	"database/sql"
	"fmt"
	"iziumov/tv-v-x/config"

	_ "github.com/lib/pq"
)

type DB struct {
	*sql.DB
}

func NewDB(conf config.DBConfig) (*DB, error) {
	dsn := fmt.Sprintf(
		"postgres://%s:%s@%s:%d/%s?sslmode=disable",
		conf.User,
		conf.Password,
		conf.Host,
		conf.Port,
		conf.Name,
	)

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("Failed to open db: %v", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("Failed to ping db: %v", err)
	}

	return &DB{
		db,
	}, nil
}
