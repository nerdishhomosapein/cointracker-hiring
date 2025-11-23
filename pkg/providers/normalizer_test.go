package providers

import (
	"conintracker-hiring/pkg/models"
	"testing"
	"time"
)

func TestNormalizeNormalTx(t *testing.T) {
	normalizer := NewEtherscanNormalizer()

	tests := []struct {
		name    string
		tx      EtherscanNormalTx
		want    *models.Transaction
		wantErr bool
	}{
		{
			name: "valid_normal_eth_transfer",
			tx: EtherscanNormalTx{
				BlockNumber:      "20000000",
				TimeStamp:        "1700000000",
				Hash:             "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
				From:             "0xa39b189482f984388a34460636fea9eb181ad1a6",
				To:               "0xd620AADaBaA20d2af700853C4504028cba7C3333",
				Value:            "1000000000000000000", // 1 ETH in wei
				GasPrice:         "50000000000",         // 50 Gwei
				GasUsed:          "21000",
				IsError:          "0",
				TxReceiptStatus:  "1",
			},
			want: &models.Transaction{
				Hash:      "0x1234567890abcdef1234567890abcdef1234567890abcdef1234567890abcdef",
				Timestamp: time.Unix(1700000000, 0),
				From:      "0xa39b189482f984388a34460636fea9eb181ad1a6",
				To:        "0xd620AADaBaA20d2af700853C4504028cba7C3333",
				Type:      models.TypeEthTransfer,
				Amount:    "1",
				GasFeeETH: "0.00105",
				BlockNumber: 20000000,
				GasUsed:   21000,
				IsError:   false,
			},
			wantErr: false,
		},
		{
			name: "failed_transaction",
			tx: EtherscanNormalTx{
				BlockNumber:     "19999999",
				TimeStamp:       "1699999990",
				Hash:            "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
				From:            "0xa39b189482f984388a34460636fea9eb181ad1a6",
				To:              "0x1111111254fb6c44bac0bed2854e76f90643097d",
				Value:           "500000000000000000",
				GasPrice:        "45000000000",
				GasUsed:         "21000",
				IsError:         "1",
				TxReceiptStatus: "0",
			},
			want: &models.Transaction{
				Hash:      "0xabcdef1234567890abcdef1234567890abcdef1234567890abcdef1234567890",
				Timestamp: time.Unix(1699999990, 0),
				From:      "0xa39b189482f984388a34460636fea9eb181ad1a6",
				To:        "0x1111111254fb6c44bac0bed2854e76f90643097d",
				Type:      models.TypeEthTransfer,
				Amount:    "0.5",
				GasFeeETH: "0.000945",
				BlockNumber: 19999999,
				GasUsed:   21000,
				IsError:   true,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizer.NormalizeNormalTx(tt.tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeNormalTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.Hash != tt.want.Hash {
					t.Errorf("Hash mismatch: got %s, want %s", got.Hash, tt.want.Hash)
				}
				if got.From != tt.want.From {
					t.Errorf("From mismatch: got %s, want %s", got.From, tt.want.From)
				}
				if got.To != tt.want.To {
					t.Errorf("To mismatch: got %s, want %s", got.To, tt.want.To)
				}
				if got.Type != tt.want.Type {
					t.Errorf("Type mismatch: got %s, want %s", got.Type, tt.want.Type)
				}
				if got.Amount != tt.want.Amount {
					t.Errorf("Amount mismatch: got %s, want %s", got.Amount, tt.want.Amount)
				}
				if got.GasFeeETH != tt.want.GasFeeETH {
					t.Errorf("GasFeeETH mismatch: got %s, want %s", got.GasFeeETH, tt.want.GasFeeETH)
				}
				if got.BlockNumber != tt.want.BlockNumber {
					t.Errorf("BlockNumber mismatch: got %d, want %d", got.BlockNumber, tt.want.BlockNumber)
				}
				if got.IsError != tt.want.IsError {
					t.Errorf("IsError mismatch: got %v, want %v", got.IsError, tt.want.IsError)
				}
			}
		})
	}
}

