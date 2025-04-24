package ammV4

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/programs/system"
	"github.com/gagliardetto/solana-go/programs/token"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/go-enols/go-log"
	"github.com/go-enols/go-raydium/core"
	"github.com/go-enols/gosolana"
)

const (
	ACCOUNT_LAYOUT_LEN = 165
	UNIT_BUDGET        = 400000
	UNIT_PRICE         = 1
)

type Client struct {
	*gosolana.Wallet
}

func NewClient(ctx context.Context, option ...gosolana.Option) (*Client, error) {
	wallet, err := gosolana.NewWallet(ctx, option...)
	return &Client{
		Wallet: wallet,
	}, err
}

func (c *Client) Buy(ctx context.Context, poolAddr string, solIn float64, slippage int) (bool, error) {
	poolAddress := solana.MustPublicKeyFromBase58(poolAddr)
	amountIn := uint64(solIn * float64(solana.LAMPORTS_PER_SOL))
	base, quote, price, err := core.GetPoolPriceByLiquidity(ctx, c.GetClient(), poolAddress)
	if err != nil {
		return false, err
	}

	log.Debugf("当前池价格 1 %s ≈ %.6f %s", base.Sysbol, price, quote.Sysbol)
	if base.Mint != solana.SystemProgramID && quote.Mint != solana.SystemProgramID && base.Mint != solana.WrappedSol && quote.Mint != solana.WrappedSol {
		log.Debug(base.Mint.String())
		log.Debug(quote.Mint.String())
		return false, errors.New("该池无法使用sol购买")
	}

	if base.Mint != solana.SystemProgramID && base.Mint != solana.WrappedSol {
		base, quote = quote, base
	}
	amountOut := solIn * price
	slippageAdjustment := 1 - (float64(slippage) / 100)
	amountOutWithSlippage := amountOut * slippageAdjustment
	minimumAmountOut := uint64(amountOutWithSlippage * math.Pow10(int(quote.Decimals)))

	_ = amountIn
	_ = base
	_ = minimumAmountOut
	log.Debugf("消耗 -> %v %s 购买 -> %v %s | 池地址: %s", solIn, base.Sysbol, solIn*price, quote.Sysbol, poolAddr)

	// 创建临时 WSOL 账户
	seed := make([]byte, 24)
	if _, err := rand.Read(seed); err != nil {
		return false, fmt.Errorf("failed to generate random seed: %v", err)
	}
	wsolAccountSeed := base64.URLEncoding.EncodeToString(seed)

	// 创建 WSOL 账户地址
	wsolAccount, _, err := solana.FindProgramAddress(
		[][]byte{
			c.PublicKey().Bytes(),
			[]byte(wsolAccountSeed),
		},
		solana.TokenProgramID,
	)
	if err != nil {
		return false, fmt.Errorf("failed to create WSOL account address: %v", err)
	}

	tokenAccount, err := GetAssociatedTokenAddress(c.PublicKey(), quote.Mint)
	if err != nil {
		return false, fmt.Errorf("failed to get token account: %v", err)
	}

	// 检查代币账户是否存在
	_, err = c.GetClient().GetAccountInfo(ctx, tokenAccount)
	createTokenAccount := err != nil
	log.Debug(wsolAccount.String())
	log.Debug(tokenAccount.String())
	log.Debug(createTokenAccount)

	// 构建交易指令
	instructions := make([]solana.Instruction, 0)
	// 添加计算预算指令
	budgetInstructions, err := MakeComputeBudgetInstruction(UNIT_BUDGET, UNIT_PRICE)
	if err != nil {
		return false, fmt.Errorf("failed to create budget instructions: %v", err)
	}
	instructions = append(instructions, budgetInstructions...)
	// 创建 WSOL 账户指令
	rentExemptBalance, err := c.GetClient().GetMinimumBalanceForRentExemption(ctx, ACCOUNT_LAYOUT_LEN, rpc.CommitmentConfirmed)
	if err != nil {
		return false, fmt.Errorf("failed to get rent exempt balance: %v", err)
	}
	wsolInstructions := MakeCreateWSOLAccountInstructions(
		c.PublicKey(),
		wsolAccount,
		rentExemptBalance+uint64(amountIn),
	)
	instructions = append(instructions, wsolInstructions...)

	// 如果需要，创建代币账户
	if createTokenAccount {
		createTokenAccountInstruction := MakeCreateAssociatedTokenAccountInstruction(
			c.PublicKey(),
			c.PublicKey(),
			quote.Mint,
		)
		instructions = append(instructions, createTokenAccountInstruction)
	}

	// 添加交换指令
	swapInstruction, err := MakeAmmV4SwapInstruction(
		amountIn,
		minimumAmountOut,
		wsolAccount,
		tokenAccount,
		poolAddress,
		base, quote,
		c.PublicKey(),
	)
	if err != nil {
		return false, fmt.Errorf("failed to create swap instruction: %v", err)
	}
	instructions = append(instructions, swapInstruction)

	// 添加关闭 WSOL 账户指令
	closeWsolInstruction := MakeCloseAccountInstruction(
		wsolAccount,
		c.PublicKey(),
		c.PublicKey(),
	)
	instructions = append(instructions, closeWsolInstruction)
	ok, err := c.SendTransaction(ctx, instructions)
	return ok, err
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

// 创建计算预算指令
func MakeComputeBudgetInstruction(units uint32, price uint64) ([]solana.Instruction, error) {
	computeBudgetProgramID, err := solana.PublicKeyFromBase58("ComputeBudget111111111111111111111111111111")
	if err != nil {
		return nil, err
	}

	unitData := make([]byte, 4)
	binary.LittleEndian.PutUint32(unitData, units)

	priceData := make([]byte, 8)
	binary.LittleEndian.PutUint64(priceData, price)

	return []solana.Instruction{
		solana.NewInstruction(
			computeBudgetProgramID,
			solana.AccountMetaSlice{},
			append([]byte{0}, unitData...), // SetComputeUnitLimit
		),
		solana.NewInstruction(
			computeBudgetProgramID,
			solana.AccountMetaSlice{},
			append([]byte{1}, priceData...), // SetComputeUnitPrice
		),
	}, nil
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

// 创建 AMM V4 交换指令
func MakeAmmV4SwapInstruction(
	amountIn uint64,
	minimumAmountOut uint64,
	tokenAccountIn solana.PublicKey,
	tokenAccountOut solana.PublicKey,
	pool solana.PublicKey,
	base, quote *core.Liquidity,
	owner solana.PublicKey,
) (solana.Instruction, error) {

	accounts := solana.AccountMetaSlice{
		{PublicKey: pool, IsSigner: false, IsWritable: true},
		{PublicKey: tokenAccountIn, IsSigner: false, IsWritable: true},
		{PublicKey: tokenAccountOut, IsSigner: false, IsWritable: true},
		{PublicKey: base.Mint, IsSigner: false, IsWritable: true},
		{PublicKey: quote.Mint, IsSigner: false, IsWritable: true},
		{PublicKey: owner, IsSigner: true, IsWritable: false},
		{PublicKey: solana.TokenProgramID, IsSigner: false, IsWritable: false},
	}

	data := make([]byte, 17)
	data[0] = 2 // Swap instruction
	binary.LittleEndian.PutUint64(data[1:9], amountIn)
	binary.LittleEndian.PutUint64(data[9:], minimumAmountOut)

	return solana.NewInstruction(core.RaydiumLiquidityProgramV4, accounts, data), nil
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
