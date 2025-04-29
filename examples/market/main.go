package main

import (
	"context"

	"github.com/go-enols/go-log"
	"github.com/go-enols/go-raydium/core"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
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

	client := wallet.GetClient()
	poolPubkey := solana.MustPublicKeyFromBase58("3bC2e2RxcfvF9oP22LvbaNsVwoS2T98q6ErCRoayQYdq")

	info, err := client.GetAccountInfoWithOpts(ctx, poolPubkey, &rpc.GetAccountInfoOpts{
		Commitment: rpc.CommitmentConfirmed,
	})
	if err != nil {
		log.Fatal(err)
	}
	data := info.GetBinary()
	if data == nil {
		log.Fatal("无法获取元数据")
	}

	res := new(core.LIQUIDITY_STATE_LAYOUT_V4)
	if err := res.UnmarshalWithDecoder(bin.NewBinDecoder(data)); err != nil {
		log.Fatal(err)
	}

	marker, err := res.MarketInfo(ctx, client)
	if err != nil {
		log.Fatal(err)
	}
	log.Debug(marker)
}