func TestNormalizeInternalTx(t *testing.T) {
	normalizer := NewEtherscanNormalizer()

	tests := []struct {
		name    string
		tx      EtherscanInternalTx
		want    *models.Transaction
		wantErr bool
	}{
		{
			name: "valid_internal_transfer",
			tx: EtherscanInternalTx{
				BlockNumber:     "19999998",
				TimeStamp:       "1699999980",
				Hash:            "0x9999999999999999999999999999999999999999999999999999999999999999",
				From:            "0xa39b189482f984388a34460636fea9eb181ad1a6",
				To:              "0x2222222254fb6c44bac0bed2854e76f90643097d",
				Value:           "250000000000000000",
				ContractAddress: "0x3333333354fb6c44bac0bed2854e76f90643097d",
				GasUsed:         "40000",
				IsError:         "0",
			},
			want: &models.Transaction{
				Hash:      "0x9999999999999999999999999999999999999999999999999999999999999999",
				Timestamp: time.Unix(1699999980, 0),
				From:      "0xa39b189482f984388a34460636fea9eb181ad1a6",
				To:        "0x2222222254fb6c44bac0bed2854e76f90643097d",
				Type:      models.TypeInternal,
				Amount:    "0.25",
				BlockNumber: 19999998,
				GasUsed:   40000,
				IsError:   false,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizer.NormalizeInternalTx(tt.tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeInternalTx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.Hash != tt.want.Hash {
					t.Errorf("Hash mismatch: got %s, want %s", got.Hash, tt.want.Hash)
				}
				if got.Type != tt.want.Type {
					t.Errorf("Type mismatch: got %s, want %s", got.Type, tt.want.Type)
				}
				if got.Amount != tt.want.Amount {
					t.Errorf("Amount mismatch: got %s, want %s", got.Amount, tt.want.Amount)
				}
			}
		})
	}
}

func TestNormalizeERC20Tx(t *testing.T) {
	normalizer := NewEtherscanNormalizer()

	tests := []struct {
		name    string
		tx      EtherscanTokenTx
		want    *models.Transaction
		wantErr bool
	}{
		{
			name: "valid_erc20_transfer",
			tx: EtherscanTokenTx{
				BlockNumber:     "19999997",
				TimeStamp:       "1699999970",
				Hash:            "0x8888888888888888888888888888888888888888888888888888888888888888",
				From:            "0xa39b189482f984388a34460636fea9eb181ad1a6",
				ContractAddress: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
				To:              "0xd620AADaBaA20d2af700853C4504028cba7C3333",
				Value:           "1000000000", // 1,000 USDC (6 decimals)
				TokenName:       "USD Coin",
				TokenSymbol:     "USDC",
				TokenDecimal:    "6",
				GasPrice:        "55000000000",
				GasUsed:         "80000",
				IsError:         "0",
				TxReceiptStatus: "1",
			},
			want: &models.Transaction{
				Hash:                 "0x8888888888888888888888888888888888888888888888888888888888888888",
				Timestamp:            time.Unix(1699999970, 0),
				From:                 "0xa39b189482f984388a34460636fea9eb181ad1a6",
				To:                   "0xd620AADaBaA20d2af700853C4504028cba7C3333",
				Type:                 models.TypeERC20Transfer,
				AssetContractAddress: "0xa0b86991c6218b36c1d19d4a2e9eb0ce3606eb48",
				AssetSymbol:          "USDC",
				Amount:               "1000.0",
				GasFeeETH:            "0.0044",
				BlockNumber:          19999997,
				GasUsed:              80000,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizer.NormalizeERC20Tx(tt.tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeERC20Tx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.Hash != tt.want.Hash {
					t.Errorf("Hash mismatch: got %s, want %s", got.Hash, tt.want.Hash)
				}
				if got.Type != tt.want.Type {
					t.Errorf("Type mismatch: got %s, want %s", got.Type, tt.want.Type)
				}
				if got.AssetSymbol != tt.want.AssetSymbol {
					t.Errorf("AssetSymbol mismatch: got %s, want %s", got.AssetSymbol, tt.want.AssetSymbol)
				}
			}
		})
	}
}

