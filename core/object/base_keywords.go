package object

import (
	"opensvc.com/opensvc/core/keywords"
	"opensvc.com/opensvc/util/converters"
	"opensvc.com/opensvc/util/key"
)

var keywordStore = keywords.Store{
	{
		Section:     "DEFAULT",
		Option:      "nodes",
		Converter:   converters.ListLowercase,
		Text:        "A node selector expression specifying the list of cluster nodes hosting service instances.",
		DefaultText: "The lowercased hostname of the evaluating node.",
		Example:     "n1 n*",
	},
	{
		Section:   "DEFAULT",
		Option:    "drpnodes",
		Converter: converters.ListLowercase,
		Text:      "The backup node where the service is activated in a DRP situation. This node is also a data synchronization target for :c-res:`sync` resources.",
		Example:   "n1 n2",
	},
	{
		Section:   "DEFAULT",
		Option:    "encapnodes",
		Converter: converters.ListLowercase,
		Text:      "The list of `containers` handled by this service and with an OpenSVC agent installed to handle the encapsulated resources. With this parameter set, parameters can be scoped with the ``@encapnodes`` suffix.",
		Example:   "n1 n2",
	},
	{
		Section: "DEFAULT",
		Option:  "app",
		Default: "default",
		Text:    "Used to identify who is responsible for this service, who is billable and provides a most useful filtering key. Better keep it a short code.",
	},
	{
		Section:    "DEFAULT",
		Option:     "topology",
		Default:    "failover",
		Candidates: []string{"failover", "flex"},
		Text:       "``failover`` the service is allowed to be up on one node at a time. ``flex`` the service can be up on :kw:`flex_target` nodes, where :kw:`flex_target` must be in the [flex_min, flex_max] range.",
	},
}

func (t Base) KeywordLookup(k key.T) keywords.Keyword {
	switch k.Section {
	case "data", "env":
		return keywords.Keyword{
			Option:   "*", // trick IsZero()
			Scopable: true,
			Required: false,
		}
	}
	return keywordStore.Lookup(k)
}
