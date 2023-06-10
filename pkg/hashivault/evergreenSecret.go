package hashivault

import (
	"net/http"
	"sync"
	"time"
)

func newEvergreen(path, vaultAddress string, sec *secret, client *http.Client, errChan chan<- error) *evergreenSecret {
	eg := &evergreenSecret{
		path:         path,
		sec:          sec,
		mux:          &sync.Mutex{},
		client:       client,
		vaultAddress: vaultAddress,
	}

	go eg.start(errChan)

	return eg
}

type evergreenSecret struct {
	path         string
	vaultAddress string
	client       *http.Client
	sec          *secret
	mux          *sync.Mutex
}

func (e *evergreenSecret) get() map[string]any {
	e.mux.Lock()
	defer e.mux.Unlock()
	return e.sec.GetData()
}

func (e *evergreenSecret) start(errChan chan<- error) {
	for {
		<-time.After(time.Duration(e.sec.LeaseDuration) * time.Second)
		e.mux.Lock()
		sec, err := get(e.path, e.vaultAddress, "", e.client)
		if err != nil {
			errChan <- err
			e.mux.Unlock()
			continue
		}
		e.sec = sec
		e.mux.Unlock()
	}
}
