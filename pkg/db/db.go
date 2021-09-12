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
type Options struct {
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

func NewDB(o *Options) (db *sql.DB, err error) {
	ds := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", o.User, o.Password, o.Host, o.Port, o.Name)
	fullDSN := fmt.Sprintf("%s?multiStatements=true", ds)
	if len(o.Options) > 0 {
		fullDSN = fmt.Sprintf("%s&%s", fullDSN, o.Options)
	}

	db, err = sql.Open(mysqlDriver, fullDSN)
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err == nil {
		db.SetMaxIdleConns(o.Threads + 1)
		return db, nil
	}

	if !strings.Contains(err.Error(), unknownDB) {
		return nil, err
	}

	tmpDB, _ := sql.Open(mysqlDriver, ds)
	defer tmpDB.Close()

	query := createDBDDL + o.Name
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
