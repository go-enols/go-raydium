package core

import (
	"log"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
)

// 奖励信息结构体
type RewardInfo struct {
	RewardState           uint8
	OpenTime              uint64
	EndTime               uint64
	LastUpdateTime        uint64
	EmissionsPerSecondX64 [16]byte // u128
	RewardTotalEmissioned uint64
	RewardClaimed         uint64
	TokenMint             solana.PublicKey
	TokenVault            solana.PublicKey
	Authority             solana.PublicKey
	RewardGrowthGlobalX64 [16]byte // u128
}

func (r *RewardInfo) UnmarshalWithDecoder(decoder *bin.Decoder) error {
	var err error
	if r.RewardState, err = decoder.ReadUint8(); err != nil {
		return err
	}
	if r.OpenTime, err = decoder.ReadUint64(bin.LE); err != nil {
		return err
	}
	if r.EndTime, err = decoder.ReadUint64(bin.LE); err != nil {
		return err
	}
	if r.LastUpdateTime, err = decoder.ReadUint64(bin.LE); err != nil {
		return err
	}
	if _, err = decoder.Read(r.EmissionsPerSecondX64[:]); err != nil {
		return err
	}
	if r.RewardTotalEmissioned, err = decoder.ReadUint64(bin.LE); err != nil {
		return err
	}
	if r.RewardClaimed, err = decoder.ReadUint64(bin.LE); err != nil {
		return err
	}
	if b, err := decoder.ReadBytes(32); err != nil {
		return err
	} else {
		copy(r.TokenMint[:], b)
	}
	if b, err := decoder.ReadBytes(32); err != nil {
		return err
	} else {
		copy(r.TokenVault[:], b)
	}
	if b, err := decoder.ReadBytes(32); err != nil {
		return err
	} else {
		copy(r.Authority[:], b)
	}
	if _, err = decoder.Read(r.RewardGrowthGlobalX64[:]); err != nil {
		return err
	}
	return nil
}

