// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package raydium_cp_swap

import (
	"errors"
	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
	ag_format "github.com/gagliardetto/solana-go/text/format"
	ag_treeout "github.com/gagliardetto/treeout"
)

// Withdraw lp for token0 ande token1
//
// # Arguments
//
// * `ctx`- The context of accounts
// * `lp_token_amount` - Amount of pool tokens to burn. User receives an output of token a and b based on the percentage of the pool tokens that are returned.
// * `minimum_token_0_amount` -  Minimum amount of token 0 to receive, prevents excessive slippage
// * `minimum_token_1_amount` -  Minimum amount of token 1 to receive, prevents excessive slippage
//
type Withdraw struct {
	LpTokenAmount       *uint64
	MinimumToken0Amount *uint64
	MinimumToken1Amount *uint64

	// [0] = [SIGNER] owner
	// ··········· Pays to mint the position
	//
	// [1] = [] authority
	//
	// [2] = [WRITE] poolState
	// ··········· Pool state account
	//
	// [3] = [WRITE] ownerLpToken
	// ··········· Owner lp token account
	//
	// [4] = [WRITE] token0Account
	// ··········· The token account for receive token_0,
	//
	// [5] = [WRITE] token1Account
	// ··········· The token account for receive token_1
	//
	// [6] = [WRITE] token0Vault
	// ··········· The address that holds pool tokens for token_0
	//
	// [7] = [WRITE] token1Vault
	// ··········· The address that holds pool tokens for token_1
	//
	// [8] = [] tokenProgram
	// ··········· token Program
	//
	// [9] = [] tokenProgram2022
	// ··········· Token program 2022
	//
	// [10] = [] vault0Mint
	// ··········· The mint of token_0 vault
	//
	// [11] = [] vault1Mint
	// ··········· The mint of token_1 vault
	//
	// [12] = [WRITE] lpMint
	// ··········· Pool lp token mint
	//
	// [13] = [] memoProgram
	// ··········· memo program
	ag_solanago.AccountMetaSlice `bin:"-"`
}

// NewWithdrawInstructionBuilder creates a new `Withdraw` instruction builder.
func NewWithdrawInstructionBuilder() *Withdraw {
	nd := &Withdraw{
		AccountMetaSlice: make(ag_solanago.AccountMetaSlice, 14),
	}
	return nd
}

// SetLpTokenAmount sets the "lpTokenAmount" parameter.
func (inst *Withdraw) SetLpTokenAmount(lpTokenAmount uint64) *Withdraw {
	inst.LpTokenAmount = &lpTokenAmount
	return inst
}

// SetMinimumToken0Amount sets the "minimumToken0Amount" parameter.
func (inst *Withdraw) SetMinimumToken0Amount(minimumToken0Amount uint64) *Withdraw {
	inst.MinimumToken0Amount = &minimumToken0Amount
	return inst
}

// SetMinimumToken1Amount sets the "minimumToken1Amount" parameter.
func (inst *Withdraw) SetMinimumToken1Amount(minimumToken1Amount uint64) *Withdraw {
	inst.MinimumToken1Amount = &minimumToken1Amount
	return inst
}

// SetOwnerAccount sets the "owner" account.
// Pays to mint the position
func (inst *Withdraw) SetOwnerAccount(owner ag_solanago.PublicKey) *Withdraw {
	inst.AccountMetaSlice[0] = ag_solanago.Meta(owner).SIGNER()
	return inst
}

// GetOwnerAccount gets the "owner" account.
// Pays to mint the position
func (inst *Withdraw) GetOwnerAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(0)
}

// SetAuthorityAccount sets the "authority" account.
func (inst *Withdraw) SetAuthorityAccount(authority ag_solanago.PublicKey) *Withdraw {
	inst.AccountMetaSlice[1] = ag_solanago.Meta(authority)
	return inst
}

// GetAuthorityAccount gets the "authority" account.
func (inst *Withdraw) GetAuthorityAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(1)
}

// SetPoolStateAccount sets the "poolState" account.
// Pool state account
func (inst *Withdraw) SetPoolStateAccount(poolState ag_solanago.PublicKey) *Withdraw {
	inst.AccountMetaSlice[2] = ag_solanago.Meta(poolState).WRITE()
	return inst
}

// GetPoolStateAccount gets the "poolState" account.
// Pool state account
func (inst *Withdraw) GetPoolStateAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(2)
}

