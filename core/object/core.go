package object

import (
	"context"
	"fmt"
	"os"
	"sync"

	"github.com/google/uuid"

	"github.com/opensvc/om3/core/env"
	"github.com/opensvc/om3/core/instance"
	"github.com/opensvc/om3/core/naming"
	"github.com/opensvc/om3/core/rawconfig"
	"github.com/opensvc/om3/core/xconfig"
	"github.com/opensvc/om3/util/compliance"
	"github.com/opensvc/om3/util/funcopt"
	"github.com/opensvc/om3/util/plog"
	"github.com/opensvc/om3/util/progress"
)

type (
	// core is the base struct embedded in all kinded objects.
	core struct {
		sync.Mutex

		path naming.Path

		// private
		volatile bool
		log      *plog.Logger

		// caches
		id         uuid.UUID
		configFile string
		configData any
		config     *xconfig.T
		node       *Node
		paths      struct {
			varDir string
			logDir string
			tmpDir string
		}

		// method plugs
		postCommit func() error
	}

	compliancer interface {
		NewCompliance() (*compliance.T, error)
	}

	volatiler interface {
		IsVolatile() bool
		SetVolatile(v bool)
	}

	// Core is implemented by all object kinds.
	Core interface {
		Configurer
		compliancer
		volatiler
		Path() naming.Path
		FQDN() string
		Status(context.Context) (instance.Status, error)
		FreshStatus(context.Context) (instance.Status, error)
		Nodes() ([]string, error)
	}
)

func (t *core) PostCommit() error {
	if t.postCommit == nil {
		return nil
	}
	return t.postCommit()
}

func (t *core) SetPostCommit(fn func() error) {
	t.postCommit = fn
}

// List returns the stringified path as data
func (t *core) List() (string, error) {
	return t.path.String(), nil
}

func (t *core) init(referrer xconfig.Referrer, id any, opts ...funcopt.O) error {
	if parsed, err := toPathType(id); err != nil {
		return err
	} else {
		t.path = parsed
	}
	if err := funcopt.Apply(t, opts...); err != nil {
		return err
	}
	t.log = plog.NewDefaultLogger().
		Attr("obj_path", t.path.String()).
		Attr("obj_kind", t.path.Kind.String()).
		Attr("obj_name", t.path.Name).
		Attr("obj_namespace", t.path.Namespace).
		WithPrefix(fmt.Sprintf("instance: %s: ", t.path))
	if v := os.Getenv(env.ActionOrchestrationIdVar); v != "" {
		t.log = t.log.Attr("ORCHESTRATION_ID", v)
	}
	if err := t.loadConfig(referrer); err != nil {
		return err
	}
	return nil
}

func (t core) String() string {
	return t.path.String()
}

func (t *core) SetVolatile(v bool) {
	t.volatile = v
}

func (t core) IsVolatile() bool {
	return t.volatile
}

func (t *core) Path() naming.Path {
	return t.path
}

// ConfigFile returns the absolute path of an opensvc object configuration
// file.
func (t core) ConfigFile() string {
	if t.configFile == "" {
		t.configFile = t.path.ConfigFile()
	}
	return t.configFile
}

// Node returns a cache Node struct pointer. If none is already cached,
// allocate a new Node{} and cache it.
func (t *core) Node() (*Node, error) {
	if t.node != nil {
		return t.node, nil
	}
	if n, err := NewNode(); err != nil {
		return nil, err
	} else {
		t.node = n
		return t.node, nil
	}
}

func (t core) Log() *plog.Logger {
	return t.log
}

func (t core) ProgressKey() []string {
	p := rawconfig.Colorize.Bold(t.path.String())
	return []string{p}
}

func (t core) Progress(ctx context.Context, cols ...any) {
	if view := progress.ViewFromContext(ctx); view != nil {
		key := t.ProgressKey()
		view.Info(key, cols)
	}
}
