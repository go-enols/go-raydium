package raydium

import "github.com/gagliardetto/solana-go"

const (
	ACCOUNT_LAYOUT_LEN uint64 = 165
)

var (
	WSOL            = solana.MustPublicKeyFromBase58("So11111111111111111111111111111111111111112")
	AmmV3ProgramId  = solana.MustPublicKeyFromBase58("CAMMCzo5YL8w4VFF8KVHrK22GGUsp5VTaW7grrKgrWqK")
	AmmV4ProgramId  = solana.MustPublicKeyFromBase58("675kPX9MHTjS2zt1qfr1NYHuzeLXfQM9H24wFSUt1Mp8")
	CPMMProgramId   = solana.MustPublicKeyFromBase58("CPMMoo8L3F4NbTegBCKVNunggL7H1ZpdTHKxQB5qKP1C")
	AmmV4Authority  = solana.MustPublicKeyFromBase58("5Q544fKrFoe6tsEbD7S8EmxGTJYAKtTVhAW5Q5pge4j1")
	CpmmAuthority   = solana.MustPublicKeyFromBase58("GpMZbSM2GgvTKHJirzeGfMFoaZ8UR2X7F4v8vHTvxFbL")
	OpenBookProgram = solana.MustPublicKeyFromBase58("srmqPvymJeFKQ4zGQed1GFppgkRHL9kaELCbyksJtPX")
)
