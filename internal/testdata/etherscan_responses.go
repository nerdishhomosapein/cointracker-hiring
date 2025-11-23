package testdata

// NormalTxResponse is a sample Etherscan response for normal ETH transfers
const NormalTxResponse = `{
  "status": "1",
  "message": "OK",
  "result": [
    {
      "blockNumber": "20000000",
      "timeStamp": "1700000000",
      "hash": "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
      "nonce": "42",
      "blockHash": "0xblockhash1234567890abcdef1234567890abcdef1234567890abcdef123456",
      "transactionIndex": "15",
      "from": "0xa39b189482f984388a34460636fea9eb181ad1a6",
      "to": "0xd620AADaBaA20d2af700853C4504028cba7C3333",
      "value": "1000000000000000000",
      "gas": "21000",
      "gasPrice": "50000000000",
      "isError": "0",
      "txreceipt_status": "1",
      "input": "0x",
      "contractAddress": "",
      "cumulativeGasUsed": "5000000",
      "gasUsed": "21000",
      "confirmations": "100000",
      "methodId": "0x",
      "functionName": ""
    },
    {
      "blockNumber": "19999999",
      "timeStamp": "1699999990",
      "hash": "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
      "nonce": "41",
      "blockHash": "0xblockhash_another_1234567890abcdef1234567890abcdef1234567890ab",
      "transactionIndex": "10",
      "from": "0xa39b189482f984388a34460636fea9eb181ad1a6",
      "to": "0x1111111254fb6c44bac0bed2854e76f90643097d",
      "value": "500000000000000000",
      "gas": "21000",
      "gasPrice": "45000000000",
      "isError": "0",
      "txreceipt_status": "1",
      "input": "0x",
      "contractAddress": "",
      "cumulativeGasUsed": "4800000",
      "gasUsed": "21000",
      "confirmations": "100001",
      "methodId": "0x",
      "functionName": ""
    }
  ]
}`

// InternalTxResponse is a sample Etherscan response for internal transfers
const InternalTxResponse = `{
  "status": "1",
  "message": "OK",
  "result": [
    {
      "blockNumber": "19999998",
      "timeStamp": "1699999980",
      "hash": "0x9999999999999999999999999999999999999999999999999999999999999999",
      "from": "0xa39b189482f984388a34460636fea9eb181ad1a6",
      "to": "0x2222222254fb6c44bac0bed2854e76f90643097d",
      "value": "250000000000000000",
      "contractAddress": "0x3333333354fb6c44bac0bed2854e76f90643097d",
      "input": "0x123456",
      "type": "call",
      "gas": "50000",
      "gasUsed": "40000",
      "traceId": "1",
      "isError": "0",
      "errCode": ""
    }
  ]
}`

// ERC20TokenTxResponse is a sample Etherscan response for ERC-20 token transfers
const ERC20TokenTxResponse = `{
  "status": "1",
  "message": "OK",
  "result": [
    {
      "blockNumber": "19999997",
      "timeStamp": "1699999970",
      "hash": "0x8888888888888888888888888888888888888888888888888888888888888888",
      "nonce": "40",
      "blockHash": "0xblockhash_erc20_1234567890abcdef1234567890abcdef1234567890",
      "from": "0xa39b189482f984388a34460636fea9eb181ad1a6",
      "contractAddress": "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
      "to": "0xd620AADaBaA20d2af700853C4504028cba7C3333",
      "value": "1000000000",
      "tokenName": "USD Coin",
      "tokenSymbol": "USDC",
      "tokenDecimal": "6",
      "transactionIndex": "20",
      "gas": "100000",
      "gasPrice": "55000000000",
      "gasUsed": "80000",
      "cumulativeGasUsed": "6000000",
      "input": "0xa9059cbb",
      "confirmations": "99999",
      "isError": "0",
      "txreceipt_status": "1"
    },
    {
      "blockNumber": "19999996",
      "timeStamp": "1699999960",
      "hash": "0x7777777777777777777777777777777777777777777777777777777777777777",
      "nonce": "39",
      "blockHash": "0xblockhash_erc20_another_1234567890abcdef1234567890abcdef12345",
      "from": "0xa39b189482f984388a34460636fea9eb181ad1a6",
      "contractAddress": "0xdac17f958d2ee523a2206206994597c13d831ec7",
      "to": "0x1111111254fb6c44bac0bed2854e76f90643097d",
      "value": "5000000000",
      "tokenName": "Tether USD",
      "tokenSymbol": "USDT",
      "tokenDecimal": "6",
      "transactionIndex": "18",
      "gas": "100000",
      "gasPrice": "50000000000",
      "gasUsed": "75000",
      "cumulativeGasUsed": "5950000",
      "input": "0xa9059cbb",
      "confirmations": "100000",
      "isError": "0",
      "txreceipt_status": "1"
    }
  ]
}`

