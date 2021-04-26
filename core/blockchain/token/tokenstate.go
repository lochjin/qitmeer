/*
 * Copyright (c) 2017-2020 The qitmeer developers
 */

package token

import (
	"fmt"
	"github.com/Qitmeer/qitmeer/core/dbnamespace"
	"github.com/Qitmeer/qitmeer/core/json"
	"github.com/Qitmeer/qitmeer/core/serialization"
	"github.com/Qitmeer/qitmeer/core/types"
	"github.com/Qitmeer/qitmeer/database"
)

// tokenState specifies the token balance of the current block.
// the updates are written in the same order as the tx in the block, which is
// used to verify the correctness of the token balance
type TokenState struct {
	PrevStateID uint32
	Types       TokenTypesMap
	Balances    TokenBalancesMap
	Updates     []BalanceUpdate
}

func (ts *TokenState) GetTokenBalances() []json.TokenBalance {
	tbs := []json.TokenBalance{}
	for k, v := range ts.Balances {
		tb := json.TokenBalance{CoinId: uint16(k), CoinName: k.Name(), Balance: v.Balance, LockedMeer: v.LockedMeer}
		tbs = append(tbs, tb)
	}
	return tbs
}

// serializeTokeState function will serialize the token state into byte slice
func serializeTokeState(ts TokenState) ([]byte, error) {
	// total number of bytes to serialize
	serializeSize := serialization.SerializeSizeVLQ(uint64(ts.PrevStateID))

	serializeSize += serialization.SerializeSizeVLQ(uint64(len(ts.Balances)))
	for id, b := range ts.Balances {
		// sanity check
		if id == types.MEERID || b.Balance < 0 || b.LockedMeer < 0 {
			return nil, fmt.Errorf("invalid token balance {%v, %v}", id, b)
		}
		serializeSize += serialization.SerializeSizeVLQ(uint64(id))
		serializeSize += serialization.SerializeSizeVLQ(uint64(b.Balance))
		serializeSize += serialization.SerializeSizeVLQ(uint64(b.LockedMeer))
	}
	serializeSize += serialization.SerializeSizeVLQ(uint64(len(ts.Updates)))
	for _, v := range ts.Updates {
		if v.Typ != TokenMint && v.Typ != TokenUnMint {
			return nil, fmt.Errorf("invalid token balance update type %v", v.Typ)
		}
		if v.MeerAmount < 0 || v.TokenAmount.Value < 0 || !types.IsKnownCoinID(v.TokenAmount.Id) {
			return nil, fmt.Errorf("invalid token balance update %v", v)
		}
		serializeSize += 1 // balanceUpdateType takes 1 byte
		serializeSize += serialization.SerializeSizeVLQ(uint64(v.MeerAmount))
		serializeSize += serialization.SerializeSizeVLQ(uint64(v.TokenAmount.Id))
		serializeSize += serialization.SerializeSizeVLQ(uint64(v.TokenAmount.Value))
	}
	serialized := make([]byte, serializeSize)
	offset := 0
	offset = serialization.PutVLQ(serialized, uint64(ts.PrevStateID))

	offset += serialization.PutVLQ(serialized[offset:], uint64(len(ts.Balances)))
	for id, b := range ts.Balances {
		offset += serialization.PutVLQ(serialized[offset:], uint64(id))
		offset += serialization.PutVLQ(serialized[offset:], uint64(b.Balance))
		offset += serialization.PutVLQ(serialized[offset:], uint64(b.LockedMeer))
	}

	offset += serialization.PutVLQ(serialized[offset:], uint64(len(ts.Updates)))
	for _, v := range ts.Updates {
		offset += serialization.PutVLQ(serialized[offset:], uint64(v.Typ))
		offset += serialization.PutVLQ(serialized[offset:], uint64(v.MeerAmount))
		offset += serialization.PutVLQ(serialized[offset:], uint64(v.TokenAmount.Id))
		offset += serialization.PutVLQ(serialized[offset:], uint64(v.TokenAmount.Value))
	}
	return serialized, nil
}

