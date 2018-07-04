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

	"github.com/ethereum/go-ethereum/consensus/ethash"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/core/vm"
	"github.com/ethereum/go-ethereum/params"
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

func visitBlock(db ethdb.Database) {
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

	receipts := rawdb.ReadReceipts(db, currentBlock.Hash(), currentBlock.NumberU64())
	fmt.Printf("%v receipts\n", len(receipts))
	for i, rcpt := range receipts {
		fmt.Printf("%v: %v\n", i, rcpt.TxHash.Hex())
	}
	body := rawdb.ReadBody(db, currentBlock.Hash(), currentBlock.NumberU64())
	fmt.Printf("%v transactions\n", len(body.Transactions))
	for i, tr := range body.Transactions {
		trjson, _ := tr.MarshalJSON()
		fmt.Printf("%v: %v\n", i, string(trjson))
	}

}

/* TODO: full chain verification */
func verifyChain(db ethdb.Database) {
	engine := ethash.NewFaker()
	chain, err := core.NewBlockChain(db, nil, params.TestChainConfig, engine, vm.Config{})
	if err != nil {
		fmt.Printf("failed to create new chain manager: %v", err)
	}
	fmt.Println("chain=", chain)
	/*
		var results <-chan error
		_, results = engine.VerifyHeaders(chain, []*types.Header{headers[i]}, []bool{true})

			// Wait for the verification result
				select {
				case result := <-results:
					if (result == nil) != valid {
						t.Errorf("test %d.%d: validity mismatch: have %v, want %v", i, j, result, valid)
					}
				case <-time.After(time.Second):
					t.Fatalf("test %d.%d: verification timeout", i, j)
				}
				// Make sure no more data is returned
				select {
				case result := <-results:
					t.Fatalf("test %d.%d: unexpected result returned: %v", i, j, result)
				case <-time.After(25 * time.Millisecond):
				}
			}
			chain.InsertChain(blocks[i : i+1])
		}*/
}

func main() {
	db := opendb()
	visitBlock(db)
	//verifyChain(db)
}
