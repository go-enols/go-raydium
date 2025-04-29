package main

import (
	"context"

	"github.com/go-enols/go-log"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/go-enols/go-raydium/ammV4"
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
	wallet, err := ammV4.NewClient(ctx, option)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug(wallet.Buy(ctx, "3bC2e2RxcfvF9oP22LvbaNsVwoS2T98q6ErCRoayQYdq", 0.1, 1))
}
