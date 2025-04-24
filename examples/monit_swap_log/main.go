package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-enols/go-log"

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

func testSwapIngo(data raydium.SwapTransferInfo) {
	log.Println("-----------------------------------------")
	log.Println("监听到一个关键交易Hash", data.Tx.Signatures)

	for Owner, value := range data.Data {
		for mint, v := range value.Data {
			log.Debugf("%s 进行了 Swap %s -> %f", Owner, mint, v.OutAmount-v.InputAmount)
		}
	}

	log.Printf("创建时间 | %s", data.Rtx.BlockTime.Time().Format(time.DateTime))
	log.Printf("发现时间 | %s", time.Now().Format(time.DateTime))
}

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	option := gosolana.Option{
		RpcUrl: NetWork.RPC,
		WsUrl:  NetWork.WS,
		// Proxy:   Proxy,
		// WsProxy: Proxy,
	}
	// 创建Solana钱包实例
	wallet, err := gosolana.NewWallet(ctx, option)
	if err != nil {
		log.Fatal(err)
	}

	script := raydium.NewMonitSwapTransferScipt(ctx, wallet.GetClient(), testSwapIngo)
	raydiumClient := core.NewClient(ctx, option)
	raydiumClient.UseLog(script.RaydiumSwapLogs)
	go raydiumClient.Start(ctx, solana.MustPublicKeyFromBase58("池地址"), rpc.CommitmentProcessed)

	var stopChan = make(chan os.Signal, 1)
	signal.Notify(stopChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	<-stopChan // wait for SIGINT

	fmt.Printf("正在退出...\n")
}
