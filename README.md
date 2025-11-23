# Cointracker - Ethereum Transaction Exporter

A command-line tool to fetch transaction history for Ethereum wallet addresses and export them to structured CSV files.

## Features

- **Multi-transaction type support**: Normal ETH transfers, internal contract interactions, and token transfers (ERC-20, ERC-721, ERC-1155)
- **Etherscan integration**: Fetches data from Etherscan API with built-in rate limiting
- **CSV export**: Generates structured CSV files with all relevant transaction details
- **Comprehensive data**: Captures transaction hash, timestamp, from/to addresses, transaction type, asset details, amounts, and gas fees
- **TDD approach**: Fully tested with unit tests, integration tests, and test fixtures

## Installation

### Prerequisites
- Go 1.24.2 or higher

### Build

```bash
git clone <repository>
cd cointracker_assignment
go build -o cointracker main.go
```

## Usage

### Basic Command

```bash
./cointracker fetch \
  --address 0xa39b189482f984388a34460636fea9eb181ad1a6 \
  --api-key YOUR_ETHERSCAN_API_KEY \
  --output transactions.csv
```

### Using Environment Variables

```bash
export ETHERSCAN_API_KEY=YOUR_API_KEY
./cointracker fetch --address 0xa39b189482f984388a34460636fea9eb181ad1a6
```

### Options

```
Global Flags:
  --api-key string      Etherscan API key (can also be set via ETHERSCAN_API_KEY env var)

Fetch Command Flags:
  -a, --address string    Ethereum wallet address (required)
  -o, --output string     Output CSV file path (default: transactions.csv)
  -p, --provider string   Data provider (default: etherscan)
  --start-page int        Starting page for pagination (default: 1)
  --end-page int          Ending page for pagination (default: 1)
```

## CSV Output Format

The exported CSV file includes the following columns:

| Column | Description |
|--------|-------------|
| Transaction Hash | Unique transaction identifier |
| Date & Time | Transaction confirmation timestamp (RFC3339) |
| From Address | Sender's Ethereum address |
| To Address | Recipient's Ethereum address |
| Transaction Type | ETH, ERC-20, ERC-721, ERC-1155, or Internal |
| Asset Contract Address | Token/NFT contract address (if applicable) |
| Asset Symbol / Name | Token symbol or NFT collection name |
| Token ID | Unique identifier for NFTs |
| Value / Amount | Quantity transferred |
| Gas Fee (ETH) | Total transaction gas cost in ETH |

## Example Transactions

### Sample Ethereum Addresses

For testing and validation, you can use these addresses:

1. **0xa39b189482f984388a34460636fea9eb181ad1a6** - Standard address
2. **0xd620AADaBaA20d2af700853C4504028cba7C3333** - Address with token transfers
3. **0xfb50526f49894b78541b776f5aaefe43e3bd8590** - Large address (160,000+ transactions)

## Testing

Run the test suite:

```bash
# All tests
go test ./...

# Provider tests
go test ./pkg/providers -v

# Output tests
go test ./pkg/output -v

# Integration tests
go test ./pkg -v
```

## Architecture

### Packages

- **pkg/models**: Core transaction model and types
- **pkg/providers**: Etherscan API client and transaction fetcher
- **pkg/output**: CSV export functionality
- **cmd**: CLI commands and orchestration

### Data Flow

1. **Fetch**: CLI → EtherscanClient → HTTP requests to Etherscan API
2. **Normalize**: Raw API responses → EtherscanNormalizer → Normalized Transaction model
3. **Export**: Normalized transactions → CSVWriter → CSV file

## Rate Limiting

The tool includes built-in rate limiting to respect Etherscan API rate limits:
- Default rate limit delay: 200ms between requests
- Automatic retry on network errors
- Clear error messages for rate limit violations

## Error Handling

The tool provides clear error messages for:
- Invalid Ethereum address format
- Missing API key
- Network/API failures
- File I/O errors
- Invalid transactions (skipped gracefully)

## Limitations

- Currently supports Etherscan only (adapter interface for future providers)
- Pagination supports up to 10,000 transactions per page
- ETH amounts in wei, tokens with decimal precision handling
- Gas fees calculated from gasUsed × gasPrice

## Future Enhancements

- Support for additional providers (Alchemy, Blockscout, Infura)
- JSON/XLSX/PDF export formats
- Filtering by transaction type, date range, or amount
- Resume capability for interrupted exports
- Multi-wallet batch processing

## License

MIT

## Support

For issues or questions, please refer to the test files for usage examples and expected behavior.
