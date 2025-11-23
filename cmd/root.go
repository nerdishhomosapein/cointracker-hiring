package cmd

import (
	"github.com/spf13/cobra"
)

var (
	version = "0.1.0"
	apiKey  string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "cointracker",
	Short:   "ETH transaction exporter - Export Ethereum wallet transactions to CSV",
	Long:    `Cointracker is a CLI tool that fetches transaction history for Ethereum wallet addresses and exports them to structured CSV files.`,
	Version: version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVar(&apiKey, "api-key", "", "Etherscan API key (can also be set via ETHERSCAN_API_KEY env var)")
}
