// Code generated by https://github.com/gagliardetto/anchor-go. DO NOT EDIT.

package amm_v3

import (
	"fmt"
	"math"
	"math/big"
	"encoding/binary"  

	ag_binary "github.com/gagliardetto/binary"
	ag_solanago "github.com/gagliardetto/solana-go"
)

type AmmConfig struct {
	// Bump to identify PDA
	Bump  uint8
	Index uint16

	// Address of the protocol owner
	Owner ag_solanago.PublicKey

	// The protocol fee
	ProtocolFeeRate uint32

	// The trade fee, denominated in hundredths of a bip (10^-6)
	TradeFeeRate uint32

	// The tick spacing
	TickSpacing uint16

	// The fund fee, denominated in hundredths of a bip (10^-6)
	FundFeeRate uint32
	PaddingU32  uint32
	FundOwner   ag_solanago.PublicKey
	Padding     [3]uint64
}

var AmmConfigDiscriminator = [8]byte{218, 244, 33, 104, 203, 203, 43, 111}

func (obj AmmConfig) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Write account discriminator:
	err = encoder.WriteBytes(AmmConfigDiscriminator[:], false)
	if err != nil {
		return err
	}
	// Serialize `Bump` param:
	err = encoder.Encode(obj.Bump)
	if err != nil {
		return err
	}
	// Serialize `Index` param:
	err = encoder.Encode(obj.Index)
	if err != nil {
		return err
	}
	// Serialize `Owner` param:
	err = encoder.Encode(obj.Owner)
	if err != nil {
		return err
	}
	// Serialize `ProtocolFeeRate` param:
	err = encoder.Encode(obj.ProtocolFeeRate)
	if err != nil {
		return err
	}
	// Serialize `TradeFeeRate` param:
	err = encoder.Encode(obj.TradeFeeRate)
	if err != nil {
		return err
	}
	// Serialize `TickSpacing` param:
	err = encoder.Encode(obj.TickSpacing)
	if err != nil {
		return err
	}
	// Serialize `FundFeeRate` param:
	err = encoder.Encode(obj.FundFeeRate)
	if err != nil {
		return err
	}
	// Serialize `PaddingU32` param:
	err = encoder.Encode(obj.PaddingU32)
	if err != nil {
		return err
	}
	// Serialize `FundOwner` param:
	err = encoder.Encode(obj.FundOwner)
	if err != nil {
		return err
	}
	// Serialize `Padding` param:
	err = encoder.Encode(obj.Padding)
	if err != nil {
		return err
	}
	return nil
}

func (obj *AmmConfig) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Read and check account discriminator:
	{
		discriminator, err := decoder.ReadTypeID()
		if err != nil {
			return err
		}
		if !discriminator.Equal(AmmConfigDiscriminator[:]) {
			return fmt.Errorf(
				"wrong discriminator: wanted %s, got %s",
				"[218 244 33 104 203 203 43 111]",
				fmt.Sprint(discriminator[:]))
		}
	}
	// Deserialize `Bump`:
	err = decoder.Decode(&obj.Bump)
	if err != nil {
		return err
	}
	// Deserialize `Index`:
	err = decoder.Decode(&obj.Index)
	if err != nil {
		return err
	}
	// Deserialize `Owner`:
	err = decoder.Decode(&obj.Owner)
	if err != nil {
		return err
	}
	// Deserialize `ProtocolFeeRate`:
	err = decoder.Decode(&obj.ProtocolFeeRate)
	if err != nil {
		return err
	}
	// Deserialize `TradeFeeRate`:
	err = decoder.Decode(&obj.TradeFeeRate)
	if err != nil {
		return err
	}
	// Deserialize `TickSpacing`:
	err = decoder.Decode(&obj.TickSpacing)
	if err != nil {
		return err
	}
	// Deserialize `FundFeeRate`:
	err = decoder.Decode(&obj.FundFeeRate)
	if err != nil {
		return err
	}
	// Deserialize `PaddingU32`:
	err = decoder.Decode(&obj.PaddingU32)
	if err != nil {
		return err
	}
	// Deserialize `FundOwner`:
	err = decoder.Decode(&obj.FundOwner)
	if err != nil {
		return err
	}
	// Deserialize `Padding`:
	err = decoder.Decode(&obj.Padding)
	if err != nil {
		return err
	}
	return nil
}

type OperationState struct {
	// Bump to identify PDA
	Bump uint8

	// Address of the operation owner
	OperationOwners [10]ag_solanago.PublicKey

	// The mint address of whitelist to emmit reward
	WhitelistMints [100]ag_solanago.PublicKey
}

var OperationStateDiscriminator = [8]byte{19, 236, 58, 237, 81, 222, 183, 252}

func (obj OperationState) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Write account discriminator:
	err = encoder.WriteBytes(OperationStateDiscriminator[:], false)
	if err != nil {
		return err
	}
	// Serialize `Bump` param:
	err = encoder.Encode(obj.Bump)
	if err != nil {
		return err
	}
	// Serialize `OperationOwners` param:
	err = encoder.Encode(obj.OperationOwners)
	if err != nil {
		return err
	}
	// Serialize `WhitelistMints` param:
	err = encoder.Encode(obj.WhitelistMints)
	if err != nil {
		return err
	}
	return nil
}

