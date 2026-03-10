package function

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"time"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
	"github.com/go-playground/validator/v10"
)

// EventStorer persists validated events if they are not already stored.
type EventStorer interface {
	StoreIfNotExist(evt *EventPayload, path string) error
}

// Dependencies aggregates all external systems required by the application.
// Start empty and add fields only when needed.
type Dependencies struct {
	EventStorer EventStorer
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
	fmt.Fprintln(w)
}

// buildEventPath deterministically builds the storage path for an event based on its timestamp and ULID.
func buildEventPath(evt *EventPayload) (string, error) {
	t, err := time.Parse(time.RFC3339, evt.TimeUTC)
	if err != nil {
		return "", err
	}

	year, month, day := t.Date()
	hour := t.Hour()

	filename := fmt.Sprintf("%s.ndjson.gz", evt.EventULID)

	return filepath.ToSlash(
		filepath.Join(
			"events",
			fmt.Sprintf("y=%04d", year),
			fmt.Sprintf("m=%02d", int(month)),
			fmt.Sprintf("d=%02d", day),
			fmt.Sprintf("hour=%02d", hour),
			filename,
		),
	), nil
}

// FakeEventStorer is an in-memory, test-double implementation of EventStorer.
// It stores gzip-compressed NDJSON representations of events keyed by path.
type FakeEventStorer struct {
	data map[string][]byte
}

// NewFakeEventStorer constructs a FakeEventStorer.
func NewFakeEventStorer() *FakeEventStorer {
	return &FakeEventStorer{
		data: make(map[string][]byte),
	}
}

// StoreIfNotExist stores the event as gzip-compressed NDJSON at the given path
// if that path is not already present.
func (s *FakeEventStorer) StoreIfNotExist(evt *EventPayload, path string) error {
	if _, exists := s.data[path]; exists {
		return nil
	}

	var buf bytes.Buffer

	gw := gzip.NewWriter(&buf)
	encoder := json.NewEncoder(gw)

	if err := encoder.Encode(evt); err != nil {
		gw.Close()
		return err
	}

	if err := gw.Close(); err != nil {
		return err
	}

	s.data[path] = buf.Bytes()
	return nil
}

// init is the composition root for the Cloud Function.
// It constructs production Dependencies and Application, then registers the HTTP function.
func init() {
	deps := Dependencies{
		EventStorer: NewFakeEventStorer(),
	}
	app := NewApplication(deps)

	functions.HTTP("Ingest", app.Ingest)
}

