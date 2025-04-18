package core

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/go-enols/gosolana"
)

type Liquidity struct {
	Address  solana.PublicKey
	Mint     solana.PublicKey
	Decimals int
	Value    uint64 // CAMM 类型的池子 存在128位超大整数所以不写入这里
	Name     string
	Sysbol   string
}

func (l *Liquidity) MintAddress() solana.PublicKey {
	return l.Mint
}

func (l *Liquidity) MintAddressString() string {
	return l.Mint.String()
}

// 查询池子的价格（自动识别类型并查余额）
func GetPoolPriceByLiquidity(ctx context.Context, client *rpc.Client, poolPubkey solana.PublicKey) (base, quote *Liquidity, price float64, err error) {
	result, err := ParsePoolAccountByRPC(ctx, client, poolPubkey)
	if err != nil {
		return nil, nil, 0, err
	}

	var vault0, vault1 solana.PublicKey
	var mint0, mint1 solana.PublicKey
	var decimals0, decimals1 int

	switch result.Type {
	case PoolTypeV4:
		log.Println("池类型: RaydiumV4")
		vault0 = result.V4.BaseVault
		vault1 = result.V4.QuoteVault
		mint0 = result.V4.BaseMint
		mint1 = result.V4.QuoteMint
		decimals0 = int(result.V4.BaseDecimal)
		decimals1 = int(result.V4.QuoteDecimal)
	case PoolTypeCPMM:
		log.Println("池类型: RaydiumCPMM")
		vault0 = result.CPMM.Token0Vault
		vault1 = result.CPMM.Token1Vault
		mint0 = result.CPMM.Token0Mint
		mint1 = result.CPMM.Token1Mint
		decimals0 = int(result.CPMM.Mint0Decimals)
		decimals1 = int(result.CPMM.Mint1Decimals)
	case PoolTypeCAMM:
		log.Println("池类型: RaydiumCAMM")
		// cAMM直接用 sqrtPriceX64 算价格，无需查余额
		price = calcCammPrice(result.CAMM)
		meta0, err := gosolana.GetTokenMetaOnChain(ctx, client, result.CAMM.TokenMint0)
		if err != nil {
			return nil, nil, 0, err
		}
		meta1, err := gosolana.GetTokenMetaOnChain(ctx, client, result.CAMM.TokenMint1)
		if err != nil {
			return nil, nil, 0, err
		}
		base = &Liquidity{
			Mint:     result.CAMM.TokenMint0,
			Decimals: int(result.CAMM.MintDecimals0),
			Name:     meta0.Name,
			Sysbol:   meta0.Symbol,
		}
		quote = &Liquidity{
			Mint:     result.CAMM.TokenMint1,
			Decimals: int(result.CAMM.MintDecimals1),
			Name:     meta1.Name,
			Sysbol:   meta1.Symbol,
		}
		return base, quote, price, nil
	default:
		return nil, nil, 0, errors.New("unsupported pool type")
	}

	// 查询两个 vault 的余额
	amounts, err := GetMultipleAccountsBalances(ctx, client, []solana.PublicKey{vault0, vault1})
	if err != nil || len(amounts) != 2 || amounts[0] == nil || amounts[1] == nil {
		return nil, nil, 0, errors.New("failed to get vault balances")
	}
	amount0 := amounts[0].Parsed.Info.TokenAmount.UIAmount
	amount1 := amounts[1].Parsed.Info.TokenAmount.UIAmount
	meta0, err := gosolana.GetTokenMetaOnChain(ctx, client, mint0)
	if err != nil {
		return nil, nil, 0, err
	}
	meta1, err := gosolana.GetTokenMetaOnChain(ctx, client, mint1)
	if err != nil {
		return nil, nil, 0, err
	}
	base = &Liquidity{
		Mint:     mint0,
		Address:  vault0,
		Value:    uint64(amount0 * pow10f(decimals0)),
		Name:     meta0.Name,
		Sysbol:   meta0.Symbol,
		Decimals: decimals0, // 保留此行（正确来源
	}
	quote = &Liquidity{
		Mint:     mint1,
		Address:  vault1,
		Value:    uint64(amount1 * pow10f(decimals1)),
		Name:     meta1.Name,
		Sysbol:   meta1.Symbol,
		Decimals: decimals1,
	}

	// 价格 = quote/base
	if amount0 == 0 {
		return base, quote, 0, errors.New("token0 vault is empty")
	}

	if base.Mint.String() != WSOL {
		base, quote = quote, base
	}
	price = SolForTokens(1, amount1, amount0, 0.25)
	return base, quote, price, nil
}

// cAMM池子价格计算（严格Uniswap V3风格）
func calcCammPrice(camm *CAMM_STATE_LAYOUT) float64 {
	pc := NewPrecisionCalculator()
	priceRat := pc.CalculatePrice(camm)
	priceFloat, _ := priceRat.Float64()
	return priceFloat
}

// 查询多个 SPL Token 账户余额（返回 UiAmount 数组，顺序与输入一致）
func GetMultipleAccountsBalances(ctx context.Context, client *rpc.Client, accounts []solana.PublicKey) ([]*AmmV4Account, error) {
	balancesResponse, err := client.GetMultipleAccountsWithOpts(
		ctx, accounts, &rpc.GetMultipleAccountsOpts{
			Encoding:   solana.EncodingJSONParsed,
			Commitment: rpc.CommitmentProcessed,
		},
	)
	if err != nil || len(balancesResponse.Value) != 2 || balancesResponse.Value[0] == nil || balancesResponse.Value[1] == nil {
		return nil, errors.New("failed to get vault balances")
	}
	amount0 := new(AmmV4Account)
	if err := json.Unmarshal(balancesResponse.Value[0].Data.GetRawJSON(), amount0); err != nil {
		log.Printf("解析 Quote 账户数据失败: %v, 原始数据: %s", err, string(balancesResponse.Value[0].Data.GetRawJSON()))
		return nil, err
	}
	amount1 := new(AmmV4Account)
	if err := json.Unmarshal(balancesResponse.Value[1].Data.GetRawJSON(), amount1); err != nil {
		log.Printf("解析 Base 账户数据失败: %v, 原始数据: %s", err, string(balancesResponse.Value[1].Data.GetRawJSON()))
		return nil, err
	}
	return []*AmmV4Account{
		amount0, amount1,
	}, nil
}

// 浮点型10的n次方
func pow10f(n int) float64 {
	res := 1.0
	for i := 0; i < n; i++ {
		res *= 10
	}
	return res
}