func (obj *OperationState) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Read and check account discriminator:
	{
		discriminator, err := decoder.ReadTypeID()
		if err != nil {
			return err
		}
		if !discriminator.Equal(OperationStateDiscriminator[:]) {
			return fmt.Errorf(
				"wrong discriminator: wanted %s, got %s",
				"[19 236 58 237 81 222 183 252]",
				fmt.Sprint(discriminator[:]))
		}
	}
	// Deserialize `Bump`:
	err = decoder.Decode(&obj.Bump)
	if err != nil {
		return err
	}
	// Deserialize `OperationOwners`:
	err = decoder.Decode(&obj.OperationOwners)
	if err != nil {
		return err
	}
	// Deserialize `WhitelistMints`:
	err = decoder.Decode(&obj.WhitelistMints)
	if err != nil {
		return err
	}
	return nil
}

type ObservationState struct {
	// Whether the ObservationState is initialized
	Initialized bool

	// recent update epoch
	RecentEpoch uint64

	// the most-recently updated index of the observations array
	ObservationIndex uint16

	// belongs to which pool
	PoolId ag_solanago.PublicKey

	// observation array
	Observations [100]Observation

	// padding for feature update
	Padding [4]uint64
}

var ObservationStateDiscriminator = [8]byte{122, 174, 197, 53, 129, 9, 165, 132}

func (obj ObservationState) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Write account discriminator:
	err = encoder.WriteBytes(ObservationStateDiscriminator[:], false)
	if err != nil {
		return err
	}
	// Serialize `Initialized` param:
	err = encoder.Encode(obj.Initialized)
	if err != nil {
		return err
	}
	// Serialize `RecentEpoch` param:
	err = encoder.Encode(obj.RecentEpoch)
	if err != nil {
		return err
	}
	// Serialize `ObservationIndex` param:
	err = encoder.Encode(obj.ObservationIndex)
	if err != nil {
		return err
	}
	// Serialize `PoolId` param:
	err = encoder.Encode(obj.PoolId)
	if err != nil {
		return err
	}
	// Serialize `Observations` param:
	err = encoder.Encode(obj.Observations)
	if err != nil {
		return err
	}
	// Serialize `Padding` param:
	err = encoder.Encode(obj.Padding)
	if err != nil {
		return err
	}
	return nil
}

func (obj *ObservationState) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Read and check account discriminator:
	{
		discriminator, err := decoder.ReadTypeID()
		if err != nil {
			return err
		}
		if !discriminator.Equal(ObservationStateDiscriminator[:]) {
			return fmt.Errorf(
				"wrong discriminator: wanted %s, got %s",
				"[122 174 197 53 129 9 165 132]",
				fmt.Sprint(discriminator[:]))
		}
	}
	// Deserialize `Initialized`:
	err = decoder.Decode(&obj.Initialized)
	if err != nil {
		return err
	}
	// Deserialize `RecentEpoch`:
	err = decoder.Decode(&obj.RecentEpoch)
	if err != nil {
		return err
	}
	// Deserialize `ObservationIndex`:
	err = decoder.Decode(&obj.ObservationIndex)
	if err != nil {
		return err
	}
	// Deserialize `PoolId`:
	err = decoder.Decode(&obj.PoolId)
	if err != nil {
		return err
	}
	// Deserialize `Observations`:
	err = decoder.Decode(&obj.Observations)
	if err != nil {
		return err
	}
	// Deserialize `Padding`:
	err = decoder.Decode(&obj.Padding)
	if err != nil {
		return err
	}
	return nil
}

type PersonalPositionState struct {
	// Bump to identify PDA
	Bump [1]uint8

	// Mint address of the tokenized position
	NftMint ag_solanago.PublicKey

	// The ID of the pool with which this token is connected
	PoolId ag_solanago.PublicKey

	// The lower bound tick of the position
	TickLowerIndex int32

	// The upper bound tick of the position
	TickUpperIndex int32

	// The amount of liquidity owned by this position
	Liquidity ag_binary.Uint128

	// The token_0 fee growth of the aggregate position as of the last action on the individual position
	FeeGrowthInside0LastX64 ag_binary.Uint128

	// The token_1 fee growth of the aggregate position as of the last action on the individual position
	FeeGrowthInside1LastX64 ag_binary.Uint128

	// The fees owed to the position owner in token_0, as of the last computation
	TokenFeesOwed0 uint64

	// The fees owed to the position owner in token_1, as of the last computation
	TokenFeesOwed1 uint64
	RewardInfos    [3]PositionRewardInfo
	RecentEpoch    uint64
	Padding        [7]uint64
}

var PersonalPositionStateDiscriminator = [8]byte{70, 111, 150, 126, 230, 15, 25, 117}

func (obj PersonalPositionState) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Write account discriminator:
	err = encoder.WriteBytes(PersonalPositionStateDiscriminator[:], false)
	if err != nil {
		return err
	}
	// Serialize `Bump` param:
	err = encoder.Encode(obj.Bump)
	if err != nil {
		return err
	}
	// Serialize `NftMint` param:
	err = encoder.Encode(obj.NftMint)
	if err != nil {
		return err
	}
	// Serialize `PoolId` param:
	err = encoder.Encode(obj.PoolId)
	if err != nil {
		return err
	}
	// Serialize `TickLowerIndex` param:
	err = encoder.Encode(obj.TickLowerIndex)
	if err != nil {
		return err
	}
	// Serialize `TickUpperIndex` param:
	err = encoder.Encode(obj.TickUpperIndex)
	if err != nil {
		return err
	}
	// Serialize `Liquidity` param:
	err = encoder.Encode(obj.Liquidity)
	if err != nil {
		return err
	}
	// Serialize `FeeGrowthInside0LastX64` param:
	err = encoder.Encode(obj.FeeGrowthInside0LastX64)
	if err != nil {
		return err
	}
	// Serialize `FeeGrowthInside1LastX64` param:
	err = encoder.Encode(obj.FeeGrowthInside1LastX64)
	if err != nil {
		return err
	}
	// Serialize `TokenFeesOwed0` param:
	err = encoder.Encode(obj.TokenFeesOwed0)
	if err != nil {
		return err
	}
	// Serialize `TokenFeesOwed1` param:
	err = encoder.Encode(obj.TokenFeesOwed1)
	if err != nil {
		return err
	}
	// Serialize `RewardInfos` param:
	err = encoder.Encode(obj.RewardInfos)
	if err != nil {
		return err
	}
	// Serialize `RecentEpoch` param:
	err = encoder.Encode(obj.RecentEpoch)
	if err != nil {
		return err
	}
	// Serialize `Padding` param:
	err = encoder.Encode(obj.Padding)
	if err != nil {
		return err
	}
	return nil
}

