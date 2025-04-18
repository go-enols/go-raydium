# go-raydium

本项目为 Raydium（Solana 上的 AMM/DEX）池子与市场的 Go 语言解析与工具库，支持自动识别、解析、查询 Raydium V4、CPMM、cAMM 等多种池子类型，并可获取池子价格、流动性、元数据等信息。适合做链上数据分析、机器人、行情监控等二次开发。

## 目录结构

```
core/
  amm.go         // Raydium V4池结构体与解析
  camm.go        // cAMM池结构体与解析
  client.go      // 客户端与WebSocket监听
  cpmm.go        // CPMM池结构体与解析
  layout.go      // 池子类型自动识别与解析入口
  liquidity.go   // 池子价格、流动性、元数据查询
  types.go       // 通用类型定义
  util.go        // AMM公式、精度工具
script.go        // 监听日志、自动发现池子、价格输出等脚本
```

## 主要功能

- **池子类型自动识别**：支持 Raydium V4、CPMM、cAMM 等主流池型，自动解析链上账户数据。
- **池子价格与流动性查询**：自动获取池子的 base/quote 资产余额，支持 AMM 恒定乘积公式、Uniswap V3 价格公式。
- **元数据解析**：自动查询并解析 SPL Token 的 name、symbol、uri 等链上元数据。
- **WebSocket 日志监听**：支持监听 Raydium/OpenBook 相关链上日志，自动发现新池子并输出行情。
- **高精度计算**：支持大整数、128位小端字节序、AMM滑点、手续费等高精度处理。

## 快速开始

1. **安装依赖**

```bash
go get github.com/go-enols/go-raydium
go mod tidy
```

2. **监听并输出新池子价格**

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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

	script := raydium.NewScript(ctx, wallet.GetClient())
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
```

3. **解析交易Hash获取池子信息**

```go
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
```
4. **根据池子地址获取池子价格**
```go
package main

import (
	"context"
	"log"
	"time"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
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

	base, quote, price, err := core.GetPoolPriceByLiquidity(ctx, wallet.GetClient(), solana.MustPublicKeyFromBase58("74iTFH46SHuzD6YRVpFCGu911XMv2oqThqMyyZK9w7vX"))
	if err != nil {
		log.Fatal(err)
		return
	}
	log.Println("-----------------------------------------")
	log.Printf("%s | 合约地址 %s | Token %s", base.Sysbol, base.Address, base.Mint)
	log.Printf("%s | 合约地址 %s | Token %s", quote.Sysbol, quote.Address, quote.Mint)
	log.Printf("当前价格 | 1 %s=%.9f %s", quote.Sysbol, price, base.Sysbol)
	log.Printf("发现时间 | %s", time.Now().Format(time.DateTime))
}
```
## 关键实现说明

- **池子类型识别**：`core/layout.go` 的 `ParsePoolAccountByRPC` 会根据账户长度和字段特征自动判断池子类型。
- **价格计算**：
  - V4/CPMM：优先用 AMM 恒定乘积公式（见 `util.go`），兼容滑点和手续费。
  - cAMM：用 sqrtPriceX64 解析，严格按 Uniswap V3 公式高精度计算。
- **元数据获取**：通过 PDA 推导和链上查询，自动获取 SPL Token 的 name/symbol/uri。
- **WebSocket监听**：`core/client.go` 支持异步日志监听，`script.go` 可自动发现新池子并输出行情。

## 依赖

- [solana-go](https://github.com/gagliardetto/solana-go)
- [gosolana](https://github.com/go-enols/gosolana)
- 其它见 go.mod

## 适用场景

- Solana 链上行情机器人
- Raydium 池子监控与分析
- 自动发现新池子与套利
- DEX 数据可视化与研究

## 贡献与反馈

如有问题、建议或需求，欢迎提 Issue 或 PR！

---

**本项目仅供学习与研究，链上操作请注意风险。**