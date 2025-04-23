package raydium

import (
	"encoding/json"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/go-enols/go-raydium/core"
)

type TxCandidate struct {
	Signature solana.Signature
	Metadata  *json.RawMessage
}

type PoolCreateInfo struct {
	PoolAddress solana.PublicKey
	Base        *core.Liquidity
	Quote       *core.Liquidity
	Rtx         *rpc.GetTransactionResult
	Tx          *solana.Transaction
	Price       float64
}

type SwapTransferInfo struct {
	Data core.SwapInfo
	Rtx  *rpc.GetTransactionResult
	Tx   *solana.Transaction
}
