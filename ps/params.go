package ps

// Parameter is the parameter exported by psadm.
// This should be sufficient for import and export.
type Parameter struct {
	Description string `yaml:"description"`
	KMSKeyID    string `yaml:"kmskeyid"`
	Decrypted   bool   `yaml:"decrypted,omitempty"`
	Name        string `yaml:"name"`
	Type        string `yaml:"type"`
	Value       string `yaml:"value"`
}
