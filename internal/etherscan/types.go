package etherscan

// Types for internal testing and normalization

// NormalTx represents a normal transaction from Etherscan
type NormalTx struct {
	Hash             string `json:"hash"`
	BlockNumber      string `json:"blockNumber"`
	TimeStamp        string `json:"timeStamp"`
	From             string `json:"from"`
	To               string `json:"to"`
	Value            string `json:"value"`
	GasPrice         string `json:"gasPrice"`
	GasUsed          string `json:"gasUsed"`
	Nonce            string `json:"nonce"`
	TransactionIndex string `json:"transactionIndex"`
	ContractAddress  string `json:"contractAddress"`
}

// InternalTx represents an internal transaction from Etherscan
type InternalTx struct {
	Hash            string `json:"hash"`
	BlockNumber     string `json:"blockNumber"`
	TimeStamp       string `json:"timeStamp"`
	From            string `json:"from"`
	To              string `json:"to"`
	Value           string `json:"value"`
	ContractAddress string `json:"contractAddress"`
	Gas             string `json:"gas"`
	GasUsed         string `json:"gasUsed"`
	IsError         string `json:"isError"`
	Type            string `json:"type"`
	TraceID         string `json:"traceId"`
}

// TokenTx represents an ERC-20 token transaction from Etherscan
type TokenTx struct {
	Hash            string `json:"hash"`
	BlockNumber     string `json:"blockNumber"`
	TimeStamp       string `json:"timeStamp"`
	From            string `json:"from"`
	To              string `json:"to"`
	Value           string `json:"value"`
	TokenName       string `json:"tokenName"`
	TokenSymbol     string `json:"tokenSymbol"`
	TokenDecimal    string `json:"tokenDecimal"`
	ContractAddress string `json:"contractAddress"`
	GasPrice        string `json:"gasPrice"`
	GasUsed         string `json:"gasUsed"`
}

// ERC721Tx represents an ERC-721 NFT transaction from Etherscan
type ERC721Tx struct {
	Hash            string `json:"hash"`
	BlockNumber     string `json:"blockNumber"`
	TimeStamp       string `json:"timeStamp"`
	From            string `json:"from"`
	To              string `json:"to"`
	TokenID         string `json:"tokenID"`
	TokenName       string `json:"tokenName"`
	TokenSymbol     string `json:"tokenSymbol"`
	ContractAddress string `json:"contractAddress"`
	GasPrice        string `json:"gasPrice"`
	GasUsed         string `json:"gasUsed"`
}

// ERC1155Tx represents an ERC-1155 token transaction from Etherscan
type ERC1155Tx struct {
	Hash            string `json:"hash"`
	BlockNumber     string `json:"blockNumber"`
	TimeStamp       string `json:"timeStamp"`
	From            string `json:"from"`
	To              string `json:"to"`
	TokenID         string `json:"tokenID"`
	TokenValue      string `json:"tokenValue"`
	TokenName       string `json:"tokenName"`
	TokenSymbol     string `json:"tokenSymbol"`
	ContractAddress string `json:"contractAddress"`
	GasPrice        string `json:"gasPrice"`
	GasUsed         string `json:"gasUsed"`
}