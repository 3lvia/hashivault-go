package main

import (
	"github.com/3lvia/hashivault-go/pkg/hashivault"
	"log"
)

func main() {
	v, errChan, err := hashivault.New()
	if err != nil {
		log.Fatal(err)
	}

	go func(ec <-chan error) {
		for err := range ec {
			log.Println(err)
		}
	}(errChan)

	secret, err := v.GetSecret("kunde/kv/data/appinsights/kunde")
	if err != nil {
		log.Fatal(err)
	}

	mapOfSecrets := secret()
	_ = mapOfSecrets
}
