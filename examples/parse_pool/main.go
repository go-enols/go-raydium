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
