package models

// VaultSetting defines a setting that comes from the vault
// swagger:model
type VaultSetting struct {
	// Key is the name of the setting.
	//
	// required: true
	Key string `json:"key"`
	// Value is the value of the setting.
	//
	// required: true
	Value string `json:"value"`
}
