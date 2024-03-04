//go:build linux

package lvm2

import (
	"encoding/json"
	"errors"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog"

	"github.com/opensvc/om3/util/command"
	"github.com/opensvc/om3/util/device"
	"github.com/opensvc/om3/util/fcache"
	"github.com/opensvc/om3/util/funcopt"
	"github.com/opensvc/om3/util/hostname"
	"github.com/opensvc/om3/util/sizeconv"
)

type (
	VG struct {
		driver
		VGName string
		log    *zerolog.Logger
	}
	VGAttrIndex uint8
	VGAttrs     string
	VGAttr      rune
)

const (
	VGAttrIndexPermissions VGAttrIndex = 0
	VGAttrIndexResizeable  VGAttrIndex = iota
	VGAttrIndexExported
	VGAttrIndexPartial
	VGAttrIndexAllocationPolicy
	VGAttrIndexClusteredOrShared
)

const (
	// State attrs field

	VGAttrStateWriteable       VGAttr = 'w'
	VGAttrStateReadOnly        VGAttr = 'r'
	VGAttrStateResizeable      VGAttr = 'z'
	VGAttrStateExported        VGAttr = 'x'
	VGAttrStatePartial         VGAttr = 'p'
	VGAttrStateAllocContiguous VGAttr = 'c'
	VGAttrStateAllocCling      VGAttr = 'l'
	VGAttrStateAllocNormal     VGAttr = 'n'
	VGAttrStateAllocAnywhere   VGAttr = 'a'
	VGAttrStateClustered       VGAttr = 'c'
	VGAttrStateShared          VGAttr = 's'
)

func NewVG(vg string, opts ...funcopt.O) *VG {
	t := VG{
		VGName: vg,
	}
	_ = funcopt.Apply(&t, opts...)
	return &t
}

func (t VG) FQN() string {
	return t.VGName
}

func (t *VG) Activate() error {
	return t.change([]string{"-ay"})
}

func (t *VG) Deactivate() error {
	return t.change([]string{"-an"})
}

func (t *VG) change(args []string) error {
	cmd := command.New(
		command.WithName("vgchange"),
		command.WithArgs(append(args, t.VGName)),
		command.WithLogger(t.Log()),
		command.WithCommandLogLevel(zerolog.InfoLevel),
		command.WithStdoutLogLevel(zerolog.InfoLevel),
		command.WithStderrLogLevel(zerolog.ErrorLevel),
	)
	cmd.Run()
	fcache.Clear("vgs")
	fcache.Clear("vgs-device")
	if cmd.ExitCode() != 0 {
		return fmt.Errorf("%s error %d", cmd, cmd.ExitCode())
	}
	return nil
}

func (t *VG) AddNodeTag() error {
	return t.AddTag("@" + hostname.Hostname())
}

func (t *VG) DelNodeTag() error {
	return t.DelTag("@" + hostname.Hostname())
}

func (t *VG) DelTag(s string) error {
	cmd := command.New(
		command.WithName("vgchange"),
		command.WithVarArgs("--deltag", s, t.VGName),
		command.WithLogger(t.Log()),
		command.WithCommandLogLevel(zerolog.InfoLevel),
		command.WithStdoutLogLevel(zerolog.InfoLevel),
		command.WithStderrLogLevel(zerolog.ErrorLevel),
	)
	cmd.Run()
	fcache.Clear("vgs")
	fcache.Clear("vgs-device")
	if cmd.ExitCode() != 0 {
		return fmt.Errorf("%s error %d", cmd, cmd.ExitCode())
	}
	return nil
}

func (t *VG) AddTag(s string) error {
	cmd := command.New(
		command.WithName("vgchange"),
		command.WithVarArgs("--addtag", s, t.VGName),
		command.WithLogger(t.Log()),
		command.WithCommandLogLevel(zerolog.InfoLevel),
		command.WithStdoutLogLevel(zerolog.InfoLevel),
		command.WithStderrLogLevel(zerolog.ErrorLevel),
	)
	cmd.Run()
	fcache.Clear("vgs")
	fcache.Clear("vgs-device")
	if cmd.ExitCode() != 0 {
		return fmt.Errorf("%s error %d", cmd, cmd.ExitCode())
	}
	return nil
}

