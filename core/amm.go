package core

import (
	"fmt"

	bin "github.com/gagliardetto/binary"
	"github.com/gagliardetto/solana-go"
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

// Serum/OpenBook V3 市场账户结构体（字段根据实际需求补充）
type MarketStateLayoutV3 struct {
	OwnAddress   solana.PublicKey
	VaultSigner  solana.PublicKey
	BaseMint     solana.PublicKey
	QuoteMint    solana.PublicKey
	BaseVault    solana.PublicKey
	QuoteVault   solana.PublicKey
	Bids         solana.PublicKey
	Asks         solana.PublicKey
	EventQueue   solana.PublicKey
	RequestQueue solana.PublicKey
	// ...补充其他字段...
}

func (m *MarketStateLayoutV3) UnmarshalWithDecoder(decoder *bin.Decoder) error {
	var err error
	readPk := func(pk *solana.PublicKey) error {
		b, err := decoder.ReadBytes(32)
		if err != nil {
			return err
		}
		copy(pk[:], b)
		return nil
	}
	if err = readPk(&m.OwnAddress); err != nil {
		return fmt.Errorf("read OwnAddress: %w", err)
	}
	if err = readPk(&m.VaultSigner); err != nil {
		return fmt.Errorf("read VaultSigner: %w", err)
	}
	if err = readPk(&m.BaseMint); err != nil {
		return fmt.Errorf("read BaseMint: %w", err)
	}
	if err = readPk(&m.QuoteMint); err != nil {
		return fmt.Errorf("read QuoteMint: %w", err)
	}
	if err = readPk(&m.BaseVault); err != nil {
		return fmt.Errorf("read BaseVault: %w", err)
	}
	if err = readPk(&m.QuoteVault); err != nil {
		return fmt.Errorf("read QuoteVault: %w", err)
	}
	if err = readPk(&m.Bids); err != nil {
		return fmt.Errorf("read Bids: %w", err)
	}
	if err = readPk(&m.Asks); err != nil {
		return fmt.Errorf("read Asks: %w", err)
	}
	if err = readPk(&m.EventQueue); err != nil {
		return fmt.Errorf("read EventQueue: %w", err)
	}
	if err = readPk(&m.RequestQueue); err != nil {
		return fmt.Errorf("read RequestQueue: %w", err)
	}
	// ...继续解析其他字段...
	return nil
}
