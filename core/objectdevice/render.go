package objectdevice

import (
	"opensvc.com/opensvc/core/rawconfig"
	"opensvc.com/opensvc/util/render/tree"
)

func (t L) Render() string {
	tree := tree.New()
	tree.AddColumn().AddText("Object").SetColor(rawconfig.Node.Color.Bold)
	tree.AddColumn().AddText("Resource").SetColor(rawconfig.Node.Color.Bold)
	tree.AddColumn().AddText("Driver").SetColor(rawconfig.Node.Color.Bold)
	tree.AddColumn().AddText("Role").SetColor(rawconfig.Node.Color.Bold)
	tree.AddColumn().AddText("Device").SetColor(rawconfig.Node.Color.Bold)
	for _, e := range t {
		n := tree.AddNode()
		n.AddColumn().AddText(e.ObjectPath.String()).SetColor(rawconfig.Node.Color.Primary)
		n.AddColumn().AddText(e.RID).SetColor(rawconfig.Node.Color.Primary)
		n.AddColumn().AddText(e.DriverGroup.String() + "." + e.DriverName).SetColor(rawconfig.Node.Color.Secondary)
		n.AddColumn().AddText(e.Role.String()).SetColor(rawconfig.Node.Color.Secondary)
		n.AddColumn().AddText(e.Device.String())
	}
	return tree.Render()
}