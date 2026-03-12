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
	Analyser    *Analyser
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

// Telemetry handles POST requests (store event) and GET requests (return report).
// Other methods receive 405 Method Not Allowed.
func (a *Application) Telemetry(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		a.serveReport(w)
		return
	case http.MethodPost:
		a.serveIngest(w, r)
		return
	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
}

func (a *Application) serveReport(w http.ResponseWriter) {
	if a.deps.EventGetter == nil || a.deps.Analyser == nil {
		http.Error(w, "report not configured", http.StatusInternalServerError)
		return
	}
	events, err := a.deps.EventGetter.GetStoredEvents()
	if err != nil {
		log.Printf("failed to get stored events: %v", err)
		http.Error(w, "failed to get events", http.StatusInternalServerError)
		return
	}
	report, err := a.deps.Analyser.Analyse(events)
	if err != nil {
		log.Printf("failed to analyse events: %v", err)
		http.Error(w, "failed to analyse", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(report); err != nil {
		log.Printf("failed to encode report: %v", err)
	}
}

func (a *Application) serveIngest(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if err := r.Body.Close(); err != nil {
			log.Printf("failed to close request body: %v", err)
		}
	}()

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

