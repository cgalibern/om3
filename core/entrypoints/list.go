package entrypoints

import (
	"sort"

	"opensvc.com/opensvc/core/objectselector"
	"opensvc.com/opensvc/core/output"
	"opensvc.com/opensvc/core/rawconfig"
)

// List is the struct exposing the object selection printing method.
type List struct {
	ObjectSelector string
	Color          string
	Format         string
	Server         string
	Local          bool
}

// Do prints the formatted object selection
func (t List) Do() error {
	selection := objectselector.NewSelection(
		t.ObjectSelector,
		objectselector.SelectionWithLocal(t.Local),
		objectselector.SelectionWithServer(t.Server),
	)
	data := make([]string, 0)
	paths, err := selection.Expand()
	if err != nil {
		return err
	}
	for _, path := range paths {
		data = append(data, path.String())
	}
	sort.Strings(data)
	human := func() string {
		s := ""
		for _, r := range data {
			s += r + "\n"
		}
		return s
	}
	output.Renderer{
		Format:        t.Format,
		Color:         t.Color,
		Data:          data,
		HumanRenderer: human,
		Colorize:      rawconfig.Colorize,
	}.Print()
	return nil
}