// SetOwnerLpTokenAccount sets the "ownerLpToken" account.
// Owner lp token account
func (inst *Withdraw) SetOwnerLpTokenAccount(ownerLpToken ag_solanago.PublicKey) *Withdraw {
	inst.AccountMetaSlice[3] = ag_solanago.Meta(ownerLpToken).WRITE()
	return inst
}

// GetOwnerLpTokenAccount gets the "ownerLpToken" account.
// Owner lp token account
func (inst *Withdraw) GetOwnerLpTokenAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(3)
}

// SetToken0AccountAccount sets the "token0Account" account.
// The token account for receive token_0,
func (inst *Withdraw) SetToken0AccountAccount(token0Account ag_solanago.PublicKey) *Withdraw {
	inst.AccountMetaSlice[4] = ag_solanago.Meta(token0Account).WRITE()
	return inst
}

// GetToken0AccountAccount gets the "token0Account" account.
// The token account for receive token_0,
func (inst *Withdraw) GetToken0AccountAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(4)
}

// SetToken1AccountAccount sets the "token1Account" account.
// The token account for receive token_1
func (inst *Withdraw) SetToken1AccountAccount(token1Account ag_solanago.PublicKey) *Withdraw {
	inst.AccountMetaSlice[5] = ag_solanago.Meta(token1Account).WRITE()
	return inst
}

// GetToken1AccountAccount gets the "token1Account" account.
// The token account for receive token_1
func (inst *Withdraw) GetToken1AccountAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(5)
}

// SetToken0VaultAccount sets the "token0Vault" account.
// The address that holds pool tokens for token_0
func (inst *Withdraw) SetToken0VaultAccount(token0Vault ag_solanago.PublicKey) *Withdraw {
	inst.AccountMetaSlice[6] = ag_solanago.Meta(token0Vault).WRITE()
	return inst
}

// GetToken0VaultAccount gets the "token0Vault" account.
// The address that holds pool tokens for token_0
func (inst *Withdraw) GetToken0VaultAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(6)
}

// SetToken1VaultAccount sets the "token1Vault" account.
// The address that holds pool tokens for token_1
func (inst *Withdraw) SetToken1VaultAccount(token1Vault ag_solanago.PublicKey) *Withdraw {
	inst.AccountMetaSlice[7] = ag_solanago.Meta(token1Vault).WRITE()
	return inst
}

// GetToken1VaultAccount gets the "token1Vault" account.
// The address that holds pool tokens for token_1
func (inst *Withdraw) GetToken1VaultAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(7)
}

// SetTokenProgramAccount sets the "tokenProgram" account.
// token Program
func (inst *Withdraw) SetTokenProgramAccount(tokenProgram ag_solanago.PublicKey) *Withdraw {
	inst.AccountMetaSlice[8] = ag_solanago.Meta(tokenProgram)
	return inst
}

// GetTokenProgramAccount gets the "tokenProgram" account.
// token Program
func (inst *Withdraw) GetTokenProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(8)
}

// SetTokenProgram2022Account sets the "tokenProgram2022" account.
// Token program 2022
func (inst *Withdraw) SetTokenProgram2022Account(tokenProgram2022 ag_solanago.PublicKey) *Withdraw {
	inst.AccountMetaSlice[9] = ag_solanago.Meta(tokenProgram2022)
	return inst
}

// GetTokenProgram2022Account gets the "tokenProgram2022" account.
// Token program 2022
func (inst *Withdraw) GetTokenProgram2022Account() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(9)
}

// SetVault0MintAccount sets the "vault0Mint" account.
// The mint of token_0 vault
func (inst *Withdraw) SetVault0MintAccount(vault0Mint ag_solanago.PublicKey) *Withdraw {
	inst.AccountMetaSlice[10] = ag_solanago.Meta(vault0Mint)
	return inst
}

// GetVault0MintAccount gets the "vault0Mint" account.
// The mint of token_0 vault
func (inst *Withdraw) GetVault0MintAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(10)
}

// SetVault1MintAccount sets the "vault1Mint" account.
// The mint of token_1 vault
func (inst *Withdraw) SetVault1MintAccount(vault1Mint ag_solanago.PublicKey) *Withdraw {
	inst.AccountMetaSlice[11] = ag_solanago.Meta(vault1Mint)
	return inst
}

// GetVault1MintAccount gets the "vault1Mint" account.
// The mint of token_1 vault
func (inst *Withdraw) GetVault1MintAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(11)
}

// SetLpMintAccount sets the "lpMint" account.
// Pool lp token mint
func (inst *Withdraw) SetLpMintAccount(lpMint ag_solanago.PublicKey) *Withdraw {
	inst.AccountMetaSlice[12] = ag_solanago.Meta(lpMint).WRITE()
	return inst
}

