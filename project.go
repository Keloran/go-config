package ConfigBuilder

// ProjectProperties is a type for storing project-level configuration properties
type ProjectProperties map[string]interface{}

// GetValue returns the value for the given key, or nil if not found
func (p ProjectProperties) GetValue(key string) interface{} {
	if p == nil {
		return nil
	}
	if value, ok := p[key]; ok {
		return value
	}
	return nil
}

// GetString returns the value for the given key as a string, or empty string if not found or not a string
func (p ProjectProperties) GetString(key string) string {
	if p == nil {
		return ""
	}
	if value, ok := p[key]; ok {
		if s, ok := value.(string); ok {
			return s
		}
	}
	return ""
}

// Has returns true if the key exists in the properties
func (p ProjectProperties) Has(key string) bool {
	if p == nil {
		return false
	}
	_, ok := p[key]
	return ok
}

// Set sets a value for the given key
func (p ProjectProperties) Set(key string, value interface{}) {
	if p != nil {
		p[key] = value
	}
}