func (t *VG) CachedDevicesShow() (*VGInfo, error) {
	var (
		err error
		out []byte
	)
	data := ShowData{}
	cmd := command.New(
		command.WithName("vgs"),
		command.WithVarArgs("--reportformat", "json", "-o", "devices"),
		command.WithLogger(t.Log()),
		command.WithCommandLogLevel(zerolog.DebugLevel),
		command.WithStdoutLogLevel(zerolog.DebugLevel),
		command.WithStderrLogLevel(zerolog.DebugLevel),
		command.WithBufferedStdout(),
	)
	if out, err = fcache.Output(cmd, "vgs-devices"); err != nil {
		return nil, err
	}
	if err = json.Unmarshal(out, &data); err != nil {
		return nil, err
	}
	if len(data.Report) != 1 {
		return nil, fmt.Errorf("vgs: no report")
	}
	for _, d := range data.Report[0].VG {
		if d.VGName == t.VGName {
			return &d, nil
		}
	}
	return nil, fmt.Errorf("%w: %s", ErrExist, t.VGName)
}

func (t *VG) CachedNormalShow() (*VGInfo, error) {
	var (
		err error
		out []byte
	)
	data := ShowData{}
	cmd := command.New(
		command.WithName("vgs"),
		command.WithVarArgs("--reportformat", "json", "-o", "+tags,pv_name"),
		command.WithLogger(t.Log()),
		command.WithCommandLogLevel(zerolog.DebugLevel),
		command.WithStdoutLogLevel(zerolog.DebugLevel),
		command.WithStderrLogLevel(zerolog.DebugLevel),
		command.WithBufferedStdout(),
	)
	if out, err = fcache.Output(cmd, "vgs"); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(out, &data); err != nil {
		return nil, err
	}
	if len(data.Report) != 1 {
		return nil, fmt.Errorf("vgs: no report")
	}
	for _, d := range data.Report[0].VG {
		if d.VGName == t.VGName {
			return &d, nil
		}
	}
	return nil, fmt.Errorf("%w: %s", ErrExist, t.VGName)
}

func (t *VG) Show(fields string) (*VGInfo, error) {
	data := ShowData{}
	cmd := command.New(
		command.WithName("vgs"),
		command.WithVarArgs("--reportformat", "json", "-o", fields, t.VGName),
		command.WithLogger(t.Log()),
		command.WithCommandLogLevel(zerolog.DebugLevel),
		command.WithStdoutLogLevel(zerolog.DebugLevel),
		command.WithStderrLogLevel(zerolog.DebugLevel),
		command.WithBufferedStdout(),
	)
	if err := cmd.Run(); err != nil {
		if cmd.ExitCode() == 5 {
			return nil, fmt.Errorf("%w: %s", ErrExist, t.VGName)
		}
		return nil, err
	}
	if err := json.Unmarshal(cmd.Stdout(), &data); err != nil {
		return nil, err
	}
	if len(data.Report) == 1 && len(data.Report[0].VG) == 1 {
		return &data.Report[0].VG[0], nil
	}
	return nil, fmt.Errorf("%w: %s", ErrExist, t.VGName)
}

func (t *VG) Attrs() (VGAttrs, error) {
	vgInfo, err := t.CachedNormalShow()
	switch {
	case errors.Is(err, ErrExist):
		return "", nil
	case err != nil:
		return "", err
	default:
		return VGAttrs(vgInfo.VGAttr), nil
	}
}

func (t *VG) Tags() ([]string, error) {
	vgInfo, err := t.CachedNormalShow()
	switch {
	case errors.Is(err, ErrExist):
		return []string{}, nil
	case err != nil:
		return []string{}, err
	default:
		return strings.Split(vgInfo.VGTags, ","), nil
	}
}

func (t *VG) HasTag(s string) (bool, error) {
	tags, err := t.Tags()
	if err != nil {
		return false, err
	}
	for _, tag := range tags {
		if tag == s {
			return true, nil
		}
	}
	return false, nil
}

