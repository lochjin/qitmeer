/*
 * Copyright (c) 2020.
 * Project:qitmeer
 * File:utreexomodule.go
 * Date:5/25/20 11:46 AM
 * Author:Jin
 * Email:lochjin@gmail.com
 */

package utreexo

import (
	"fmt"
	"github.com/Qitmeer/qitmeer/common/hash"
	"github.com/Qitmeer/qitmeer/config"
	"github.com/Qitmeer/qitmeer/core/blockchain"
	"github.com/Qitmeer/qitmeer/core/blockdag"
	"github.com/Qitmeer/qitmeer/core/dbnamespace"
	"github.com/Qitmeer/qitmeer/core/types"
	"github.com/Qitmeer/qitmeer/database"
	"github.com/Qitmeer/qitmeer/services/index"
	"github.com/Qitmeer/qitmeer/utreexo/accumulator"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

const (
	ForestFilePath     = "utreexo/forestfile.dat"
	MiscForestFilePath = "utreexo/miscforestfile.dat"
	UtreexoDirPath     = "utreexo"
)

type UtreexoModule struct {
	bc     *blockchain.BlockChain
	bd     *blockdag.BlockDAG
	db     database.DB
	info   *UtreexoInfo
	cfg    *config.Config
	forest *accumulator.Forest

	started  int32
	shutdown int32
	msgChan  chan interface{}
	wg       sync.WaitGroup
	quit     chan struct{}
}

func (bn *UtreexoModule) Init() error {
	if !bn.cfg.Utreexo {
		log.Info("Utreexo disable.")
		return nil
	}
	log.Info("Utreexo enable:init...")
	err := bn.db.Update(func(dbTx database.Tx) error {
		meta := dbTx.Metadata()
		infoData := meta.Get(dbnamespace.UtreexoInfoBucketName)
		if infoData != nil {
			err := bn.info.deserialize(infoData)
			if err != nil {
				return err
			}
			log.Trace(fmt.Sprintf("Utreexo info:%s %d", bn.info.GetHash().String(), bn.info.indexOrder))
		}

		if meta.Bucket(dbnamespace.UtreexoProofBucketName) == nil {
			_, err := meta.CreateBucket(dbnamespace.UtreexoProofBucketName)
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return err
	}

	os.MkdirAll(bn.getUtreexoDirPath(), os.ModePerm)

	var forest *accumulator.Forest
	if HasAccess(bn.getForestFilePath()) {
		log.Info("Has access to forestdata, resuming")
		forest, err = bn.restoreForest()
	} else {
		log.Info("Creating new forestdata")
		forest, err = bn.createForest()

		bn.info.indexHash = hash.ZeroHash
		bn.info.indexOrder = 0
	}
	if err != nil {
		return err
	}
	bn.forest = forest

	if bn.info.isValid() {
		curIB := bn.bd.GetBlockByOrder(bn.info.indexOrder)
		if !curIB.IsEqual(bn.info.GetHash()) {
			return fmt.Errorf("The Utreexo data was damaged. you can cleanup your data base by '--droputreexo'.")
		}
	}

	//
	mainIB := bn.bd.GetMainChainTip()
	if bn.info.indexOrder < mainIB.GetOrder() || !bn.info.isValid() {
		err = bn.catchUpIndex(mainIB)
		if err != nil {
			return err
		}
	}
	return nil
}

func (bn *UtreexoModule) catchUpIndex(mainIB blockdag.IBlock) error {
	i := bn.info.indexOrder + 1
	if !bn.info.isValid() {
		i = 0
	}
	for ; i <= mainIB.GetOrder(); i++ {
		ib := bn.bd.GetIBlockByOrder(i)
		if ib == nil {
			return fmt.Errorf(fmt.Sprintf("Can't find order(%d)", i))
		}
		block, err := bn.bc.FetchBlockByHash(ib.GetHash())
		if err != nil {
			return err
		}
		txs := map[int]*types.Transaction{}
		err = bn.db.View(func(dbTx database.Tx) error {

			for i, tx := range block.Transactions() {
				if index.DBHasTxIndexEntry(dbTx, tx.Hash()) {
					txs[i] = tx.Tx
				} else if i == 0 {
					return fmt.Errorf("The block is invalid")
				}
			}
			return nil
		})
		if err != nil {
			continue
		}

		err = bn.buildProofs(&addBlockMsg{block, txs, i})
		if err != nil {
			return err
		}
	}

	return nil
}

func (bn *UtreexoModule) restoreForest() (forest *accumulator.Forest, err error) {

	// Where the forestfile exists
	forestFile, err := os.OpenFile(
		bn.getForestFilePath(), os.O_RDWR, 0400)
	if err != nil {
		return nil, err
	}
	// Where the misc forest data exists
	miscForestFile, err := os.OpenFile(
		bn.getMiscForestFilePath(), os.O_RDONLY, 0400)
	if err != nil {
		return nil, err
	}

	forest, err = accumulator.RestoreForest(miscForestFile, forestFile)
	if err != nil {
		return nil, err
	}

	return
}

// createForest initializes forest
func (bn *UtreexoModule) createForest() (forest *accumulator.Forest, err error) {

	// Where the forestfile exists
	forestFile, err := os.OpenFile(
		bn.getForestFilePath(), os.O_CREATE|os.O_RDWR, 0600)
	if err != nil {
		return nil, err
	}

	// Restores all the forest data
	forest = accumulator.NewForest(forestFile)

	return
}

func (bn *UtreexoModule) getForestFilePath() string {
	return filepath.Join(bn.cfg.DataDir, ForestFilePath)
}

func (bn *UtreexoModule) getMiscForestFilePath() string {
	return filepath.Join(bn.cfg.DataDir, MiscForestFilePath)
}

func (bn *UtreexoModule) getUtreexoDirPath() string {
	return filepath.Join(bn.cfg.DataDir, UtreexoDirPath)
}

func (bn *UtreexoModule) Start() {
	if !bn.cfg.Utreexo {
		return
	}
	if atomic.AddInt32(&bn.started, 1) != 1 {
		return
	}
	log.Info("Start utreexo")
	bn.wg.Add(1)
	go bn.handler()
}

func (bn *UtreexoModule) Stop() error {
	if !bn.cfg.Utreexo {
		return nil
	}
	if atomic.AddInt32(&bn.shutdown, 1) != 1 {
		err := fmt.Errorf("Utreexo is already in the process of shutting down")
		log.Error(err.Error())
		return err
	}
	log.Info("End utreexo")

	close(bn.quit)
	bn.wg.Wait()
	return nil
}

func (bn *UtreexoModule) handler() {

out:
	for {
		select {
		case m := <-bn.msgChan:
			switch msg := m.(type) {
			case *addBlockMsg:
				err := bn.buildProofs(msg)
				if err != nil {
					log.Error(err.Error())
				}
			default:
				log.Warn(fmt.Sprintf("Invalid message type: %T", msg))
			}

		case <-bn.quit:
			break out
		}
	}

	// Drain any wait channels before going away so there is nothing left
	// waiting on this goroutine.
cleanup:
	for {
		select {
		case <-bn.msgChan:
		default:
			break cleanup
		}
	}

	bn.wg.Done()
	log.Trace("Utreexo handler done")
}

func (bn *UtreexoModule) BlockConnected(blk *types.SerializedBlock) {
	// Ignore if we are shutting down.
	if atomic.LoadInt32(&bn.shutdown) != 0 {
		return
	}

	txs := map[int]*types.Transaction{}
	err := bn.db.View(func(dbTx database.Tx) error {

		for i, tx := range blk.Transactions() {
			if index.DBHasTxIndexEntry(dbTx, tx.Hash()) {
				txs[i] = tx.Tx
			} else if i == 0 {
				return fmt.Errorf("The block is invalid")
			}
		}
		return nil
	})
	if err != nil {
		log.Error(err.Error())
		return
	}

	bn.msgChan <- &addBlockMsg{blk, txs, uint(blk.Order())}
}

func (bn *UtreexoModule) BlockWillDisconnect(blk *types.SerializedBlock) {
	// Ignore if we are shutting down.
	if atomic.LoadInt32(&bn.shutdown) != 0 {
		return
	}

	txs := map[int]*types.Transaction{}
	err := bn.db.View(func(dbTx database.Tx) error {

		for i, tx := range blk.Transactions() {
			if index.DBHasTxIndexEntry(dbTx, tx.Hash()) {
				txs[i] = tx.Tx
			} else if i == 0 {
				return fmt.Errorf("The block is invalid")
			}
		}
		return nil
	})
	if err != nil {
		log.Error(err.Error())
		return
	}

	bn.msgChan <- &removeBlockMsg{blk, txs}
}

func NewUtreexoModule(bc *blockchain.BlockChain, db database.DB, cfg *config.Config) (*UtreexoModule, error) {
	bn := UtreexoModule{bc: bc,
		bd:      bc.BlockDAG(),
		db:      db,
		info:    &UtreexoInfo{hash.ZeroHash, 0},
		cfg:     cfg,
		msgChan: make(chan interface{}),
		quit:    make(chan struct{})}
	err := bn.Init()
	return &bn, err
}
