package core

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
	"github.com/gagliardetto/solana-go/rpc"
)

type LIQUIDITY_STATE_LAYOUT_V4 struct {
	Status                 uint64
	Nonce                  uint64
	MaxOrder               uint64
	Depth                  uint64
	BaseDecimal            uint64
	QuoteDecimal           uint64
	State                  uint64
	ResetFlag              uint64
	MinSize                uint64
	VolMaxCutRatio         uint64
	AmountWaveRatio        uint64
	BaseLotSize            uint64
	QuoteLotSize           uint64
	MinPriceMultiplier     uint64
	MaxPriceMultiplier     uint64
	SystemDecimalValue     uint64
	MinSeparateNumerator   uint64
	MinSeparateDenominator uint64
	TradeFeeNumerator      uint64
	TradeFeeDenominator    uint64
	PnlNumerator           uint64
	PnlDenominator         uint64
	SwapFeeNumerator       uint64
	SwapFeeDenominator     uint64
	BaseNeedTakePnl        uint64
	QuoteNeedTakePnl       uint64
	QuoteTotalPnl          uint64
	BaseTotalPnl           uint64
	PoolOpenTime           uint64
	PunishPcAmount         uint64
	PunishCoinAmount       uint64
	OrderbookToInitTime    uint64

	SwapBaseInAmount   [16]byte
	SwapQuoteOutAmount [16]byte
	SwapBase2QuoteFee  uint64
	SwapQuoteInAmount  [16]byte
	SwapBaseOutAmount  [16]byte
	SwapQuote2BaseFee  uint64

	BaseVault       solana.PublicKey
	QuoteVault      solana.PublicKey
	BaseMint        solana.PublicKey
	QuoteMint       solana.PublicKey
	LpMint          solana.PublicKey
	OpenOrders      solana.PublicKey
	MarketID        solana.PublicKey
	MarketProgramID solana.PublicKey
	TargetOrders    solana.PublicKey
	WithdrawQueue   solana.PublicKey
	LpVault         solana.PublicKey
	Owner           solana.PublicKey

	LpReserve uint64
	Padding   [3]uint64
}

