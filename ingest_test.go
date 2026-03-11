package function

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type erroringEventStorer struct{}

func (e *erroringEventStorer) StoreIfNotExist(evt *EventPayload, path string) error {
	return errors.New("store failed")
}

type erroringEventGetter struct{}

func (e *erroringEventGetter) GetStoredEvents() ([]*EventPayload, error) {
	return nil, errors.New("get failed")
}

func newTestApplicationWithStorer(storer EventStorer) *Application {
	return NewApplication(Dependencies{
		EventStorer: storer,
		EventGetter: NewFakeEventGetter(nil),
		Analyser:    NewAnalyser(),
	})
}

func validEventJSON() string {
	return `{"schemaVersion":0,"eventULID":"01ARZ3NDEKTSV4RRFFQ69G5FAV","proxyUserID":"123e4567-e89b-42d3-a456-426614174000","timeUTC":"2025-03-10T11:22:33Z","visit":1,"event":"view","parameters":"param"}`
}

func TestIngest_MethodNotAllowed(t *testing.T) {
	app := newTestApplicationWithStorer(NewFakeEventStorer())

	req := httptest.NewRequest(http.MethodPut, "/", nil)
	rr := httptest.NewRecorder()

	app.Ingest(rr, req)

	if rr.Code != http.StatusMethodNotAllowed {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusMethodNotAllowed)
	}
}

func TestIngest_InvalidJSON_ReturnsBadRequest(t *testing.T) {
	app := newTestApplicationWithStorer(NewFakeEventStorer())

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(`{"schemaVersion":`))
	rr := httptest.NewRecorder()

	app.Ingest(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}

	if body := rr.Body.String(); !strings.Contains(body, "invalid JSON payload") {
		t.Fatalf("body = %q, want to contain %q", body, "invalid JSON payload")
	}
}

func TestIngest_InvalidPayload_ReturnsBadRequest(t *testing.T) {
	app := newTestApplicationWithStorer(NewFakeEventStorer())

	// eventULID is too short to satisfy the validation tag.
	body := `{"schemaVersion":0,"eventULID":"short","proxyUserID":"123e4567-e89b-42d3-a456-426614174000","timeUTC":"2025-03-10T11:22:33Z","visit":1,"event":"view","parameters":"param"}`

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(body))
	rr := httptest.NewRecorder()

	app.Ingest(rr, req)

	if rr.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusBadRequest)
	}

	if body := rr.Body.String(); !strings.Contains(body, "invalid payload") {
		t.Fatalf("body = %q, want to contain %q", body, "invalid payload")
	}
}

func TestIngest_ValidPayload_StoresEventAndReturnsNoContent(t *testing.T) {
	storer := NewFakeEventStorer()
	app := newTestApplicationWithStorer(storer)

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(validEventJSON()))
	rr := httptest.NewRecorder()

	app.Ingest(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNoContent)
	}

	if len(storer.data) != 1 {
		t.Fatalf("expected exactly one stored event, got %d", len(storer.data))
	}

	var storedPath string
	for k := range storer.data {
		storedPath = k
	}

	if !strings.HasPrefix(storedPath, "events/") {
		t.Fatalf("stored path = %q, want it to start with %q", storedPath, "events/")
	}
}

func TestIngest_StoreFailureReturnsInternalServerError(t *testing.T) {
	app := newTestApplicationWithStorer(&erroringEventStorer{})

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(validEventJSON()))
	rr := httptest.NewRecorder()

	app.Ingest(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusInternalServerError)
	}

	if body := rr.Body.String(); !strings.Contains(body, "failed to store event") {
		t.Fatalf("body = %q, want to contain %q", body, "failed to store event")
	}
}

func TestIngest_NilEventStorerStillSucceeds(t *testing.T) {
	app := NewApplication(Dependencies{})

	req := httptest.NewRequest(http.MethodPost, "/", strings.NewReader(validEventJSON()))
	rr := httptest.NewRecorder()

	app.Ingest(rr, req)

	if rr.Code != http.StatusNoContent {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusNoContent)
	}
}

func TestReport_GET_ReturnsReportJSON(t *testing.T) {
	events := []*EventPayload{
		{ProxyUserID: "u1", Event: EventLaunched},
		{ProxyUserID: "u2", Event: EventLoadedExample},
	}
	app := NewApplication(Dependencies{
		EventGetter: NewFakeEventGetter(events),
		Analyser:    NewAnalyser(),
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	app.Ingest(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	if ct := rr.Header().Get("Content-Type"); ct != "application/json" {
		t.Fatalf("Content-Type = %q, want application/json", ct)
	}
	expected := &Report{
		HowManyPeopleHave: HowManyPeopleHave{
			Launched:        1,
			LoadedAnExample: 1,
		},
		TotalRecoverableErrors: 0,
		TotalFatalErrors:       0,
	}
	var got Report
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	assertReportEqual(t, &got, expected)
}

func TestReport_GET_EmptyStorage_ReturnsZeroReport(t *testing.T) {
	app := NewApplication(Dependencies{
		EventGetter: NewFakeEventGetter(nil),
		Analyser:    NewAnalyser(),
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	app.Ingest(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusOK)
	}
	expected := &Report{
		HowManyPeopleHave:     HowManyPeopleHave{},
		TotalRecoverableErrors: 0,
		TotalFatalErrors:       0,
	}
	var got Report
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("decode response: %v", err)
	}
	assertReportEqual(t, &got, expected)
}

func TestReport_GET_EventGetterFails_ReturnsInternalServerError(t *testing.T) {
	app := NewApplication(Dependencies{
		EventGetter: &erroringEventGetter{},
		Analyser:    NewAnalyser(),
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	app.Ingest(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusInternalServerError)
	}
	if body := rr.Body.String(); !strings.Contains(body, "failed to get events") {
		t.Fatalf("body = %q, want to contain %q", body, "failed to get events")
	}
}

func TestReport_GET_NilEventGetter_ReturnsInternalServerError(t *testing.T) {
	app := NewApplication(Dependencies{
		EventGetter: nil,
		Analyser:    NewAnalyser(),
	})

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rr := httptest.NewRecorder()

	app.Ingest(rr, req)

	if rr.Code != http.StatusInternalServerError {
		t.Fatalf("status = %d, want %d", rr.Code, http.StatusInternalServerError)
	}
}

