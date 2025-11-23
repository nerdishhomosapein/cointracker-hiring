package providers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const (
	// Etherscan API base URL
	EtherscanBaseURL = "https://api.etherscan.io/api"
	
	// Default pagination
	DefaultPageSize = 10000
	DefaultStartBlock = 0
	DefaultEndBlock = 99999999
	
	// Rate limit delays (Etherscan free tier)
	RateLimitDelay = 200 * time.Millisecond
)

// EtherscanClient implements the Provider interface for Etherscan API
type EtherscanClient struct {
	apiKey     string
	httpClient *http.Client
	baseURL    string
	lastReq    time.Time // Track last request for rate limiting
}

// ClientConfig holds configuration for Etherscan client
type ClientConfig struct {
	APIKey      string
	HTTPClient  *http.Client
	BaseURL     string
	RateLimit   time.Duration
}

// NewEtherscanClient creates a new Etherscan API client
func NewEtherscanClient(cfg ClientConfig) *EtherscanClient {
	if cfg.HTTPClient == nil {
		cfg.HTTPClient = &http.Client{
			Timeout: 30 * time.Second,
		}
	}
	if cfg.BaseURL == "" {
		cfg.BaseURL = EtherscanBaseURL
	}
	
	return &EtherscanClient{
		apiKey:     cfg.APIKey,
		httpClient: cfg.HTTPClient,
		baseURL:    cfg.BaseURL,
		lastReq:    time.Now(),
	}
}

