package object

import (
	"context"
	"fmt"
	"sort"
	"time"

	"github.com/ssrathi/go-attr"

	"github.com/opensvc/om3/core/actionresdeps"
	"github.com/opensvc/om3/core/driver"
	"github.com/opensvc/om3/core/keywords"
	"github.com/opensvc/om3/core/manifest"
	"github.com/opensvc/om3/core/resource"
	"github.com/opensvc/om3/core/resourceid"
	"github.com/opensvc/om3/core/resourceset"
	"github.com/opensvc/om3/core/schedule"
	"github.com/opensvc/om3/core/xconfig"
	"github.com/opensvc/om3/util/funcopt"
	"github.com/opensvc/om3/util/key"
	"github.com/opensvc/om3/util/pg"
)

type (
	actor struct {
		core
		pg                 *pg.Config
		resources          resource.Drivers
		_resources         resource.Drivers
		actionResourceDeps *actionresdeps.Store
	}

	// freezer is implemented by object kinds supporting freeze and thaw.
	freezer interface {
		Freeze(ctx context.Context) error
		Unfreeze(ctx context.Context) error
		Frozen() time.Time
	}

	// resourceLister provides a method to list and filter resources
	resourceLister interface {
		Resources() resource.Drivers
		IsDesc() bool
	}

	// Actor is implemented by object kinds supporting start, stop, ...
	Actor interface {
		Core
		resourceLister
		freezer
		PG() *pg.Config
		IsVolatile() bool
		ResourceSets() resourceset.L
		ResourceByID(rid string) resource.Driver
		GetActionResDeps() *actionresdeps.Store
		ConfigureResources()
		IsDisabled() bool

		Boot(context.Context) error
		Restart(context.Context) error
		Run(context.Context) error
		Shutdown(context.Context) error
		Start(context.Context) error
		StartStandby(context.Context) error
		Stop(context.Context) error
		PRStart(context.Context) error
		PRStop(context.Context) error
		Provision(context.Context) error
		Unprovision(context.Context) error
		SetProvisioned(context.Context) error
		SetUnprovisioned(context.Context) error
		SyncFull(context.Context) error
		SyncResync(context.Context) error
		SyncUpdate(context.Context) error
		Enter(context.Context, string) error

		PrintSchedule() schedule.Table
		PushResInfo(context.Context) (resource.Infos, error)
	}
)

func (t *actor) PG() *pg.Config {
	return t.pg
}

func (t *actor) init(referrer xconfig.Referrer, p any, opts ...funcopt.O) error {
	if err := t.core.init(referrer, p, opts...); err != nil {
		return err
	}
	t.pg = t.pgConfig("")
	t.actionResourceDeps = actionresdeps.NewStore()
	t.actionResourceDeps.SetActionMap(map[string]string{
		"provision":   "start",
		"shutdown":    "stop",
		"unprovision": "stop",
		"toc":         "stop",
	})
	return nil
}

func (t *actor) ResourceSets() resourceset.L {
	l := resourceset.NewList()
	done := make(map[string]*resourceset.T)
	//
	// subsetSectionString returns the existing section name found in the
	// config file
	//   [subset#fs:g1]   (most precise)
	//   [subset#g1]      (less precise)
	//
	subsetSectionString := func(g driver.Group, name string) string {
		s := resourceset.FormatSectionName(g.String(), name)
		if t.config.HasSectionString(s) {
			return s
		}
		return "subset#" + s
	}
	//
	// configureResourceSet allocates and configures the resourceset, looking
	// for keywords in the following sections:
	//   [subset#fs:g1]   (most precise)
	//   [subset#g1]      (less precise)
	//
	// If the rset is already configured, avoid doing the job twice.
	//
	configureResourceSet := func(g driver.Group, name string) *resourceset.T {
		id := resourceset.FormatSectionName(g.String(), name)
		if rset, ok := done[id]; ok {
			return rset
		}
		k := subsetSectionString(g, name)
		rset := resourceset.New()
		rset.DriverGroup = g
		rset.Name = name
		rset.SectionName = k
		rset.ResourceLister = t
		parallelKey := key.New(k, "parallel")
		rset.Parallel = t.config.GetBool(parallelKey)
		rset.PG = t.pgConfig(k)
		rset.SetLogger(t.log)
		done[id] = rset
		l = append(l, rset)
		return rset
	}

	for _, k := range t.config.SectionStrings() {
		//
		// look for resource sections with a defined subset
		//   [fs#1]
		//   subset = g1
		//
		rid, err := resourceid.Parse(k)
		if err != nil {
			continue
		}
		subsetKey := key.New(k, "subset")
		subsetName := t.config.Get(subsetKey)
		if subsetName == "" {
			// discard section with no 'subset' keyword
			continue
		}
		configureResourceSet(rid.DriverGroup(), subsetName)
	}

	// add generic resourcesets not already found as a section
	for _, k := range driver.GroupNames() {
		if _, ok := done[k]; ok {
			continue
		}
		if rset, err := resourceset.Generic(k); err == nil {
			rset.ResourceLister = t
			rset.SetLogger(t.log)
			l = append(l, rset)
		} else {
			t.log.Debugf("%s", err)
		}
	}
	sort.Sort(l)
	return l
}

