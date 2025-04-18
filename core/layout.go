package core

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

// 池子类型枚举
type PoolType int

const (
	PoolTypeUnknown PoolType = iota
	PoolTypeV4
	PoolTypeCPMM
	PoolTypeCAMM
)

// 统一的池子解析结果
type PoolParseResult struct {
	Type PoolType
	V4   *LIQUIDITY_STATE_LAYOUT_V4
	CPMM *CPMM_STATE_LAYOUT
	CAMM *CAMM_STATE_LAYOUT
}

// 自动识别并解析池子
func ParsePoolAccountByRPC(ctx context.Context, client *rpc.Client, poolPubkey solana.PublicKey) (*PoolParseResult, error) {
	info, err := client.GetAccountInfoWithOpts(ctx, poolPubkey, &rpc.GetAccountInfoOpts{
		Commitment: rpc.CommitmentConfirmed,
	})
	if err != nil {
		return nil, err
	}
	data := info.GetBinary()
	if data == nil {
		return nil, errors.New("pool account not found")
	}
	length := len(data)
	if 700 < length && length < 800 {
		// 1. 尝试V4
		v4 := new(LIQUIDITY_STATE_LAYOUT_V4)
		if err := v4.UnmarshalWithDecoder(bin.NewBinDecoder(data)); err == nil {
			return &PoolParseResult{Type: PoolTypeV4, V4: v4}, nil
		}
	} else if length < 700 {
		// 2. 尝试CPMM
		cpmm := new(CPMM_STATE_LAYOUT)
		if err := cpmm.UnmarshalWithDecoder(bin.NewBinDecoder(data)); err == nil {
			return &PoolParseResult{Type: PoolTypeCPMM, CPMM: cpmm}, nil
		}
	} else if length > 1000 {
		// 3. 尝试cAMM
		camm := new(CAMM_STATE_LAYOUT)
		if err := camm.UnmarshalWithDecoder(bin.NewBinDecoder(data)); err == nil {
			return &PoolParseResult{Type: PoolTypeCAMM, CAMM: camm}, nil
		}
	}

	log.Printf("无法识别的池子类型，数据长度: %d", length)
	return nil, errors.New("unsupported or unknown pool layout")
}

// 等待交易被确认
func GetConfirmedTransaction(ctx context.Context, client *rpc.Client, tx solana.Signature) (*rpc.GetTransactionResult, *solana.Transaction, error) {
	for {
		select {
		case <-ctx.Done():
			return nil, nil, ctx.Err()
		default:
			rctx, rcancel := context.WithTimeout(ctx, 5*time.Second)
			rpcTx, err := client.GetTransaction(rctx, tx, &rpc.GetTransactionOpts{
				MaxSupportedTransactionVersion: &Max_Transaction_Version,
				Commitment:                     rpc.CommitmentConfirmed,
			})
			rcancel()

			if err != nil {
				log.Println("读取交易信息失败,3秒后重试 | ", err)
				time.Sleep(3 * time.Second)
				continue
			}

			if rpcTx.Meta.Err != nil {
				return nil, nil, fmt.Errorf("交易失败: %v", rpcTx.Meta.Err)
			}

			tx, err := rpcTx.Transaction.GetTransaction()
			if err != nil {
				return nil, nil, fmt.Errorf("无法交易: %v", err)
			}

			return rpcTx, tx, nil
		}
	}
}
