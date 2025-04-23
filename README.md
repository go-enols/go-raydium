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
examples/        // 示例代码
```

## 主要功能

- **池子类型自动识别**：支持 Raydium V4、CPMM、cAMM 等主流池型，自动解析链上账户数据。
- **池子价格与流动性查询**：自动获取池子的 base/quote 资产余额，支持 AMM 恒定乘积公式、Uniswap V3 价格公式。
- **元数据解析**：自动查询并解析 SPL Token 的 name、symbol、uri 等链上元数据。
- **WebSocket 日志监听**：支持监听 Raydium/OpenBook 相关链上日志，自动发现新池子并输出行情。
- **高精度计算**：支持大整数、128 位小端字节序、AMM 滑点、手续费等高精度处理。

## 快速开始

1. **安装依赖**

```bash
go get github.com/go-enols/go-raydium
go mod tidy
```

2. **示例代码**

所有常用用法和场景的示例代码已移至 [`examples`](./examples) 文件夹。  
请参考以下链接获取详细示例：

- [监听并输出新池子价格](./examples/monit_pool_create/main.go)
- [解析交易 Hash 获取池子信息](./examples/parse_transfer_log_pool/main.go)
- [监听池子获取所有的 swap 交易并输出交易数据](./examples/swap-log/main.go)

你可以直接复制、运行或根据需要修改这些示例。

## 关键实现说明

- **池子类型识别**：`core/layout.go` 的 `ParsePoolAccountByRPC` 会根据账户长度和字段特征自动判断池子类型。
- **价格计算**：
  - V4/CPMM：优先用 AMM 恒定乘积公式（见 `util.go`），兼容滑点和手续费。
  - cAMM：用 sqrtPriceX64 解析，严格按 Uniswap V3 公式高精度计算。
- **元数据获取**：通过 PDA 推导和链上查询，自动获取 SPL Token 的 name/symbol/uri。
- **WebSocket 监听**：`core/client.go` 支持异步日志监听，`script.go` 可自动发现新池子并输出行情。

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
