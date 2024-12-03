package state

import (
	db "github.com/cosmos/cosmos-db"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"log"
	"testing"
)

func TestPrune(t *testing.T) {
	var (
		dataDir           = "/Users/aminnami/Desktop/NodeData/teritori/data"
		applicationDBName = "application"
		opts              = opt.Options{
			DisableSeeksCompaction: true,
		}
		//versionsToKeep uint64 = 1
	)

	applicationDB, err := db.NewGoLevelDBWithOpts(applicationDBName, dataDir, &opts)
	if err != nil {
		t.Fatal(err)
	}

	if err := applicationDB.ForceCompact(nil, nil); err != nil {
		return
	}
	//if err := Prune(applicationDB, versionsToKeep); err != nil {
	//	t.Fatal(err)
	//}

	log.Println("prune ok")
}