// GetLpMintAccount gets the "lpMint" account.
// Pool lp token mint
func (inst *Withdraw) GetLpMintAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(12)
}

// SetMemoProgramAccount sets the "memoProgram" account.
// memo program
func (inst *Withdraw) SetMemoProgramAccount(memoProgram ag_solanago.PublicKey) *Withdraw {
	inst.AccountMetaSlice[13] = ag_solanago.Meta(memoProgram)
	return inst
}

// GetMemoProgramAccount gets the "memoProgram" account.
// memo program
func (inst *Withdraw) GetMemoProgramAccount() *ag_solanago.AccountMeta {
	return inst.AccountMetaSlice.Get(13)
}

func (inst Withdraw) Build() *Instruction {
	return &Instruction{BaseVariant: ag_binary.BaseVariant{
		Impl:   inst,
		TypeID: Instruction_Withdraw,
	}}
}

// ValidateAndBuild validates the instruction parameters and accounts;
// if there is a validation error, it returns the error.
// Otherwise, it builds and returns the instruction.
func (inst Withdraw) ValidateAndBuild() (*Instruction, error) {
	if err := inst.Validate(); err != nil {
		return nil, err
	}
	return inst.Build(), nil
}

func (inst *Withdraw) Validate() error {
	// Check whether all (required) parameters are set:
	{
		if inst.LpTokenAmount == nil {
			return errors.New("LpTokenAmount parameter is not set")
		}
		if inst.MinimumToken0Amount == nil {
			return errors.New("MinimumToken0Amount parameter is not set")
		}
		if inst.MinimumToken1Amount == nil {
			return errors.New("MinimumToken1Amount parameter is not set")
		}
	}

	// Check whether all (required) accounts are set:
	{
		if inst.AccountMetaSlice[0] == nil {
			return errors.New("accounts.Owner is not set")
		}
		if inst.AccountMetaSlice[1] == nil {
			return errors.New("accounts.Authority is not set")
		}
		if inst.AccountMetaSlice[2] == nil {
			return errors.New("accounts.PoolState is not set")
		}
		if inst.AccountMetaSlice[3] == nil {
			return errors.New("accounts.OwnerLpToken is not set")
		}
		if inst.AccountMetaSlice[4] == nil {
			return errors.New("accounts.Token0Account is not set")
		}
		if inst.AccountMetaSlice[5] == nil {
			return errors.New("accounts.Token1Account is not set")
		}
		if inst.AccountMetaSlice[6] == nil {
			return errors.New("accounts.Token0Vault is not set")
		}
		if inst.AccountMetaSlice[7] == nil {
			return errors.New("accounts.Token1Vault is not set")
		}
		if inst.AccountMetaSlice[8] == nil {
			return errors.New("accounts.TokenProgram is not set")
		}
		if inst.AccountMetaSlice[9] == nil {
			return errors.New("accounts.TokenProgram2022 is not set")
		}
		if inst.AccountMetaSlice[10] == nil {
			return errors.New("accounts.Vault0Mint is not set")
		}
		if inst.AccountMetaSlice[11] == nil {
			return errors.New("accounts.Vault1Mint is not set")
		}
		if inst.AccountMetaSlice[12] == nil {
			return errors.New("accounts.LpMint is not set")
		}
		if inst.AccountMetaSlice[13] == nil {
			return errors.New("accounts.MemoProgram is not set")
		}
	}
	return nil
}

