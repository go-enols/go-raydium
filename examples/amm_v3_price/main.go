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
