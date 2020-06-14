package utreexo

import (
	"fmt"
	"github.com/Qitmeer/qitmeer/common/hash"
	"github.com/Qitmeer/qitmeer/core/types"
	"github.com/Qitmeer/qitmeer/database"
	"github.com/Qitmeer/qitmeer/engine/txscript"
	"github.com/Qitmeer/qitmeer/services/index"
	"github.com/Qitmeer/qitmeer/utreexo/accumulator"
)

func (bn *UtreexoModule) buildProofs(msg *addBlockMsg) error {
	if !bn.info.indexHash.IsEqual(&hash.ZeroHash) {
		if bn.info.indexOrder+1 != msg.order {
			return fmt.Errorf("Must be continuous block.(%d)-->(%d)", bn.info.indexOrder, msg.order)
		}
	}

	// Get the add and remove data needed from the block & undo block
	blockAdds, delLeaves, err := bn.blockToAddDel(msg)
	if err != nil {
		return err
	}

	// use the accumulator to get inclusion proofs, and produce a block
	// proof with all data needed to verify the block
	ud, err := genUData(delLeaves, bn.forest, uint32(msg.order))
	if err != nil {
		return err
	}

	err = bn.db.Update(func(dbTx database.Tx) error {
		return DBPutUData(dbTx, &ud)
	})
	if err != nil {
		return err
	}

	ud.AccProof.SortTargets()
	_, err = bn.forest.Modify(blockAdds, ud.AccProof.Targets)
	if err != nil {
		return err
	}
	bn.info.indexOrder = msg.order
	bn.info.indexHash = *msg.blk.Hash()
	err = bn.updateDB()
	if err != nil {
		log.Error(err.Error())
		return err
	}
	return nil
}

func (bn *UtreexoModule) blockToAddDel(msg *addBlockMsg) (blockAdds []accumulator.Leaf,
	delLeaves []LeafData, err error) {

	inskip, outskip := bn.DedupeBlock(msg)
	// fmt.Printf("inskip %v outskip %v\n", inskip, outskip)
	delLeaves, err = bn.blockToDelLeaves(msg, inskip)
	if err != nil {
		return
	}

	// this is bridgenode, so don't need to deal with memorable leaves
	blockAdds = bn.BlockToAddLeaves(msg, nil, outskip, uint32(msg.order))

	return
}

func (bn *UtreexoModule) DedupeBlock(msg *addBlockMsg) (inskip []uint32, outskip []uint32) {

	var i uint32
	// wire.Outpoints are comparable with == which is nice.
	inmap := make(map[types.TxOutPoint]uint32)

	// go through txs then inputs building map
	for cbif0, tx := range msg.blk.Transactions() {
		if cbif0 == 0 { // coinbase tx can't be deduped
			i++
			continue
		}
		_, ok := msg.txs[cbif0]
		if !ok {
			i += uint32(len(tx.Tx.TxIn))
			continue
		}
		for _, in := range tx.Tx.TxIn {
			// fmt.Printf("%s into inmap\n", in.PreviousOutPoint.String())
			inmap[in.PreviousOut] = i
			i++
		}
	}

	i = 0
	// start over, go through outputs finding skips
	for cbif0, tx := range msg.blk.Transactions() {
		if cbif0 == 0 { // coinbase tx can't be deduped
			i += uint32(len(tx.Tx.TxOut))
			continue
		}
		_, ok := msg.txs[cbif0]
		if !ok {
			i += uint32(len(tx.Tx.TxOut))
			continue
		}

		txid := tx.Tx.TxHash()

		for outidx, _ := range tx.Tx.TxOut {
			op := types.TxOutPoint{Hash: txid, OutIndex: uint32(outidx)}
			// fmt.Printf("%s check for inmap... ", op.String())
			inpos, exists := inmap[op]
			if exists {
				// fmt.Printf("hit")
				inskip = append(inskip, inpos)
				outskip = append(outskip, i)
			}
			// fmt.Printf("\n")
			i++
		}
	}
	// sort inskip list, as it's built in order consumed not created
	sortUint32s(inskip)
	return
}

