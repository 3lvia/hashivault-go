package main

import (
	"context"
	"github.com/3lvia/hashivault-go/pkg/hashivault"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/sdk/trace"
	"log"
	"os"
)

func main() {
	ctx := context.Background()

	// Use stout logger, replace with more sophisticated logger in production.
	l := log.New(os.Stdout, "", log.LstdFlags)

	// Use stdout exporter, replace with more sophisticated exporter in production.
	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		log.Fatal(err)
	}

	tp := trace.NewTracerProvider(
		trace.WithSyncer(exporter),
		trace.WithSampler(trace.AlwaysSample()),
	)
	otel.SetTracerProvider(tp)

	v, errChan, err := hashivault.New(
		ctx,
		hashivault.WithOIDC(),
		hashivault.WithVaultAddress("https://vault.dev-elvia.io"),
		hashivault.WithLogger(l),
	)
	if err != nil {
		log.Fatal(err)
	}

	go func(ec <-chan error) {
		for err := range ec {
			log.Println(err)
		}
	}(errChan)

	secret, err := v.GetSecret(ctx, "kunde/kv/data/appinsights/kunde")
	if err != nil {
		log.Fatal(err)
	}

	mapOfSecrets := secret()
	_ = mapOfSecrets
}
