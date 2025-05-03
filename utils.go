package raydium

import (
	"errors"
	"fmt"
	"math"
	"math/big"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
)

// 创建关联代币账户指令
func MakeCreateAssociatedTokenAccountInstruction(
	payer solana.PublicKey,
	owner solana.PublicKey,
	mint solana.PublicKey,
) solana.Instruction {
	ata, _, _ := solana.FindProgramAddress(
		[][]byte{
			owner.Bytes(),
			solana.TokenProgramID.Bytes(),
			mint.Bytes(),
		},
		solana.SPLAssociatedTokenAccountProgramID,
	)

	accounts := solana.AccountMetaSlice{
		{PublicKey: payer, IsSigner: true, IsWritable: true},
		{PublicKey: ata, IsSigner: false, IsWritable: true},
		{PublicKey: owner, IsSigner: false, IsWritable: false},
		{PublicKey: mint, IsSigner: false, IsWritable: false},
		{PublicKey: solana.SystemProgramID, IsSigner: false, IsWritable: false},
		{PublicKey: solana.TokenProgramID, IsSigner: false, IsWritable: false},
		{PublicKey: solana.SysVarRentPubkey, IsSigner: false, IsWritable: false},
	}

	return solana.NewInstruction(
		solana.SPLAssociatedTokenAccountProgramID,
		accounts,
		[]byte{},
	)
}

// 获取关联代币账户地址
func GetAssociatedTokenAddress(wallet, mint solana.PublicKey) (solana.PublicKey, error) {
	seeds := [][]byte{
		wallet.Bytes(),
		solana.TokenProgramID.Bytes(),
		mint.Bytes(),
	}
	addr, _, err := solana.FindProgramAddress(seeds, solana.SPLAssociatedTokenAccountProgramID)
	if err != nil {
		return solana.PublicKey{}, fmt.Errorf("failed to find associated token address: %v", err)
	}
	return addr, nil
}

// 创建临时 WSOL 账户指令
func MakeCreateWSOLAccountInstructions(
	owner solana.PublicKey,
	wsolAccount solana.PublicKey,
	lamports uint64,
) []solana.Instruction {
	return []solana.Instruction{
		system.NewCreateAccountInstruction(
			lamports,
			165,
			solana.TokenProgramID,
			owner,
			wsolAccount,
		).Build(),
		token.NewInitializeAccountInstruction(
			wsolAccount,
			solana.SystemProgramID,
			owner,
			solana.TokenProgramID,
		).Build(),
	}
}

// 创建关闭账户指令
func MakeCloseAccountInstruction(
	account solana.PublicKey,
	destination solana.PublicKey,
	owner solana.PublicKey,
) solana.Instruction {
	return token.NewCloseAccountInstruction(
		account,
		destination,
		owner,
		[]solana.PublicKey{},
	).Build()
}

// 修复后的价格限制计算
func calculatePriceLimit(
	amountIn uint64, // 输入的代币数量（链上整数格式，例如 1000000 表示 1 USDC）
	decimals0 uint8, // 输入代币的小数位数（例如 USDC 是 6）
	slippage float64, // 滑点比例（例如 0.25 表示 25%）
) (bin.Uint128, error) {
	// 1. 将 amountIn 转换为实际代币单位（考虑小数位数）
	decimalsFactor := math.Pow10(int(decimals0)) // 10^decimals0
	amountInFloat := new(big.Float).SetUint64(amountIn)
	amountInActual := new(big.Float).Quo(amountInFloat, big.NewFloat(decimalsFactor))

	// 2. 应用滑点因子（例如 1 - 0.25 = 0.75）
	slippageFactor := new(big.Float).SetFloat64(1 - slippage)
	priceLimit := new(big.Float).Mul(amountInActual, slippageFactor)

	// 3. 转换回链上整数格式（根据输出代币的小数位数）
	// 假设输出代币的小数位数为 decimals1（需根据实际情况传入）
	decimals1 := uint8(9) // 例如 SOL 是 9 位
	outputFactor := math.Pow10(int(decimals1))
	priceLimit.Mul(priceLimit, big.NewFloat(outputFactor))

	// 4. 转换为整数（向下取整）
	priceLimitInt := new(big.Int)
	priceLimit.Int(priceLimitInt)
	return bigIntToUint128(priceLimitInt)
}

// 将 big.Int 转换为 binary.Uint128（处理高低位）
func bigIntToUint128(n *big.Int) (bin.Uint128, error) {
	// 验证范围
	if n.BitLen() > 128 || n.Sign() < 0 {
		return bin.Uint128{}, errors.New("数值超出 uint128 范围")
	}

	// 分解为高64位和低64位
	var (
		lo uint64
		hi uint64
	)
	// 复制 big.Int 避免修改原始值
	num := new(big.Int).Set(n)
	mask := new(big.Int).SetUint64(math.MaxUint64) // 0xFFFF_FFFF_FFFF_FFFF

	// 低64位
	lo = new(big.Int).And(num, mask).Uint64()
	// 右移64位获取高64位
	num.Rsh(num, 64)
	hi = num.Uint64()

	return bin.Uint128{Lo: lo, Hi: hi}, nil
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
