package utreexo

import "github.com/Qitmeer/qitmeer/utreexo/accumulator"

type UData struct {
	Order          uint32
	AccProof       accumulator.BatchProof
	UtxoData       []LeafData
	RememberLeaves []bool
}

func (ud *UData) ToBytes() (b []byte) {
	b = U32tB(ud.Order)
	// first stick the batch proof on the beginning
	batchBytes := ud.AccProof.ToBytes()
	b = append(b, U32tB(uint32(len(batchBytes)))...)
	b = append(b, batchBytes...)

	// next, all the leafDatas
	for _, ld := range ud.UtxoData {
		ldb := ld.ToBytes()
		b = append(b, PrefixLen16(ldb)...)
	}

	return
}
