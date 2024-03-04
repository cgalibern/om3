package oxcmd

import (
	"time"
)

type (
	// OptsGlobal contains options accepted by all actions
	OptsGlobal struct {
		Color          string
		Output         string
		ObjectSelector string
		Server         string
	}

	// OptsLogs contains options used by all log commands:
	// node logs, cluster logs, object logs
	OptsLogs struct {
		Follow bool
		Lines  int
		Filter []string
	}

	// OptsAsync contains options accepted by all actions having an orchestration
	OptsAsync struct {
		Watch bool
		Wait  bool
		Time  time.Duration
	}

	// OptsResourceSelector contains options needed to initialize a
	// resourceselector.Options struct
	OptsResourceSelector struct {
		RID    string
		Subset string
		Tag    string
	}

	// OptsLock contains options accepted by all actions using an action lock
	OptsLock struct {
		Disable bool
		Timeout time.Duration
	}

	// OptTo sets a barrier when iterating over a resource lister
	OptTo struct {
		To     string
		UpTo   string // Deprecated
		DownTo string // Deprecated
	}
)
