package goolean

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/thecatster/goolean/pkg/goolean"
)

var version = "0.0.1"

var rootCmd = &cobra.Command{
	Use:     "goolean",
	Version: version,
	Short:   "goolean - a simple CLI to solve boolean algebra",
	Long: `goolean is a CLI tool to solve boolean algebra

One can use goolean to solve boolean algebra straight from the terminal`,
	Run: func(cmd *cobra.Command, args []string) {
		goolean.Run()
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Whoops. There was an error while executing your CLI '%s'", err)
		os.Exit(1)
	}
}
