package hashivault

import (
	"github.com/3lvia/hashivault-go/internal/auth"
	"log"
	"net/http"
	"sync"
)

type tokenGetterFunc func() string

func startTokenJob(c *optionsCollector, errChan chan<- error, initializedChan chan<- struct{}, client *http.Client, l *log.Logger) tokenGetterFunc {
	if c.vaultToken != "" {
		// If the token is already set, just return it
		return func() string {
			return c.vaultToken
		}
	}

	j := &tokenJob{
		mux:          &sync.Mutex{},
		vaultAddress: c.vaultAddress,
		gitHubToken:  c.gitHubToken,
		k8sMountPath: c.k8sMountPath,
		k8sRole:      c.k8sRole,
		client:       client,
		useOICD:      c.useOIDC,
		l:            l,
	}

	go j.start(errChan, initializedChan)
	return j.token
}

type tokenJob struct {
	mux          *sync.Mutex
	vaultAddress string
	gitHubToken  string
	k8sMountPath string
	k8sRole      string
	currentToken string
	useOICD      bool
	client       *http.Client
	l            *log.Logger
}

func (j *tokenJob) start(errChannel chan<- error, initializedChan chan<- struct{}) {
	j.l.Print("starting token job")

	j.mux.Lock()
	authResponse, err := j.authenticate()
	if err != nil {
		close(initializedChan)
		errChannel <- err
		return
	}
	j.currentToken = authResponse.ClientToken()
	j.mux.Unlock()

	// signal that we're done initializing
	close(initializedChan)
	j.l.Print("token job initialized, first token acquired")

	if !authResponse.Renewable() {
		// no need to renew token, so we're done
		return
	}

	after := authResponse.After()
	for {
		<-after
		j.l.Print("renewing token")
		j.mux.Lock()
		ar, err := j.authenticate()
		if err != nil {
			errChannel <- err
			j.mux.Unlock()
			continue
		}
		j.currentToken = ar.ClientToken()
		after = ar.After()
		j.mux.Unlock()
		j.l.Print("token renewed")
	}
}

func (j *tokenJob) token() string {
	j.mux.Lock()
	defer j.mux.Unlock()
	return j.currentToken
}

func (j *tokenJob) authenticate() (auth.AuthenticationResponse, error) {
	if j.useOICD {
		j.l.Print("using OIDC authentication")
		return auth.Authenticate(j.vaultAddress, auth.MethodOICD, auth.WithClient(j.client))
	}
	if j.gitHubToken != "" {
		j.l.Print("using GitHub authentication")
		return auth.Authenticate(j.vaultAddress, auth.MethodGitHub, auth.WithGitHubToken(j.gitHubToken), auth.WithClient(j.client))
	}

	j.l.Print("using Kubernetes authentication")
	return auth.Authenticate(j.vaultAddress, auth.MethodK8s, auth.WithK8s(j.k8sMountPath, j.k8sRole), auth.WithClient(j.client))
}