func (l *LIQUIDITY_STATE_LAYOUT_V4) UnmarshalWithDecoder(decoder *bin.Decoder) error {
	var err error

	readBytes := func(dst []byte) error {
		b, err := decoder.ReadBytes(len(dst))
		if err != nil {
			return err
		}
		copy(dst, b)
		return nil
	}

	readU64 := func(field *uint64, name string) error {
		*field, err = decoder.ReadUint64(bin.LE)
		if err != nil {
			return fmt.Errorf("failed to decode %s: %w", name, err)
		}
		return nil
	}

	readPk := func(pk *solana.PublicKey, name string) error {
		b, err := decoder.ReadBytes(32)
		if err != nil {
			return fmt.Errorf("failed to decode %s: %w", name, err)
		}
		copy(pk[:], b)
		return nil
	}

	if err = readU64(&l.Status, "Status"); err != nil {
		return err
	}
	if err = readU64(&l.Nonce, "Nonce"); err != nil {
		return err
	}
	if err = readU64(&l.MaxOrder, "MaxOrder"); err != nil {
		return err
	}
	if err = readU64(&l.Depth, "Depth"); err != nil {
		return err
	}
	if err = readU64(&l.BaseDecimal, "BaseDecimal"); err != nil {
		return err
	}
	if err = readU64(&l.QuoteDecimal, "QuoteDecimal"); err != nil {
		return err
	}
	if err = readU64(&l.State, "State"); err != nil {
		return err
	}
	if err = readU64(&l.ResetFlag, "ResetFlag"); err != nil {
		return err
	}
	if err = readU64(&l.MinSize, "MinSize"); err != nil {
		return err
	}
	if err = readU64(&l.VolMaxCutRatio, "VolMaxCutRatio"); err != nil {
		return err
	}
	if err = readU64(&l.AmountWaveRatio, "AmountWaveRatio"); err != nil {
		return err
	}
	if err = readU64(&l.BaseLotSize, "BaseLotSize"); err != nil {
		return err
	}
	if err = readU64(&l.QuoteLotSize, "QuoteLotSize"); err != nil {
		return err
	}
	if err = readU64(&l.MinPriceMultiplier, "MinPriceMultiplier"); err != nil {
		return err
	}
	if err = readU64(&l.MaxPriceMultiplier, "MaxPriceMultiplier"); err != nil {
		return err
	}
	if err = readU64(&l.SystemDecimalValue, "SystemDecimalValue"); err != nil {
		return err
	}
	if err = readU64(&l.MinSeparateNumerator, "MinSeparateNumerator"); err != nil {
		return err
	}
	if err = readU64(&l.MinSeparateDenominator, "MinSeparateDenominator"); err != nil {
		return err
	}
	if err = readU64(&l.TradeFeeNumerator, "TradeFeeNumerator"); err != nil {
		return err
	}
	if err = readU64(&l.TradeFeeDenominator, "TradeFeeDenominator"); err != nil {
		return err
	}
	if err = readU64(&l.PnlNumerator, "PnlNumerator"); err != nil {
		return err
	}
	if err = readU64(&l.PnlDenominator, "PnlDenominator"); err != nil {
		return err
	}
	if err = readU64(&l.SwapFeeNumerator, "SwapFeeNumerator"); err != nil {
		return err
	}
	if err = readU64(&l.SwapFeeDenominator, "SwapFeeDenominator"); err != nil {
		return err
	}
	if err = readU64(&l.BaseNeedTakePnl, "BaseNeedTakePnl"); err != nil {
		return err
	}
	if err = readU64(&l.QuoteNeedTakePnl, "QuoteNeedTakePnl"); err != nil {
		return err
	}
	if err = readU64(&l.QuoteTotalPnl, "QuoteTotalPnl"); err != nil {
		return err
	}
	if err = readU64(&l.BaseTotalPnl, "BaseTotalPnl"); err != nil {
		return err
	}
	if err = readU64(&l.PoolOpenTime, "PoolOpenTime"); err != nil {
		return err
	}
	if err = readU64(&l.PunishPcAmount, "PunishPcAmount"); err != nil {
		return err
	}
	if err = readU64(&l.PunishCoinAmount, "PunishCoinAmount"); err != nil {
		return err
	}
	if err = readU64(&l.OrderbookToInitTime, "OrderbookToInitTime"); err != nil {
		return err
	}

	if err = readBytes(l.SwapBaseInAmount[:]); err != nil {
		return fmt.Errorf("failed to decode SwapBaseInAmount: %w", err)
	}
	if err = readBytes(l.SwapQuoteOutAmount[:]); err != nil {
		return fmt.Errorf("failed to decode SwapQuoteOutAmount: %w", err)
	}
	if err = readU64(&l.SwapBase2QuoteFee, "SwapBase2QuoteFee"); err != nil {
		return err
	}
	if err = readBytes(l.SwapQuoteInAmount[:]); err != nil {
		return fmt.Errorf("failed to decode SwapQuoteInAmount: %w", err)
	}
	if err = readBytes(l.SwapBaseOutAmount[:]); err != nil {
		return fmt.Errorf("failed to decode SwapBaseOutAmount: %w", err)
	}
	if err = readU64(&l.SwapQuote2BaseFee, "SwapQuote2BaseFee"); err != nil {
		return err
	}

	if err = readPk(&l.BaseVault, "BaseVault"); err != nil {
		return err
	}
	if err = readPk(&l.QuoteVault, "QuoteVault"); err != nil {
		return err
	}
	if err = readPk(&l.BaseMint, "BaseMint"); err != nil {
		return err
	}
	if err = readPk(&l.QuoteMint, "QuoteMint"); err != nil {
		return err
	}
	if err = readPk(&l.LpMint, "LpMint"); err != nil {
		return err
	}
	if err = readPk(&l.OpenOrders, "OpenOrders"); err != nil {
		return err
	}
	if err = readPk(&l.MarketID, "MarketID"); err != nil {
		return err
	}
	if err = readPk(&l.MarketProgramID, "MarketProgramID"); err != nil {
		return err
	}
	if err = readPk(&l.TargetOrders, "TargetOrders"); err != nil {
		return err
	}
	if err = readPk(&l.WithdrawQueue, "WithdrawQueue"); err != nil {
		return err
	}
	if err = readPk(&l.LpVault, "LpVault"); err != nil {
		return err
	}
	if err = readPk(&l.Owner, "Owner"); err != nil {
		return err
	}

	if err = readU64(&l.LpReserve, "LpReserve"); err != nil {
		return err
	}
	for i := 0; i < 3; i++ {
		if err = readU64(&l.Padding[i], fmt.Sprintf("Padding[%d]", i)); err != nil {
			return err
		}
	}

	return nil
}