// ERC721NFTResponse is a sample Etherscan response for ERC-721 NFT transfers
const ERC721NFTResponse = `{
  "status": "1",
  "message": "OK",
  "result": [
    {
      "blockNumber": "19999995",
      "timeStamp": "1699999950",
      "hash": "0x6666666666666666666666666666666666666666666666666666666666666666",
      "nonce": "38",
      "blockHash": "0xblockhash_nft_1234567890abcdef1234567890abcdef1234567890",
      "from": "0xa39b189482f984388a34460636fea9eb181ad1a6",
      "contractAddress": "0xbc4ca0eda7647a8ab7c2061c2e2ad183",
      "to": "0xd620AADaBaA20d2af700853C4504028cba7C3333",
      "tokenID": "1337",
      "tokenName": "Bored Ape Yacht Club",
      "tokenSymbol": "BAYC",
      "tokenDecimal": "0",
      "transactionIndex": "25",
      "gas": "150000",
      "gasPrice": "60000000000",
      "gasUsed": "125000",
      "cumulativeGasUsed": "6125000",
      "input": "0x23b872dd",
      "confirmations": "99998",
      "isError": "0",
      "txreceipt_status": "1"
    }
  ]
}`

// ERC1155Response is a sample Etherscan response for ERC-1155 token transfers
const ERC1155Response = `{
  "status": "1",
  "message": "OK",
  "result": [
    {
      "blockNumber": "19999994",
      "timeStamp": "1699999940",
      "hash": "0x5555555555555555555555555555555555555555555555555555555555555555",
      "nonce": "37",
      "blockHash": "0xblockhash_1155_1234567890abcdef1234567890abcdef1234567890",
      "from": "0xa39b189482f984388a34460636fea9eb181ad1a6",
      "contractAddress": "0x76be3b62873462d2142405439777e053",
      "to": "0xd620AADaBaA20d2af700853C4504028cba7C3333",
      "tokenID": "999",
      "tokenValue": "50",
      "tokenName": "Polymarket Conditional Token",
      "tokenSymbol": "POLY",
      "tokenDecimal": "0",
      "transactionIndex": "30",
      "gas": "200000",
      "gasPrice": "65000000000",
      "gasUsed": "150000",
      "cumulativeGasUsed": "6275000",
      "input": "0xf242432a",
      "confirmations": "99997",
      "isError": "0",
      "txreceipt_status": "1"
    }
  ]
}`

// ErrorResponse is a sample error response from Etherscan
const ErrorResponse = `{
  "status": "0",
  "message": "NOTOK",
  "result": "Invalid API Key"
}`

// RateLimitResponse is a sample rate-limit response
const RateLimitResponse = `{
  "status": "0",
  "message": "NOTOK",
  "result": "Max rate limit reached, please use API Key for higher rate limit"
}`

// EmptyResultResponse is a response with no transactions
const EmptyResultResponse = `{
  "status": "1",
  "message": "OK",
  "result": []
}`
