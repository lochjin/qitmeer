/*
 * Copyright (c) 2020.
 * Project:qitmeer
 * File:dbhelper.go
 * Date:5/29/20 7:56 PM
 * Author:Jin
 * Email:lochjin@gmail.com
 */

package utreexo

import (
	"github.com/Qitmeer/qitmeer/core/dbnamespace"
	"github.com/Qitmeer/qitmeer/database"
	"os"
)

func DBPutUData(dbTx database.Tx, udata *UData) error {
	bucket := dbTx.Metadata().Bucket(dbnamespace.UtreexoProofBucketName)
	var serializedID [4]byte
	dbnamespace.ByteOrder.PutUint32(serializedID[:], uint32(udata.Order))

	key := serializedID[:]

	return bucket.Put(key, udata.ToBytes())
}

func (bn *UtreexoModule) updateDB() error {
	err := bn.db.Update(func(dbTx database.Tx) error {
		meta := dbTx.Metadata()
		meta.Put(dbnamespace.UtreexoInfoBucketName, bn.info.toBytes())
		return nil
	})
	if err != nil {
		return err
	}

	miscForestFile, err := os.OpenFile(bn.getMiscForestFilePath(), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return err
	}
	err = bn.forest.WriteForest(miscForestFile)
	if err != nil {
		return err
	}
	return nil
}