func (l *LIQUIDITY_STATE_LAYOUT_V4) MarketInfo(ctx context.Context, client *rpc.Client) (*MarketStateLayoutV3, error) {
	info, err := client.GetAccountInfoWithOpts(ctx, l.MarketID, &rpc.GetAccountInfoOpts{
		Commitment: rpc.CommitmentConfirmed,
	})
	if err != nil {
		return nil, err
	}
	data := info.GetBinary()
	if data == nil {
		return nil, errors.New("pool account not found")
	}
	// 1. 尝试V4
	result := new(MarketStateLayoutV3)
	if err := result.UnmarshalBinary(data); err == nil {
		return nil, nil
	}
	return result, nil
}

type AccountFlags struct {
	Initialized  bool
	Market       bool
	OpenOrders   bool
	RequestQueue bool
	EventQueue   bool
	Bids         bool
	Asks         bool
	// 剩余57位为0
}

type MarketStateLayoutV3 struct {
	Padding0              [5]byte
	AccountFlags          AccountFlags
	OwnAddress            [32]byte
	VaultSignerNonce      uint64
	BaseMint              [32]byte
	QuoteMint             [32]byte
	BaseVault             [32]byte
	BaseDepositsTotal     uint64
	BaseFeesAccrued       uint64
	QuoteVault            [32]byte
	QuoteDepositsTotal    uint64
	QuoteFeesAccrued      uint64
	QuoteDustThreshold    uint64
	RequestQueue          [32]byte
	EventQueue            [32]byte
	Bids                  [32]byte
	Asks                  [32]byte
	BaseLotSize           uint64
	QuoteLotSize          uint64
	FeeRateBps            uint64
	ReferrerRebateAccrued uint64
	Padding1              [7]byte
}

// AccountFlags序列化为8字节（低7位为flag，高57位为0）
func (f *AccountFlags) MarshalBinary() ([]byte, error) {
	var v uint64
	if f.Initialized {
		v |= 1 << 0
	}
	if f.Market {
		v |= 1 << 1
	}
	if f.OpenOrders {
		v |= 1 << 2
	}
	if f.RequestQueue {
		v |= 1 << 3
	}
	if f.EventQueue {
		v |= 1 << 4
	}
	if f.Bids {
		v |= 1 << 5
	}
	if f.Asks {
		v |= 1 << 6
	}
	buf := make([]byte, 8)
	binary.LittleEndian.PutUint64(buf, v)
	return buf, nil
}

func (f *AccountFlags) UnmarshalBinary(data []byte) error {
	if len(data) < 8 {
		return errors.New("data too short for AccountFlags")
	}
	v := binary.LittleEndian.Uint64(data)
	f.Initialized = v&(1<<0) != 0
	f.Market = v&(1<<1) != 0
	f.OpenOrders = v&(1<<2) != 0
	f.RequestQueue = v&(1<<3) != 0
	f.EventQueue = v&(1<<4) != 0
	f.Bids = v&(1<<5) != 0
	f.Asks = v&(1<<6) != 0
	return nil
}

