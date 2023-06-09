package hashivault

// SecretsManager represents a service that is able to provide clients with a secret stored at a path.
type SecretsManager interface {
	GetSecret(path string) (Secret, error)
	SetDefaultGoogleCredentials(path, key string) error
}

type Secret interface {
	GetRequestID() string
	GetLeaseID() string
	IsRenewable() bool
	GetLeaseDuration() int
	GetData() map[string]interface{}
	GetMetadata() map[string]interface{}
}

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
