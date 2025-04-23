package main

import (
	"context"
	"log"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/go-enols/go-raydium"
	"github.com/go-enols/go-raydium/core"
	"github.com/go-enols/gosolana"
)

var (
	Proxy               = ""
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
	}
	// 创建Solana钱包实例
	wallet, err := gosolana.NewWallet(ctx, option)
	if err != nil {
		log.Fatal(err)
	}
	retx, tx, err := core.GetConfirmedTransaction(ctx, wallet.GetClient(), solana.MustSignatureFromBase58("5CVDHTjoXRw47MoQn6CYx7imiEinFyvNT7BY85LXpBSWYgNu7ErtEh8tA3rw8Za7qfkhwksbNGnuKmxupYtpYfpU"))
	if err != nil {
		log.Println("查询交易信息失败 |", err)
		return
	}
	poolAddress, err := raydium.ParseLpAddressByLogs(tx)
	if err != nil {
		return
	}
	base, quote, price, err := core.GetPoolPriceByLiquidity(ctx, wallet.GetClient(), poolAddress)
	if err != nil {
		log.Println(poolAddress.String(), "查询池子数据失败", err)
		return
	}
	log.Println("-----------------------------------------")
	log.Println("交易Hash", tx.Signatures)
	log.Println("发现一个新的raydium池子:", poolAddress.String(), base.Sysbol, "-", quote.Sysbol)
	log.Printf("%s | 合约地址 %s | Token %s", base.Sysbol, base.Address, base.Mint)
	log.Printf("%s | 合约地址 %s | Token %s", quote.Sysbol, quote.Address, quote.Mint)
	log.Printf("当前价格 | 1 %s=%.9f %s", quote.Sysbol, price, base.Sysbol)
	log.Printf("池子创建时间 | %s", retx.BlockTime.Time().Format(time.DateTime))
	log.Printf("发现时间 | %s", time.Now().Format(time.DateTime))
}
