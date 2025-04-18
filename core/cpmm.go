package core

import (
	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
)

type CPMM_STATE_LAYOUT struct {
	AmmConfig          solana.PublicKey
	PoolCreator        solana.PublicKey
	Token0Vault        solana.PublicKey
	Token1Vault        solana.PublicKey
	LpMint             solana.PublicKey
	Token0Mint         solana.PublicKey
	Token1Mint         solana.PublicKey
	Token0Program      solana.PublicKey
	Token1Program      solana.PublicKey
	ObservationKey     solana.PublicKey
	AuthBump           uint8
	Status             uint8
	LpMintDecimals     uint8
	Mint0Decimals      uint8
	Mint1Decimals      uint8
	LpSupply           uint64
	ProtocolFeesToken0 uint64
	ProtocolFeesToken1 uint64
	FundFeesToken0     uint64
	FundFeesToken1     uint64
	OpenTime           uint64
	Padding            [32]uint64
}

func (c *CPMM_STATE_LAYOUT) UnmarshalWithDecoder(decoder *bin.Decoder) error {
	var err error
	// 跳过8字节头部填充
	if _, err = decoder.ReadBytes(8); err != nil {
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
	if err = readPk(&c.PoolCreator); err != nil {
		return err
	}
	if err = readPk(&c.Token0Vault); err != nil {
		return err
	}
	if err = readPk(&c.Token1Vault); err != nil {
		return err
	}
	if err = readPk(&c.LpMint); err != nil {
		return err
	}
	if err = readPk(&c.Token0Mint); err != nil {
		return err
	}
	if err = readPk(&c.Token1Mint); err != nil {
		return err
	}
	if err = readPk(&c.Token0Program); err != nil {
		return err
	}
	if err = readPk(&c.Token1Program); err != nil {
		return err
	}
	if err = readPk(&c.ObservationKey); err != nil {
		return err
	}
	if c.AuthBump, err = decoder.ReadUint8(); err != nil {
		return err
	}
	if c.Status, err = decoder.ReadUint8(); err != nil {
		return err
	}
	if c.LpMintDecimals, err = decoder.ReadUint8(); err != nil {
		return err
	}
	if c.Mint0Decimals, err = decoder.ReadUint8(); err != nil {
		return err
	}
	if c.Mint1Decimals, err = decoder.ReadUint8(); err != nil {
		return err
	}
	if c.LpSupply, err = decoder.ReadUint64(bin.LE); err != nil {
		return err
	}
	if c.ProtocolFeesToken0, err = decoder.ReadUint64(bin.LE); err != nil {
		return err
	}
	if c.ProtocolFeesToken1, err = decoder.ReadUint64(bin.LE); err != nil {
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
	for i := 0; i < 32; i++ {
		if c.Padding[i], err = decoder.ReadUint64(bin.LE); err != nil {
			return err
		}
	}
	return nil
}
