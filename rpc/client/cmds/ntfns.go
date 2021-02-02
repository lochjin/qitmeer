package cmds

import (
	"github.com/Qitmeer/qitmeer/core/json"
	"github.com/Qitmeer/qitmeer/core/types"
)

const (
	BlockConnectedNtfnMethod    = "blockConnected"
	BlockDisconnectedNtfnMethod = "blockDisconnected"
	BlockAcceptedNtfnMethod     = "blockAccepted"
	ReorganizationNtfnMethod    = "reorganization"
	TxAcceptedNtfnMethod        = "txaccepted"
	TxAcceptedVerboseNtfnMethod = "txacceptedverbose"
)

type BlockConnectedNtfn struct {
	Hash  string
	Order int64
	Time  int64
	Txs   []string
}

func NewBlockConnectedNtfn(hash string, order int64, time int64, txs []string) *BlockConnectedNtfn {
	return &BlockConnectedNtfn{
		Hash:  hash,
		Order: order,
		Time:  time,
		Txs:   txs,
	}
}

type BlockDisconnectedNtfn struct {
	Hash  string
	Order int64
	Time  int64
	Txs   []string
}

func NewBlockDisconnectedNtfn(hash string, order int64, time int64, txs []string) *BlockDisconnectedNtfn {
	return &BlockDisconnectedNtfn{
		Hash:  hash,
		Order: order,
		Time:  time,
		Txs:   txs,
	}
}

type BlockAcceptedNtfn struct {
	Hash  string
	Order int64
	Time  int64
	Txs   []string
}

func NewBlockAcceptedNtfn(hash string, order int64, time int64, txs []string) *BlockAcceptedNtfn {
	return &BlockAcceptedNtfn{
		Hash:  hash,
		Order: order,
		Time:  time,
		Txs:   txs,
	}
}

type ReorganizationNtfn struct {
	Hash  string
	Order int64
	Olds  []string
}

func NewReorganizationNtfn(hash string, order int64, olds []string) *ReorganizationNtfn {
	return &ReorganizationNtfn{
		Hash:  hash,
		Order: order,
		Olds:  olds,
	}
}

type TxAcceptedNtfn struct {
	TxID    string
	Amounts types.AmountGroup
}

func NewTxAcceptedNtfn(txHash string, amounts types.AmountGroup) *TxAcceptedNtfn {
	return &TxAcceptedNtfn{
		TxID:    txHash,
		Amounts: amounts,
	}
}

type TxAcceptedVerboseNtfn struct {
	Tx json.DecodeRawTransactionResult
}

func NewTxAcceptedVerboseNtfn(tx json.DecodeRawTransactionResult) *TxAcceptedVerboseNtfn {
	return &TxAcceptedVerboseNtfn{
		Tx: tx,
	}
}

func init() {
	flags := UFWebsocketOnly | UFNotification

	MustRegisterCmd(BlockConnectedNtfnMethod, (*BlockConnectedNtfn)(nil), flags, NotifyNameSpace)
	MustRegisterCmd(BlockDisconnectedNtfnMethod, (*BlockDisconnectedNtfn)(nil), flags, NotifyNameSpace)
	MustRegisterCmd(BlockAcceptedNtfnMethod, (*BlockAcceptedNtfn)(nil), flags, NotifyNameSpace)
	MustRegisterCmd(ReorganizationNtfnMethod, (*ReorganizationNtfn)(nil), flags, NotifyNameSpace)
	MustRegisterCmd(TxAcceptedNtfnMethod, (*TxAcceptedNtfn)(nil), flags, NotifyNameSpace)
	MustRegisterCmd(TxAcceptedVerboseNtfnMethod, (*TxAcceptedVerboseNtfn)(nil), flags, NotifyNameSpace)
}
