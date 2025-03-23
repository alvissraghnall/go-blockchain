package cli

import (
    "fmt"
    "os"

    "github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
    Use:   "gochain",
    Short: "A CLI for interacting with gochain",
    Long:  `A CLI for managing the GOCHAIN blockchain network, including full nodes and miners.`,
}

func Execute() {
    if err := rootCmd.Execute(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
