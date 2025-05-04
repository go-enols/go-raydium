package main

import (
	"context"

	"github.com/go-enols/go-log"
	"github.com/go-enols/go-raydium"
	"github.com/go-enols/go-raydium/raydium_cp_swap"

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
	poolAddress := solana.MustPublicKeyFromBase58("5kuvtU1KT8gP92u4zPH9sprJzE9qPUJBw1tmThiVZKBf")

	out, err := wallet.GetClient().GetAccountInfo(ctx, poolAddress)
	if err != nil {
		log.Fatal("获取价格失败", err)
	}
	data := new(raydium_cp_swap.PoolState)
	if err := bin.NewBinDecoder(out.GetBinary()).Decode(data); err != nil {
		log.Error("解析交易池子数据失败", err)
		return
	}

	priceInfo, err := data.PriceInfo(ctx, wallet.GetClient())
	if err != nil {
		log.Fatal("获取价格信息失败", err)
		return
	}

	// 1 SOl = ? Token
	price := raydium.TokensForSol(1, priceInfo.BaseAmount, priceInfo.QuoteAmount, 0.25)
	log.Debug(priceInfo)
	log.Debugf("当前价格 | 1 %s = %.6f %s", priceInfo.BaseName, price, priceInfo.QuoteName)

	// 1 Token = ? Sol
	// price := raydium.SolForTokens(1, priceInfo.BaseAmount, priceInfo.QuoteAmount, 0.25)
	// log.Debug(priceInfo)
	// log.Debugf("当前价格 | 1 %s = %.6f %s", priceInfo.QuoteName, price, priceInfo.BaseName)
}
