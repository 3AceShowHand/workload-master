package cmd

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/spf13/cobra"
)

func RegisterSQLCoverage(root *cobra.Command) {
	cmd := &cobra.Command{
		Use: "sql_coverage",
		Run: func(cmd *cobra.Command, args []string) {
			runSQLCoverage()
		},
	}
	root.AddCommand(cmd)
}

func runSQLCoverage() {

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
