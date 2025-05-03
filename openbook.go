package raydium

import (
	"bytes"
	"encoding/binary"
	"errors"
)

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

type OpenBook struct {
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
func (m *OpenBook) MarshalBinary() ([]byte, error) {
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

var OpenBookDiscriminator = [5]byte{115, 101, 114, 117, 109}

// MarketStateV3反序列化
func (m *OpenBook) UnmarshalBinary(data []byte) error {
	if len(data) < 5+8+32*8+8*8+7 {
		return errors.New("data too short for MarketStateV3")
	}
	offset := 0
	copy(m.Padding0[:], data[offset:offset+5])
	if m.Padding0 != OpenBookDiscriminator {
		return errors.New("没有找到openbook的特定标识符")
	}
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
