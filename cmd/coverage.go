package cmd

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/3AceShowHand/workload-master/pkg/db"
	"github.com/pingcap/log"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

func NewCmdSQLCoverage() *cobra.Command {
	cmd := &cobra.Command{
		Use: "sql_coverage",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := executeSQLCoverage(); err != nil {
				log.Error("execute SQL Coverage failed", zap.Error(err))
				return err
			}
			return nil
		},
	}

	return cmd
}

func executeSQLCoverage() error {
	database, err := db.NewDB(&db.Options{
		Host:     "",
		Port:     0,
		User:     "",
		Password: "",
		Name:     "",
		Threads:  0,
		Options:  nil,
	})
	if err != nil {
		return err
	}
	task := &SQLCoverageTask{
		db:      database,
		dbName:  "",
		threads: 0,
	}

	ctx, cancel := context.WithTimeout(globalCtx, time.Duration())
	defer cancel()

}

var (
	sqlFormat = map[string]string{
		"insert":  "",
		"delete":  "",
		"update":  "",
		"replace": "",
	}
)

type SQLCoverageTask struct {
	db      *sql.DB
	dbName  string
	threads int
}

func (t *SQLCoverageTask) Name() string {
	return fmt.Sprintf("SQL Coverage Task")
}

func (t *SQLCoverageTask) CleanUp(ctx context.Context, threadID int) error {
	return nil
}

func (t *SQLCoverageTask) Prepare(ctx context.Context, threadID int) error {
	return nil
}

func (t *SQLCoverageTask) Run(ctx context.Context, threadID int) error {
	return nil
}

func (t *SQLCoverageTask) DBName() string {
	return t.dbName
}
