package providers

// EtherscanNormalTx represents a normal ETH transfer response from Etherscan
type EtherscanNormalTx struct {
	BlockNumber      string `json:"blockNumber"`
	TimeStamp        string `json:"timeStamp"`
	Hash             string `json:"hash"`
	Nonce            string `json:"nonce"`
	BlockHash        string `json:"blockHash"`
	TransactionIndex string `json:"transactionIndex"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	Gas              string `json:"gas"`
	GasPrice         string `json:"gasPrice"`
	IsError          string `json:"isError"`
	TxReceiptStatus  string `json:"txreceipt_status"`
	Input            string `json:"input"`
	ContractAddress  string `json:"contractAddress"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	GasUsed          string `json:"gasUsed"`
	Confirmations    string `json:"confirmations"`
	MethodId         string `json:"methodId"`
	FunctionName     string `json:"functionName"`
}

// EtherscanInternalTx represents an internal transaction response from Etherscan
type EtherscanInternalTx struct {
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
	IsError         string `json:"isError"`
	ErrCode         string `json:"errCode"`
}

// EtherscanTokenTx represents a token transfer response from Etherscan (ERC-20/721/1155)
type EtherscanTokenTx struct {
	BlockNumber       string `json:"blockNumber"`
	TimeStamp         string `json:"timeStamp"`
	Hash              string `json:"hash"`
	Nonce             string `json:"nonce"`
	BlockHash         string `json:"blockHash"`
	From              string `json:"from"`
	ContractAddress   string `json:"contractAddress"`
	To                string `json:"to"`
	Value             string `json:"value"`
	TokenName         string `json:"tokenName"`
	TokenSymbol       string `json:"tokenSymbol"`
	TokenDecimal      string `json:"tokenDecimal"`
	TransactionIndex  string `json:"transactionIndex"`
	Gas               string `json:"gas"`
	GasPrice          string `json:"gasPrice"`
	GasUsed           string `json:"gasUsed"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	Input             string `json:"input"`
	Confirmations     string `json:"confirmations"`
	IsError           string `json:"isError"`
	TxReceiptStatus   string `json:"txreceipt_status"`
	TokenID           string `json:"tokenID"`   // For NFTs (ERC-721, ERC-1155)
	TokenValue        string `json:"tokenValue"` // For ERC-1155
}

// EtherscanResponse is the common response wrapper
type EtherscanResponse struct {
	Status  string          `json:"status"`
	Message string          `json:"message"`
	Result  []interface{}   `json:"result"`
}

// NormalTxResponse wraps Etherscan normal transaction results
type NormalTxResponse struct {
	Status  string                `json:"status"`
	Message string                `json:"message"`
	Result  []EtherscanNormalTx    `json:"result"`
}

// InternalTxResponse wraps Etherscan internal transaction results
type InternalTxResponse struct {
	Status  string                 `json:"status"`
	Message string                 `json:"message"`
	Result  []EtherscanInternalTx   `json:"result"`
}

// TokenTxResponse wraps Etherscan token transfer results
type TokenTxResponse struct {
	Status  string              `json:"status"`
	Message string              `json:"message"`
	Result  []EtherscanTokenTx   `json:"result"`
}
