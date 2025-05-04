package main

import (
	"context"

	"github.com/go-enols/go-log"
	"github.com/go-enols/go-raydium/amm_v3"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/go-enols/gosolana"
)

var (
	Proxy               = "http://127.0.0.1:7890"
	NetWork rpc.Cluster = rpc.Cluster{
		RPC: "https://mainnet.helius-rpc.com/?api-key=ce5ee933-a6ba-46b3-8e00-3f08bb2c49b1",
		WS:  "wss://mainnet.helius-rpc.com/?api-key=ce5ee933-a6ba-46b3-8e00-3f08bb2c49b1",
	}
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	option := gosolana.Option{
		RpcUrl: NetWork.RPC,
		WsUrl:  NetWork.WS,
		// Proxy:   Proxy,
		// WsProxy: Proxy,
		Pkey: "26HX8sewDP8Y6xTE3v4DtR5HHB5D4ua1MUxPEHyUA2j3SrFt4FDLwaXTWZg7BoeGGooyojtftUjR8CTMCczhQyrD",
	}
	wallet, err := gosolana.NewWallet(ctx, option)
	if err != nil {
		log.Fatal(err)
	}
	poolAddress := solana.MustPublicKeyFromBase58("zYTvAuRgJ5vxoCqqkR2PBkK7MkajznfGNVtPockGQ7L")
	out, err := wallet.GetClient().GetAccountInfo(ctx, poolAddress)
	if err != nil {
		log.Fatal("获取价格失败", err)
	}
	poolAmmV3 := new(amm_v3.PoolState)
	if err := bin.NewBinDecoder(out.GetBinary()).Decode(poolAmmV3); err != nil {
		log.Error("解析交易池子数据失败", err)
		return
	}

	log.Debugf("当前价格 | %.6f", poolAmmV3.Price())

}
