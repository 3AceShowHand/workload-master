package db

import (
	"database/sql"
	"fmt"
	"strings"
)

// Config is the config for database
type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	Name     string

	// connection params
	Params []string
}

const (
	unknownDB   = "Unknown database"
	createDBDDL = "CREATE DATABASE IF NOT EXISTS "
	mysqlDriver = "mysql"
)

func NewDB(cfg *Config) *sql.DB {

}

func CloseDB(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}

func openDB() {
	// TODO: support other drivers
	var (
		tmpDB *sql.DB
		err   error
		ds    = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", user, password, host, port, dbName)
	)
	// allow multiple statements in one query to allow q15 on the TPC-H
	fullDsn := fmt.Sprintf("%s?multiStatements=true", ds)
	if len(connParams) > 0 {
		fullDsn = fmt.Sprintf("%s&%s", fullDsn, connParams)
	}
	globalDB, err = sql.Open(mysqlDriver, fullDsn)
	if err != nil {
		panic(err)
	}
	if err := globalDB.Ping(); err != nil {
		errString := err.Error()
		if strings.Contains(errString, unknownDB) {
			tmpDB, _ = sql.Open(mysqlDriver, ds)
			defer tmpDB.Close()
			if _, err := tmpDB.Exec(createDBDDL + dbName); err != nil {
				panic(fmt.Errorf("failed to create database, err %v\n", err))
			}
		} else {
			globalDB = nil
		}
	} else {
		globalDB.SetMaxIdleConns(threads + acThreads + 1)
	}
}

// func closeDB() {
// 	if globalDB != nil {
// 		globalDB.Close()
// 	}
// 	globalDB = nil
// }