func (t *actor) getConfiguringResourceByID(rid string) resource.Driver {
	for _, r := range t._resources {
		if r.RID() == rid {
			return r
		}
	}
	return nil
}

func (t *actor) getConfiguredResourceByID(rid string) resource.Driver {
	for _, r := range t.resources {
		if r.RID() == rid {
			return r
		}
	}
	return nil
}

func (t *actor) ResourceByID(rid string) resource.Driver {
	if r := t.getConfiguredResourceByID(rid); r != nil {
		return r
	}
	return t.getConfiguringResourceByID(rid)
}

func listResources(i interface{}) resource.Drivers {
	if lister, ok := i.(resourceLister); ok {
		return lister.Resources()
	}
	return resource.Drivers{}
}

func listResourceSets(i interface{}) resourceset.L {
	if actor, ok := i.(Actor); ok {
		return actor.ResourceSets()
	}
	return resourceset.L{}
}

func (t *actor) ResourcesByDrivergroups(drvgrps []driver.Group) resource.Drivers {
	return resourcesByDrivergroups(t, drvgrps)
}

func resourcesByDrivergroups(i interface{}, drvgrps []driver.Group) resource.Drivers {
	l := make([]resource.Driver, 0)
	for _, r := range listResources(i) {
		drvgrp := r.ID().DriverGroup()
		for _, d := range drvgrps {
			if drvgrp == d {
				l = append(l, r)
				break
			}
		}
	}
	return l
}

func (t *actor) Resources() resource.Drivers {
	if t.resources != nil {
		return t.resources
	}
	t.ConfigureResources()
	return t.resources
}

func (t *actor) ConfigureResources() {
	t.Lock()
	defer t.Unlock()
	begin := time.Now()
	postponed := make(map[string][]resource.Driver)
	t._resources = make(resource.Drivers, 0)
	for _, k := range t.config.SectionStrings() {
		rid, err := resourceid.Parse(k)
		if err != nil {
			continue
		}
		driverGroup := rid.DriverGroup()
		if driverGroup == driver.GroupUnknown {
			t.log.Attr("rid", k).Attr("f", "listResources").Debugf("%s: unknown driver group in rid %s", t.path, k)
			continue
		}
		typeKey := key.New(k, "type")
		driverName := t.config.Get(typeKey)
		driverID := driver.NewID(driverGroup, driverName)
		factory := resource.NewResourceFunc(driverID)
		if factory == nil {
			t.log.Attr("driver", driverID.String()).Debugf("%s: unknown driver %s", t.path, driverID)
			continue
		}
		r := factory()
		rBegin := time.Now()
		if err := t.configureResource(r, k); err != nil {
			switch o := err.(type) {
			case xconfig.ErrPostponedRef:
				if _, ok := postponed[o.RID]; !ok {
					postponed[o.RID] = make([]resource.Driver, 0)
				}
				postponed[o.RID] = append(postponed[o.RID], r)
			default:
				t.log.Attr("rid", k).Errorf("%s: configure resource %s: %s", t.path, k, err)
			}
			continue
		}
		dur := time.Now().Sub(rBegin)
		t.log.Attr("rid", k).Attr("duration", dur).Debugf("%s: resource %s configured in %s", t.path, k, dur)
		t._resources = append(t._resources, r)
	}
	for _, resources := range postponed {
		for _, r := range resources {
			rBegin := time.Now()
			if err := t.ReconfigureResource(r); err != nil {
				t.log.Attr("rid", r.RID()).Errorf("%s: configure postponed resource %s: %s", t.path, r.RID(), err)
				continue
			}
			dur := time.Now().Sub(rBegin)
			t.log.Attr("rid", r.RID()).Attr("duration", dur).Debugf("%s: postponed resource %s configured in %s", t.path, r.RID(), dur)
			t._resources = append(t._resources, r)
		}
	}
	t.resources = t._resources
	t._resources = nil
	dur := time.Now().Sub(begin)
	t.log.Attr("duration", dur).Debugf("%s: all resources configured in %s", t.path, dur)
	return
}

