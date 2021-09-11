package cmd

import "github.com/spf13/cobra"

func RegisterSQLCoverage(root *cobra.Command) {
	cmd := &cobra.Command{
		Use: "sql_coverage",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	root.AddCommand(cmd)
}
