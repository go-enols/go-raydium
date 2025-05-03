package raydium

import "github.com/gagliardetto/solana-go"

const (
	ACCOUNT_LAYOUT_LEN uint64 = 165
)

var (
	WSOL            = solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")
	AmmAuthority    = solana.MustPublicKeyFromBase58("5Q544fKrFoe6tsEbD7S8EmxGTJYAKtTVhAW5Q5pge4j1")
	CpmmAuthority   = solana.MustPublicKeyFromBase58("GpMZbSM2GgvTKHJirzeGfMFoaZ8UR2X7F4v8vHTvxFbL")
	OpenBookProgram = solana.MustPublicKeyFromBase58("srmqPvymJeFKQ4zGQed1GFppgkRHL9kaELCbyksJtPX")
)
