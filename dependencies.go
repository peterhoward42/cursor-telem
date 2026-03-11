package function

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-playground/validator/v10"
)

// EventStorer persists validated events if they are not already stored.
type EventStorer interface {
	StoreIfNotExist(evt *EventPayload, path string) error
}

// EventGetter retrieves stored events for analysis and reporting.
type EventGetter interface {
	GetStoredEvents() ([]*EventPayload, error)
}

// Dependencies aggregates all external systems required by the application.
// Start empty and add fields only when needed.
type Dependencies struct {
	EventStorer EventStorer
	EventGetter EventGetter
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

	path, err := buildEventPath(&payload)
	if err != nil {
		log.Printf("failed to build event path: %v", err)
		http.Error(w, "invalid event time", http.StatusBadRequest)
		return
	}

	if a.deps.EventStorer != nil {
		if err := a.deps.EventStorer.StoreIfNotExist(&payload, path); err != nil {
			log.Printf("failed to store event: %v", err)
			http.Error(w, "failed to store event", http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusNoContent)
}

