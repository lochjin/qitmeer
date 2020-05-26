package utreexo

import (
	"crypto/sha256"
	"fmt"
	"github.com/Qitmeer/qitmeer/common/hash"
	"github.com/Qitmeer/qitmeer/core/types"
)

// LeafData is all the data that goes into a leaf in the utreexo accumulator
type LeafData struct {
	BlockHash [hash.HashSize]byte
	Outpoint  types.TxOutPoint
	Order     uint32
	Coinbase  bool
	Amt       int64
	PkScript  []byte
}

// turn a LeafData into bytes
func (l *LeafData) ToString() (s string) {
	s = l.Outpoint.String()
	// s += fmt.Sprintf(" bh %x ", l.BlockHash)
	s += fmt.Sprintf("o %d ", l.Order)
	s += fmt.Sprintf("cb %v ", l.Coinbase)
	s += fmt.Sprintf("amt %d ", l.Amt)
	s += fmt.Sprintf("pks %x ", l.PkScript)
	s += fmt.Sprintf("%x", l.LeafHash())
	return
}

func LeafDataFromBytes(b []byte) (LeafData, error) {
	var l LeafData
	if len(b) < 80 {
		return l, fmt.Errorf("%x for leafdata, need 80 bytes", b)
	}
	copy(l.BlockHash[:], b[0:32])
	copy(l.Outpoint.Hash[:], b[32:64])
	l.Outpoint.OutIndex = BtU32(b[64:68])
	l.Order = BtU32(b[68:72])
	if l.Order&1 == 1 {
		l.Coinbase = true
	}
	l.Order >>= 1
	l.Amt = BtI64(b[72:80])
	l.PkScript = b[80:]

	return l, nil
}

// turn a LeafData into bytes
func (l *LeafData) ToBytes() (b []byte) {
	b = append(l.BlockHash[:], l.Outpoint.Hash[:]...)
	b = append(b, U32tB(l.Outpoint.OutIndex)...)
	hcb := l.Order << 1
	if l.Coinbase {
		hcb |= 1
	}
	b = append(b, U32tB(hcb)...)
	b = append(b, I64tB(l.Amt)...)
	b = append(b, l.PkScript...)
	return
}

// turn a LeafData into a LeafHash
func (l *LeafData) LeafHash() [hash.HashSize]byte {
	return sha256.Sum256(l.ToBytes())
}
