package function

import (
	"context"
	"log"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

// init is the composition root for the Cloud Function.
// It constructs production Dependencies and Application, then registers the HTTP function.
func init() {
	ctx := context.Background()

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		log.Fatalf("failed to create storage client: %v", err)
	}

	deps := Dependencies{
		EventStorer: NewGCSEventStorer(storageClient, "drawexact-telemetry"),
		EventGetter: NewFakeEventGetter(nil),
	}
	app := NewApplication(deps)

	functions.HTTP("Ingest", app.Ingest)
}

