package main

import (
	"context"

	"github.com/go-enols/go-log"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/go-enols/go-raydium/core"
	"github.com/go-enols/gosolana"
)

var (
	Proxy               = "http://127.0.0.1:7890"
	NetWork rpc.Cluster = rpc.MainNetBeta
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	option := gosolana.Option{
		RpcUrl:  NetWork.RPC,
		WsUrl:   NetWork.WS,
		Proxy:   Proxy,
		WsProxy: Proxy,
		Pkey:    "26HX8sewDP8Y6xTE3v4DtR5HHB5D4ua1MUxPEHyUA2j3SrFt4FDLwaXTWZg7BoeGGooyojtftUjR8CTMCczhQyrD",
	}
	// 创建Solana钱包实例
	wallet, err := gosolana.NewWallet(ctx, option)
	if err != nil {
		log.Fatal(err)
	}

	amountIn := 0.1
	poolPublick := solana.MustPublicKeyFromBase58("3bC2e2RxcfvF9oP22LvbaNsVwoS2T98q6ErCRoayQYdq")
	base, quote, amountOut, err := core.GetPoolPriceByLiquidity(
		ctx,
		wallet.GetClient(),
		poolPublick,
		amountIn,
	)
	if err != nil {
		log.Fatal(err)
	}

	log.Debugf("池地址 %s", poolPublick.String())
	log.Debugf("%f %s = %f %s", amountIn, base.Sysbol, amountOut, quote.Sysbol)
	log.Debugf("%s 供应量 %f | 价值 %f %s", base.Mint.String(), base.Value, base.Cost, base.Sysbol)
	log.Debugf("%s 供应量 %f | 价值 %f %s", quote.Mint.String(), quote.Value, quote.Cost, base.Sysbol)
	// 2025-04-29 09:34:19 |DEBUG   | main.main:47 - 池地址 3bC2e2RxcfvF9oP22LvbaNsVwoS2T98q6ErCRoayQYdq
	// 2025-04-29 09:34:19 |DEBUG   | main.main:48 - 0.100000 SOL = 423.584214 jellyjelly
	// 2025-04-29 09:34:19 |DEBUG   | main.main:49 - So11111111111111111111111111111111111111112 供应量 10139.361460 | 价值 10139.361460 SOL
	// 2025-04-29 09:34:19 |DEBUG   | main.main:50 - FeR8VBqNRSUD5NtXAj2n3j1dAHkZHfyDktKuLXD4pump 供应量 43060611.281756 | 价值 10165.773394 SOL

}
