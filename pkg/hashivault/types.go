package hashivault

// SecretsManager represents a service that is able to provide clients with a secrets identified by paths.
type SecretsManager interface {
	// GetSecret returns a function that returns a map of secrets. The point is that the returned function will always
	// return the latest version of the secret. Therefore, clients should save a reference to the function rather than
	// saving the actual secrets, and invoke the func just-in-time as the secret is needed. The returned function is
	// safe to use concurrently.
	GetSecret(path string) (EvergreenSecretsFunc, error)

	// SetDefaultGoogleCredentials fetches the Google credentials from the given path and key and sets them as the
	// default credentials for the current process. This means saving the credentials to disk and setting the
	// environment variable GOOGLE_APPLICATION_CREDENTIALS to point to the saved file.
	SetDefaultGoogleCredentials(path, key string) error
}

// EvergreenSecretsFunc is a function that returns a map of secrets. The point is that the returned function will always
// return the latest version of the secret. Therefore, clients should save a reference to the function rather than
// saving the actual secrets, and invoke the func just-in-time as the secret is needed. The returned function is
// safe to use concurrently.
type EvergreenSecretsFunc func() map[string]any

// secret contains all data and metadata from a Vault secret
type secret struct {
	RequestID     string                            `json:"request_id"`
	LeaseID       string                            `json:"lease_id"`
	Renewable     bool                              `json:"renewable"`
	LeaseDuration int                               `json:"lease_duration"`
	Data          map[string]map[string]interface{} `json:"data"`
}

func (s *secret) GetRequestID() string {
	return s.RequestID
}

func (s *secret) GetLeaseID() string {
	return s.LeaseID
}

func (s *secret) IsRenewable() bool {
	return s.Renewable
}

func (s *secret) GetLeaseDuration() int {
	return s.LeaseDuration
}

func (s *secret) GetData() map[string]interface{} {
	return s.Data["data"]
}

func (s *secret) GetMetadata() map[string]interface{} {
	return s.Data["metadata"]
}
