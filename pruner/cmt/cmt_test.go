package cmt

import (
	dbm "github.com/cometbft/cometbft-db"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"log"
	"testing"
)

func TestPrune(t *testing.T) {
	var (
		dataDir          = "/Users/aminnami/Desktop/NodeData/teritori/data"
		blockStoreDBName = "blockstore"
		stateDBName      = "state"
		opts             = opt.Options{
			DisableSeeksCompaction: true,
		}
		blocksToKeep uint64 = 1
	)

	blockStoreDB, err := dbm.NewGoLevelDBWithOpts(blockStoreDBName, dataDir, &opts)
	if err != nil {
		t.Fatal(err)
	}

	stateDB, err := dbm.NewGoLevelDBWithOpts(stateDBName, dataDir, &opts)
	if err != nil {
		t.Fatal(err)
	}

	if err := Prune(blockStoreDB, stateDB, blocksToKeep); err != nil {
		t.Fatal(err)
	}

	log.Println("prune ok")
}

func TestPruneTxIndex(t *testing.T) {
	var (
		dataDir = "/Users/aminnami/Desktop/NodeData/teritori/data"
	)

	if err := PruneTxIndex(dataDir); err != nil {
		t.Fatal(err)
	}

	log.Println("prune ok")
}
