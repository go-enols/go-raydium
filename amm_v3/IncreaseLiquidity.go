// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package amm_v3

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// #[deprecated(note = "Use `increase_liquidity_v2` instead.")]
// Increases liquidity with a exist position, with amount paid by `payer`
//
// # Arguments
//
// * `ctx` - The context of accounts
// * `liquidity` - The desired liquidity to be added, can't be zero
// * `amount_0_max` - The max amount of token_0 to spend, which serves as a slippage check
// * `amount_1_max` - The max amount of token_1 to spend, which serves as a slippage check
//
type IncreaseLiquidity struct {
	Liquidity  *ag_binary.Uint128
	Amount0Max *uint64
	Amount1Max *uint64

	// [0] = [SIGNER] nftOwner
	// ··········· Pays to mint the position
	//
	// [1] = [] nftAccount
	// ··········· The token account for nft
	//
	// [2] = [WRITE] poolState
	//
	// [3] = [WRITE] protocolPosition
	//
	// [4] = [WRITE] personalPosition
	// ··········· Increase liquidity for this position
	//
	// [5] = [WRITE] tickArrayLower
	// ··········· Stores init state for the lower tick
	//
	// [6] = [WRITE] tickArrayUpper
	// ··········· Stores init state for the upper tick
	//
	// [7] = [WRITE] tokenAccount0
	// ··········· The payer's token account for token_0
	//
	// [8] = [WRITE] tokenAccount1
	// ··········· The token account spending token_1 to mint the position
	//
	// [9] = [WRITE] tokenVault0
	// ··········· The address that holds pool tokens for token_0
	//
	// [10] = [WRITE] tokenVault1
	// ··········· The address that holds pool tokens for token_1
	//
	// [11] = [] tokenProgram
	// ··········· Program to create mint account and mint tokens
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewIncreaseLiquidityInstructionBuilder creates a new `IncreaseLiquidity` instruction builder.
func NewIncreaseLiquidityInstructionBuilder() *IncreaseLiquidity {
	nd := &IncreaseLiquidity{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 12),
	}
	return nd
}

// SetLiquidity sets the "liquidity" parameter.
func (inst *IncreaseLiquidity) SetLiquidity(liquidity ag_binary.Uint128) *IncreaseLiquidity {
	inst.Liquidity = &liquidity
	return inst
}

// SetAmount0Max sets the "amount0Max" parameter.
func (inst *IncreaseLiquidity) SetAmount0Max(amount0Max uint64) *IncreaseLiquidity {
	inst.Amount0Max = &amount0Max
	return inst
}

// SetAmount1Max sets the "amount1Max" parameter.
func (inst *IncreaseLiquidity) SetAmount1Max(amount1Max uint64) *IncreaseLiquidity {
	inst.Amount1Max = &amount1Max
	return inst
}

// SetNftOwnerAccount sets the "nftOwner" account.
// Pays to mint the position
func (inst *IncreaseLiquidity) SetNftOwnerAccount(nftOwner ag_solanago.PublicKey) *IncreaseLiquidity {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(nftOwner).SIGNER()
	return inst
}

// GetNftOwnerAccount gets the "nftOwner" account.
// Pays to mint the position
func (inst *IncreaseLiquidity) GetNftOwnerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetNftAccountAccount sets the "nftAccount" account.
// The token account for nft
func (inst *IncreaseLiquidity) SetNftAccountAccount(nftAccount ag_solanago.PublicKey) *IncreaseLiquidity {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(nftAccount)
	return inst
}

// GetNftAccountAccount gets the "nftAccount" account.
// The token account for nft
func (inst *IncreaseLiquidity) GetNftAccountAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetPoolStateAccount sets the "poolState" account.
func (inst *IncreaseLiquidity) SetPoolStateAccount(poolState ag_solanago.PublicKey) *IncreaseLiquidity {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(poolState).WRITE()
	return inst
}

// GetPoolStateAccount gets the "poolState" account.
func (inst *IncreaseLiquidity) GetPoolStateAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetProtocolPositionAccount sets the "protocolPosition" account.
func (inst *IncreaseLiquidity) SetProtocolPositionAccount(protocolPosition ag_solanago.PublicKey) *IncreaseLiquidity {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(protocolPosition).WRITE()
	return inst
}

// GetProtocolPositionAccount gets the "protocolPosition" account.
func (inst *IncreaseLiquidity) GetProtocolPositionAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

// SetPersonalPositionAccount sets the "personalPosition" account.
// Increase liquidity for this position
func (inst *IncreaseLiquidity) SetPersonalPositionAccount(personalPosition ag_solanago.PublicKey) *IncreaseLiquidity {
	inst.AccountMetaSlice[4] = ag_solanago.Meta(personalPosition).WRITE()
	return inst
}

// GetPersonalPositionAccount gets the "personalPosition" account.
// Increase liquidity for this position
func (inst *IncreaseLiquidity) GetPersonalPositionAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(4)
}

