package common

import (
	"time"

	"github.com/opslevel/opslevel-go/v2023"
)

// SyncCache Performs a one-time sync of the opslevel-go caches
func SyncCache(client *opslevel.Client) {
	opslevel.Cache.CacheTiers(client)
	opslevel.Cache.CacheLifecycles(client)
	opslevel.Cache.CacheTeams(client)
}

// SyncCaches Runs a goroutine that will periodically sync the opslevel-go caches
func SyncCaches(client *opslevel.Client, resync time.Duration) {
	ticker := time.NewTicker(resync)
	go func() {
		for {
			<-ticker.C
			// has a mutex lock that will block TryGet in ReconcileService goroutine
			SyncCache(client)
		}
	}()
}