func (inst *Withdraw) EncodeToTree(parent ag_treeout.Branches) {
	parent.Child(ag_format.Program(ProgramName, ProgramID)).
		//
		ParentFunc(func(programBranch ag_treeout.Branches) {
			programBranch.Child(ag_format.Instruction("Withdraw")).
				//
				ParentFunc(func(instructionBranch ag_treeout.Branches) {

					// Parameters of the instruction:
					instructionBranch.Child("Params[len=3]").ParentFunc(func(paramsBranch ag_treeout.Branches) {
						paramsBranch.Child(ag_format.Param("      LpTokenAmount", *inst.LpTokenAmount))
						paramsBranch.Child(ag_format.Param("MinimumToken0Amount", *inst.MinimumToken0Amount))
						paramsBranch.Child(ag_format.Param("MinimumToken1Amount", *inst.MinimumToken1Amount))
					})

					// Accounts of the instruction:
					instructionBranch.Child("Accounts[len=14]").ParentFunc(func(accountsBranch ag_treeout.Branches) {
						accountsBranch.Child(ag_format.Meta("           owner", inst.AccountMetaSlice.Get(0)))
						accountsBranch.Child(ag_format.Meta("       authority", inst.AccountMetaSlice.Get(1)))
						accountsBranch.Child(ag_format.Meta("       poolState", inst.AccountMetaSlice.Get(2)))
						accountsBranch.Child(ag_format.Meta("    ownerLpToken", inst.AccountMetaSlice.Get(3)))
						accountsBranch.Child(ag_format.Meta("          token0", inst.AccountMetaSlice.Get(4)))
						accountsBranch.Child(ag_format.Meta("          token1", inst.AccountMetaSlice.Get(5)))
						accountsBranch.Child(ag_format.Meta("     token0Vault", inst.AccountMetaSlice.Get(6)))
						accountsBranch.Child(ag_format.Meta("     token1Vault", inst.AccountMetaSlice.Get(7)))
						accountsBranch.Child(ag_format.Meta("    tokenProgram", inst.AccountMetaSlice.Get(8)))
						accountsBranch.Child(ag_format.Meta("tokenProgram2022", inst.AccountMetaSlice.Get(9)))
						accountsBranch.Child(ag_format.Meta("      vault0Mint", inst.AccountMetaSlice.Get(10)))
						accountsBranch.Child(ag_format.Meta("      vault1Mint", inst.AccountMetaSlice.Get(11)))
						accountsBranch.Child(ag_format.Meta("          lpMint", inst.AccountMetaSlice.Get(12)))
						accountsBranch.Child(ag_format.Meta("     memoProgram", inst.AccountMetaSlice.Get(13)))
					})
				})
		})
}

func (obj Withdraw) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Serialize `LpTokenAmount` param:
	err = encoder.Encode(obj.LpTokenAmount)
	if err != nil {
		return err
	}
	// Serialize `MinimumToken0Amount` param:
	err = encoder.Encode(obj.MinimumToken0Amount)
	if err != nil {
		return err
	}
	// Serialize `MinimumToken1Amount` param:
	err = encoder.Encode(obj.MinimumToken1Amount)
	if err != nil {
		return err
	}
	return nil
}
func (obj *Withdraw) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Deserialize `LpTokenAmount`:
	err = decoder.Decode(&obj.LpTokenAmount)
	if err != nil {
		return err
	}
	// Deserialize `MinimumToken0Amount`:
	err = decoder.Decode(&obj.MinimumToken0Amount)
	if err != nil {
		return err
	}
	// Deserialize `MinimumToken1Amount`:
	err = decoder.Decode(&obj.MinimumToken1Amount)
	if err != nil {
		return err
	}
	return nil
}

// NewWithdrawInstruction declares a new Withdraw instruction with the provided parameters and accounts.
func NewWithdrawInstruction(
	// Parameters:
	lpTokenAmount uint64,
	minimumToken0Amount uint64,
	minimumToken1Amount uint64,
	// Accounts:
	owner ag_solanago.PublicKey,
	authority ag_solanago.PublicKey,
	poolState ag_solanago.PublicKey,
	ownerLpToken ag_solanago.PublicKey,
	token0Account ag_solanago.PublicKey,
	token1Account ag_solanago.PublicKey,
	token0Vault ag_solanago.PublicKey,
	token1Vault ag_solanago.PublicKey,
	tokenProgram ag_solanago.PublicKey,
	tokenProgram2022 ag_solanago.PublicKey,
	vault0Mint ag_solanago.PublicKey,
	vault1Mint ag_solanago.PublicKey,
	lpMint ag_solanago.PublicKey,
	memoProgram ag_solanago.PublicKey) *Withdraw {
	return NewWithdrawInstructionBuilder().
		SetLpTokenAmount(lpTokenAmount).
		SetMinimumToken0Amount(minimumToken0Amount).
		SetMinimumToken1Amount(minimumToken1Amount).
		SetOwnerAccount(owner).
		SetAuthorityAccount(authority).
		SetPoolStateAccount(poolState).
		SetOwnerLpTokenAccount(ownerLpToken).
		SetToken0AccountAccount(token0Account).
		SetToken1AccountAccount(token1Account).
		SetToken0VaultAccount(token0Vault).
		SetToken1VaultAccount(token1Vault).
		SetTokenProgramAccount(tokenProgram).
		SetTokenProgram2022Account(tokenProgram2022).
		SetVault0MintAccount(vault0Mint).
		SetVault1MintAccount(vault1Mint).
		SetLpMintAccount(lpMint).
		SetMemoProgramAccount(memoProgram)
}
