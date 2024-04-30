package auth

// Authenticator is a struct that holds a list of API keys
type Authenticator struct {
	APIKeys []string
}

// NewAPIAuth is a constructor function that returns an Authenticator
func NewAPIAuth(keys []string) *Authenticator {
	return &Authenticator{
		APIKeys: keys,
	}
}

// Authenticate is a method that checks if a given key is in the list of API keys
func (a *Authenticator) Authenticate(key string) bool {
	for _, k := range a.APIKeys {
		if k == key {
			return true
		}
	}
	return false
}

// Set is a method that sets the API keys
func (a *Authenticator) Set(keys []string) {
	a.APIKeys = keys
}
