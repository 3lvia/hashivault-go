package hashivault

import (
	"github.com/3lvia/hashivault-go/internal/auth"
	"net/http"
	"sync"
)

type tokenGetterFunc func() string

func startTokenJob(vaultAddress, vaultToken, gitHubToken, k8sMountPath, k8sRole string, errChan chan<- error, client *http.Client) tokenGetterFunc {
	if vaultToken != "" {
		return func() string {
			return vaultToken
		}
	}

	j := &tokenJob{
		mux:          &sync.Mutex{},
		vaultAddress: vaultAddress,
		gitHubToken:  gitHubToken,
		k8sMountPath: k8sMountPath,
		k8sRole:      k8sRole,
		client:       client,
	}
	go j.start(errChan)
	return j.token
}

type tokenJob struct {
	mux          *sync.Mutex
	vaultAddress string
	gitHubToken  string
	k8sMountPath string
	k8sRole      string
	currentToken string
	client       *http.Client
}

func (j *tokenJob) start(errChannel chan<- error) {
	j.mux.Lock()
	authResponse, err := j.authenticate()
	if err != nil {
		errChannel <- err
		return
	}
	j.currentToken = authResponse.ClientToken()
	j.mux.Unlock()

	if !authResponse.Renewable() {
		// no need to renew token, so we're done
		return
	}

	after := authResponse.After()
	for {
		<-after
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
	}
}

func (j *tokenJob) token() string {
	j.mux.Lock()
	defer j.mux.Unlock()
	return j.currentToken
}

func (j *tokenJob) authenticate() (auth.AuthenticationResponse, error) {
	if j.gitHubToken != "" {
		return auth.Authenticate(j.vaultAddress, auth.MethodGitHub, auth.WithGitHubToken(j.gitHubToken), auth.WithClient(j.client))
	}

	return auth.Authenticate(j.vaultAddress, auth.MethodK8s, auth.WithK8s(j.k8sMountPath, j.k8sRole), auth.WithClient(j.client))
}
