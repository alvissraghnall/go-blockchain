package cli

import (
    "fmt"

    "github.com/spf13/cobra"
)

var startNodeCmd = &cobra.Command{
    Use:   "startnode",
    Short: "Start the full node",
    Run: func(cmd *cobra.Command, args []string) {
        fmt.Println("Starting full node...")
        // TODO: Add logic to start the node
    },
}

func init() {
    rootCmd.AddCommand(startNodeCmd)
}
