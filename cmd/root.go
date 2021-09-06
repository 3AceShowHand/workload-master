package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command {
	Use: "workload-master",
	Short: "workload-master master all different kind of workload.",
	Long: "workload-master can run different kind of database workload, such as go-tpc, go-ycsb, sysbench",

	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand()
}