// cAMM池子结构体
type CAMM_STATE_LAYOUT struct {
	Bump                   [1]uint8
	AmmConfig              solana.PublicKey
	Owner                  solana.PublicKey
	TokenMint0             solana.PublicKey
	TokenMint1             solana.PublicKey
	TokenVault0            solana.PublicKey
	TokenVault1            solana.PublicKey
	ObservationKey         solana.PublicKey
	MintDecimals0          uint8
	MintDecimals1          uint8
	TickSpacing            uint16
	Liquidity              [16]byte // u128
	SqrtPriceX64           [16]byte // u128
	TickCurrent            int32
	Padding3               uint16
	Padding4               uint16
	FeeGrowthGlobal0X64    [16]byte // u128
	FeeGrowthGlobal1X64    [16]byte // u128
	ProtocolFeesToken0     uint64
	ProtocolFeesToken1     uint64
	SwapInAmountToken0     [16]byte // u128
	SwapOutAmountToken1    [16]byte // u128
	SwapInAmountToken1     [16]byte // u128
	SwapOutAmountToken0    [16]byte // u128
	Status                 uint8
	Padding                [7]uint8
	RewardInfos            [3]RewardInfo
	TickArrayBitmap        [16]uint64
	TotalFeesToken0        uint64
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

func (c *CAMM_STATE_LAYOUT) UnmarshalWithDecoder(decoder *bin.Decoder) error {
	var err error
	// 跳过8字节头部填充
	if _, err = decoder.ReadBytes(8); err != nil {
		return err
	}
	if _, err = decoder.Read(c.Bump[:]); err != nil {
		return err
	}
	readPk := func(pk *solana.PublicKey) error {
		b, err := decoder.ReadBytes(32)
		if err != nil {
			return err
		}
		copy(pk[:], b)
		return nil
	}
	if err = readPk(&c.AmmConfig); err != nil {
		return err
	}
	if err = readPk(&c.Owner); err != nil {
		return err
	}
	if err = readPk(&c.TokenMint0); err != nil {
		return err
	}
	if err = readPk(&c.TokenMint1); err != nil {
		return err
	}
	if err = readPk(&c.TokenVault0); err != nil {
		return err
	}
	if err = readPk(&c.TokenVault1); err != nil {
		return err
	}
	if err = readPk(&c.ObservationKey); err != nil {
		return err
	}
	log.Println(c.ObservationKey.String())
	if c.MintDecimals0, err = decoder.ReadUint8(); err != nil {
		return err
	}
	if c.MintDecimals1, err = decoder.ReadUint8(); err != nil {
		return err
	}
	if c.TickSpacing, err = decoder.ReadUint16(bin.LE); err != nil {
		return err
	}
	if _, err = decoder.Read(c.Liquidity[:]); err != nil {
		return err
	}
	if _, err = decoder.Read(c.SqrtPriceX64[:]); err != nil {
		return err
	}
	if c.TickCurrent, err = decoder.ReadInt32(bin.LE); err != nil {
		return err
	}
	if c.Padding3, err = decoder.ReadUint16(bin.LE); err != nil {
		return err
	}
	if c.Padding4, err = decoder.ReadUint16(bin.LE); err != nil {
		return err
	}
	if _, err = decoder.Read(c.FeeGrowthGlobal0X64[:]); err != nil {
		return err
	}
	if _, err = decoder.Read(c.FeeGrowthGlobal1X64[:]); err != nil {
		return err
	}
	if c.ProtocolFeesToken0, err = decoder.ReadUint64(bin.LE); err != nil {
		return err
	}
	if c.ProtocolFeesToken1, err = decoder.ReadUint64(bin.LE); err != nil {
		return err
	}
	if _, err = decoder.Read(c.SwapInAmountToken0[:]); err != nil {
		return err
	}
	if _, err = decoder.Read(c.SwapOutAmountToken1[:]); err != nil {
		return err
	}
	if _, err = decoder.Read(c.SwapInAmountToken1[:]); err != nil {
		return err
	}
	if _, err = decoder.Read(c.SwapOutAmountToken0[:]); err != nil {
		return err
	}
	if c.Status, err = decoder.ReadUint8(); err != nil {
		return err
	}
	if _, err = decoder.Read(c.Padding[:]); err != nil {
		return err
	}
	for i := 0; i < 3; i++ {
		if err = c.RewardInfos[i].UnmarshalWithDecoder(decoder); err != nil {
			return err
		}
	}
	for i := 0; i < 16; i++ {
		if c.TickArrayBitmap[i], err = decoder.ReadUint64(bin.LE); err != nil {
			return err
		}
	}
	if c.TotalFeesToken0, err = decoder.ReadUint64(bin.LE); err != nil {
		return err
	}
	if c.TotalFeesClaimedToken0, err = decoder.ReadUint64(bin.LE); err != nil {
		return err
	}
	if c.TotalFeesToken1, err = decoder.ReadUint64(bin.LE); err != nil {
		return err
	}
	if c.TotalFeesClaimedToken1, err = decoder.ReadUint64(bin.LE); err != nil {
		return err
	}
	if c.FundFeesToken0, err = decoder.ReadUint64(bin.LE); err != nil {
		return err
	}
	if c.FundFeesToken1, err = decoder.ReadUint64(bin.LE); err != nil {
		return err
	}
	if c.OpenTime, err = decoder.ReadUint64(bin.LE); err != nil {
		return err
	}
	if c.RecentEpoch, err = decoder.ReadUint64(bin.LE); err != nil {
		return err
	}
	for i := 0; i < 24; i++ {
		if c.Padding1[i], err = decoder.ReadUint64(bin.LE); err != nil {
			return err
		}
	}
	for i := 0; i < 32; i++ {
		if c.Padding2[i], err = decoder.ReadUint64(bin.LE); err != nil {
			return err
		}
	}
	return nil
}
