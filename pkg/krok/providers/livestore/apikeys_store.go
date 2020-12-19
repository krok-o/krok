package livestore

// ApiKeysStore is a postgres based store for ApiKeys.
type ApiKeysStore struct {
	ApiKeysDependencies
	Config
}

// ApiKeysDependencies ApiKeys specific dependencies.
type ApiKeysDependencies struct {
	Dependencies
	Connector *Connector
}

// NewApiKeysStore creates a new ApiKeysStore
func NewApiKeysStore(cfg Config, deps ApiKeysDependencies) *ApiKeysStore {
	return &ApiKeysStore{Config: cfg, ApiKeysDependencies: deps}
}
