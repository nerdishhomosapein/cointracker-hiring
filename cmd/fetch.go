package cmd

import (
	"conintracker-hiring/pkg/output"
	"conintracker-hiring/pkg/providers"
	"context"
	"fmt"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

var (
	address    string
	outputFile string
	startPage  int
	endPage    int
	provider   string
)

// fetchCmd represents the fetch command
var fetchCmd = &cobra.Command{
	Use:   "fetch",
	Short: "Fetch transaction history for an Ethereum wallet address",
	Long:  `Fetches all transactions (normal, internal, token transfers) for a given Ethereum address and exports to CSV.`,
	RunE:  runFetch,
}

func init() {
	rootCmd.AddCommand(fetchCmd)

	// Command-specific flags
	fetchCmd.Flags().StringVarP(&address, "address", "a", "", "Ethereum wallet address (required)")
	fetchCmd.Flags().StringVarP(&outputFile, "output", "o", "", "Output CSV file path (default: transactions.csv)")
	fetchCmd.Flags().IntVar(&startPage, "start-page", 1, "Starting page for pagination")
	fetchCmd.Flags().IntVar(&endPage, "end-page", 1, "Ending page for pagination")
	fetchCmd.Flags().StringVarP(&provider, "provider", "p", "etherscan", "Data provider (currently only 'etherscan' supported)")

	// Mark required flags
	fetchCmd.MarkFlagRequired("address")
}

func runFetch(cmd *cobra.Command, args []string) error {
	// Validate address format
	if !isValidEthereumAddress(address) {
		return fmt.Errorf("invalid Ethereum address format: %s", address)
	}

	// Get API key from flag or environment variable
	etherscanKey := apiKey
	if etherscanKey == "" {
		etherscanKey = os.Getenv("ETHERSCAN_API_KEY")
	}
	if etherscanKey == "" {
		return fmt.Errorf("Etherscan API key is required (set via --api-key flag or ETHERSCAN_API_KEY env var)")
	}

	// Set default output file
	if outputFile == "" {
		outputFile = "transactions.csv"
	}

	// Create output file
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	// Print progress
	fmt.Printf("Fetching transactions for address: %s\n", address)
	fmt.Printf("Output file: %s\n\n", outputFile)

	// Create Etherscan client
	client := providers.NewEtherscanClient(providers.ClientConfig{
		APIKey: etherscanKey,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	})

	// Create normalizer and fetcher
	normalizer := providers.NewEtherscanNormalizer()
	fetcher := providers.NewTransactionFetcher(client, normalizer)

	// Fetch transactions
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	fmt.Println("Fetching transactions...")
	txs, err := fetcher.FetchAllTransactions(ctx, address, startPage, endPage)
	if err != nil {
		return fmt.Errorf("failed to fetch transactions: %w", err)
	}

	fmt.Printf("Found %d transactions\n", len(txs))

	if len(txs) == 0 {
		fmt.Println("No transactions found for this address")
		return nil
	}

	// Write to CSV
	fmt.Println("Writing to CSV...")
	csvWriter, err := output.NewCSVWriter(output.CSVConfig{Writer: file})
	if err != nil {
		return fmt.Errorf("failed to create CSV writer: %w", err)
	}

	if err := csvWriter.WriteTransactions(txs); err != nil {
		csvWriter.Close()
		return fmt.Errorf("failed to write transactions to CSV: %w", err)
	}

	if err := csvWriter.Close(); err != nil {
		return fmt.Errorf("failed to close CSV writer: %w", err)
	}

	// Print summary
	fmt.Println("\n✓ Successfully exported transactions to CSV")
	fmt.Printf("Total transactions: %d\n", len(txs))

	// Count by type
	typeCounts := make(map[string]int)
	for _, tx := range txs {
		typeCounts[string(tx.Type)]++
	}

	fmt.Println("\nTransaction breakdown:")
	for txType, count := range typeCounts {
		fmt.Printf("  %s: %d\n", txType, count)
	}

	return nil
}

// isValidEthereumAddress validates Ethereum address format
func isValidEthereumAddress(addr string) bool {
	// Ethereum addresses are 42 characters long (0x + 40 hex chars)
	if len(addr) != 42 {
		return false
	}

	if !strings.HasPrefix(addr, "0x") {
		return false
	}

	// Check if remaining 40 characters are valid hex
	validHex := regexp.MustCompile(`^0x[0-9a-fA-F]{40}$`)
	return validHex.MatchString(addr)
}