func (obj *PersonalPositionState) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Read and check account discriminator:
	{
		discriminator, err := decoder.ReadTypeID()
		if err != nil {
			return err
		}
		if !discriminator.Equal(PersonalPositionStateDiscriminator[:]) {
			return fmt.Errorf(
				"wrong discriminator: wanted %s, got %s",
				"[70 111 150 126 230 15 25 117]",
				fmt.Sprint(discriminator[:]))
		}
	}
	// Deserialize `Bump`:
	err = decoder.Decode(&obj.Bump)
	if err != nil {
		return err
	}
	// Deserialize `NftMint`:
	err = decoder.Decode(&obj.NftMint)
	if err != nil {
		return err
	}
	// Deserialize `PoolId`:
	err = decoder.Decode(&obj.PoolId)
	if err != nil {
		return err
	}
	// Deserialize `TickLowerIndex`:
	err = decoder.Decode(&obj.TickLowerIndex)
	if err != nil {
		return err
	}
	// Deserialize `TickUpperIndex`:
	err = decoder.Decode(&obj.TickUpperIndex)
	if err != nil {
		return err
	}
	// Deserialize `Liquidity`:
	err = decoder.Decode(&obj.Liquidity)
	if err != nil {
		return err
	}
	// Deserialize `FeeGrowthInside0LastX64`:
	err = decoder.Decode(&obj.FeeGrowthInside0LastX64)
	if err != nil {
		return err
	}
	// Deserialize `FeeGrowthInside1LastX64`:
	err = decoder.Decode(&obj.FeeGrowthInside1LastX64)
	if err != nil {
		return err
	}
	// Deserialize `TokenFeesOwed0`:
	err = decoder.Decode(&obj.TokenFeesOwed0)
	if err != nil {
		return err
	}
	// Deserialize `TokenFeesOwed1`:
	err = decoder.Decode(&obj.TokenFeesOwed1)
	if err != nil {
		return err
	}
	// Deserialize `RewardInfos`:
	err = decoder.Decode(&obj.RewardInfos)
	if err != nil {
		return err
	}
	// Deserialize `RecentEpoch`:
	err = decoder.Decode(&obj.RecentEpoch)
	if err != nil {
		return err
	}
	// Deserialize `Padding`:
	err = decoder.Decode(&obj.Padding)
	if err != nil {
		return err
	}
	return nil
}

type PoolState struct {
	// Bump to identify PDA
	Bump      [1]uint8
	AmmConfig ag_solanago.PublicKey
	Owner     ag_solanago.PublicKey

	// Token pair of the pool, where token_mint_0 address < token_mint_1 address
	TokenMint0 ag_solanago.PublicKey
	TokenMint1 ag_solanago.PublicKey

	// Token pair vault
	TokenVault0 ag_solanago.PublicKey
	TokenVault1 ag_solanago.PublicKey

	// observation account key
	ObservationKey ag_solanago.PublicKey

	// mint0 and mint1 decimals
	MintDecimals0 uint8
	MintDecimals1 uint8

	// The minimum number of ticks between initialized ticks
	TickSpacing uint16

	// The currently in range liquidity available to the pool.
	Liquidity ag_binary.Uint128

	// The current price of the pool as a sqrt(token_1/token_0) Q64.64 value
	SqrtPriceX64 ag_binary.Uint128

	// The current tick of the pool, i.e. according to the last tick transition that was run.
	TickCurrent int32
	Padding3    uint16
	Padding4    uint16

	// The fee growth as a Q64.64 number, i.e. fees of token_0 and token_1 collected per
	// unit of liquidity for the entire life of the pool.
	FeeGrowthGlobal0X64 ag_binary.Uint128
	FeeGrowthGlobal1X64 ag_binary.Uint128

	// The amounts of token_0 and token_1 that are owed to the protocol.
	ProtocolFeesToken0 uint64
	ProtocolFeesToken1 uint64

	// The amounts in and out of swap token_0 and token_1
	SwapInAmountToken0  ag_binary.Uint128
	SwapOutAmountToken1 ag_binary.Uint128
	SwapInAmountToken1  ag_binary.Uint128
	SwapOutAmountToken0 ag_binary.Uint128

	// Bitwise representation of the state of the pool
	// bit0, 1: disable open position and increase liquidity, 0: normal
	// bit1, 1: disable decrease liquidity, 0: normal
	// bit2, 1: disable collect fee, 0: normal
	// bit3, 1: disable collect reward, 0: normal
	// bit4, 1: disable swap, 0: normal
	Status uint8

	// Leave blank for future use
	Padding     [7]uint8
	RewardInfos [3]RewardInfo

	// Packed initialized tick array state
	TickArrayBitmap [16]uint64

	// except protocol_fee and fund_fee
	TotalFeesToken0 uint64

	// except protocol_fee and fund_fee
	TotalFeesClaimedToken0 uint64
	TotalFeesToken1        uint64
	TotalFeesClaimedToken1 uint64
	FundFeesToken0         uint64
	FundFeesToken1         uint64
	OpenTime               uint64
	RecentEpoch            uint64
	Padding1               [24]uint64
	Padding2               [32]uint64
}