// MarketStateV3序列化
func (m *MarketStateLayoutV3) MarshalBinary() ([]byte, error) {
	buf := new(bytes.Buffer)
	buf.Write(m.Padding0[:])
	af, _ := m.AccountFlags.MarshalBinary()
	buf.Write(af)
	buf.Write(m.OwnAddress[:])
	binary.Write(buf, binary.LittleEndian, m.VaultSignerNonce)
	buf.Write(m.BaseMint[:])
	buf.Write(m.QuoteMint[:])
	buf.Write(m.BaseVault[:])
	binary.Write(buf, binary.LittleEndian, m.BaseDepositsTotal)
	binary.Write(buf, binary.LittleEndian, m.BaseFeesAccrued)
	buf.Write(m.QuoteVault[:])
	binary.Write(buf, binary.LittleEndian, m.QuoteDepositsTotal)
	binary.Write(buf, binary.LittleEndian, m.QuoteFeesAccrued)
	binary.Write(buf, binary.LittleEndian, m.QuoteDustThreshold)
	buf.Write(m.RequestQueue[:])
	buf.Write(m.EventQueue[:])
	buf.Write(m.Bids[:])
	buf.Write(m.Asks[:])
	binary.Write(buf, binary.LittleEndian, m.BaseLotSize)
	binary.Write(buf, binary.LittleEndian, m.QuoteLotSize)
	binary.Write(buf, binary.LittleEndian, m.FeeRateBps)
	binary.Write(buf, binary.LittleEndian, m.ReferrerRebateAccrued)
	buf.Write(m.Padding1[:])
	return buf.Bytes(), nil
}

// MarketStateV3反序列化
func (m *MarketStateLayoutV3) UnmarshalBinary(data []byte) error {
	if len(data) < 5+8+32*8+8*8+7 {
		return errors.New("data too short for MarketStateV3")
	}
	offset := 0
	copy(m.Padding0[:], data[offset:offset+5])
	offset += 5
	if err := m.AccountFlags.UnmarshalBinary(data[offset : offset+8]); err != nil {
		return err
	}
	offset += 8
	copy(m.OwnAddress[:], data[offset:offset+32])
	offset += 32
	m.VaultSignerNonce = binary.LittleEndian.Uint64(data[offset : offset+8])
	offset += 8
	copy(m.BaseMint[:], data[offset:offset+32])
	offset += 32
	copy(m.QuoteMint[:], data[offset:offset+32])
	offset += 32
	copy(m.BaseVault[:], data[offset:offset+32])
	offset += 32
	m.BaseDepositsTotal = binary.LittleEndian.Uint64(data[offset : offset+8])
	offset += 8
	m.BaseFeesAccrued = binary.LittleEndian.Uint64(data[offset : offset+8])
	offset += 8
	copy(m.QuoteVault[:], data[offset:offset+32])
	offset += 32
	m.QuoteDepositsTotal = binary.LittleEndian.Uint64(data[offset : offset+8])
	offset += 8
	m.QuoteFeesAccrued = binary.LittleEndian.Uint64(data[offset : offset+8])
	offset += 8
	m.QuoteDustThreshold = binary.LittleEndian.Uint64(data[offset : offset+8])
	offset += 8
	copy(m.RequestQueue[:], data[offset:offset+32])
	offset += 32
	copy(m.EventQueue[:], data[offset:offset+32])
	offset += 32
	copy(m.Bids[:], data[offset:offset+32])
	offset += 32
	copy(m.Asks[:], data[offset:offset+32])
	offset += 32
	m.BaseLotSize = binary.LittleEndian.Uint64(data[offset : offset+8])
	offset += 8
	m.QuoteLotSize = binary.LittleEndian.Uint64(data[offset : offset+8])
	offset += 8
	m.FeeRateBps = binary.LittleEndian.Uint64(data[offset : offset+8])
	offset += 8
	m.ReferrerRebateAccrued = binary.LittleEndian.Uint64(data[offset : offset+8])
	offset += 8
	copy(m.Padding1[:], data[offset:offset+7])
	return nil
}
