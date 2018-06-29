package main

import (
	_ "bytes"
	// "encoding/hex"
	"fmt"
	"log"
	// "github.com/davecgh/go-spew/spew"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/rawdb"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethdb"
	// "github.com/ethereum/go-ethereum/rlp"
	"os"
	// "time"
)

func opendb() ethdb.Database {
	netdn := "testnet"
	dirname := fmt.Sprintf("%s/.ethereum/%s/geth/chaindata", os.Getenv("HOME"), netdn)
	db, err := ethdb.NewLDBDatabase(dirname, 0, 0)
	if err != nil {
		fmt.Printf("ldb %s open failed\n", dirname)
		return nil
	}
	// defer db.Close()
	return db
}

func GetBlock(db ethdb.Database, hash common.Hash, number uint64) *types.Block {
	// Short circuit if the block's already in the cache, retrieve otherwise
	block := rawdb.ReadBlock(db, hash, number)
	if block == nil {
		return nil
	}
	return block
}

func GetBlockByNumber(db ethdb.Database, number uint64) *types.Block {
	hash := rawdb.ReadCanonicalHash(db, number)
	if hash == (common.Hash{}) {
		return nil
	}
	return GetBlock(db, hash, number)
}

func GetBlockByHash(db ethdb.Database, hash common.Hash) *types.Block {
	number := rawdb.ReadHeaderNumber(db, hash)
	return GetBlock(db, hash, *number)
}

func loadTrieState(db ethdb.Database) {
	head := rawdb.ReadHeadBlockHash(db)
	fmt.Println("headblock hash=", head)
	if head == (common.Hash{}) {
		// Corrupt or empty database, init from scratch
		log.Println("Empty database, resetting chain")
	}
	// Make sure the entire head block is available
	currentBlock := GetBlockByHash(db, head)
	if currentBlock == nil { // Corrupt or empty database, init from scratch
		log.Println("Head block missing, resetting chain", "hash", head)
		return
	}
	fmt.Println("currentBlock=", currentBlock)
}

func main() {
	db := opendb()
	loadTrieState(db)
}
