package raydium

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/go-enols/go-log"

	"github.com/go-enols/go-raydium/core"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/go-enols/gosolana/ws"
)

type MonitPoolCreateScipt struct {
	RpcClient *rpc.Client      // http rpc
	signer    chan TxCandidate // 监听到池子创建之后的rpc通道
	callback  []func(PoolCreateInfo)
}

func NewMonitPoolCreateScipt(ctx context.Context, rpc *rpc.Client, CallBack ...func(PoolCreateInfo)) *MonitPoolCreateScipt {
	script := &MonitPoolCreateScipt{
		RpcClient: rpc,
		signer:    make(chan TxCandidate, 1000), // 设置通道缓存，如果数据量过大可以适量增加
		callback:  CallBack,
	}
	go script.signerProcess(ctx)
	return script
}

// 解析Raydium 日志
func (s *MonitPoolCreateScipt) RaydiumLogs(logs *ws.LogResult) {
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
func (s *MonitPoolCreateScipt) OpenBookLogs(logs *ws.LogResult) {
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
func (s *MonitPoolCreateScipt) signerProcess(ctx context.Context) {
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

			data := PoolCreateInfo{
				Quote:       quote,
				Base:        base,
				PoolAddress: poolAddress,
				Rtx:         retx,
				Tx:          tx,
				Price:       price,
			}
			for _, call := range s.callback {
				go call(data)
			}
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

type MonitSwapTransferScipt struct {
	RpcClient *rpc.Client      // http rpc
	signer    chan TxCandidate // 监听到有新的swap交易之后的rpc通道
	callback  []func(SwapTransferInfo)
}

func NewMonitSwapTransferScipt(ctx context.Context, rpc *rpc.Client, CallBack ...func(SwapTransferInfo)) *MonitSwapTransferScipt {
	script := &MonitSwapTransferScipt{
		RpcClient: rpc,
		signer:    make(chan TxCandidate, 1000), // 设置通道缓存，如果数据量过大可以适量增加
		callback:  CallBack,
	}
	go script.signerProcess(ctx)
	return script
}

// 解析Raydium Swap 日志
func (m *MonitSwapTransferScipt) RaydiumSwapLogs(logs *ws.LogResult) {
	// Raw input data; bytes[2:6] -> possible timestamp for ido openning; bytes[7:15] -> possible pc amount, last 8 bytes, possible coin amount
	// Find possible Purchase IDO instruction logs:
	for i := range logs.Value.Logs {
		curLog := logs.Value.Logs[i]
		if !strings.Contains(strings.ToLower(curLog), "swap") {
			continue // Search further, not IDO log.
		}
		m.signer <- TxCandidate{
			Signature: logs.Value.Signature,
		}
		break
	}
}

// 处理每一个过来的签名直接查询出对应的数据
func (s *MonitSwapTransferScipt) signerProcess(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		case sign := <-s.signer:
			log.Debug(sign.Signature.String())
			// 获取签名对应的交易信息
			retx, tx, err := core.GetConfirmedTransaction(ctx, s.RpcClient, sign.Signature)
			if err != nil {
				log.Println("查询交易信息失败 |", err)
				continue
			}
			data, err := core.SwapTransferLog(ctx, s.RpcClient, retx, tx)
			if err != nil {
				log.Errorf("无法解析 Swap 交易日志 | %s", err)
				continue
			}
			for _, call := range s.callback {
				go call(SwapTransferInfo{
					Data: data,
					Tx:   tx,
					Rtx:  retx,
				})
			}
		}
	}
}
