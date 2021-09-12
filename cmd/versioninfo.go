package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func printVersion() {
	fmt.Println("Git Commit Hash:", commit)
	fmt.Println("UTC Build Time:", date)
	fmt.Println("Release version:", version)
}

func NewCmdVersion() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Output version information",
		Run: func(cmd *cobra.Command, args []string) {
			printVersion()
		},
	}
}
