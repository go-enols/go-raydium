package raydium

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"math"

	"github.com/go-enols/go-raydium/amm_v3"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	associatedtokenaccount "github.com/gagliardetto/solana-go/programs/associated-token-account"
	"github.com/gagliardetto/solana-go/rpc"
	"github.com/go-enols/go-log"
	"github.com/go-enols/gosolana"
)

func init() {
	// 设置AmmV3的程序地址
	amm_v3.SetProgramID(AmmV3ProgramId)
}

type V3Client struct {
	*gosolana.Wallet
}

func NewV3Client(wallet *gosolana.Wallet) *V3Client {
	return &V3Client{
		Wallet: wallet,
	}
}

func (v *V3Client) Swap(ctx context.Context, solIn float64, poolPub solana.PublicKey, isBuy bool, slippage float64, poolDatas ...[]byte) (*solana.Transaction, error) {
	var poolData []byte
	if len(poolDatas) > 0 {
		poolData = poolDatas[0]
	} else {
		out, err := v.GetClient().GetAccountInfo(ctx, poolPub)
		if err != nil {
			log.Fatal(err)
		}
		poolData = out.GetBinary()
	}
	poolAmmV3 := new(amm_v3.PoolState)
	if err := bin.NewBinDecoder(poolData).Decode(poolAmmV3); err != nil {
		log.Errorf("解析交易池子数据失败 | %s", err)
		return nil, err
	}
	if !poolAmmV3.TokenMint0.Equals(WSOL) && !poolAmmV3.TokenMint1.Equals(WSOL) {
		return nil, errors.New("不是SOL的交易,暂不支持")
	}
	// 确保SOL是在第一位需要兑换的代币是第二位
	if !poolAmmV3.TokenMint0.Equals(WSOL) {
		poolAmmV3.TokenMint0, poolAmmV3.TokenMint1 = poolAmmV3.TokenMint1, poolAmmV3.TokenMint0
		poolAmmV3.MintDecimals0, poolAmmV3.MintDecimals1 = poolAmmV3.MintDecimals1, poolAmmV3.MintDecimals0
		poolAmmV3.FundFeesToken0, poolAmmV3.FundFeesToken1 = poolAmmV3.FundFeesToken1, poolAmmV3.FundFeesToken0
		poolAmmV3.TotalFeesClaimedToken0, poolAmmV3.TotalFeesClaimedToken1 = poolAmmV3.TotalFeesClaimedToken1, poolAmmV3.TotalFeesClaimedToken0
		poolAmmV3.TokenVault0, poolAmmV3.TokenVault1 = poolAmmV3.TokenVault1, poolAmmV3.TokenVault0
		poolAmmV3.ProtocolFeesToken0, poolAmmV3.ProtocolFeesToken1 = poolAmmV3.ProtocolFeesToken1, poolAmmV3.ProtocolFeesToken0
		poolAmmV3.SwapInAmountToken0, poolAmmV3.SwapInAmountToken1 = poolAmmV3.SwapInAmountToken1, poolAmmV3.SwapInAmountToken0
		poolAmmV3.SwapOutAmountToken0, poolAmmV3.SwapOutAmountToken1 = poolAmmV3.SwapOutAmountToken1, poolAmmV3.SwapOutAmountToken0
	}
	var amountIn float64
	if isBuy {
		amountIn = solIn * math.Pow10(int(poolAmmV3.MintDecimals0))
	} else {
		amountIn = solIn * math.Pow10(int(poolAmmV3.MintDecimals1))
	}
	data, err := v.CreateInstruction(
		uint64(amountIn),
		poolPub,
		poolAmmV3,
		isBuy,
		slippage,
	)
	if err != nil {
		return nil, err
	}
	recentBlockHash, err := v.GetClient().GetLatestBlockhash(ctx, rpc.CommitmentFinalized)
	if err != nil {
		log.Errorf("获取Hash失败 | %s", err)
		return nil, err
	}
	// 构造交易
	tx, err := solana.NewTransaction(
		data,
		recentBlockHash.Value.Blockhash,
		solana.TransactionPayer(v.PublicKey()),
	)
	if err != nil {
		log.Errorf("构建交易失败 | %s", err)
		return nil, err
	}

	return tx, nil
}

func (v *V3Client) CreateInstruction(
	amountIn uint64,
	poolAddr solana.PublicKey,
	poolAmmV3 *amm_v3.PoolState,
	isBuy bool,
	slippage float64,
) ([]solana.Instruction, error) {
	var rseult []solana.Instruction
	solIn := amountIn
	if !isBuy {
		solIn = 0
	}
	crteteAccount, closeAccount, account, err := v.createAccountInstruction(solIn)
	if err != nil {
		return nil, err
	}
	log.Debugf("创建 Base 过渡账户 | %s", account.String())
	swap, err := v.createSwapInstruction(
		account, amountIn, poolAddr, poolAmmV3, isBuy, slippage,
	)
	if err != nil {
		return nil, err
	}
	rseult = append(rseult, crteteAccount...)
	rseult = append(rseult, swap...)
	rseult = append(rseult, closeAccount)
	return rseult, nil
}

