package object

import (
	"path/filepath"
	"time"

	"github.com/opensvc/om3/core/resource"
	"github.com/opensvc/om3/core/schedule"
	"github.com/opensvc/om3/util/file"
	"github.com/opensvc/om3/util/hostname"
	"github.com/opensvc/om3/util/key"
)

// PrintSchedule display the object scheduling table
func (t *actor) PrintSchedule() schedule.Table {
	return t.Schedules()
}

func (t *actor) lastRunFile(action, rid, base string) string {
	base = "last_" + base
	if rid != "" {
		base = base + "_" + rid
	}
	return filepath.Join(t.VarDir(), "scheduler", base)
}

func (t *actor) lastSuccessFile(action, rid, base string) string {
	return filepath.Join(t.lastRunFile(action, rid, base) + ".success")
}

func (t *actor) loadLast(action, rid, base string) time.Time {
	fpath := t.lastRunFile(action, rid, base)
	return file.ModTime(fpath)
}

func (t *actor) newScheduleEntry(action, keyStr, rid, base string, reqCol, reqProv bool) schedule.Entry {
	k := key.Parse(keyStr)
	def, err := t.config.GetStringStrict(k)
	if err != nil {
		panic(err)
	}
	return schedule.Entry{
		Node:               hostname.Hostname(),
		Path:               t.path,
		Action:             action,
		LastRunAt:          t.loadLast(action, rid, base),
		Key:                k.String(),
		Schedule:           def,
		LastRunFile:        t.lastRunFile(action, rid, base),
		LastSuccessFile:    t.lastSuccessFile(action, rid, base),
		RequireCollector:   reqCol,
		RequireProvisioned: reqProv,
	}
}

func (t *actor) Schedules() schedule.Table {
	table := schedule.NewTable(
		t.newScheduleEntry("status", "status_schedule", "", "status", false, false),
		t.newScheduleEntry("compliance_auto", "comp_schedule", "", "comp_check", true, true),
	)
	needResMon := false
	type scheduleOptioner interface {
		ScheduleOptions() resource.ScheduleOptions
	}
	for _, r := range listResources(t) {
		if !needResMon && r.IsMonitored() {
			needResMon = true
		}
		if i, ok := r.(scheduleOptioner); ok {
			opts := i.ScheduleOptions()
			rid := r.RID()
			e := t.newScheduleEntry(opts.Action, key.T{Section: rid, Option: opts.Option}.String(), rid, opts.Base, opts.RequireCollector, opts.RequireProvisioned)
			table = table.Add(e)
		}
	}
	if needResMon {
		e := t.newScheduleEntry("resource_monitor", "monitor_schedule", "", "resource_monitor", false, true)
		table = table.Add(e)
	}
	if len(listResources(t)) > 0 {
		e := t.newScheduleEntry("push_resinfo", "resinfo_schedule", "", "push_resinfo", true, false)
		table = table.Add(e)
	}
	return table
}
