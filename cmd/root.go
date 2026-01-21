package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	appErrors "github.com/tmustier/economist-cli/internal/errors"
)

var (
	debugMode bool
)

var rootCmd = &cobra.Command{
	Use:           "economist",
	Short:         "CLI tool to read The Economist",
	Long:          `A command-line interface to browse and read articles from The Economist.`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		if appErrors.IsUserError(err) {
			fmt.Fprintln(os.Stderr, err)
		} else {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolVar(&debugMode, "debug", false, "Enable debug output")
	rootCmd.AddCommand(headlinesCmd)
	rootCmd.AddCommand(readCmd)
	rootCmd.AddCommand(loginCmd)
	rootCmd.AddCommand(sectionsCmd)
}
