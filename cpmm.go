package raydium

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"math"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/go-enols/go-log"
	"github.com/go-enols/go-raydium/raydium_cp_swap"
	"github.com/go-enols/gosolana"
)

type CpmmClient struct {
	*gosolana.Wallet
}

func NewCpmmClient(wallet *gosolana.Wallet) *CpmmClient {
	return &CpmmClient{
		Wallet: wallet,
	}
}

func (c *CpmmClient) Swap(ctx context.Context, solIn float64, poolPub solana.PublicKey, isBuy bool, slippage float64, poolDatas ...[]byte) (*solana.Transaction, error) {
	var poolData = new(raydium_cp_swap.PoolState)
	if len(poolDatas) > 0 {
		if err := bin.NewBinDecoder(poolDatas[0]).Decode(poolData); err != nil {
			return nil, err
		}
	} else {
		out, err := c.GetClient().GetAccountInfo(ctx, poolPub)
		if err != nil {
			log.Errorf("获取池数据失败 | %s", err)
			return nil, err
		}
		if err := bin.NewBinDecoder(out.GetBinary()).Decode(poolData); err != nil {
			return nil, err
		}
	}
	data, err := c.CreateInstruction(solIn, poolPub, poolData, isBuy, slippage)
	if err != nil {
		log.Printf("创建指令失败 | %s", err)
		return nil, err
	}
	recentBlockHash, err := c.GetClient().GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		log.Printf("获取Hash失败 | %s", err)
		return nil, err
	}
	// 构造交易
	tx, err := solana.NewTransaction(
		data,
		recentBlockHash.Value.Blockhash,
		solana.TransactionPayer(c.PublicKey()),
	)
	if err != nil {
		log.Printf("构建交易失败 | %s", err)
		return nil, err
	}
	return tx, nil
}

func (c *CpmmClient) CreateInstruction(
	amountIn float64,
	poolAddr solana.PublicKey,
	poolCpmm *raydium_cp_swap.PoolState,
	isBuy bool,
	slippage float64,
	opt ...struct {
		Limit uint32
		Price uint64
	},
) ([]solana.Instruction, error) {
	var result []solana.Instruction
	var price uint64 = 1_000_000
	var limit uint32 = 100_000
	if len(opt) > 0 {
		price = opt[0].Price
		limit = opt[0].Limit
	}

	swap, err := c.createSwapInstruction(
		amountIn, poolAddr, poolCpmm, isBuy, slippage,
	)
	if err != nil {
		return nil, err
	}
	result = append(result, c.createSetComputeUnitPriceInstruction(price))
	result = append(result, c.createSetComputeUnitLimitInstruction(limit))
	result = append(result, swap...)
	return result, nil
}

func (c *CpmmClient) createSetComputeUnitPriceInstruction(price uint64) *computebudget.Instruction {
	return computebudget.NewSetComputeUnitPriceInstruction(price).Build()
}
func (c *CpmmClient) createSetComputeUnitLimitInstruction(limit uint32) *computebudget.Instruction {
	return computebudget.NewSetComputeUnitLimitInstruction(limit).Build()
}

