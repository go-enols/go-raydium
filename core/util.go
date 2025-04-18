package core

import (
	"encoding/binary"
	"math"
	"math/big"
)

// SolForTokens 计算用SOL换取代币的数量（AMM恒定乘积公式）
func SolForTokens(solAmount, baseVaultBalance, quoteVaultBalance, swapFee float64) float64 {
	// 计算有效使用的SOL（扣除手续费）
	effectiveSolUsed := solAmount - (solAmount * (swapFee / 100))

	// 计算恒定乘积
	constantProduct := baseVaultBalance * quoteVaultBalance

	// 计算新的基础代币储备
	updatedBase := constantProduct / (quoteVaultBalance + effectiveSolUsed)

	// 计算实际获得的代币数量
	tokensReceived := baseVaultBalance - updatedBase

	// 四舍五入到9位小数
	return tokensReceived
}

// TokensForSol 计算用代币换取SOL的数量（AMM恒定乘积公式）
func TokensForSol(tokenAmount, baseVaultBalance, quoteVaultBalance, swapFee float64) float64 {
	// 计算有效卖出的代币数量（扣除手续费）
	effectiveTokensSold := tokenAmount * (1 - (swapFee / 100))

	// 计算恒定乘积
	constantProduct := baseVaultBalance * quoteVaultBalance

	// 计算新的SOL储备
	updatedQuote := constantProduct / (baseVaultBalance + effectiveTokensSold)

	// 计算实际获得的SOL数量
	solReceived := quoteVaultBalance - updatedQuote

	// 四舍五入到9位小数
	return solReceived
}

// 精度计算核心模块
type PrecisionCalculator struct {
	power128 *big.Int     // 2^128缓存
	decimals [19]*big.Int // 10^0到10^18的缓存
}

func NewPrecisionCalculator() *PrecisionCalculator {
	pc := &PrecisionCalculator{
		power128: new(big.Int).Exp(big.NewInt(2), big.NewInt(128), nil),
		decimals: [19]*big.Int{},
	}
	for i := 0; i <= 18; i++ {
		pc.decimals[i] = new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(i)), nil)
	}
	return pc
}

// 小端序转换核心方法
func (pc *PrecisionCalculator) leBytesToUint128(b [16]byte) *big.Int {
	var lo, hi uint64
	lo = binary.LittleEndian.Uint64(b[:8])
	hi = binary.LittleEndian.Uint64(b[8:])
	return new(big.Int).Or(
		new(big.Int).Lsh(new(big.Int).SetUint64(hi), 64),
		new(big.Int).SetUint64(lo),
	)
}

// 分层精度计算流程
func (pc *PrecisionCalculator) CalculatePrice(camm *CAMM_STATE_LAYOUT) *big.Rat {
	// 第一层：原始数据转换
	sqrtPrice := pc.leBytesToUint128(camm.SqrtPriceX64)

	// 第二层：数学运算层
	sqrtPriceSquared := new(big.Int).Mul(sqrtPrice, sqrtPrice)

	// 第三层：精度调整层
	decimals0 := int(math.Min(float64(camm.MintDecimals0), 18))
	decimals1 := int(math.Min(float64(camm.MintDecimals1), 18))

	numerator := new(big.Int).Mul(sqrtPriceSquared, pc.decimals[decimals0])
	denominator := new(big.Int).Mul(pc.power128, pc.decimals[decimals1])

	// 第四层：最终比率生成
	return new(big.Rat).SetFrac(numerator, denominator)
}
func (pc *PrecisionCalculator) StringPrice(camm *CAMM_STATE_LAYOUT) string {
	rat := pc.CalculatePrice(camm)
	return rat.FloatString(int(math.Max(float64(camm.MintDecimals0), float64(camm.MintDecimals1))))
}
func (pc *PrecisionCalculator) BatchCalculate(pools []*CAMM_STATE_LAYOUT) []*big.Rat {
	results := make([]*big.Rat, len(pools))
	for i, p := range pools {
		results[i] = pc.CalculatePrice(p)
	}
	return results
}