// genDels generates txs to be deleted from the Utreexo forest. These are TxIns
func (bn *UtreexoModule) blockToDelLeaves(msg *addBlockMsg, skiplist []uint32) (
	delLeaves []LeafData, err error) {

	var blockInIdx uint32
	for txinblock, tx := range msg.blk.Transactions() {
		if txinblock == 0 {
			blockInIdx++ // coinbase tx always has 1 input
			continue
		}
		_, ok := msg.txs[txinblock]
		if !ok {
			blockInIdx += uint32(len(tx.Tx.TxIn))
			continue
		}
		// loop through inputs
		for _, txin := range tx.Tx.TxIn {
			// check if on skiplist.  If so, don't make leaf
			if len(skiplist) > 0 && skiplist[0] == blockInIdx {
				skiplist = skiplist[1:]
				blockInIdx++
				continue
			}

			// build leaf
			var l LeafData
			var preTx *types.Transaction
			var preBlockH *hash.Hash
			err := bn.db.View(func(dbTx database.Tx) error {
				dtx, blockH, erro := index.DBFetchTxAndBlock(dbTx, &txin.PreviousOut.Hash)
				if erro != nil {
					return erro
				}
				preTx = dtx
				preBlockH = blockH
				return nil
			})

			if err != nil {
				log.Error(err.Error())
				blockInIdx++
				continue
			}
			ib := bn.bd.GetBlock(preBlockH)
			if ib == nil {
				log.Error(err.Error())
				blockInIdx++
				continue
			}
			l.Outpoint = txin.PreviousOut
			l.Order = uint32(ib.GetOrder())
			l.Coinbase = preTx.IsCoinBase()
			// TODO get blockhash from headers here -- empty for now
			// l.BlockHash = getBlockHashByHeight(l.CbHeight >> 1)
			l.Amt = int64(preTx.TxOut[txin.PreviousOut.OutIndex].Amount)
			l.PkScript = preTx.TxOut[txin.PreviousOut.OutIndex].GetPkScript()
			delLeaves = append(delLeaves, l)
			blockInIdx++
		}
	}
	return
}

func (bn *UtreexoModule) BlockToAddLeaves(msg *addBlockMsg,
	remember []bool, skiplist []uint32,
	order uint32) (leaves []accumulator.Leaf) {

	var txonum uint32
	// bh := bl.Blockhash
	for coinbaseif0, tx := range msg.blk.Transactions() {
		_, ok := msg.txs[coinbaseif0]
		if !ok {
			txonum += uint32(len(tx.Tx.TxOut))
			continue
		}
		// cache txid aka txhash
		txid := tx.Tx.TxHash()
		for i, out := range tx.Tx.TxOut {
			// Skip all the OP_RETURNs
			if txscript.IsUnspendable(out.GetPkScript()) {
				txonum++
				continue
			}

			// Skip txos on the skip list
			if len(skiplist) > 0 && skiplist[0] == txonum {
				skiplist = skiplist[1:]
				txonum++
				continue
			}

			var l LeafData
			// TODO put blockhash back in -- leaving empty for now!
			// l.BlockHash = bh
			l.Outpoint.Hash = txid
			l.Outpoint.OutIndex = uint32(i)
			l.Order = order
			if coinbaseif0 == 0 {
				l.Coinbase = true
			}
			l.Amt = int64(out.Amount)
			l.PkScript = out.PkScript
			uleaf := accumulator.Leaf{Hash: l.LeafHash()}
			if uint32(len(remember)) > txonum {
				uleaf.Remember = remember[txonum]
			}
			leaves = append(leaves, uleaf)
			// fmt.Printf("add %s\n", l.ToString())
			// fmt.Printf("add %s -> %x\n", l.Outpoint.String(), l.LeafHash())
			txonum++
		}
	}
	return
}

func genUData(delLeaves []LeafData, f *accumulator.Forest, order uint32) (
	ud UData, err error) {

	ud.Order = order
	ud.UtxoData = delLeaves
	// make slice of hashes from leafdata
	delHashes := make([]accumulator.Hash, len(ud.UtxoData))
	for i, _ := range ud.UtxoData {
		delHashes[i] = ud.UtxoData[i].LeafHash()
		// fmt.Printf("del %s -> %x\n",
		// ud.UtxoData[i].Outpoint.String(), delHashes[i][:4])
	}
	// generate block proof. Errors if the tx cannot be proven
	// Should never error out with genproofs as it takes
	// blk*.dat files which have already been vetted by Bitcoin Core
	ud.AccProof, err = f.ProveBatch(delHashes)
	if err != nil {
		err = fmt.Errorf("genUData failed at block %d %s %s",
			order, f.Stats(), err.Error())
		return
	}

	if len(ud.AccProof.Targets) != len(delLeaves) {
		err = fmt.Errorf("genUData %d targets but %d leafData",
			len(ud.AccProof.Targets), len(delLeaves))
		return
	}

	return
}
