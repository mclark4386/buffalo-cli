package plugins

// Plugin is the most basic interface a plugin can implement.
type Plugin interface {
	// Name is the name of the plugin.
	// This will also be used for the cli sub-command
	// 	"pop" | "heroku" | "auth" | etc...
	Name() string
}

type Plugins []Plugin

// Len is the number of elements in the collection.
func (plugs Plugins) Len() int {
	return len(plugs)
}

// Less reports whether the element with
// index i should sort before the element with index j.
func (plugs Plugins) Less(i int, j int) bool {
	return plugs[i].Name() < plugs[j].Name()
}

// Swap swaps the elements with indexes i and j.
func (plugs Plugins) Swap(i int, j int) {
	plugs[i], plugs[j] = plugs[j], plugs[i]
}