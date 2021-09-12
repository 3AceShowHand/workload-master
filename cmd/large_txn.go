package cmd

import "github.com/spf13/cobra"

func NewCmdLargeTxn() *cobra.Command {
	cmd := &cobra.Command{
		Use: "large_txn",
		Run: func(cmd *cobra.Command, args []string) {
			runLargeTxn()
		},
	}

	return cmd
}

func runLargeTxn() {

}
