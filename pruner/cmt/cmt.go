package cmt

import (
	"fmt"
	dbm "github.com/cometbft/cometbft-db"
	cmtstate "github.com/cometbft/cometbft/state"
	cmtstore "github.com/cometbft/cometbft/store"
	"log"
	"os"
	"path/filepath"
)

func Prune(blockStoreDB, stateDB dbm.DB, blocksToKeep uint64) error {
	blockStore := cmtstore.NewBlockStore(blockStoreDB)
	stateStore := cmtstate.NewStore(stateDB, cmtstate.StoreOptions{})
	base := blockStore.Base()
	pruneHeight := blockStore.Height() - int64(blocksToKeep)

	state, err := stateStore.Load()
	if err != nil {
		return fmt.Errorf("prune state error: %v", err)
	}

	log.Println("pruning block store")

	blocksPruned, evidenceHeight, err := blockStore.PruneBlocks(pruneHeight, state)
	if err != nil {
		return fmt.Errorf("prune blocks %w", err)
	}

	log.Printf("%d block pruned", blocksPruned)

	log.Println("pruning state store")
	err = stateStore.PruneStates(base, pruneHeight, evidenceHeight)
	if err != nil {
		return fmt.Errorf("prune states %w", err)
	}

	return nil
}

func PruneTxIndex(dataDir string) error {
	txIndexPath := filepath.Join(dataDir, "tx_index.db")

	err := os.RemoveAll(txIndexPath)
	if err != nil {
		return fmt.Errorf("Failed to remove directory: %v\n", err)
	}

	log.Println("tx_index.db removed")
	return nil
}