// executeRequest performs an HTTP request with rate limiting and error handling
func (c *EtherscanClient) executeRequest(ctx context.Context, params url.Values) (map[string]interface{}, error) {
	// Rate limiting: wait if necessary
	timeSinceLastReq := time.Since(c.lastReq)
	if timeSinceLastReq < RateLimitDelay {
		select {
		case <-time.After(RateLimitDelay - timeSinceLastReq):
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
	c.lastReq = time.Now()

	// Build URL
	u, _ := url.Parse(c.baseURL)
	u.RawQuery = params.Encode()

	// Create request
	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Execute request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Parse JSON
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	// Check for API errors
	if status, ok := result["status"].(string); ok {
		if status == "0" {
			if message, ok := result["message"].(string); ok {
				if message == "NOTOK" {
					if resultMsg, ok := result["result"].(string); ok {
										return nil, fmt.Errorf("etherscan error: %s", resultMsg)
									}
								}
							}
						}
					}

	return result, nil
}

// buildParams creates base query parameters for Etherscan API
func (c *EtherscanClient) buildParams(action, module string, address string) url.Values {
	params := url.Values{}
	params.Set("apikey", c.apiKey)
	params.Set("module", module)
	params.Set("action", action)
	params.Set("address", address)
	return params
}

// FetchNormalTransactions fetches normal ETH transfers from Etherscan
func (c *EtherscanClient) FetchNormalTransactions(ctx context.Context, address string, startPage, endPage int) ([]EtherscanNormalTx, error) {
	params := c.buildParams("txlist", "account", address)
	params.Set("startblock", strconv.Itoa(DefaultStartBlock))
	params.Set("endblock", strconv.Itoa(DefaultEndBlock))
	params.Set("page", strconv.Itoa(startPage))
	params.Set("offset", strconv.Itoa(endPage - startPage + 1))
	params.Set("sort", "asc")

	result, err := c.executeRequest(ctx, params)
	if err != nil {
		return nil, err
	}

	// Parse results
	var txs []EtherscanNormalTx
	if resultData, ok := result["result"].([]interface{}); ok {
		for _, item := range resultData {
			if itemMap, ok := item.(map[string]interface{}); ok {
				// Convert map to JSON and back to typed struct
				jsonData, _ := json.Marshal(itemMap)
				var tx EtherscanNormalTx
				if err := json.Unmarshal(jsonData, &tx); err == nil {
					txs = append(txs, tx)
				}
			}
		}
	}

	return txs, nil
}

// FetchInternalTransactions fetches internal contract interactions from Etherscan
func (c *EtherscanClient) FetchInternalTransactions(ctx context.Context, address string, startPage, endPage int) ([]EtherscanInternalTx, error) {
	params := c.buildParams("txlistinternal", "account", address)
	params.Set("startblock", strconv.Itoa(DefaultStartBlock))
	params.Set("endblock", strconv.Itoa(DefaultEndBlock))
	params.Set("page", strconv.Itoa(startPage))
	params.Set("offset", strconv.Itoa(endPage - startPage + 1))
	params.Set("sort", "asc")

	result, err := c.executeRequest(ctx, params)
	if err != nil {
		return nil, err
	}

	// Parse results
	var txs []EtherscanInternalTx
	if resultData, ok := result["result"].([]interface{}); ok {
		for _, item := range resultData {
			if itemMap, ok := item.(map[string]interface{}); ok {
				jsonData, _ := json.Marshal(itemMap)
				var tx EtherscanInternalTx
				if err := json.Unmarshal(jsonData, &tx); err == nil {
					txs = append(txs, tx)
				}
			}
		}
	}

	return txs, nil
}

// FetchTokenTransfers fetches ERC-20 token transfers from Etherscan
func (c *EtherscanClient) FetchTokenTransfers(ctx context.Context, address string, startPage, endPage int) ([]EtherscanTokenTx, error) {
	params := c.buildParams("tokentx", "account", address)
	params.Set("startblock", strconv.Itoa(DefaultStartBlock))
	params.Set("endblock", strconv.Itoa(DefaultEndBlock))
	params.Set("page", strconv.Itoa(startPage))
	params.Set("offset", strconv.Itoa(endPage - startPage + 1))
	params.Set("sort", "asc")

	result, err := c.executeRequest(ctx, params)
	if err != nil {
		return nil, err
	}

	// Parse results
	var txs []EtherscanTokenTx
	if resultData, ok := result["result"].([]interface{}); ok {
		for _, item := range resultData {
			if itemMap, ok := item.(map[string]interface{}); ok {
				jsonData, _ := json.Marshal(itemMap)
				var tx EtherscanTokenTx
				if err := json.Unmarshal(jsonData, &tx); err == nil {
					txs = append(txs, tx)
				}
			}
		}
	}

	return txs, nil
}

// FetchNFTTransfers fetches ERC-721 NFT transfers from Etherscan
func (c *EtherscanClient) FetchNFTTransfers(ctx context.Context, address string, startPage, endPage int) ([]EtherscanTokenTx, error) {
	params := c.buildParams("tokennfttx", "account", address)
	params.Set("startblock", strconv.Itoa(DefaultStartBlock))
	params.Set("endblock", strconv.Itoa(DefaultEndBlock))
	params.Set("page", strconv.Itoa(startPage))
	params.Set("offset", strconv.Itoa(endPage - startPage + 1))
	params.Set("sort", "asc")

	result, err := c.executeRequest(ctx, params)
	if err != nil {
		return nil, err
	}

	// Parse results
	var txs []EtherscanTokenTx
	if resultData, ok := result["result"].([]interface{}); ok {
		for _, item := range resultData {
			if itemMap, ok := item.(map[string]interface{}); ok {
				jsonData, _ := json.Marshal(itemMap)
				var tx EtherscanTokenTx
				if err := json.Unmarshal(jsonData, &tx); err == nil {
					txs = append(txs, tx)
				}
			}
		}
	}

	return txs, nil
}

// FetchERC1155Transfers fetches ERC-1155 multi-token transfers from Etherscan
func (c *EtherscanClient) FetchERC1155Transfers(ctx context.Context, address string, startPage, endPage int) ([]EtherscanTokenTx, error) {
	params := c.buildParams("token1155tx", "account", address)
	params.Set("startblock", strconv.Itoa(DefaultStartBlock))
	params.Set("endblock", strconv.Itoa(DefaultEndBlock))
	params.Set("page", strconv.Itoa(startPage))
	params.Set("offset", strconv.Itoa(endPage - startPage + 1))
	params.Set("sort", "asc")

	result, err := c.executeRequest(ctx, params)
	if err != nil {
		return nil, err
	}

	// Parse results
	var txs []EtherscanTokenTx
	if resultData, ok := result["result"].([]interface{}); ok {
		for _, item := range resultData {
			if itemMap, ok := item.(map[string]interface{}); ok {
				jsonData, _ := json.Marshal(itemMap)
				var tx EtherscanTokenTx
				if err := json.Unmarshal(jsonData, &tx); err == nil {
					txs = append(txs, tx)
				}
			}
		}
	}

	return txs, nil
}
