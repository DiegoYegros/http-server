package database

import (
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb" // SQL Server driver
	"httpserver/config"
)

var DB *sql.DB

func InitDB(cfg *config.Config) error {
	if cfg.Database.Driver == "" || cfg.Database.ConnectionString == "" {
		return nil
	}
	var err error
	DB, err = sql.Open(cfg.Database.Driver, cfg.Database.ConnectionString)
	if err != nil {
		return fmt.Errorf("Couldn't open database: %w", err)
	}
	if err = DB.Ping(); err != nil {
		return fmt.Errorf("Couldn't connect to database: %w", err)
	}
	return nil
}

func CloseDB() {
	if DB != nil {
		DB.Close()
	}
}

func IsDBConnected() bool {
	return DB != nil && DB.Ping() == nil
}
