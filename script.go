package raydium

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/go-enols/go-raydium/core"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/go-enols/gosolana/ws"
)

type TxCandidate struct {
	Signature solana.Signature
	Metadata  *json.RawMessage
}

type Scipt struct {
	RpcClient *rpc.Client // http rpc
	signer    chan TxCandidate
}

func NewScript(ctx context.Context, rpc *rpc.Client) *Scipt {
	script := &Scipt{
		RpcClient: rpc,
		signer:    make(chan TxCandidate, 1000), // 设置通道缓存，如果数据量过大可以适量增加
	}
	go script.signerProcess(ctx)
	return script
}

// 解析Raydium 日志
func (s *Scipt) RaydiumLogs(logs *ws.LogResult) {
	// Raw input data; bytes[2:6] -> possible timestamp for ido openning; bytes[7:15] -> possible pc amount, last 8 bytes, possible coin amount
	// Find possible Purchase IDO instruction logs:
	for i := range logs.Value.Logs {
		curLog := logs.Value.Logs[i]

		// Parse IDO info from log
		_, after, found := strings.Cut(curLog, " InitializeInstruction2 ")
		if !found {
			continue // Search further, not IDO log.
		}

		// Add quotes to keys.
		splitted := strings.Split(after, " ")
		for i, s := range splitted {
			if strings.Contains(s, ":") {
				splitted[i] = "\"" + s[:len(s)-1] + "\":"
			}
		}

		metadata := json.RawMessage(strings.Join(splitted, " "))
		if !json.Valid(metadata) {
			continue // Search further, invalid JSON.
		}

		s.signer <- TxCandidate{
			Signature: logs.Value.Signature,
			Metadata:  &metadata,
		}
		break
	}
}

// 解析OpenBook日志
func (s *Scipt) OpenBookLogs(logs *ws.LogResult) {
	// Find possible InitMarket instruction logs:
	for i := range logs.Value.Logs {
		curLog := logs.Value.Logs[i]
		if !strings.Contains(curLog, "Program 11111111111111111111111111111111 success") {
			continue // Search further.
		}

		if i+1 >= len(logs.Value.Logs) {
			break // No more logs.
		}

		nextLog := logs.Value.Logs[i+1]
		if !strings.Contains(nextLog, "Program srmqPvymJeFKQ4zGQed1GFppgkRHL9kaELCbyksJtPX invoke [1]") {
			continue // Search further.
		}
		s.signer <- TxCandidate{
			Signature: logs.Value.Signature,
			Metadata:  nil,
		}
		break
	}
}

// 处理每一个过来的签名直接查询出对应的数据
func (s *Scipt) signerProcess(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case sign := <-s.signer:
			// 获取签名对应的交易信息
			retx, tx, err := core.GetConfirmedTransaction(ctx, s.RpcClient, sign.Signature)
			if err != nil {
				log.Println("查询交易信息失败 |", err)
				continue
			}
			poolAddress, err := ParseLpAddressByLogs(tx)
			if err != nil {
				continue
			}

			base, quote, price, err := core.GetPoolPriceByLiquidity(ctx, s.RpcClient, poolAddress)
			if err != nil {
				log.Println(poolAddress.String(), "查询池子数据失败", err)
				continue
			}
			log.Println("-----------------------------------------")
			log.Println("监听到一个关键交易Hash", tx.Signatures)
			log.Println("发现一个新的raydium池子:", poolAddress.String(), base.Sysbol, "-", quote.Sysbol)
			log.Printf("%s | 合约地址 %s | Token %s", base.Sysbol, base.Address, base.Mint)
			log.Printf("%s | 合约地址 %s | Token %s", quote.Sysbol, quote.Address, quote.Mint)
			log.Printf("当前价格 | 1 %s=%.9f %s", quote.Sysbol, price, base.Sysbol)
			log.Printf("池子创建时间 | %s", retx.BlockTime.Time().Format(time.DateTime))
			log.Printf("发现时间 | %s", time.Now().Format(time.DateTime))

		}
	}
}

// 从交易日志中解析出池子地址
func ParseLpAddressByLogs(tx *solana.Transaction) (solana.PublicKey, error) {
	safeIndex := func(idx uint16) solana.PublicKey {
		if idx >= uint16(len(tx.Message.AccountKeys)) {
			return solana.PublicKey{}
		}
		return tx.Message.AccountKeys[idx]
	}
	log.Println(tx.Signatures)
	for _, instr := range tx.Message.Instructions {
		program, err := tx.Message.Program(instr.ProgramIDIndex)
		if err != nil {
			continue // Program account index out of range.
		}

		if program.String() != core.OpenBookDex.String() {
			continue // Not called by serum.
		}

		if len(instr.Accounts) < 10 {
			log.Println("如果小于10条表示这个指令有问题,或者openbook的版本更新了!!!!")
			// 如果小于10条表示这个指令有问题，或者openbook的版本更新了
			continue // Not enough accounts for InitializeMarket instruction.
		}

		const BaseMintIndex = 7
		const QuoteMinIndex = 8

		if safeIndex(instr.Accounts[QuoteMinIndex]) != solana.WrappedSol && safeIndex(instr.Accounts[BaseMintIndex]) != solana.WrappedSol {
			return solana.PublicKey{}, fmt.Errorf("找到一个新的市场,但是不是SOL的市场,已忽略")
		}

		return safeIndex(instr.Accounts[0]), nil
	}
	for _, instr := range tx.Message.Instructions {
		program, err := tx.Message.Program(instr.ProgramIDIndex)
		if err != nil {
			continue // Program account index out of range.
		}

		if program.String() != core.RaydiumLiquidityProgramV4.String() {
			continue // Not called by serum.
		}

		if len(instr.Accounts) < 21 {
			continue // Not enough accounts for Purchase IDO instruction.
		}
		return safeIndex(instr.Accounts[4]), nil
	}
	return solana.PublicKey{}, errors.New("不是openbook中raydium订单")
}
