/*
 * Copyright (c) 2020.
 * Project:qitmeer
 * File:utreexoinfo.go
 * Date:5/25/20 3:56 PM
 * Author:Jin
 * Email:lochjin@gmail.com
 */

package utreexo

import (
	"fmt"
	"github.com/Qitmeer/qitmeer/common/hash"
	"github.com/Qitmeer/qitmeer/core/dbnamespace"
)

type UtreexoInfo struct {
	indexHash  hash.Hash
	indexOrder uint
}

func (ui *UtreexoInfo) deserialize(serializedData []byte) error {
	expectedMinLen := hash.HashSize + 4
	if len(serializedData) < expectedMinLen {
		return fmt.Errorf("corrupt size; min %v got %v", expectedMinLen, len(serializedData))
	}

	copy(ui.indexHash[:], serializedData[0:hash.HashSize])
	offset := uint32(hash.HashSize)
	ui.indexOrder = uint(dbnamespace.ByteOrder.Uint32(serializedData[offset : offset+4]))
	offset += 4

	return nil
}

func (ui *UtreexoInfo) GetHash() *hash.Hash {
	return &ui.indexHash
}

func (ui *UtreexoInfo) isValid() bool {
	return !ui.indexHash.IsEqual(&hash.ZeroHash)
}
