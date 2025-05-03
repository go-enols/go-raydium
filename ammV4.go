package raydium

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/binary"
	"errors"
	"fmt"
	"math"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	computebudget "github.com/gagliardetto/solana-go/programs/compute-budget"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/go-enols/go-log"

	"github.com/go-enols/go-raydium/raydium_amm"
	"github.com/go-enols/gosolana"
)

type V4Client struct {
	*gosolana.Wallet
}

func NewV4Client(wallet *gosolana.Wallet) *V4Client {
	return &V4Client{
		Wallet: wallet,
	}
}

func (v *V4Client) Swap(ctx context.Context, amount float64, poolPub solana.PublicKey, isBuy bool, slippage float64, opt ...struct {
	Pool   []byte
	Market []byte
}) (*solana.Transaction, error) {
	var poolData []byte
	var marketData []byte
	if len(opt) > 0 {
		poolData = opt[0].Pool
		marketData = opt[0].Market
	} else {
		out, err := v.GetClient().GetAccountInfo(ctx, poolPub)
		if err != nil {
			log.Fatal(err)
		}
		poolData = out.GetBinary()
	}

	poolAmmV4 := new(raydium_amm.AmmInfo)

	if err := bin.NewBinDecoder(poolData).Decode(poolAmmV4); err != nil {
		log.Errorf("解析交易池子数据失败 | %s", err)
		return nil, err
	}
	openBook := new(OpenBook)
	if len(marketData) == 0 {
		log.Debug(poolAmmV4.Market.String())
		out, err := v.GetClient().GetAccountInfo(ctx, poolAmmV4.Market)
		if err != nil {
			log.Fatal(err)
		}
		if err := openBook.UnmarshalBinary(out.GetBinary()); err != nil {
			log.Errorf("解析市场数据失败 | %s", err)
			return nil, err
		}
	} else {
		if err := openBook.UnmarshalBinary(marketData); err != nil {
			log.Errorf("解析市场数据失败 | %s", err)
			return nil, err
		}
	}

	data, err := v.CreateInstruction(amount, poolPub, poolAmmV4, openBook, isBuy, slippage)
	if err != nil {
		log.Printf("创建指令失败 | %s", err)
		return nil, err
	}
	recentBlockHash, err := v.GetClient().GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		log.Printf("获取Hash失败 | %s", err)
		return nil, err
	}
	// 构造交易
	tx, err := solana.NewTransaction(
		data,
		recentBlockHash.Value.Blockhash,
		solana.TransactionPayer(v.PublicKey()),
	)
	if err != nil {
		log.Printf("构建交易失败 | %s", err)
		return nil, err
	}
	return tx, nil
}

func (v *V4Client) CreateInstruction(
	amountIn float64,
	poolAddr solana.PublicKey,
	poolAmmV4 *raydium_amm.AmmInfo,
	market *OpenBook,
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

	swap, err := v.createSwapInstruction(
		amountIn, poolAddr, poolAmmV4, market, isBuy, slippage,
	)
	if err != nil {
		return nil, err
	}
	result = append(result, v.createSetComputeUnitPriceInstruction(price))
	result = append(result, v.createSetComputeUnitLimitInstruction(limit))
	result = append(result, swap...)
	return result, nil
}

func (v *V4Client) createSetComputeUnitPriceInstruction(price uint64) *computebudget.Instruction {
	return computebudget.NewSetComputeUnitPriceInstruction(price).Build()
}
func (v *V4Client) createSetComputeUnitLimitInstruction(limit uint32) *computebudget.Instruction {
	return computebudget.NewSetComputeUnitLimitInstruction(limit).Build()
}

func (v *V4Client) createSwapInstruction(
	amount float64,
	poolAddr solana.PublicKey,
	poolAmmV4 *raydium_amm.AmmInfo,
	market *OpenBook,
	isBuy bool,
	slippage float64,
) ([]solana.Instruction, error) {
	var result []solana.Instruction
	log.Debug("LP Value | ", poolAddr)
	log.Debug("LP Mint | ", poolAmmV4.LpMint)
	var vault0, vault1 solana.PublicKey = poolAmmV4.TokenCoin, poolAmmV4.TokenPc
	var mint0, mint1 solana.PublicKey = poolAmmV4.CoinMint, poolAmmV4.PcMint

	if !mint0.Equals(WSOL) && !mint1.Equals(WSOL) {
		return nil, errors.New("该池子不存在WSOL代币,无法交易")
	}

	priceInfo, err := poolAmmV4.PriceInfo(context.TODO(), v.GetClient())
	if err != nil {
		return nil, err
	}

	if !mint0.Equals(WSOL) {
		priceInfo.BaseAmount, priceInfo.QuoteAmount = priceInfo.QuoteAmount, priceInfo.BaseAmount
		priceInfo.BaseDecimals, priceInfo.QuoteDecimals = priceInfo.QuoteDecimals, priceInfo.BaseDecimals
		priceInfo.BaseName, priceInfo.QuoteName = priceInfo.QuoteName, priceInfo.BaseName
		mint0, mint1 = mint1, mint0
		vault0, vault1 = vault1, vault0
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
	}
	crteteAccount, closeAccount, account, err := v.createAccountInstruction(amountIn)
	if err != nil {
		return nil, err
	}
	log.Debugf("创建 Base 过渡账户 | %s", account.String())
	tokenCreate, tokenAccount, err := v.checkTokenAccount(mint1)
	if err != nil {
		return nil, err
	}
	result = append(result, crteteAccount...)
	if tokenCreate != nil {
		result = append(result, tokenCreate)
	}
	// 将uint64转换为[]byte
	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, market.VaultSignerNonce)
	serumVaultSigner, err := solana.CreateProgramAddress([][]byte{
		poolAmmV4.Market.Bytes(),
		nonceBytes,
	}, OpenBookProgram)
	if err != nil {
		return nil, fmt.Errorf("failed to create WSOL account address: %v", err)
	}
	// 输入token的账户和输出token的账户
	var inputTokenAccount, outputTokenAccount solana.PublicKey
	if isBuy {
		inputTokenAccount, outputTokenAccount = account, tokenAccount
	} else {
		inputTokenAccount, outputTokenAccount = tokenAccount, account
	}

	result = append(result, raydium_amm.NewSwapBaseInInstruction(
		amountIn,
		miniAmountOut,
		solana.TokenProgramID,
		poolAddr,
		AmmAuthority,
		poolAmmV4.OpenOrders,
		poolAmmV4.TargetOrders,
		vault0,
		vault1,
		OpenBookProgram,
		poolAmmV4.Market,
		market.Bids,
		market.Asks,
		market.EventQueue,
		market.BaseVault,
		market.QuoteVault,
		serumVaultSigner,
		inputTokenAccount,
		outputTokenAccount,
		v.PublicKey(),
	).Build())
	result = append(result, closeAccount)

	return result, nil

}

// 构建创建临时WSOL账户的指令
func (v *V4Client) createAccountInstruction(amountIn uint64) ([]solana.Instruction, solana.Instruction, solana.PublicKey, error) {

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
func (v *V4Client) checkTokenAccount(mint solana.PublicKey) (solana.Instruction, solana.PublicKey, error) {
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