var PoolStateDiscriminator = [8]byte{247, 237, 227, 245, 215, 195, 222, 70}

func (obj PoolState) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Write account discriminator:
	err = encoder.WriteBytes(PoolStateDiscriminator[:], false)
	if err != nil {
		return err
	}
	// Serialize `Bump` param:
	err = encoder.Encode(obj.Bump)
	if err != nil {
		return err
	}
	// Serialize `AmmConfig` param:
	err = encoder.Encode(obj.AmmConfig)
	if err != nil {
		return err
	}
	// Serialize `Owner` param:
	err = encoder.Encode(obj.Owner)
	if err != nil {
		return err
	}
	// Serialize `TokenMint0` param:
	err = encoder.Encode(obj.TokenMint0)
	if err != nil {
		return err
	}
	// Serialize `TokenMint1` param:
	err = encoder.Encode(obj.TokenMint1)
	if err != nil {
		return err
	}
	// Serialize `TokenVault0` param:
	err = encoder.Encode(obj.TokenVault0)
	if err != nil {
		return err
	}
	// Serialize `TokenVault1` param:
	err = encoder.Encode(obj.TokenVault1)
	if err != nil {
		return err
	}
	// Serialize `ObservationKey` param:
	err = encoder.Encode(obj.ObservationKey)
	if err != nil {
		return err
	}
	// Serialize `MintDecimals0` param:
	err = encoder.Encode(obj.MintDecimals0)
	if err != nil {
		return err
	}
	// Serialize `MintDecimals1` param:
	err = encoder.Encode(obj.MintDecimals1)
	if err != nil {
		return err
	}
	// Serialize `TickSpacing` param:
	err = encoder.Encode(obj.TickSpacing)
	if err != nil {
		return err
	}
	// Serialize `Liquidity` param:
	err = encoder.Encode(obj.Liquidity)
	if err != nil {
		return err
	}
	// Serialize `SqrtPriceX64` param:
	err = encoder.Encode(obj.SqrtPriceX64)
	if err != nil {
		return err
	}
	// Serialize `TickCurrent` param:
	err = encoder.Encode(obj.TickCurrent)
	if err != nil {
		return err
	}
	// Serialize `Padding3` param:
	err = encoder.Encode(obj.Padding3)
	if err != nil {
		return err
	}
	// Serialize `Padding4` param:
	err = encoder.Encode(obj.Padding4)
	if err != nil {
		return err
	}
	// Serialize `FeeGrowthGlobal0X64` param:
	err = encoder.Encode(obj.FeeGrowthGlobal0X64)
	if err != nil {
		return err
	}
	// Serialize `FeeGrowthGlobal1X64` param:
	err = encoder.Encode(obj.FeeGrowthGlobal1X64)
	if err != nil {
		return err
	}
	// Serialize `ProtocolFeesToken0` param:
	err = encoder.Encode(obj.ProtocolFeesToken0)
	if err != nil {
		return err
	}
	// Serialize `ProtocolFeesToken1` param:
	err = encoder.Encode(obj.ProtocolFeesToken1)
	if err != nil {
		return err
	}
	// Serialize `SwapInAmountToken0` param:
	err = encoder.Encode(obj.SwapInAmountToken0)
	if err != nil {
		return err
	}
	// Serialize `SwapOutAmountToken1` param:
	err = encoder.Encode(obj.SwapOutAmountToken1)
	if err != nil {
		return err
	}
	// Serialize `SwapInAmountToken1` param:
	err = encoder.Encode(obj.SwapInAmountToken1)
	if err != nil {
		return err
	}
	// Serialize `SwapOutAmountToken0` param:
	err = encoder.Encode(obj.SwapOutAmountToken0)
	if err != nil {
		return err
	}
	// Serialize `Status` param:
	err = encoder.Encode(obj.Status)
	if err != nil {
		return err
	}
	// Serialize `Padding` param:
	err = encoder.Encode(obj.Padding)
	if err != nil {
		return err
	}
	// Serialize `RewardInfos` param:
	err = encoder.Encode(obj.RewardInfos)
	if err != nil {
		return err
	}
	// Serialize `TickArrayBitmap` param:
	err = encoder.Encode(obj.TickArrayBitmap)
	if err != nil {
		return err
	}
	// Serialize `TotalFeesToken0` param:
	err = encoder.Encode(obj.TotalFeesToken0)
	if err != nil {
		return err
	}
	// Serialize `TotalFeesClaimedToken0` param:
	err = encoder.Encode(obj.TotalFeesClaimedToken0)
	if err != nil {
		return err
	}
	// Serialize `TotalFeesToken1` param:
	err = encoder.Encode(obj.TotalFeesToken1)
	if err != nil {
		return err
	}
	// Serialize `TotalFeesClaimedToken1` param:
	err = encoder.Encode(obj.TotalFeesClaimedToken1)
	if err != nil {
		return err
	}
	// Serialize `FundFeesToken0` param:
	err = encoder.Encode(obj.FundFeesToken0)
	if err != nil {
		return err
	}
	// Serialize `FundFeesToken1` param:
	err = encoder.Encode(obj.FundFeesToken1)
	if err != nil {
		return err
	}
	// Serialize `OpenTime` param:
	err = encoder.Encode(obj.OpenTime)
	if err != nil {
		return err
	}
	// Serialize `RecentEpoch` param:
	err = encoder.Encode(obj.RecentEpoch)
	if err != nil {
		return err
	}
	// Serialize `Padding1` param:
	err = encoder.Encode(obj.Padding1)
	if err != nil {
		return err
	}
	// Serialize `Padding2` param:
	err = encoder.Encode(obj.Padding2)
	if err != nil {
		return err
	}
	return nil
}

