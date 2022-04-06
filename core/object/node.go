package object

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"opensvc.com/opensvc/core/xconfig"
	"opensvc.com/opensvc/util/funcopt"
)

type (
	// Node is the node struct.
	Node struct {
		//private
		log      zerolog.Logger
		volatile bool

		// caches
		id           uuid.UUID
		configFile   string
		config       *xconfig.T
		mergedConfig *xconfig.T
		paths        NodePaths
	}
)

// NewNode allocates a node.
func NewNode(opts ...funcopt.O) *Node {
	t := &Node{}
	t.init(opts...)
	return t
}

func (t *Node) init(opts ...funcopt.O) error {
	if err := funcopt.Apply(t, opts...); err != nil {
		return err
	}

	// log.Logger is configured in cmd/root.go
	t.log = log.Logger

	/*
			t.log = logging.Configure(logging.Config{
				ConsoleLoggingEnabled: true,
				EncodeLogsAsJSON:      true,
				FileLoggingEnabled:    true,
				Directory:             t.LogDir(),
				Filename:              "node.log",
				MaxSize:               5,
				MaxBackups:            1,
				MaxAge:                30,
				WithCaller:            logging.WithCaller,
			}).
				With().
				Str("n", hostname.Hostname()).
				Str("sid", xsession.ID).
				Logger()
		}
	*/

	if err := t.loadConfig(); err != nil {
		return err
	}
	return nil
}

func (t Node) String() string {
	return fmt.Sprintf("node")
}

func (t Node) IsVolatile() bool {
	return t.volatile
}

func (t Node) SetStandardConfigFile() {
	return
}
