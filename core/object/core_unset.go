package object

import (
	"opensvc.com/opensvc/core/actioncontext"
	"opensvc.com/opensvc/core/xconfig"
	"opensvc.com/opensvc/util/key"
)

// OptsUnset is the options of the Unset object method.
type OptsUnset struct {
	OptsLock
	Keywords []string `flag:"kws"`
}

// Unset gets a keyword value
func (t *core) Unset(options OptsUnset) error {
	props := actioncontext.Unset
	unlock, err := t.lockAction(props, options.OptsLock)
	if err != nil {
		return err
	}
	defer unlock()
	return unset(t.config, options)
}

func (t *core) UnsetKeys(kws ...key.T) error {
	return unsetKeys(t.config, kws...)
}

func unset(cf *xconfig.T, options OptsUnset) error {
	kws := make([]key.T, 0)
	for _, kw := range options.Keywords {
		kws = append(kws, key.Parse(kw))
	}
	return unsetKeys(cf, kws...)
}

func unsetKeys(cf *xconfig.T, kws ...key.T) error {
	changes := 0
	for _, k := range kws {
		changes += cf.Unset(k)
	}
	if changes > 0 {
		return cf.Commit()
	}
	return nil
}