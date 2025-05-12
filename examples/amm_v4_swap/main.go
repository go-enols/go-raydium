package main

import (
	"context"

	"github.com/go-enols/go-log"
	"github.com/go-enols/go-raydium"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/go-enols/gosolana"
)

var (
	NetWork rpc.Cluster = rpc.MainNetBeta
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	option := gosolana.Option{
		RpcUrl: NetWork.RPC,
		WsUrl:  NetWork.WS,
	}
	wallet, err := gosolana.NewWallet(ctx, option)
	if err != nil {
		log.Fatal(err)
	}
	poolAddress := solana.MustPublicKeyFromBase58("AgFnRLUScRD2E4nWQxW73hdbSN7eKEUb2jHX7tx9YTYc")

	// ammBuild := amm_v3.NewSwapV2InstructionBuilder().SetAmount(amount uint64)
	client := raydium.NewV4Client(wallet)

	tx, err := client.Swap(ctx, 21.45, poolAddress, false, 0.01)
	if err != nil {
		log.Errorf("创建交易失败 | %s", err)
		return
	}
	log.Debug(tx.String())
	// return
	// 签名交易
	out, err := tx.Sign(func(key solana.PublicKey) *solana.PrivateKey {
		if key.Equals(wallet.PublicKey()) {
			return &wallet.PrivateKey
		}
		return nil
	})
	if err != nil {
		log.Print("签名交易失败", err.Error())
		return
	}
	log.Printf("签名交易输出 | %v", out)
	// 7. 发送交易
	sig, err := wallet.GetClient().SendTransactionWithOpts(
		context.Background(),
		tx,
		rpc.TransactionOpts{
			SkipPreflight:       false,
			PreflightCommitment: rpc.CommitmentFinalized,
		},
	)
	if err != nil {
		log.Errorf("发送交易失败 %s", err)
		return
	}
	log.Printf("[INFO] Transaction Signature: %s", sig)
	log.Printf("[INFO] 交易详情 | %v", tx) // 打印交易详情
	result, err := wallet.GetTransaction(ctx, sig)
	if err != nil {
		log.Printf("[ERROR] 获取交易状态失败 | %s", err)
		return
	}
	if result {
		log.Debug("交易成功")
	} else {
		log.Error("交易失败")
	}
}
