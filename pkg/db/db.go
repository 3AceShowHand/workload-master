package db

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/pingcap/log"
	"go.uber.org/zap"

	// mysql package
	_ "github.com/go-sql-driver/mysql"
)

// Config is the config for database
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string

	Threads int
	// connection params
	Options []string
}

const (
	unknownDB   = "Unknown database"
	createDBDDL = "CREATE DATABASE IF NOT EXISTS "
	mysqlDriver = "mysql"
)

func NewDB(cfg *Config) (db *sql.DB, err error) {
	ds := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)
	fullDSN := fmt.Sprintf("%s?multiStatements=true", ds)
	if len(cfg.Options) > 0 {
		fullDSN = fmt.Sprintf("%s&%s", fullDSN, cfg.Options)
	}

	db, err = sql.Open(mysqlDriver, fullDSN)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err == nil {
		db.SetMaxIdleConns(cfg.Threads + 1)
		return db, nil
	}

	if !strings.Contains(err.Error(), unknownDB) {
		return nil, err
	}

	tmpDB, _ := sql.Open(mysqlDriver, ds)
	defer tmpDB.Close()

	query := createDBDDL + cfg.Name
	if _, err := tmpDB.Exec(query); err != nil {
		log.Info("failed to create database", zap.Error(err), zap.String("query", query))
		return nil, err
	}

	return db, nil
}

func CloseDB(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}