// SetTickArrayLowerAccount sets the "tickArrayLower" account.
// Stores init state for the lower tick
func (inst *IncreaseLiquidity) SetTickArrayLowerAccount(tickArrayLower ag_solanago.PublicKey) *IncreaseLiquidity {
	inst.AccountMetaSlice[5] = ag_solanago.Meta(tickArrayLower).WRITE()
	return inst
}

// GetTickArrayLowerAccount gets the "tickArrayLower" account.
// Stores init state for the lower tick
func (inst *IncreaseLiquidity) GetTickArrayLowerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(5)
}

// SetTickArrayUpperAccount sets the "tickArrayUpper" account.
// Stores init state for the upper tick
func (inst *IncreaseLiquidity) SetTickArrayUpperAccount(tickArrayUpper ag_solanago.PublicKey) *IncreaseLiquidity {
	inst.AccountMetaSlice[6] = ag_solanago.Meta(tickArrayUpper).WRITE()
	return inst
}

// GetTickArrayUpperAccount gets the "tickArrayUpper" account.
// Stores init state for the upper tick
func (inst *IncreaseLiquidity) GetTickArrayUpperAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(6)
}

// SetTokenAccount0Account sets the "tokenAccount0" account.
// The payer's token account for token_0
func (inst *IncreaseLiquidity) SetTokenAccount0Account(tokenAccount0 ag_solanago.PublicKey) *IncreaseLiquidity {
	inst.AccountMetaSlice[7] = ag_solanago.Meta(tokenAccount0).WRITE()
	return inst
}

// GetTokenAccount0Account gets the "tokenAccount0" account.
// The payer's token account for token_0
func (inst *IncreaseLiquidity) GetTokenAccount0Account() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(7)
}

// SetTokenAccount1Account sets the "tokenAccount1" account.
// The token account spending token_1 to mint the position
func (inst *IncreaseLiquidity) SetTokenAccount1Account(tokenAccount1 ag_solanago.PublicKey) *IncreaseLiquidity {
	inst.AccountMetaSlice[8] = ag_solanago.Meta(tokenAccount1).WRITE()
	return inst
}

// GetTokenAccount1Account gets the "tokenAccount1" account.
// The token account spending token_1 to mint the position
func (inst *IncreaseLiquidity) GetTokenAccount1Account() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(8)
}

// SetTokenVault0Account sets the "tokenVault0" account.
// The address that holds pool tokens for token_0
func (inst *IncreaseLiquidity) SetTokenVault0Account(tokenVault0 ag_solanago.PublicKey) *IncreaseLiquidity {
	inst.AccountMetaSlice[9] = ag_solanago.Meta(tokenVault0).WRITE()
	return inst
}

// GetTokenVault0Account gets the "tokenVault0" account.
// The address that holds pool tokens for token_0
func (inst *IncreaseLiquidity) GetTokenVault0Account() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(9)
}

// SetTokenVault1Account sets the "tokenVault1" account.
// The address that holds pool tokens for token_1
func (inst *IncreaseLiquidity) SetTokenVault1Account(tokenVault1 ag_solanago.PublicKey) *IncreaseLiquidity {
	inst.AccountMetaSlice[10] = ag_solanago.Meta(tokenVault1).WRITE()
	return inst
}

// GetTokenVault1Account gets the "tokenVault1" account.
// The address that holds pool tokens for token_1
func (inst *IncreaseLiquidity) GetTokenVault1Account() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(10)
}

// SetTokenProgramAccount sets the "tokenProgram" account.
// Program to create mint account and mint tokens
func (inst *IncreaseLiquidity) SetTokenProgramAccount(tokenProgram ag_solanago.PublicKey) *IncreaseLiquidity {
	inst.AccountMetaSlice[11] = ag_solanago.Meta(tokenProgram)
	return inst
}

// GetTokenProgramAccount gets the "tokenProgram" account.
// Program to create mint account and mint tokens
func (inst *IncreaseLiquidity) GetTokenProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(11)
}

func (inst IncreaseLiquidity) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_IncreaseLiquidity,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst IncreaseLiquidity) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *IncreaseLiquidity) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.Liquidity == nil {
			return errors.New("Liquidity parameter is not set")
		}
		if inst.Amount0Max == nil {
			return errors.New("Amount0Max parameter is not set")
		}
		if inst.Amount1Max == nil {
			return errors.New("Amount1Max parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.NftOwner is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.NftAccount is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.PoolState is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.ProtocolPosition is not set")
		}
		if inst.AccountMetaSlice[4] == nil {
			return errors.New("accounts.PersonalPosition is not set")
		}
		if inst.AccountMetaSlice[5] == nil {
			return errors.New("accounts.TickArrayLower is not set")
		}
		if inst.AccountMetaSlice[6] == nil {
			return errors.New("accounts.TickArrayUpper is not set")
		}
		if inst.AccountMetaSlice[7] == nil {
			return errors.New("accounts.TokenAccount0 is not set")
		}
		if inst.AccountMetaSlice[8] == nil {
			return errors.New("accounts.TokenAccount1 is not set")
		}
		if inst.AccountMetaSlice[9] == nil {
			return errors.New("accounts.TokenVault0 is not set")
		}
		if inst.AccountMetaSlice[10] == nil {
			return errors.New("accounts.TokenVault1 is not set")
		}
		if inst.AccountMetaSlice[11] == nil {
			return errors.New("accounts.TokenProgram is not set")
		}
	}
	return nil
}

