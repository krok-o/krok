package krok

// Execute defines a Krok plugin and how it should behave and what it should be capable of.
// This interface is returned by the plugin loader.
type Execute = func(payload string, opts ...interface{}) (string, bool, error)
