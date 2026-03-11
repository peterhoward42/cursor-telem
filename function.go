package function

import (
	"context"
	"fmt"

	"cloud.google.com/go/storage"
	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

// Register is the composition root for the Cloud Function.
// It constructs production Dependencies and Application, then registers the HTTP function.
// Call it explicitly from the process entrypoint (e.g. main); do not rely on init().
func Register() error {
	ctx := context.Background()

	storageClient, err := storage.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("create storage client: %w", err)
	}

	deps := Dependencies{
		EventStorer: NewGCSEventStorer(storageClient, "drawexact-telemetry"),
		EventGetter: NewFakeEventGetter(nil),
		Analyser:    NewAnalyser(),
	}
	app := NewApplication(deps)

	functions.HTTP("Ingest", app.Ingest)
	return nil
}