func (c *CpmmClient) createSwapInstruction(
	amount float64,
	poolAddr solana.PublicKey,
	poolCpmm *raydium_cp_swap.PoolState,
	isBuy bool,
	slippage float64,
) ([]solana.Instruction, error) {
	var result []solana.Instruction
	var vault0, vault1 solana.PublicKey = poolCpmm.Token0Vault, poolCpmm.Token1Vault
	var mint0, mint1 solana.PublicKey = poolCpmm.Token0Mint, poolCpmm.Token1Mint
	var inputTokenProgram, outputTokenProgram = poolCpmm.Token0Program, poolCpmm.Token1Program

	priceInfo, err := poolCpmm.PriceInfo(context.TODO(), c.GetClient())
	if err != nil {
		return nil, err
	}
	if !poolCpmm.Token0Mint.Equals(WSOL) {
		priceInfo.BaseAmount, priceInfo.QuoteAmount = priceInfo.QuoteAmount, priceInfo.BaseAmount
		priceInfo.BaseDecimals, priceInfo.QuoteDecimals = priceInfo.QuoteDecimals, priceInfo.BaseDecimals
		priceInfo.BaseName, priceInfo.QuoteName = priceInfo.QuoteName, priceInfo.BaseName
		mint0, mint1 = mint1, mint0
		vault0, vault1 = vault1, vault0
		inputTokenProgram, outputTokenProgram = outputTokenProgram, inputTokenProgram
	}
	var price float64
	var amountIn, miniAmountOut uint64
	var amountOut float64
	if isBuy {
		price = SolForTokens(1, priceInfo.BaseAmount, priceInfo.QuoteAmount, 0.25)
		amountIn = uint64(amount * math.Pow10(int(priceInfo.BaseDecimals)))
		amountOut = amount / price * (1 - slippage)
		miniAmountOut = uint64(amountOut * math.Pow10(int(priceInfo.QuoteDecimals)))
		log.Debugf("BUY %.6f %s -> %.6f %s ", amount, priceInfo.BaseName, amountOut, priceInfo.QuoteName)
	} else {
		price = TokensForSol(1, priceInfo.BaseAmount, priceInfo.QuoteAmount, 0.25)
		amountIn = uint64(amount * math.Pow10(int(priceInfo.QuoteDecimals)))
		amountOut = amount / price * (1 - slippage)
		miniAmountOut = uint64(amountOut * math.Pow10(int(priceInfo.BaseDecimals)))
		log.Debugf("SHELL %.6f %s -> %.6f %s ", amount, priceInfo.QuoteName, amountOut, priceInfo.BaseName)
		vault0, vault1 = vault1, vault0                                               //反转tokenValue账户
		inputTokenProgram, outputTokenProgram = outputTokenProgram, inputTokenProgram // 反转程序账户
		mint0, mint1 = mint1, mint0                                                   // 反转Mint账户
	}

	crteteAccount, closeAccount, account, err := c.createAccountInstruction(amountIn)
	if err != nil {
		return nil, err
	}
	log.Debugf("创建 Base 过渡账户 | %s", account.String())
	tokenCreate, tokenAccount, err := c.checkTokenAccount(mint1)
	if err != nil {
		return nil, err
	}
	result = append(result, crteteAccount...)
	if tokenCreate != nil {
		result = append(result, tokenCreate)
	}

	// 输入token的账户和输出token的账户
	var inputTokenAccount, outputTokenAccount solana.PublicKey
	if isBuy {
		inputTokenAccount, outputTokenAccount = account, tokenAccount
	} else {
		inputTokenAccount, outputTokenAccount = tokenAccount, account
	}

	result = append(
		result,
		raydium_cp_swap.NewSwapBaseInputInstruction(
			amountIn, miniAmountOut,
			c.PublicKey(),
			CpmmAuthority,
			poolCpmm.AmmConfig,
			poolAddr,
			inputTokenAccount,
			outputTokenAccount,
			vault0,
			vault1,
			inputTokenProgram,
			outputTokenProgram,
			mint0,
			mint1,
			poolCpmm.ObservationKey,
		).Build(),
	)
	result = append(result, closeAccount)

	return result, nil
}

// 构建创建临时WSOL账户的指令
func (v *CpmmClient) createAccountInstruction(amountIn uint64) ([]solana.Instruction, solana.Instruction, solana.PublicKey, error) {

	// 创建临时 WSOL 账户
	seed := make([]byte, 24)
	if _, err := rand.Read(seed); err != nil {
		return nil, nil, solana.PublicKey{}, fmt.Errorf("failed to generate random seed: %v", err)
	}
	wsolAccountSeed := base64.URLEncoding.EncodeToString(seed)

	// 创建 WSOL 账户地址
	wsolAccount, _, err := solana.FindProgramAddress(
		[][]byte{
			v.PublicKey().Bytes(),
			[]byte(wsolAccountSeed),
		},
		solana.TokenProgramID,
	)
	if err != nil {
		return nil, nil, solana.PublicKey{}, fmt.Errorf("failed to create WSOL account address: %v", err)
	}

	// 获取租金余额
	rentExemptBalance, err := v.GetClient().GetMinimumBalanceForRentExemption(
		context.TODO(),
		ACCOUNT_LAYOUT_LEN,
		rpc.CommitmentConfirmed,
	)
	if err != nil {
		return nil, nil, solana.PublicKey{}, fmt.Errorf("failed to get rent exempt balance: %v", err)
	}

	return MakeCreateWSOLAccountInstructions(
		v.PublicKey(),
		wsolAccount,
		rentExemptBalance+uint64(amountIn),
	), MakeCloseAccountInstruction(wsolAccount, v.PublicKey(), v.PublicKey()), wsolAccount, nil
}

// 检查是否有对应mint地址的账户,如果没有则创建
func (v *CpmmClient) checkTokenAccount(mint solana.PublicKey) (solana.Instruction, solana.PublicKey, error) {
	tokenAccount, err := GetAssociatedTokenAddress(v.PublicKey(), mint)
	if err != nil {
		return nil, solana.PublicKey{}, fmt.Errorf("failed to get token account: %v", err)
	}

	// 检查代币账户是否存在
	_, err = v.GetClient().GetAccountInfo(context.TODO(), tokenAccount)
	createTokenAccount := err != nil
	// 如果代币账户不存在，创建代币账户
	if createTokenAccount {
		createTokenAccountInstruction := MakeCreateAssociatedTokenAccountInstruction(
			v.PublicKey(),
			v.PublicKey(),
			mint,
		)
		log.Debugf("创建 Quote 过渡账户 | %s", tokenAccount.String())
		return createTokenAccountInstruction, tokenAccount, nil
	}
	return nil, tokenAccount, nil
}