func (t *VG) HasNodeTag() (bool, error) {
	return t.HasTag(hostname.Hostname())
}

func (t VGAttrs) Attr(index VGAttrIndex) VGAttr {
	if len(t) < int(index)+1 {
		return ' '
	}
	return VGAttr(t[index])
}

func (t *VG) Exists() (bool, error) {
	_, err := t.CachedNormalShow()
	switch {
	case errors.Is(err, ErrExist):
		return false, nil
	case err != nil:
		return false, err
	default:
		return true, nil
	}
}

func (t *VG) IsActive() (bool, error) {
	/*
		if attrs, err := t.Attrs(); err != nil {
			return false, err
		} else {
			return attrs.Attr(VGAttrIndexState) == VGAttrStateActive, nil
		}
	*/
	return false, nil
}

func (t *VG) Devices() (device.L, error) {
	l := make(device.L, 0)
	data, err := t.CachedDevicesShow()
	if err != nil {
		return nil, err
	}
	for _, s := range strings.Fields(data.Devices) {
		path := strings.Split(s, "(")[0]
		dev := device.New(path, device.WithLogger(t.Log()))
		l = append(l, dev)
	}
	return l, nil
}

func (t *VG) Create(size string, pvs []string, options []string) error {
	if i, err := sizeconv.FromSize(size); err == nil {
		// default unit is not "B", explicitely tell
		size = fmt.Sprintf("%dB", i)
	}
	args := make([]string, 0)
	args = append(args, t.VGName)
	args = append(args, pvs...)
	args = append(args, options...)
	args = append(args, "--yes")
	cmd := command.New(
		command.WithName("vgcreate"),
		command.WithArgs(args),
		command.WithLogger(t.Log()),
		command.WithCommandLogLevel(zerolog.InfoLevel),
		command.WithStdoutLogLevel(zerolog.InfoLevel),
		command.WithStderrLogLevel(zerolog.ErrorLevel),
	)
	cmd.Run()
	fcache.Clear("vgs")
	fcache.Clear("vgs-device")
	if cmd.ExitCode() != 0 {
		return fmt.Errorf("%s error %d", cmd, cmd.ExitCode())
	}
	return nil
}

func (t *VG) Wipe() error {
	return nil
}

func (t *VG) Remove(args []string) error {
	cmd := command.New(
		command.WithName("vgremove"),
		command.WithArgs(append(args, t.VGName)),
		command.WithLogger(t.Log()),
		command.WithCommandLogLevel(zerolog.InfoLevel),
		command.WithStdoutLogLevel(zerolog.InfoLevel),
		command.WithStderrLogLevel(zerolog.ErrorLevel),
	)
	cmd.Run()
	fcache.Clear("vgs")
	fcache.Clear("vgs-device")
	if cmd.ExitCode() != 0 {
		return fmt.Errorf("%s error %d", cmd, cmd.ExitCode())
	}
	return nil
}

func (t *VG) PVs() (device.L, error) {
	l := make(device.L, 0)
	vgInfo, err := t.CachedNormalShow()
	switch {
	case errors.Is(err, ErrExist):
		return l, nil
	case err != nil:
		return l, err
	}
	for _, s := range strings.Split(vgInfo.PVName, ",") {
		l = append(l, device.New(s, device.WithLogger(t.Log())))
	}
	return l, nil
}

func (t *VG) ActiveLVs() (device.L, error) {
	l := make(device.L, 0)
	pattern := fmt.Sprintf("/dev/mapper/%s-*", t.VGName)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return l, err
	}
	for _, p := range matches {
		switch {
		case strings.Contains(p, "_rimage_"), strings.Contains(p, "_rmeta_"):
			continue
		case strings.Contains(p, "_mimage_"), strings.Contains(p, "_mlog_"), strings.HasSuffix(p, "_mlog"):
			continue
		}
		l = append(l, device.New(p, device.WithLogger(t.Log())))
	}
	return l, nil
}
