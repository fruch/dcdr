package consul

import (
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/watch"
	"github.com/vsco/dcdr/cli/api/stores"
	"github.com/vsco/dcdr/cli/printer"
	"github.com/vsco/dcdr/config"
)

type ConsulWatcher struct {
	config *config.Config
	cb     func(kvb stores.KVBytes)
}

func New(cfg *config.Config) (cw *ConsulWatcher) {
	cw = &ConsulWatcher{
		config: cfg,
	}

	return
}

func (cw *ConsulWatcher) Register(cb func(kvb stores.KVBytes)) {
	cw.cb = cb
}

func (cw *ConsulWatcher) Updated(kvs interface{}) {
	kvp := kvs.(api.KVPairs)
	kvb, err := stores.KvPairsToKvBytes(kvp)

	if err != nil {
		printer.LogErr("%v", err)
		return
	}

	cw.cb(kvb)
}

func (cw *ConsulWatcher) Watch() {
	params := map[string]interface{}{
		"type":   "keyprefix",
		"prefix": cw.config.Namespace,
	}

	wp, err := watch.Parse(params)
	defer wp.Stop()

	if err != nil {
		printer.LogErr("%v", err)
	}

	wp.Handler = func(idx uint64, data interface{}) {
		cw.Updated(data)
	}

	if err := wp.Run(""); err != nil {
		printer.LogErr("Error querying Consul agent: %s", err)
	}
}
