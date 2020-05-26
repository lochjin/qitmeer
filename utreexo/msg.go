package utreexo

import (
	"github.com/Qitmeer/qitmeer/core/blockdag"
	"github.com/Qitmeer/qitmeer/core/types"
)

type addBlockMsg struct {
	blk   *types.SerializedBlock
	order uint
}

type removeBlockMsg struct {
	blk *types.SerializedBlock
	ib  blockdag.IBlock
}
