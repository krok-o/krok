package krok

// Plugin defines a Krok plugin and how it should behave and what it should be capable of.
// This interface is returned by the plugin loader.
type Plugin interface {
	// Execute defines the ability to run a command on a hook.
	// Opts defines a variable number of arguments that can be passed to the command.
	// These could come from the Vault, things like, auth information or environment properties, etc.
	Execute(payload string, opts ...interface{}) (string, bool, error)
}
