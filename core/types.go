package core

import (
	"github.com/gagliardetto/solana-go"
	"github.com/go-enols/gosolana/ws"
)

var (
	Max_Transaction_Version uint64 = 1
	WSOL                           = "So11111111111111111111111111111111111111112"

	// OpenBookDex program ID. 账本
	OpenBookDex solana.PublicKey = solana.MustPublicKeyFromBase58("srmqPvymJeFKQ4zGQed1GFppgkRHL9kaELCbyksJtPX")

	// Raydium program ID. 流动性
	RaydiumLiquidityProgramV4 solana.PublicKey = solana.MustPublicKeyFromBase58("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8") // This program calls Raydium Purchase IDO to create a new pair.

	// Raydium Authority V4 program ID
	RaydiumAuthorityProgramV4 solana.PublicKey = solana.MustPublicKeyFromBase58("5Q544fKrFoe6tsEbD7S8EmxGTJYAKtTVhAW5Q5pge4j1") // This is also a wallet that holds tokens and do swaps.
)

type LogApusic func(*ws.LogResult)

type AmmV4Account struct {
	Program string `json:"program"`
	Parsed  struct {
		Info struct {
			IsNative    bool   `json:"isNative"`
			Mint        string `json:"mint"`
			Owner       string `json:"owner"`
			State       string `json:"state"`
			TokenAmount struct {
				Amount         string  `json:"amount"`
				Decimals       int     `json:"decimals"`
				UIAmount       float64 `json:"uiAmount"`
				UIAmountString string  `json:"uiAmountString"`
			} `json:"tokenAmount"`
		} `json:"info"`
		Type string `json:"type"`
	} `json:"parsed"`
	Space int `json:"space"`
}
