package models

// VaultSetting defines a setting that comes from the vault
type VaultSetting struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
