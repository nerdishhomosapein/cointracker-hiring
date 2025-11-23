package etherscan

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
)

// Client represents an Etherscan API client for testing purposes
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient creates a new Etherscan client
func NewClient(baseURL, apiKey string, httpClient *http.Client) *Client {
	if httpClient == nil {
		httpClient = &http.Client{}
	}
	return &Client{
		baseURL:    baseURL,
		apiKey:     apiKey,
		httpClient: httpClient,
	}
}

// Transaction represents a normal Ethereum transaction
type Transaction struct {
	BlockNumber      string `json:"blockNumber"`
	TimeStamp        string `json:"timeStamp"`
	Hash             string `json:"hash"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	GasPrice         string `json:"gasPrice"`
	GasUsed          string `json:"gasUsed"`
	Nonce            string `json:"nonce"`
	TransactionIndex string `json:"transactionIndex"`
	Input            string `json:"input"`
	ContractAddress  string `json:"contractAddress"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	Confirmations    string `json:"confirmations"`
}

// InternalTransaction represents an internal transaction
type InternalTransaction struct {
	BlockNumber     string `json:"blockNumber"`
	TimeStamp       string `json:"timeStamp"`
	Hash            string `json:"hash"`
	From            string `json:"from"`
	To              string `json:"to"`
	Value           string `json:"value"`
	ContractAddress string `json:"contractAddress"`
	Input           string `json:"input"`
	Type            string `json:"type"`
	Gas             string `json:"gas"`
	GasUsed         string `json:"gasUsed"`
	TraceId         string `json:"traceId"`
	ErrCode         string `json:"errCode"`
	IsError         string `json:"isError"`
}

// TokenTransaction represents an ERC-20 token transaction
type TokenTransaction struct {
	BlockNumber     string `json:"blockNumber"`
	TimeStamp       string `json:"timeStamp"`
	Hash            string `json:"hash"`
	From            string `json:"from"`
	To              string `json:"to"`
	Value           string `json:"value"`
	TokenName       string `json:"tokenName"`
	TokenSymbol     string `json:"tokenSymbol"`
	TokenDecimal    string `json:"tokenDecimal"`
	ContractAddress string `json:"contractAddress"`
	TokenID         string `json:"tokenID,omitempty"`
	TokenValue      string `json:"tokenValue,omitempty"`
	Gas             string `json:"gas"`
	GasPrice        string `json:"gasPrice"`
	GasUsed         string `json:"gasUsed"`
}

// EtherscanResponse represents the API response structure
type EtherscanResponse[T any] struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Result  []T    `json:"result"`
}

// makeRequest makes a generic request to Etherscan API
func (c *Client) makeRequest(ctx context.Context, params map[string]string, result interface{}) error {
	u, err := url.Parse(c.baseURL)
	if err != nil {
		return fmt.Errorf("invalid base URL: %w", err)
	}

	query := u.Query()
	for key, value := range params {
		query.Set(key, value)
	}
	query.Set("apikey", c.apiKey)
	u.RawQuery = query.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	// First, unmarshal to check status
	var baseResp struct {
		Status  string `json:"status"`
		Message string `json:"message"`
	}
	if err := json.Unmarshal(body, &baseResp); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if baseResp.Status != "1" {
		return fmt.Errorf("API error: %s", baseResp.Message)
	}

	// Now unmarshal the full response
	if err := json.Unmarshal(body, result); err != nil {
		return fmt.Errorf("failed to unmarshal result: %w", err)
	}

	return nil
}

// GetNormalTx fetches normal transactions for an address
func (c *Client) GetNormalTx(ctx context.Context, address string, page, offset int, sort string) ([]Transaction, error) {
	params := map[string]string{
		"module":  "account",
		"action":  "txlist",
		"address": address,
		"page":    strconv.Itoa(page),
		"offset":  strconv.Itoa(offset),
		"sort":    sort,
	}

	var resp EtherscanResponse[Transaction]
	if err := c.makeRequest(ctx, params, &resp); err != nil {
		return nil, err
	}

	return resp.Result, nil
}

// GetInternalTx fetches internal transactions for an address
func (c *Client) GetInternalTx(ctx context.Context, address string, page, offset int, sort string) ([]InternalTransaction, error) {
	params := map[string]string{
		"module":  "account",
		"action":  "txlistinternal",
		"address": address,
		"page":    strconv.Itoa(page),
		"offset":  strconv.Itoa(offset),
		"sort":    sort,
	}

	var resp EtherscanResponse[InternalTransaction]
	if err := c.makeRequest(ctx, params, &resp); err != nil {
		return nil, err
	}

	return resp.Result, nil
}

// GetTokenTx fetches ERC-20 token transactions for an address
func (c *Client) GetTokenTx(ctx context.Context, address string, page, offset int, sort string) ([]TokenTransaction, error) {
	params := map[string]string{
		"module":  "account",
		"action":  "tokentx",
		"address": address,
		"page":    strconv.Itoa(page),
		"offset":  strconv.Itoa(offset),
		"sort":    sort,
	}

	var resp EtherscanResponse[TokenTransaction]
	if err := c.makeRequest(ctx, params, &resp); err != nil {
		return nil, err
	}

	return resp.Result, nil
}

// GetERC721Tx fetches ERC-721 (NFT) transactions for an address
func (c *Client) GetERC721Tx(ctx context.Context, address string, page, offset int, sort string) ([]TokenTransaction, error) {
	params := map[string]string{
		"module":  "account",
		"action":  "tokennfttx",
		"address": address,
		"page":    strconv.Itoa(page),
		"offset":  strconv.Itoa(offset),
		"sort":    sort,
	}

	var resp EtherscanResponse[TokenTransaction]
	if err := c.makeRequest(ctx, params, &resp); err != nil {
		return nil, err
	}

	return resp.Result, nil
}

// GetERC1155Tx fetches ERC-1155 token transactions for an address
func (c *Client) GetERC1155Tx(ctx context.Context, address string, page, offset int, sort string) ([]TokenTransaction, error) {
	params := map[string]string{
		"module":  "account",
		"action":  "token1155tx",
		"address": address,
		"page":    strconv.Itoa(page),
		"offset":  strconv.Itoa(offset),
		"sort":    sort,
	}

	var resp EtherscanResponse[TokenTransaction]
	if err := c.makeRequest(ctx, params, &resp); err != nil {
		return nil, err
	}

	return resp.Result, nil
}