// deserializeTokenState function will deserializes token state from the byte slice
func deserializeTokenState(data []byte) (*TokenState, error) {
	prevStateID, offset := serialization.DeserializeVLQ(data)
	if offset == 0 {
		return nil, fmt.Errorf("unexpected end of data while reading prevStateID")
	}
	// Deserialize the balance.
	var balances map[types.CoinID]TokenBalance
	numOfBalances, bytesRead := serialization.DeserializeVLQ(data[offset:])
	if bytesRead == 0 {
		return nil, fmt.Errorf("unexpected end of data while reading number of balances")
	}
	offset += bytesRead

	if numOfBalances > 0 {
		balances = make(map[types.CoinID]TokenBalance, numOfBalances)
		for i := uint64(0); i < numOfBalances; i++ {
			// token id
			derId, bytesRead := serialization.DeserializeVLQ(data[offset:])
			if bytesRead == 0 {
				return nil, fmt.Errorf("unexpected end of data while reading token id at balances{%d}", i)
			}
			offset += bytesRead

			// token balance
			balance, bytesRead := serialization.DeserializeVLQ(data[offset:])
			if bytesRead == 0 {
				return nil, fmt.Errorf("unexpected end of data while reading balance at balances{%d}", i)
			}
			offset += bytesRead

			// locked meer
			lockedMeer, bytesRead := serialization.DeserializeVLQ(data[offset:])
			if bytesRead == 0 {
				return nil, fmt.Errorf("unexpected end of data while reading balance at balances{%d}", i)
			}
			offset += bytesRead

			id := types.CoinID(uint16(derId))
			balances[id] = TokenBalance{int64(balance), int64(lockedMeer)}
		}
	}
	updates := []BalanceUpdate{}
	numOfUpdates, bytesRead := serialization.DeserializeVLQ(data[offset:])
	if bytesRead == 0 {
		return nil, fmt.Errorf("unexpected end of data while reading number of balances")
	}
	offset += bytesRead
	if numOfUpdates > 0 {
		updates = make([]BalanceUpdate, numOfUpdates)
		for i := uint64(0); i < numOfUpdates; i++ {
			//type
			updateType, bytesRead := serialization.DeserializeVLQ(data[offset:])
			if bytesRead == 0 {
				return nil, fmt.Errorf("unexpected end of data while reading balance update type at update[%d]", i)
			}
			offset += bytesRead
			//meerAmount
			meerAmount, bytesRead := serialization.DeserializeVLQ(data[offset:])
			if bytesRead == 0 {
				return nil, fmt.Errorf("unexpected end of data while reading meer amount at update[%d]", i)
			}
			offset += bytesRead
			//tokenId
			tokenId, bytesRead := serialization.DeserializeVLQ(data[offset:])
			if bytesRead == 0 {
				return nil, fmt.Errorf("unexpected end of data while reading token id at update[%d]", i)
			}
			offset += bytesRead
			//tokenAmount
			tokenAmount, bytesRead := serialization.DeserializeVLQ(data[offset:])
			if bytesRead == 0 {
				return nil, fmt.Errorf("unexpected end of data while reading token amount at update[%d]", i)
			}
			offset += bytesRead

			updates[i] = BalanceUpdate{
				Typ:         BalanceUpdateType(updateType),
				MeerAmount:  int64(meerAmount),
				TokenAmount: types.Amount{Value: int64(tokenAmount), Id: types.CoinID(uint16(tokenId))},
			}
		}
	}
	return &TokenState{PrevStateID: uint32(prevStateID), Balances: balances, Updates: updates}, nil
}

// dbPutTokenState put a token balance record into the token state database.
// the key is the provided block hash
func DBPutTokenState(dbTx database.Tx, bid uint32, ts TokenState) error {
	// Serialize the current token state.
	serializedData, err := serializeTokeState(ts)
	if err != nil {
		return err
	}
	// Store the current token balance record into the token state database.
	meta := dbTx.Metadata()
	bucket := meta.Bucket(dbnamespace.TokenBucketName)
	var serializedID [4]byte
	dbnamespace.ByteOrder.PutUint32(serializedID[:], bid)
	return bucket.Put(serializedID[:], serializedData)
}

// dbFetchTokenState fetch the token balance record from the token state database.
// the key is the input block hash.
func DBFetchTokenState(dbTx database.Tx, bid uint32) (*TokenState, error) {
	// if it is genesis hash, return empty tokenState directly
	if bid == 0 {
		return &TokenState{}, nil
	}
	// Fetch record from the token state database by block hash
	meta := dbTx.Metadata()
	bucket := meta.Bucket(dbnamespace.TokenBucketName)

	var serializedID [4]byte
	dbnamespace.ByteOrder.PutUint32(serializedID[:], bid)
	v := bucket.Get(serializedID[:])
	if v == nil {
		return nil, fmt.Errorf("tokenstate db can't find record from block id : %v", bid)
	}
	// deserialize the fetched token state record
	return deserializeTokenState(v)
}

func DBRemoveTokenState(dbTx database.Tx, id uint32) error {
	bucket := dbTx.Metadata().Bucket(dbnamespace.TokenBucketName)
	var serializedID [4]byte
	dbnamespace.ByteOrder.PutUint32(serializedID[:], id)

	key := serializedID[:]
	return bucket.Delete(key)
}

func CheckUnMintUpdate(update *BalanceUpdate) error {
	if update.Typ != TokenUnMint {
		return fmt.Errorf("checkUnMintUpdate : wrong update type %v", update.Typ)
	}
	if err := checkUpdateCommon(update); err != nil {
		return err
	}
	return nil
}

func CheckMintUpdate(update *BalanceUpdate) error {
	if update.Typ != TokenMint {
		return fmt.Errorf("checkUnMintUpdate : wrong update type %v", update.Typ)
	}
	if err := checkUpdateCommon(update); err != nil {
		return err
	}
	return nil
}

func checkUpdateCommon(update *BalanceUpdate) error {
	if !types.IsKnownCoinID(update.TokenAmount.Id) {
		return fmt.Errorf("checkUpdateCommon : unknown token id %v", update.TokenAmount.Id.Name())
	}
	if update.TokenAmount.Value <= 0 {
		return fmt.Errorf("checkUpdateCommon : wrong token amount : %v", update.TokenAmount.Value)
	}
	if update.MeerAmount <= 0 {
		return fmt.Errorf("checkUpdateCommon : wrong meer amount : %v", update.MeerAmount)
	}
	return nil
}