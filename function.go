package function

import "github.com/GoogleCloudPlatform/functions-framework-go/functions"

// init is the composition root for the Cloud Function.
// It constructs production Dependencies and Application, then registers the HTTP function.
func init() {
	deps := Dependencies{
		EventStorer: NewFakeEventStorer(),
		EventGetter: NewFakeEventGetter(nil),
	}
	app := NewApplication(deps)

	functions.HTTP("Ingest", app.Ingest)
}