func (obj *PoolState) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Read and check account discriminator:
	{
		discriminator, err := decoder.ReadTypeID()
		if err != nil {
			return err
		}
		if !discriminator.Equal(PoolStateDiscriminator[:]) {
			return fmt.Errorf(
				"wrong discriminator: wanted %s, got %s",
				"[247 237 227 245 215 195 222 70]",
				fmt.Sprint(discriminator[:]))
		}
	}
	// Deserialize `Bump`:
	err = decoder.Decode(&obj.Bump)
	if err != nil {
		return err
	}
	// Deserialize `AmmConfig`:
	err = decoder.Decode(&obj.AmmConfig)
	if err != nil {
		return err
	}
	// Deserialize `Owner`:
	err = decoder.Decode(&obj.Owner)
	if err != nil {
		return err
	}
	// Deserialize `TokenMint0`:
	err = decoder.Decode(&obj.TokenMint0)
	if err != nil {
		return err
	}
	// Deserialize `TokenMint1`:
	err = decoder.Decode(&obj.TokenMint1)
	if err != nil {
		return err
	}
	// Deserialize `TokenVault0`:
	err = decoder.Decode(&obj.TokenVault0)
	if err != nil {
		return err
	}
	// Deserialize `TokenVault1`:
	err = decoder.Decode(&obj.TokenVault1)
	if err != nil {
		return err
	}
	// Deserialize `ObservationKey`:
	err = decoder.Decode(&obj.ObservationKey)
	if err != nil {
		return err
	}
	// Deserialize `MintDecimals0`:
	err = decoder.Decode(&obj.MintDecimals0)
	if err != nil {
		return err
	}
	// Deserialize `MintDecimals1`:
	err = decoder.Decode(&obj.MintDecimals1)
	if err != nil {
		return err
	}
	// Deserialize `TickSpacing`:
	err = decoder.Decode(&obj.TickSpacing)
	if err != nil {
		return err
	}
	// Deserialize `Liquidity`:
	err = decoder.Decode(&obj.Liquidity)
	if err != nil {
		return err
	}
	// Deserialize `SqrtPriceX64`:
	err = decoder.Decode(&obj.SqrtPriceX64)
	if err != nil {
		return err
	}
	// Deserialize `TickCurrent`:
	err = decoder.Decode(&obj.TickCurrent)
	if err != nil {
		return err
	}
	// Deserialize `Padding3`:
	err = decoder.Decode(&obj.Padding3)
	if err != nil {
		return err
	}
	// Deserialize `Padding4`:
	err = decoder.Decode(&obj.Padding4)
	if err != nil {
		return err
	}
	// Deserialize `FeeGrowthGlobal0X64`:
	err = decoder.Decode(&obj.FeeGrowthGlobal0X64)
	if err != nil {
		return err
	}
	// Deserialize `FeeGrowthGlobal1X64`:
	err = decoder.Decode(&obj.FeeGrowthGlobal1X64)
	if err != nil {
		return err
	}
	// Deserialize `ProtocolFeesToken0`:
	err = decoder.Decode(&obj.ProtocolFeesToken0)
	if err != nil {
		return err
	}
	// Deserialize `ProtocolFeesToken1`:
	err = decoder.Decode(&obj.ProtocolFeesToken1)
	if err != nil {
		return err
	}
	// Deserialize `SwapInAmountToken0`:
	err = decoder.Decode(&obj.SwapInAmountToken0)
	if err != nil {
		return err
	}
	// Deserialize `SwapOutAmountToken1`:
	err = decoder.Decode(&obj.SwapOutAmountToken1)
	if err != nil {
		return err
	}
	// Deserialize `SwapInAmountToken1`:
	err = decoder.Decode(&obj.SwapInAmountToken1)
	if err != nil {
		return err
	}
	// Deserialize `SwapOutAmountToken0`:
	err = decoder.Decode(&obj.SwapOutAmountToken0)
	if err != nil {
		return err
	}
	// Deserialize `Status`:
	err = decoder.Decode(&obj.Status)
	if err != nil {
		return err
	}
	// Deserialize `Padding`:
	err = decoder.Decode(&obj.Padding)
	if err != nil {
		return err
	}
	// Deserialize `RewardInfos`:
	err = decoder.Decode(&obj.RewardInfos)
	if err != nil {
		return err
	}
	// Deserialize `TickArrayBitmap`:
	err = decoder.Decode(&obj.TickArrayBitmap)
	if err != nil {
		return err
	}
	// Deserialize `TotalFeesToken0`:
	err = decoder.Decode(&obj.TotalFeesToken0)
	if err != nil {
		return err
	}
	// Deserialize `TotalFeesClaimedToken0`:
	err = decoder.Decode(&obj.TotalFeesClaimedToken0)
	if err != nil {
		return err
	}
	// Deserialize `TotalFeesToken1`:
	err = decoder.Decode(&obj.TotalFeesToken1)
	if err != nil {
		return err
	}
	// Deserialize `TotalFeesClaimedToken1`:
	err = decoder.Decode(&obj.TotalFeesClaimedToken1)
	if err != nil {
		return err
	}
	// Deserialize `FundFeesToken0`:
	err = decoder.Decode(&obj.FundFeesToken0)
	if err != nil {
		return err
	}
	// Deserialize `FundFeesToken1`:
	err = decoder.Decode(&obj.FundFeesToken1)
	if err != nil {
		return err
	}
	// Deserialize `OpenTime`:
	err = decoder.Decode(&obj.OpenTime)
	if err != nil {
		return err
	}
	// Deserialize `RecentEpoch`:
	err = decoder.Decode(&obj.RecentEpoch)
	if err != nil {
		return err
	}
	// Deserialize `Padding1`:
	err = decoder.Decode(&obj.Padding1)
	if err != nil {
		return err
	}
	// Deserialize `Padding2`:
	err = decoder.Decode(&obj.Padding2)
	if err != nil {
		return err
	}
	return nil
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
func (pc *PrecisionCalculator) leBytesToUint128(b ag_binary.Uint128) *big.Int {
	return new(big.Int).Or(
		new(big.Int).Lsh(new(big.Int).SetUint64(b.Hi), 64),
		new(big.Int).SetUint64(b.Lo),
	)
}

