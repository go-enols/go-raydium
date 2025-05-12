package main

import (
	"context"

	"github.com/go-enols/go-log"
	"github.com/go-enols/go-raydium"
	"github.com/go-enols/go-raydium/raydium_amm"

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
	poolAddress := solana.MustPublicKeyFromBase58("Bzc9NZfMqkXR6fz1DBph7BDf9BroyEf6pnzESP7v5iiw")

	out, err := wallet.GetClient().GetAccountInfo(ctx, poolAddress)
	if err != nil {
		log.Fatal("获取价格失败", err)
	}
	data := new(raydium_amm.AmmInfo)
	if err := bin.NewBinDecoder(out.GetBinary()).Decode(data); err != nil {
		log.Error("解析交易池子数据失败", err)
		return
	}

	priceInfo, err := data.PriceInfo(ctx, wallet.GetClient())
	if err != nil {
		log.Fatal("获取价格信息失败", err)
		return
	}

	log.Debug(priceInfo)
	// 1 SOL = ? Token
	price := raydium.TokensForSol(1, priceInfo.BaseAmount, priceInfo.QuoteAmount, 0.25)
	log.Debugf("当前价格 | 1 %s = %.6f %s", priceInfo.BaseName, price, priceInfo.QuoteName)

	// 1 Token = ? Sol
	// price := raydium.SolForTokens(1, priceInfo.BaseAmount, priceInfo.QuoteAmount, 0.25)
	// log.Debugf("当前价格 | 1 %s = %.6f %s", priceInfo.QuoteName, price, priceInfo.BaseName)
}
