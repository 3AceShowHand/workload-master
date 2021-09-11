package cmd

import "github.com/spf13/cobra"

func RegisterLargeTxn(root *cobra.Command) {
	cmd := &cobra.Command{
		Use: "large_txn",
		Run: func(cmd *cobra.Command, args []string) {

		},
	}
	root.AddCommand(cmd)
}
