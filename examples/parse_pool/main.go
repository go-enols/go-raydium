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
	ammv3 := solana.MustPublicKeyFromBase58("zYTvAuRgJ5vxoCqqkR2PBkK7MkajznfGNVtPockGQ7L")
	ammv4 := solana.MustPublicKeyFromBase58("Bzc9NZfMqkXR6fz1DBph7BDf9BroyEf6pnzESP7v5iiw")
	cpmm := solana.MustPublicKeyFromBase58("5kuvtU1KT8gP92u4zPH9sprJzE9qPUJBw1tmThiVZKBf")

	var res *raydium.ParsePoolResult
	res, _ = raydium.ParsePool(wallet.GetClient(), ammv3)
	log.Debug("zYTvAuRgJ5vxoCqqkR2PBkK7MkajznfGNVtPockGQ7L", res.Types)
	res, _ = raydium.ParsePool(wallet.GetClient(), ammv4)
	log.Debug("Bzc9NZfMqkXR6fz1DBph7BDf9BroyEf6pnzESP7v5iiw", res.Types)
	res, _ = raydium.ParsePool(wallet.GetClient(), cpmm)
	log.Debug("5kuvtU1KT8gP92u4zPH9sprJzE9qPUJBw1tmThiVZKBf", res.Types)

}
