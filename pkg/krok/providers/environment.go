package providers

// EnvironmentConverter provides an option to parse the environment.
// This is needed in case we are running in a docker swarm
// where secrets come from a mounted file instead of an environment
// variable.
type EnvironmentConverter interface {
	// LoadValueFromFile provides the ability to load a secret from a docker
	// mounted secret file if the value contains `/run/secret`.
	LoadValueFromFile(f string) (string, error)
}