func (t *actor) ReconfigureResource(r resource.Driver) error {
	return t.configureResource(r, r.RID())
}

func (t *actor) configureResource(r resource.Driver, rid string) error {
	r.SetRID(rid)
	m := r.Manifest()
	getDNS := func() ([]string, error) {
		n, err := t.Node()
		if err != nil {
			return nil, err
		}
		return n.Nameservers()
	}
	getCNIConfig := func() (string, error) {
		n, err := t.Node()
		if err != nil {
			return "", err
		}
		return n.CNIConfig()
	}
	getCNIPlugins := func() (string, error) {
		n, err := t.Node()
		if err != nil {
			return "", err
		}
		return n.CNIPlugins()
	}

	setAttr := func(c manifest.Context) error {
		switch {
		case c.Ref == "object.path":
			if err := attr.SetValue(r, c.Attr, t.path); err != nil {
				return err
			}
		case c.Ref == "object.drpnodes":
			if nodes, err := t.DRPNodes(); err != nil {
				return err
			} else if err := attr.SetValue(r, c.Attr, nodes); err != nil {
				return err
			}
		case c.Ref == "object.nodes":
			if nodes, err := t.Nodes(); err != nil {
				return err
			} else if err := attr.SetValue(r, c.Attr, nodes); err != nil {
				return err
			}
		case c.Ref == "object.peers":
			if nodes, err := t.Peers(); err != nil {
				return err
			} else if err := attr.SetValue(r, c.Attr, nodes); err != nil {
				return err
			}
		case c.Ref == "object.id":
			if err := attr.SetValue(r, c.Attr, t.ID()); err != nil {
				return err
			}
		case c.Ref == "object.topology":
			if err := attr.SetValue(r, c.Attr, t.Topology()); err != nil {
				return err
			}
		case c.Ref == "node.dns":
			if dns, err := getDNS(); err != nil {
				return err
			} else if err := attr.SetValue(r, c.Attr, dns); err != nil {
				return err
			}
		case c.Ref == "cni.config":
			if s, err := getCNIConfig(); err != nil {
				return err
			} else if err := attr.SetValue(r, c.Attr, s); err != nil {
				return err
			}
		case c.Ref == "cni.plugins":
			if s, err := getCNIPlugins(); err != nil {
				return err
			} else if err := attr.SetValue(r, c.Attr, s); err != nil {
				return err
			}
		}
		return nil
	}
	for _, attr := range m.Attrs {
		switch o := attr.(type) {
		case keywords.Keyword:
			k := key.New(rid, o.Option)
			val, err := t.config.EvalKeywordAs(k, o, "")
			if err != nil {
				if o.Required {
					return err
				}
				r.Log().Debugf("%s: %s keyword eval: %s", t.path, k, err)
				continue
			}
			if err := o.SetValue(r, val); err != nil {
				return fmt.Errorf("%s.%s: %w", rid, o.Option, err)
			}
		case manifest.Context:
			if err := setAttr(o); err != nil {
				return fmt.Errorf("%s: %w", o.Attr, err)
			}
		}
	}
	r.SetObject(t)
	r.SetPG(t.pgConfig(rid))
	if i, ok := r.(resource.Configurer); ok {
		if err := i.Configure(); err != nil {
			return err
		}
	}
	if i, ok := r.(resource.ActionResourceDepser); ok {
		deps := i.ActionResourceDeps()
		t.actionResourceDeps.RegisterSlice(deps)
	}
	//r.Log().Debug().Msgf("%s: configured resource: %+v", t.path, r)
	return nil
}

func (t *actor) GetActionResDeps() *actionresdeps.Store {
	return t.actionResourceDeps
}

// IsDesc is a requirement of the ResourceLister interface. actor Resources() is always ascending.
func (t *actor) IsDesc() bool {
	return false
}

// IsDisabled returns true if the object config evaluates DEFAULT.disable as true.
// CRM actions are skipped on a disabled instance.
func (t *actor) IsDisabled() bool {
	k := key.Parse("disable")
	return t.config.GetBool(k)
}