// 构建一下创建临时账户和并初始化的指令
//
// 返回创建指令和关闭指令以及创建的账户
func (v *V3Client) createAccountInstruction(amountIn uint64) ([]solana.Instruction, solana.Instruction, solana.PublicKey, error) {

	// 创建临时 WSOL 账户
	seed := make([]byte, 24)
	if _, err := rand.Read(seed); err != nil {
		return nil, nil, solana.PublicKey{}, fmt.Errorf("failed to generate random seed: %v", err)
	}
	wsolAccountSeed := base64.URLEncoding.EncodeToString(seed)

	wsolAccount, err := solana.CreateWithSeed(v.PublicKey(), wsolAccountSeed, solana.TokenProgramID)
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
		wsolAccountSeed,
		wsolAccount,
		rentExemptBalance+uint64(amountIn),
	), MakeCloseAccountInstruction(wsolAccount, v.PublicKey(), v.PublicKey()), wsolAccount, nil
}

//	other_amount_threshold（滑点保护参数）
//	other_amount_threshold 是用于滑点保护的参数，根据交易方向的不同，它有两种不同的含义：
//	当 is_base_input = true 时（输入固定金额）：它表示最小接收数量，即愿意接受的最少输出代币数量
//	当 is_base_input = false 时（输出固定金额）：它表示最大支付数量，即愿意支付的最多输入代币数量
//	计算方法：
//	获取 other_amount_threshold 的过程主要通过以下步骤：
//	首先使用 get_out_put_amount_and_remaining_accounts 函数计算预期的输出金额： main.rs:843-854
//	然后根据交易类型应用滑点调整： main.rs:859-874
//	滑点计算使用 amount_with_slippage 函数，该函数根据滑点方向进行向上或向下取整： utils.rs:56-62
//	sqrt_price_limit_x64（价格限制参数）
//	sqrt_price_limit_x64 是价格限制参数，它以 Q64.64 定点数格式表示（一种特殊的数值表示法），用于设置交易可接受的价格限制。
//	计算方法：
//	用户可以提供一个限价，然后通过 price_to_sqrt_price_x64 函数将其转换为 Q64.64 格式： main.rs:833-841
//	price_to_sqrt_price_x64 函数实现如下，它考虑了两种代币的小数位数： utils.rs:269-272
//	如果不提供价格限制（参数为0或None），则根据交易方向默认设置为可能的最小/最大价格： swap.rs:611-618
//	注意：在实际使用时，sqrt_price_limit_x64 必须遵守特定的边界条件，确保它比当前价格高或低（取决于交易方向）： swap.rs:144-153
//	代码实际应用
//	在 swap 指令中，这两个参数组合使用来执行交易并保护用户： swap.rs:755-792
//	Notes
//	 other_amount_threshold 和 sqrt_price_limit_x64 都是重要的交易保护参数
//	 正确设置这些参数有助于防止在高波动市场中的意外滑点
//	 对于大多数用户来说，通常只需要设置一个合理的滑点百分比，客户端会自动计算 other_amount_threshold
//	 sqrt_price_limit_x64 可以省略（设为0），除非有特定的价格限制需求
//
// 构建一下交换指令
func (v *V3Client) createSwapInstruction(
	createTempAccount solana.PublicKey,
	amountIn uint64,
	poolAddr solana.PublicKey,
	poolAmmV3 *amm_v3.PoolState,
	isBuy bool,
	slippage float64, // 先预留一下
) ([]solana.Instruction, error) {
	var result []solana.Instruction
	price := poolAmmV3.Price()
	log.Debugf("当前价格：%.6f", price)

	log.Debugf("限价 | %.6f | 交换金额: %d ", price, amountIn)
	inter, tokenAccount, err := v.checkTokenAccount(poolAmmV3.TokenMint1)
	if err != nil {
		return nil, err
	}
	if inter != nil {
		result = append(result, inter)
	}
	tickArray, _, err := poolAmmV3.GetCurrentTickArrayAddress(AmmV3ProgramId, poolAddr)
	if err != nil {
		return nil, err
	}
	if !isBuy {
		createTempAccount, tokenAccount = tokenAccount, createTempAccount
		poolAmmV3.TokenVault0, poolAmmV3.TokenVault1 = poolAmmV3.TokenVault1, poolAmmV3.TokenVault0
	}
	log.Debug(createTempAccount.String())
	log.Debug(tokenAccount.String())
	result = append(result, amm_v3.NewSwapInstruction(
		amountIn,
		0,
		bin.Uint128{},
		true,
		v.PublicKey(),
		poolAmmV3.AmmConfig,
		poolAddr,
		createTempAccount,
		tokenAccount,
		poolAmmV3.TokenVault0,
		poolAmmV3.TokenVault1,
		poolAmmV3.ObservationKey,
		solana.TokenProgramID,
		tickArray,
	).Build())
	return result, nil
}

// 检查是否有对应mint地址的账户,如果没有则创建
func (v *V3Client) checkTokenAccount(mint solana.PublicKey) (solana.Instruction, solana.PublicKey, error) {
	out, err := v.GetClient().GetTokenAccountsByOwner(
		context.TODO(),
		v.PublicKey(), &rpc.GetTokenAccountsConfig{
			Mint: &mint,
		}, &rpc.GetTokenAccountsOpts{
			Commitment: rpc.CommitmentProcessed,
		})
	if err != nil {
		return nil, solana.PublicKey{}, err
	}
	if len(out.Value) > 0 {
		return nil, out.Value[0].Pubkey, nil
	} else {
		tokenAccount, err := GetAssociatedTokenAddress(v.PublicKey(), mint)
		if err != nil {
			return nil, solana.PublicKey{}, err
		}
		log.Debug(tokenAccount.String())
		return associatedtokenaccount.NewCreateInstruction(
			v.PublicKey(),
			v.PublicKey(),
			mint,
		).Build(), tokenAccount, nil
	}
}