// 分层精度计算流程
func (pc *PrecisionCalculator) CalculatePrice(camm *PoolState) *big.Rat {
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
// SwapDirection 表示交换方向  
type SwapDirection bool  
  
const (  
	Coin2PC SwapDirection = true // 从 Coin 到 PC  
	PC2Coin SwapDirection = false                     // 从 PC 到 Coin  
)  
  
func (obj *PoolState) Price() float64{
	pc := NewPrecisionCalculator()
	priceRat := pc.CalculatePrice(obj)
	priceFloat, ok := priceRat.Float64()
	if ok {
		return 0
	}
	return priceFloat
}
// 计算带滑点的金额  
func (obj *PoolState)AmountWithSlippage(amount uint64, slippage float64, isMaxIn bool) uint64 {  
    amountBig := new(big.Float).SetUint64(amount)  
    slippageBig := new(big.Float).SetFloat64(slippage)  
      
    var result *big.Float  
    if isMaxIn {  
        // 计算最大输入量(增加滑点)  
        onePlusSlippage := new(big.Float).Add(big.NewFloat(1.0), slippageBig)  
        result = new(big.Float).Mul(amountBig, onePlusSlippage)  
    } else {  
        // 计算最小输出量(减少滑点)  
        oneMinusSlippage := new(big.Float).Sub(big.NewFloat(1.0), slippageBig)  
        result = new(big.Float).Mul(amountBig, oneMinusSlippage)  
    }  
      
    // 转换回uint64  
    resultUint64, _ := result.Uint64()  
    return resultUint64  
}
// 计算10的n次方  
func multiplier(decimals uint8) *big.Float {  
    return new(big.Float).SetFloat64(math.Pow(10, float64(decimals)))  
}  
  
// 将价格转换为sqrtPriceX64格式  
func (obj *PoolState)PriceToSqrtPriceX64() *big.Int {  
    priceBig := new(big.Float).SetFloat64(obj.Price()*(1-0.05))  
      
    // 考虑代币小数位数调整价格  
    multiplier1 := multiplier(obj.MintDecimals0)  
    multiplier0 := multiplier(obj.MintDecimals1)  
      
    // price * 10^decimals1 / 10^decimals0  
    priceWithDecimals := new(big.Float).Mul(priceBig, multiplier1)  
    priceWithDecimals = new(big.Float).Quo(priceWithDecimals, multiplier0)  
      
    // 计算平方根  
    sqrtPrice := new(big.Float).Sqrt(priceWithDecimals)  
      
    // 转换为Q64.64格式  
    // 2^64 = 18446744073709551616  
    q64Multiplier := new(big.Float).SetFloat64(18446744073709551616.0)  
    sqrtPriceX64Float := new(big.Float).Mul(sqrtPrice, q64Multiplier)  
      
    // 转换为big.Int  
    sqrtPriceX64Int, _ := new(big.Int).SetString(sqrtPriceX64Float.Text('f', 0), 10)  
    return sqrtPriceX64Int  
}



type ProtocolPositionState struct {
	// Bump to identify PDA
	Bump uint8

	// The ID of the pool with which this token is connected
	PoolId ag_solanago.PublicKey

	// The lower bound tick of the position
	TickLowerIndex int32

	// The upper bound tick of the position
	TickUpperIndex int32

	// The amount of liquidity owned by this position
	Liquidity ag_binary.Uint128

	// The token_0 fee growth per unit of liquidity as of the last update to liquidity or fees owed
	FeeGrowthInside0LastX64 ag_binary.Uint128

	// The token_1 fee growth per unit of liquidity as of the last update to liquidity or fees owed
	FeeGrowthInside1LastX64 ag_binary.Uint128

	// The fees owed to the position owner in token_0
	TokenFeesOwed0 uint64

	// The fees owed to the position owner in token_1
	TokenFeesOwed1 uint64

	// The reward growth per unit of liquidity as of the last update to liquidity
	RewardGrowthInside [3]ag_binary.Uint128
	RecentEpoch        uint64
	Padding            [7]uint64
}

var ProtocolPositionStateDiscriminator = [8]byte{100, 226, 145, 99, 146, 218, 160, 106}

func (obj ProtocolPositionState) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Write account discriminator:
	err = encoder.WriteBytes(ProtocolPositionStateDiscriminator[:], false)
	if err != nil {
		return err
	}
	// Serialize `Bump` param:
	err = encoder.Encode(obj.Bump)
	if err != nil {
		return err
	}
	// Serialize `PoolId` param:
	err = encoder.Encode(obj.PoolId)
	if err != nil {
		return err
	}
	// Serialize `TickLowerIndex` param:
	err = encoder.Encode(obj.TickLowerIndex)
	if err != nil {
		return err
	}
	// Serialize `TickUpperIndex` param:
	err = encoder.Encode(obj.TickUpperIndex)
	if err != nil {
		return err
	}
	// Serialize `Liquidity` param:
	err = encoder.Encode(obj.Liquidity)
	if err != nil {
		return err
	}
	// Serialize `FeeGrowthInside0LastX64` param:
	err = encoder.Encode(obj.FeeGrowthInside0LastX64)
	if err != nil {
		return err
	}
	// Serialize `FeeGrowthInside1LastX64` param:
	err = encoder.Encode(obj.FeeGrowthInside1LastX64)
	if err != nil {
		return err
	}
	// Serialize `TokenFeesOwed0` param:
	err = encoder.Encode(obj.TokenFeesOwed0)
	if err != nil {
		return err
	}
	// Serialize `TokenFeesOwed1` param:
	err = encoder.Encode(obj.TokenFeesOwed1)
	if err != nil {
		return err
	}
	// Serialize `RewardGrowthInside` param:
	err = encoder.Encode(obj.RewardGrowthInside)
	if err != nil {
		return err
	}
	// Serialize `RecentEpoch` param:
	err = encoder.Encode(obj.RecentEpoch)
	if err != nil {
		return err
	}
	// Serialize `Padding` param:
	err = encoder.Encode(obj.Padding)
	if err != nil {
		return err
	}
	return nil
}

