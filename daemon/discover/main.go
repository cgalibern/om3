// Package discover implements object discovery for daemon
//
// It watches config filesystem to create initial instance config worker when
// config file is created. Instance config worker is then responsible for
// watching instance config updates
//
// When is discovers that another remote config is available and no instance
// config worker is running, it fetches remote instance config to local config
// directory.
//
// It is responsible for initial aggregated worker creation.
package discover

import (
	"context"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/rs/zerolog"

	"opensvc.com/opensvc/daemon/daemonlogctx"
	"opensvc.com/opensvc/daemon/msgbus"
	"opensvc.com/opensvc/util/hostname"
)

type (
	discover struct {
		cfgCmdC    chan *msgbus.Msg
		svcaggCmdC chan *msgbus.Msg
		ctx        context.Context
		log        zerolog.Logger

		// cfgMTime is a map of local instance config file time, indexed by object
		// path string representation.
		// More recent remote config files are fetched.
		cfgMTime map[string]time.Time

		svcAggCancel map[string]context.CancelFunc
		svcAgg       map[string]map[string]struct{}

		remoteNodeCtx        map[string]context.Context
		remoteNodeCancel     map[string]context.CancelFunc
		remoteCfgFetchCancel map[string]context.CancelFunc

		// fetcherUpdated map[svc] updated timestamp of svc config being fetched
		fetcherUpdated map[string]time.Time

		// fetcherFrom map[svc] node
		fetcherFrom map[string]string

		// fetcherCancel map[svc] cancel func for svc fetcher
		fetcherCancel map[string]context.CancelFunc

		// fetcherNodeCancel map[node]map[svc] cancel func for node
		fetcherNodeCancel map[string]map[string]context.CancelFunc

		localhost string
		fsWatcher *fsnotify.Watcher
	}
)

var (
	dropCmdTimeout = 100 * time.Millisecond
)

// Start function starts file system watcher on config directory
// then listen for config file creation to create
func Start(ctx context.Context) (func(), error) {
	var wg sync.WaitGroup
	ctx, cancel := context.WithCancel(ctx)
	d := discover{
		cfgCmdC:    make(chan *msgbus.Msg),
		svcaggCmdC: make(chan *msgbus.Msg),
		cfgMTime:   make(map[string]time.Time),
		ctx:        ctx,
		log:        daemonlogctx.Logger(ctx).With().Str("name", "daemon.discover").Logger(),

		svcAgg: make(map[string]map[string]struct{}),

		fetcherFrom:       make(map[string]string),
		fetcherCancel:     make(map[string]context.CancelFunc),
		fetcherNodeCancel: make(map[string]map[string]context.CancelFunc),
		fetcherUpdated:    make(map[string]time.Time),
		localhost:         hostname.Hostname(),
	}
	wg.Add(2)
	go func() {
		defer wg.Done()
		d.cfg()
	}()
	go func() {
		defer wg.Done()
		d.agg()
	}()

	stopFSWatcher, err := d.fsWatcherStart()
	if err != nil {
		d.log.Error().Err(err).Msg("start")
		return stopFSWatcher, err
	}

	cancelAndWait := func() {
		stopFSWatcher()
		cancel() // stop cfg and agg via context cancel
		wg.Wait()
	}
	return cancelAndWait, nil
}