func TestNormalizeERC721Tx(t *testing.T) {
	normalizer := NewEtherscanNormalizer()

	tests := []struct {
		name    string
		tx      EtherscanTokenTx
		want    *models.Transaction
		wantErr bool
	}{
		{
			name: "valid_erc721_transfer",
			tx: EtherscanTokenTx{
				BlockNumber:     "19999995",
				TimeStamp:       "1699999950",
				Hash:            "0x6666666666666666666666666666666666666666666666666666666666666666",
				From:            "0xa39b189482f984388a34460636fea9eb181ad1a6",
				ContractAddress: "0xbc4ca0eda7647a8ab7c2061c2e2ad183",
				To:              "0xd620AADaBaA20d2af700853C4504028cba7C3333",
				TokenID:         "1337",
				TokenName:       "Bored Ape Yacht Club",
				TokenSymbol:     "BAYC",
				TokenDecimal:    "0",
				GasPrice:        "60000000000",
				GasUsed:         "125000",
				IsError:         "0",
				TxReceiptStatus: "1",
			},
			want: &models.Transaction{
				Hash:                 "0x6666666666666666666666666666666666666666666666666666666666666666",
				Timestamp:            time.Unix(1699999950, 0),
				From:                 "0xa39b189482f984388a34460636fea9eb181ad1a6",
				To:                   "0xd620AADaBaA20d2af700853C4504028cba7C3333",
				Type:                 models.TypeERC721Transfer,
				AssetContractAddress: "0xbc4ca0eda7647a8ab7c2061c2e2ad183",
				AssetSymbol:          "BAYC",
				TokenID:              "1337",
				Amount:               "1",
				GasFeeETH:            "0.0075",
				BlockNumber:          19999995,
				GasUsed:              125000,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizer.NormalizeERC721Tx(tt.tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeERC721Tx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.Type != tt.want.Type {
					t.Errorf("Type mismatch: got %s, want %s", got.Type, tt.want.Type)
				}
				if got.TokenID != tt.want.TokenID {
					t.Errorf("TokenID mismatch: got %s, want %s", got.TokenID, tt.want.TokenID)
				}
				if got.Amount != tt.want.Amount {
					t.Errorf("Amount for NFT mismatch: got %s, want %s", got.Amount, tt.want.Amount)
				}
			}
		})
	}
}

func TestNormalizeERC1155Tx(t *testing.T) {
	normalizer := NewEtherscanNormalizer()

	tests := []struct {
		name    string
		tx      EtherscanTokenTx
		want    *models.Transaction
		wantErr bool
	}{
		{
			name: "valid_erc1155_transfer",
			tx: EtherscanTokenTx{
				BlockNumber:     "19999994",
				TimeStamp:       "1699999940",
				Hash:            "0x5555555555555555555555555555555555555555555555555555555555555555",
				From:            "0xa39b189482f984388a34460636fea9eb181ad1a6",
				ContractAddress: "0x76be3b62873462d2142405439777e053",
				To:              "0xd620AADaBaA20d2af700853C4504028cba7C3333",
				TokenID:         "999",
				TokenValue:      "50",
				TokenName:       "Polymarket Conditional Token",
				TokenSymbol:     "POLY",
				TokenDecimal:    "0",
				GasPrice:        "65000000000",
				GasUsed:         "150000",
				IsError:         "0",
				TxReceiptStatus: "1",
			},
			want: &models.Transaction{
				Hash:                 "0x5555555555555555555555555555555555555555555555555555555555555555",
				Timestamp:            time.Unix(1699999940, 0),
				From:                 "0xa39b189482f984388a34460636fea9eb181ad1a6",
				To:                   "0xd620AADaBaA20d2af700853C4504028cba7C3333",
				Type:                 models.TypeERC1155Transfer,
				AssetContractAddress: "0x76be3b62873462d2142405439777e053",
				AssetSymbol:          "POLY",
				TokenID:              "999",
				Amount:               "50",
				GasFeeETH:            "0.00975",
				BlockNumber:          19999994,
				GasUsed:              150000,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := normalizer.NormalizeERC1155Tx(tt.tx)
			if (err != nil) != tt.wantErr {
				t.Errorf("NormalizeERC1155Tx() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if got.Type != tt.want.Type {
					t.Errorf("Type mismatch: got %s, want %s", got.Type, tt.want.Type)
				}
				if got.Amount != tt.want.Amount {
					t.Errorf("Amount mismatch: got %s, want %s", got.Amount, tt.want.Amount)
				}
			}
		})
	}
}
