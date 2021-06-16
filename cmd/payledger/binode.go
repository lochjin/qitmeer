/*
 * Copyright (c) 2020.
 * Project:qitmeer
 * File:binode.go
 * Date:6/6/20 9:28 AM
 * Author:Jin
 * Email:lochjin@gmail.com
 */

package main

import (
	"fmt"
	"github.com/Qitmeer/qitmeer/core/blockchain"
	"github.com/Qitmeer/qitmeer/core/blockdag"
	"github.com/Qitmeer/qitmeer/database"
	"github.com/Qitmeer/qitmeer/log"
	"github.com/Qitmeer/qitmeer/params"
	"github.com/Qitmeer/qitmeer/services/index"
	"github.com/Qitmeer/qitmeer/services/mining"
	"path"
)

type BINode struct {
	name string
	bc   *blockchain.BlockChain
	db   database.DB
	cfg  *Config
}

func (node *BINode) init(cfg *Config) error {
	node.cfg = cfg

	// Load the block database.
	db, err := LoadBlockDB(cfg.DbType, cfg.DataDir, true)
	if err != nil {
		log.Error("load block database", "error", err)
		return err
	}

	node.db = db
	//
	var indexes []index.Indexer
	txIndex := index.NewTxIndex(db)
	indexes = append(indexes, txIndex)
	// index-manager
	indexManager := index.NewManager(db, indexes, params.ActiveNetParams.Params)

	bc, err := blockchain.New(&blockchain.Config{
		DB:           db,
		ChainParams:  params.ActiveNetParams.Params,
		TimeSource:   blockchain.NewMedianTime(),
		DAGType:      cfg.DAGType,
		BlockVersion: mining.BlockVersion(params.ActiveNetParams.Params.Net),
		IndexManager: indexManager,
	})
	if err != nil {
		log.Error(err.Error())
		return err
	}
	node.bc = bc
	node.name = path.Base(cfg.DataDir)

	log.Info(fmt.Sprintf("Load Data:%s", cfg.DataDir))

	return node.statistics()
}

func (node *BINode) exit() {
	if node.db != nil {
		log.Info(fmt.Sprintf("Gracefully shutting down the database:%s", node.name))
		node.db.Close()
	}
}

func (node *BINode) BlockChain() *blockchain.BlockChain {
	return node.bc
}

func (node *BINode) DB() database.DB {
	return node.db
}

func (node *BINode) statistics() error {
	bd := node.bc.BlockDAG()

	endID := uint(0)

	// 1:find end
	mainTip := bd.GetMainChainTip()
	for cur := mainTip; cur != nil; cur = node.bc.BlockDAG().GetBlockById(cur.GetMainParent()) {
		block, err := node.bc.FetchBlockByHash(cur.GetHash())
		if err != nil {
			return err
		}
		if block.Transactions()[0].Tx.TxOut[0].Amount != 0 {
			endID = cur.GetID()
			break
		}
	}

	if endID == 0 {
		return fmt.Errorf("End block id error\n")
	}

	total := endID + 1
	validCount := 1
	subsidyCount := 0
	subsidy := uint64(0)

	// all blue blocks
	mainTip = bd.GetBlockById(endID)
	log.Info(fmt.Sprintf("Find mining main tip:id=%d hash=%s", endID, mainTip.GetHash()))

	bluesArr := make([]bool, total)

	for cur := mainTip; cur != nil; cur = node.bc.BlockDAG().GetBlockById(cur.GetMainParent()) {
		bluesArr[cur.GetID()] = true
		for k := range cur.(*blockdag.PhantomBlock).GetBlueDiffAnticone().GetMap() {
			bluesArr[k] = true
		}
	}

	// genesis
	genblock, err := node.bc.FetchBlockByHash(bd.GetGenesisHash())
	if err != nil {
		return err
	}
	gencoinbaseAmout := uint64(0)
	for _, v := range genblock.Transactions()[0].Tx.TxOut {
		gencoinbaseAmout += v.Amount
	}
	fmt.Printf("id:%8d hash:%s txsvalid:%t coinbaseAmout:%18d coinbaseValid:%t isBlue:%t\n", 0, bd.GetGenesisHash(), true, gencoinbaseAmout, true, true)

	coinbaseValidNum := 1

	for i := uint(1); i < total; i++ {
		ib := node.bc.BlockDAG().GetBlockById(i)
		if ib == nil {
			return fmt.Errorf("No block:%d\n", i)
		}
		if !ib.IsOrdered() {
			return fmt.Errorf("Block is not ordered:%s (%d)\n", ib.GetHash(), ib.GetID())
		}
		block, err := node.bc.FetchBlockByHash(ib.GetHash())
		if err != nil {
			return err
		}

		coinbaseAmout := block.Transactions()[0].Tx.TxOut[0].Amount
		coinbaseValid := false

		if !knownInvalid(byte(ib.GetStatus())) {
			validCount++
		}

		txfullHash := block.Transactions()[0].Tx.TxHashFull()

		if isTxValid(node.db, block.Transactions()[0].Hash(), &txfullHash, ib.GetHash()) {
			if bluesArr[i] {
				subsidyCount++
				subsidy += block.Transactions()[0].Tx.TxOut[0].Amount

			}
			coinbaseValid = true
		}

		if coinbaseValid {
			coinbaseValidNum++
		}
		fmt.Printf("id:%8d hash:%s txsvalid:%t coinbaseAmout:%18d coinbaseValid:%t isBlue:%t\n", i, ib.GetHash(), !knownInvalid(byte(ib.GetStatus())), coinbaseAmout, coinbaseValid, bluesArr[i])
	}

	blues := mainTip.(*blockdag.PhantomBlock).GetBlueNum() + 1
	reds := mainTip.GetOrder() + 1 - blues
	unconfirmed := total - (mainTip.GetOrder() + 1)

	fmt.Println()

	resultStr := fmt.Sprintf("数据库分析:(包含创世区块; 主链Tip Hash: %s (%d)\n总区块数:%d\n合法的区块数:%d\nCoinbase数据库记录数:%d\n包含重复Coinbase交易的区块数:%d\n红色块数:%d\n蓝色块数:%d\n合法蓝色区块:%d (不包括创世块)\n铸币量:%d * 12000000000 = %d\n流通量:%d + %d = %d\n",
		mainTip.GetHash(), endID, total, validCount, coinbaseValidNum, validCount-coinbaseValidNum, reds, blues, subsidyCount, subsidyCount, subsidy, gencoinbaseAmout, subsidy, gencoinbaseAmout+subsidy)

	if unconfirmed > 0 {
		resultStr = fmt.Sprintf("%s Unconfirmed:%d", resultStr, unconfirmed)
	}
	fmt.Println(resultStr)
	log.Info(resultStr)

	return nil
}
