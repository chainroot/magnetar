package cmd

import (
	"fmt"
	"github.com/chainroot/magnetar/pruner/cmt"
	"github.com/chainroot/magnetar/pruner/state"
	dbm "github.com/cometbft/cometbft-db"
	db "github.com/cosmos/cosmos-db"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
	"path/filepath"
)

func pruneCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "prune [path_to_home]",
		Short: "prune data from the application store and block store",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {

			fmt.Println("cometbft", cometBft)
			fmt.Println("cosmossdk", cosmosSdk)
			fmt.Println("blocks", keepBlocks)
			fmt.Println("versions", keepVersions)
			fmt.Println("txIndex", txIndex)
			fmt.Println("compact", compact)

			dbDir := rootify(args[0], homePath)
			fmt.Println("dbDir", dbDir)

			var group errgroup.Group

			blockStoreDb, err := dbm.NewGoLevelDB("blockstore", dbDir)
			if err != nil {
				return fmt.Errorf("opening block store: %w", err)
			}
			defer blockStoreDb.Close()

			stateDb, err := dbm.NewGoLevelDB("state", dbDir)
			if err != nil {
				return fmt.Errorf("opening state: %w", err)
			}
			defer stateDb.Close()

			applicationDb, err := db.NewGoLevelDB("application", dbDir, nil)
			if err != nil {
				return fmt.Errorf("opening state: %w", err)
			}
			defer applicationDb.Close()

			group.Go(func() error {
				if cometBft {
					return cmt.Prune(blockStoreDb, stateDb, keepBlocks)
				}
				return nil
			})

			group.Go(func() error {
				if cosmosSdk {
					return state.Prune(applicationDb, keepVersions)
				}
				return nil
			})

			group.Go(func() error {
				if txIndex {
					return cmt.PruneTxIndex(dataDir)
				}
				return nil
			})

			if err := group.Wait(); err != nil {
				return err
			}

			if compact {
				if err := blockStoreDb.Compact(nil, nil); err != nil {
					return err
				}

				if err := stateDb.Compact(nil, nil); err != nil {
					return err
				}

				if err := applicationDb.ForceCompact(nil, nil); err != nil {
					return err
				}
			}

			return nil
		},
	}
	return cmd
}

func rootify(path, root string) string {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(root, path)
}
