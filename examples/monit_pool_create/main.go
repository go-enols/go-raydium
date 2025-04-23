package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gagliardetto/solana-go/rpc"
	"github.com/go-enols/go-raydium"
	"github.com/go-enols/go-raydium/core"
	"github.com/go-enols/gosolana"
)

var (
	Proxy               = ""
	NetWork rpc.Cluster = rpc.MainNetBeta
)

func logInfo(data raydium.PoolCreateInfo) {
	log.Println("-----------------------------------------")
	log.Println("监听到一个关键交易Hash", data.Tx.Signatures)
	log.Println("发现一个新的raydium池子:", data.PoolAddress.String(), data.Base.Sysbol, "-", data.Quote.Sysbol)
	log.Printf("%s | 合约地址 %s | Token %s", data.Base.Sysbol, data.Base.Address, data.Base.Mint)
	log.Printf("%s | 合约地址 %s | Token %s", data.Quote.Sysbol, data.Quote.Address, data.Quote.Mint)
	log.Printf("当前价格 | 1 %s=%.9f %s", data.Quote.Sysbol, data.Price, data.Quote.Sysbol)
	log.Printf("池子创建时间 | %s", data.Rtx.BlockTime.Time().Format(time.DateTime))
	log.Printf("发现时间 | %s", time.Now().Format(time.DateTime))
}

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

	script := raydium.NewMonitPoolCreateScipt(ctx, wallet.GetClient(), logInfo)
	raydiumClient := core.NewClient(ctx, option)
	raydiumClient.UseLog(script.RaydiumLogs)
	go raydiumClient.Start(ctx, core.RaydiumLiquidityProgramV4, rpc.CommitmentProcessed)

	openbook := core.NewClient(ctx, option)
	openbook.UseLog(script.OpenBookLogs)
	go openbook.Start(ctx, core.RaydiumLiquidityProgramV4, rpc.CommitmentProcessed)

	var stopChan = make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stopChan // wait for SIGINT

	fmt.Printf("正在退出...\n")
}