func (obj *ProtocolPositionState) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Read and check account discriminator:
	{
		discriminator, err := decoder.ReadTypeID()
		if err != nil {
			return err
		}
		if !discriminator.Equal(ProtocolPositionStateDiscriminator[:]) {
			return fmt.Errorf(
				"wrong discriminator: wanted %s, got %s",
				"[100 226 145 99 146 218 160 106]",
				fmt.Sprint(discriminator[:]))
		}
	}
	// Deserialize `Bump`:
	err = decoder.Decode(&obj.Bump)
	if err != nil {
		return err
	}
	// Deserialize `PoolId`:
	err = decoder.Decode(&obj.PoolId)
	if err != nil {
		return err
	}
	// Deserialize `TickLowerIndex`:
	err = decoder.Decode(&obj.TickLowerIndex)
	if err != nil {
		return err
	}
	// Deserialize `TickUpperIndex`:
	err = decoder.Decode(&obj.TickUpperIndex)
	if err != nil {
		return err
	}
	// Deserialize `Liquidity`:
	err = decoder.Decode(&obj.Liquidity)
	if err != nil {
		return err
	}
	// Deserialize `FeeGrowthInside0LastX64`:
	err = decoder.Decode(&obj.FeeGrowthInside0LastX64)
	if err != nil {
		return err
	}
	// Deserialize `FeeGrowthInside1LastX64`:
	err = decoder.Decode(&obj.FeeGrowthInside1LastX64)
	if err != nil {
		return err
	}
	// Deserialize `TokenFeesOwed0`:
	err = decoder.Decode(&obj.TokenFeesOwed0)
	if err != nil {
		return err
	}
	// Deserialize `TokenFeesOwed1`:
	err = decoder.Decode(&obj.TokenFeesOwed1)
	if err != nil {
		return err
	}
	// Deserialize `RewardGrowthInside`:
	err = decoder.Decode(&obj.RewardGrowthInside)
	if err != nil {
		return err
	}
	// Deserialize `RecentEpoch`:
	err = decoder.Decode(&obj.RecentEpoch)
	if err != nil {
		return err
	}
	// Deserialize `Padding`:
	err = decoder.Decode(&obj.Padding)
	if err != nil {
		return err
	}
	return nil
}

type TickArrayState struct {
	PoolId               ag_solanago.PublicKey
	StartTickIndex       int32
	Ticks                [60]TickState
	InitializedTickCount uint8
	RecentEpoch          uint64
	Padding              [107]uint8
}

var TickArrayStateDiscriminator = [8]byte{192, 155, 85, 205, 49, 249, 129, 42}

func (obj TickArrayState) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Write account discriminator:
	err = encoder.WriteBytes(TickArrayStateDiscriminator[:], false)
	if err != nil {
		return err
	}
	// Serialize `PoolId` param:
	err = encoder.Encode(obj.PoolId)
	if err != nil {
		return err
	}
	// Serialize `StartTickIndex` param:
	err = encoder.Encode(obj.StartTickIndex)
	if err != nil {
		return err
	}
	// Serialize `Ticks` param:
	err = encoder.Encode(obj.Ticks)
	if err != nil {
		return err
	}
	// Serialize `InitializedTickCount` param:
	err = encoder.Encode(obj.InitializedTickCount)
	if err != nil {
		return err
	}
	// Serialize `RecentEpoch` param:
	err = encoder.Encode(obj.RecentEpoch)
	if err != nil {
		return err
	}
	// Serialize `Padding` param:
	err = encoder.Encode(obj.Padding)
	if err != nil {
		return err
	}
	return nil
}

