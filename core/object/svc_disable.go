package object

import (
	"context"

	"github.com/opensvc/om3/core/actioncontext"
	"github.com/opensvc/om3/core/keyop"
	"github.com/opensvc/om3/core/resourceselector"
	"github.com/opensvc/om3/core/statusbus"
	"github.com/opensvc/om3/util/key"
)

// Disable unsets disable=true from the svc config
func (t *svc) Disable(ctx context.Context) error {
	ctx = actioncontext.WithProps(ctx, actioncontext.Disable)
	ctx, stop := statusbus.WithContext(ctx, t.path)
	defer stop()
	defer t.postActionStatusEval(ctx)
	unlock, err := t.lockAction(ctx)
	if err != nil {
		return err
	}
	defer unlock()
	var kops keyop.L
	rs := resourceselector.FromContext(ctx, t)
	for _, r := range rs.Resources() {
		kops = append(kops, keyop.T{
			Key:   key.T{r.RID(), "disable"},
			Op:    keyop.Set,
			Value: "true",
		})
	}
	kops = append(kops, keyop.T{
		Key:   key.T{"DEFAULT", "disable"},
		Op:    keyop.Set,
		Value: "true",
	})
	return t.config.Set(kops...)
}
