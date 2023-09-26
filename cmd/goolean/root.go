package goolean

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "goolean",
	Short: "goolean - a simple CLI to solve boolean algebra",
	Long: `goolean is a CLI tool to solve boolean algebra

One can use goolean to solve boolean algebra straight from the terminal`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