func (inst *IncreaseLiquidity) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("IncreaseLiquidity")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=3]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param(" Liquidity", *inst.Liquidity))
						paramsBranch.Child(ag_format.Param("Amount0Max", *inst.Amount0Max))
						paramsBranch.Child(ag_format.Param("Amount1Max", *inst.Amount1Max))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=12]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("        nftOwner", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("             nft", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("       poolState", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("protocolPosition", inst.AccountMetaSlice.Get(3)))
						accountsBranch.Child(ag_format.Meta("personalPosition", inst.AccountMetaSlice.Get(4)))
						accountsBranch.Child(ag_format.Meta("  tickArrayLower", inst.AccountMetaSlice.Get(5)))
						accountsBranch.Child(ag_format.Meta("  tickArrayUpper", inst.AccountMetaSlice.Get(6)))
						accountsBranch.Child(ag_format.Meta("   tokenAccount0", inst.AccountMetaSlice.Get(7)))
						accountsBranch.Child(ag_format.Meta("   tokenAccount1", inst.AccountMetaSlice.Get(8)))
						accountsBranch.Child(ag_format.Meta("     tokenVault0", inst.AccountMetaSlice.Get(9)))
						accountsBranch.Child(ag_format.Meta("     tokenVault1", inst.AccountMetaSlice.Get(10)))
						accountsBranch.Child(ag_format.Meta("    tokenProgram", inst.AccountMetaSlice.Get(11)))
					})
				})
		})
}

func (obj IncreaseLiquidity) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `Liquidity` param:
	err = encoder.Encode(obj.Liquidity)
	if err != nil {
		return err
	}
	// Serialize `Amount0Max` param:
	err = encoder.Encode(obj.Amount0Max)
	if err != nil {
		return err
	}
	// Serialize `Amount1Max` param:
	err = encoder.Encode(obj.Amount1Max)
	if err != nil {
		return err
	}
	return nil
}
func (obj *IncreaseLiquidity) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `Liquidity`:
	err = decoder.Decode(&obj.Liquidity)
	if err != nil {
		return err
	}
	// Deserialize `Amount0Max`:
	err = decoder.Decode(&obj.Amount0Max)
	if err != nil {
		return err
	}
	// Deserialize `Amount1Max`:
	err = decoder.Decode(&obj.Amount1Max)
	if err != nil {
		return err
	}
	return nil
}

// NewIncreaseLiquidityInstruction declares a new IncreaseLiquidity instruction with the provided parameters and accounts.
func NewIncreaseLiquidityInstruction(
	// Parameters:
	liquidity ag_binary.Uint128,
	amount0Max uint64,
	amount1Max uint64,
	// Accounts:
	nftOwner ag_solanago.PublicKey,
	nftAccount ag_solanago.PublicKey,
	poolState ag_solanago.PublicKey,
	protocolPosition ag_solanago.PublicKey,
	personalPosition ag_solanago.PublicKey,
	tickArrayLower ag_solanago.PublicKey,
	tickArrayUpper ag_solanago.PublicKey,
	tokenAccount0 ag_solanago.PublicKey,
	tokenAccount1 ag_solanago.PublicKey,
	tokenVault0 ag_solanago.PublicKey,
	tokenVault1 ag_solanago.PublicKey,
	tokenProgram ag_solanago.PublicKey) *IncreaseLiquidity {
	return NewIncreaseLiquidityInstructionBuilder().
		SetLiquidity(liquidity).
		SetAmount0Max(amount0Max).
		SetAmount1Max(amount1Max).
		SetNftOwnerAccount(nftOwner).
		SetNftAccountAccount(nftAccount).
		SetPoolStateAccount(poolState).
		SetProtocolPositionAccount(protocolPosition).
		SetPersonalPositionAccount(personalPosition).
		SetTickArrayLowerAccount(tickArrayLower).
		SetTickArrayUpperAccount(tickArrayUpper).
		SetTokenAccount0Account(tokenAccount0).
		SetTokenAccount1Account(tokenAccount1).
		SetTokenVault0Account(tokenVault0).
		SetTokenVault1Account(tokenVault1).
		SetTokenProgramAccount(tokenProgram)
}
