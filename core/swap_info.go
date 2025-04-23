package core

import (
	"context"
	"errors"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// string 交易人地址
type SwapInfo map[string]struct {
	Data map[string]struct { // string 代币地址
		InputAmount float64 // 交易输入金额
		OutAmount   float64 // 交易输出金额
		Name        string  // 代币名称 可能不需要？先预留具体获取方法可以使用gosolana中的GetTokenMetaOnChain方法获取
	}
}

// 更具交易日志解析这个是不是某个池子的交易，如果是就返回交易方向代币数量以及账户
//
// rtx, tx 请使用`GetConfirmedTransaction`方法获取
func SwapTransferLog(ctx context.Context, client *rpc.Client, rtx *rpc.GetTransactionResult, tx *solana.Transaction) (SwapInfo, error) {
	// 可能不需要解析？
	// isRaydium := false
	// for _, inst := range tx.Message.Instructions {
	// 	programId, err := tx.Message.Program(inst.ProgramIDIndex)
	// 	if err != nil {
	// 		log.Error("没有找到对应的程序id,可能是错误的tx")
	// 		return nil, err
	// 	}
	// 	log.Debug(programId.String())
	// 	if programId != RaydiumLiquidityProgramV4 {
	// 		continue
	// 	}
	// 	isRaydium = true
	// 	break
	// }
	// if !isRaydium {
	// 	return nil, errors.New("不是 Raydium 订单")
	// }

	data, err := swapInfoByTransfer(rtx)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func swapInfoByTransfer(tx *rpc.GetTransactionResult) (SwapInfo, error) {
	if tx == nil || tx.Meta == nil {
		return nil, errors.New("交易数据为空")
	}

	result := make(SwapInfo)

	preTokenBalances := tx.Meta.PreTokenBalances
	postTokenBalances := tx.Meta.PostTokenBalances

	// 用于唯一标识每个账户和代币组合
	type key struct {
		Owner string // 账户地址
		Mint  string // 代币地址
	}
	preMap := make(map[key]float64)  // 存储每个账户+代币的交易前余额
	postMap := make(map[key]float64) // 存储每个账户+代币的交易后余额

	// 遍历交易前余额，填充 preMap
	for _, pre := range preTokenBalances {
		if pre.UiTokenAmount.UiAmount == nil {
			continue // 如果余额为空则跳过
		}
		k := key{Owner: pre.Owner.String(), Mint: pre.Mint.String()}
		preMap[k] = *pre.UiTokenAmount.UiAmount
	}
	// 遍历交易后余额，填充 postMap
	for _, post := range postTokenBalances {
		if post.UiTokenAmount.UiAmount == nil {
			continue // 如果余额为空则跳过
		}
		k := key{Owner: post.Owner.String(), Mint: post.Mint.String()}
		postMap[k] = *post.UiTokenAmount.UiAmount
	}

	// 合并所有出现过的账户+代币组合，避免 Pre/Post 不一一对应导致遗漏
	seen := make(map[key]struct{})
	for k := range preMap {
		seen[k] = struct{}{}
	}
	for k := range postMap {
		seen[k] = struct{}{}
	}

	// 遍历所有账户+代币组合，计算变化量
	for k := range seen {
		preAmt := preMap[k]       // 交易前余额，若没有则为0
		postAmt := postMap[k]     // 交易后余额，若没有则为0
		inputAmount := float64(0) // 转入金额
		outAmount := float64(0)   // 转出金额

		// 余额增加为转入，减少为转出
		if postAmt > preAmt {
			inputAmount = postAmt - preAmt
		} else if preAmt > postAmt {
			outAmount = preAmt - postAmt
		}
		// 如果没有变化则跳过
		if inputAmount == 0 && outAmount == 0 {
			continue
		}
		// 初始化 SwapInfo 的结构
		if _, ok := result[k.Owner]; !ok {
			result[k.Owner] = struct {
				Data map[string]struct {
					InputAmount float64
					OutAmount   float64
					Name        string
				}
			}{Data: make(map[string]struct {
				InputAmount float64
				OutAmount   float64
				Name        string
			})}
		}
		// 填充每个账户的交易信息
		result[k.Owner].Data[k.Mint] = struct {
			InputAmount float64
			OutAmount   float64
			Name        string
		}{
			InputAmount: inputAmount,
			OutAmount:   outAmount,
			Name:        "", // 代币名称可后续补充
		}
	}

	return result, nil
}
