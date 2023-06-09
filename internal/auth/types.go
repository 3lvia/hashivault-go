package auth

import "time"

// gitToken holds github authentication information to be formatted to a bytes buffer
type gitToken struct {
	Token string `json:"token"`
}

// k8sToken holds kubernetes authentication information to be formatted to a bytes buffer
type k8sToken struct {
	JWT  string `json:"jwt"`
	Role string `json:"role"`
}

type AuthenticationResponse interface {
	ClientToken() string
	LeaseDurationSeconds() int
	Renewable() bool
	After() <-chan time.Time
}

type authenticationResponse struct {
	RequestID     string             `json:"request_id"`
	LeaseID       string             `json:"lease_id"`
	Renew         bool               `json:"renewable"`
	LeaseDuration int                `json:"lease_duration"`
	Data          interface{}        `json:"data"`
	WrapInfo      interface{}        `json:"wrap_info"`
	Warnings      interface{}        `json:"warnings"`
	Auth          authenticationData `json:"auth"`
}

func (a authenticationResponse) ClientToken() string {
	return a.Auth.ClientToken
}

func (a authenticationResponse) LeaseDurationSeconds() int {
	return a.Auth.LeaseDuration
}

func (a authenticationResponse) Renewable() bool {
	return a.Auth.Renewable
}

func (a authenticationResponse) After() <-chan time.Time {
	return time.After(time.Duration(a.Auth.LeaseDuration) * time.Second)
}

type authenticationData struct {
	ClientToken    string                 `json:"client_token"`
	Accessor       string                 `json:"accessor"`
	Policies       []string               `json:"policies"`
	TokenPolicies  []string               `json:"token_policies"`
	Metadata       map[string]interface{} `json:"metadata"`
	LeaseDuration  int                    `json:"lease_duration"`
	Renewable      bool                   `json:"renewable"`
	EntityID       string                 `json:"entity_id"`
	TokenType      string                 `json:"token_type"`
	Orphan         bool                   `json:"orphan"`
	MFARequirement interface{}            `json:"mfa_requirement"`
	NumUses        int                    `json:"num_uses"`
}
