package hashivault

import "github.com/3lvia/hashivault-go/internal/auth"

func New() (SecretsManager, error) {
	v := &vault{}

	ghToken := ""
	addr := "https://vault.dev-elvia.io"

	token, err := auth.Authenticate(addr, auth.WithGitHubToken(ghToken))
	if err != nil {
		return nil, err
	}
	_ = token

	return v, nil
}