func (obj *TickArrayState) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Read and check account discriminator:
	{
		discriminator, err := decoder.ReadTypeID()
		if err != nil {
			return err
		}
		if !discriminator.Equal(TickArrayStateDiscriminator[:]) {
			return fmt.Errorf(
				"wrong discriminator: wanted %s, got %s",
				"[192 155 85 205 49 249 129 42]",
				fmt.Sprint(discriminator[:]))
		}
	}
	// Deserialize `PoolId`:
	err = decoder.Decode(&obj.PoolId)
	if err != nil {
		return err
	}
	// Deserialize `StartTickIndex`:
	err = decoder.Decode(&obj.StartTickIndex)
	if err != nil {
		return err
	}
	// Deserialize `Ticks`:
	err = decoder.Decode(&obj.Ticks)
	if err != nil {
		return err
	}
	// Deserialize `InitializedTickCount`:
	err = decoder.Decode(&obj.InitializedTickCount)
	if err != nil {
		return err
	}
	// Deserialize `RecentEpoch`:
	err = decoder.Decode(&obj.RecentEpoch)
	if err != nil {
		return err
	}
	// Deserialize `Padding`:
	err = decoder.Decode(&obj.Padding)
	if err != nil {
		return err
	}
	return nil
}

type TickArrayBitmapExtension struct {
	PoolId ag_solanago.PublicKey

	// Packed initialized tick array state for start_tick_index is positive
	PositiveTickArrayBitmap [14][8]uint64

	// Packed initialized tick array state for start_tick_index is negitive
	NegativeTickArrayBitmap [14][8]uint64
}

var TickArrayBitmapExtensionDiscriminator = [8]byte{60, 150, 36, 219, 97, 128, 139, 153}

func (obj TickArrayBitmapExtension) MarshalWithEncoder(encoder *ag_binary.Encoder) (err error) {
	// Write account discriminator:
	err = encoder.WriteBytes(TickArrayBitmapExtensionDiscriminator[:], false)
	if err != nil {
		return err
	}
	// Serialize `PoolId` param:
	err = encoder.Encode(obj.PoolId)
	if err != nil {
		return err
	}
	// Serialize `PositiveTickArrayBitmap` param:
	err = encoder.Encode(obj.PositiveTickArrayBitmap)
	if err != nil {
		return err
	}
	// Serialize `NegativeTickArrayBitmap` param:
	err = encoder.Encode(obj.NegativeTickArrayBitmap)
	if err != nil {
		return err
	}
	return nil
}

func (obj *TickArrayBitmapExtension) UnmarshalWithDecoder(decoder *ag_binary.Decoder) (err error) {
	// Read and check account discriminator:
	{
		discriminator, err := decoder.ReadTypeID()
		if err != nil {
			return err
		}
		if !discriminator.Equal(TickArrayBitmapExtensionDiscriminator[:]) {
			return fmt.Errorf(
				"wrong discriminator: wanted %s, got %s",
				"[60 150 36 219 97 128 139 153]",
				fmt.Sprint(discriminator[:]))
		}
	}
	// Deserialize `PoolId`:
	err = decoder.Decode(&obj.PoolId)
	if err != nil {
		return err
	}
	// Deserialize `PositiveTickArrayBitmap`:
	err = decoder.Decode(&obj.PositiveTickArrayBitmap)
	if err != nil {
		return err
	}
	// Deserialize `NegativeTickArrayBitmap`:
	err = decoder.Decode(&obj.NegativeTickArrayBitmap)
	if err != nil {
		return err
	}
	return nil
}

// Constants matching Raydium CLMM protocol  
const (  
	TICK_ARRAY_SIZE = 60  
	TICK_ARRAY_SEED = "tick_array"  
)  

  
// GetArrayStartIndex calculates the start index of a tick array based on a tick index  
// This matches the Rust implementation in tick_array.rs  
func GetArrayStartIndex(tickIndex int32, tickSpacing uint16) int32 {  
	ticksInArray := int32(tickSpacing) * TICK_ARRAY_SIZE  
	start := tickIndex / ticksInArray  
	if tickIndex < 0 && tickIndex%ticksInArray != 0 {  
		start = start - 1  
	}  
	return start * ticksInArray  
}  
  
// GetTickArrayAddress calculates the PDA address for a tick array  
func GetTickArrayAddress(programID ag_solanago.PublicKey, poolID ag_solanago.PublicKey, tickIndex int32, tickSpacing uint16) (ag_solanago.PublicKey, uint8, error) {  
	// Calculate tick_array_start_index  
	tickArrayStartIndex := GetArrayStartIndex(tickIndex, tickSpacing)  
	  
	// Convert to byte array - using big-endian byte order (BE)  
	tickArrayStartIndexBytes := make([]byte, 4)  
	binary.BigEndian.PutUint32(tickArrayStartIndexBytes, uint32(tickArrayStartIndex))  
	  
	// Generate PDA address  
	seeds := [][]byte{  
		[]byte(TICK_ARRAY_SEED),  
		poolID[:],  
		tickArrayStartIndexBytes,  
	}  
	  
	return ag_solanago.FindProgramAddress(seeds, programID)  
}  
  

// GetCurrentTickArrayAddress gets the tick array address for the current tick in the pool  
func (pool *PoolState) GetCurrentTickArrayAddress(programID ag_solanago.PublicKey, poolID ag_solanago.PublicKey) (ag_solanago.PublicKey, uint8, error) {  
	return GetTickArrayAddress(programID, poolID, pool.TickCurrent, pool.TickSpacing)  
} 