package state

import (
	cosmoslog "cosmossdk.io/log"
	"cosmossdk.io/store/metrics"
	"cosmossdk.io/store/types"
	"fmt"
	"github.com/chainroot/magnetar/internal/rootmulti"
	cosmosdb "github.com/cosmos/cosmos-db"
	"log"
	"os"
)

func Prune(applicationDB cosmosdb.DB, versionsToKeep uint64) error {
	fmt.Println("pruning application state")

	appStore := rootmulti.NewStore(applicationDB, cosmoslog.NewLogger(os.Stderr), metrics.NewNoOpMetrics())
	ver := rootmulti.GetLatestVersion(applicationDB)

	var storeNames []string
	if ver != 0 {
		cInfo, err := appStore.GetCommitInfo(ver)
		if err != nil {
			return err
		}

		for _, storeInfo := range cInfo.StoreInfos {
			// we only want to prune the stores with actual data.
			// sometimes in-memory stores get leaked to disk without data.
			// if that happens, the store's computed hash is empty as well.
			if len(storeInfo.CommitId.Hash) > 0 {
				storeNames = append(storeNames, storeInfo.Name)
			} else {
				log.Println("skipping", storeInfo.Name, "store due to empty hash")
			}
		}
	}

	keys := types.NewKVStoreKeys(storeNames...)
	for _, value := range keys {
		appStore.MountStoreWithDB(value, types.StoreTypeIAVL, nil)
	}

	if err := appStore.LoadLatestVersion(); err != nil {
		return fmt.Errorf("load latest version: %w", err)
	}

	versions := appStore.GetAllVersions()

	v64 := make([]int64, len(versions))
	for i := 0; i < len(versions); i++ {
		v64[i] = int64(versions[i])
	}

	if err := appStore.PruneStores(int64(len(v64)) - int64(versionsToKeep)); err != nil {
		return fmt.Errorf("prune stores: %w", err)
	}

	return nil
}
