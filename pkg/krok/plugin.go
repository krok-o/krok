package krok

// Plugin defines a Krok plugin and how it should behave and what it should be capable of.
// This interface is returned by the plugin loader.
type Plugin interface {
	Execute(payload string) (string, bool, error)
}
