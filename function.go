package function

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/go-playground/validator/v10"
)

// Dependencies aggregates all external systems required by the application.
// Start empty and add fields only when needed.
type Dependencies struct {
}

// Application owns all request handling logic.
type Application struct {
	deps     Dependencies
	validate *validator.Validate
}

// NewApplication constructs an Application with its explicit dependencies.
func NewApplication(deps Dependencies) *Application {
	return &Application{
		deps:     deps,
		validate: validator.New(),
	}
}

// EventPayload is the validated JSON request payload.
type EventPayload struct {
	SchemaVersion int    `json:"schemaVersion" validate:"max=1,min=0"`
	EventULID     string `json:"eventULID" validate:"len=26"`
	ProxyUserID   string `json:"proxyUserID" validate:"uuid4"`
	TimeUTC       string `json:"timeUTC" validate:"datetime=2006-01-02T15:04:05Z07:00"`
	Visit         int    `json:"visit" validate:"max=100000,min=1"`
	Event         string `json:"event" validate:"max=40,min=4"`
	Parameters    string `json:"parameters" validate:"max=80"`
}

// Ingest handles POST requests containing a single EventPayload JSON body.
func (a *Application) Ingest(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()

	var payload EventPayload

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&payload); err != nil {
		log.Printf("failed to decode request body: %v", err)
		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	if err := a.validate.Struct(payload); err != nil {
		log.Printf("payload validation failed: %v", err)
		http.Error(w, "invalid payload", http.StatusBadRequest)
		return
	}

	// For now, accept valid payloads without additional processing.
	w.WriteHeader(http.StatusNoContent)
	fmt.Fprintln(w)
}

// init is the composition root for the Cloud Function.
// It constructs production Dependencies and Application, then registers the HTTP function.
func init() {
	deps := Dependencies{}
	app := NewApplication(deps)

	functions.HTTP("Ingest", app.Ingest)
}

