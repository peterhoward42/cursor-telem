package function

import "testing"

func TestBuildEventPath_Success(t *testing.T) {
	evt := &EventPayload{
		EventULID: "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		TimeUTC:   "2025-03-10T11:22:33Z",
	}

	got, err := buildEventPath(evt)
	if err != nil {
		t.Fatalf("buildEventPath() error = %v, want nil", err)
	}

	want := "events/y=2025/m=03/d=10/hour=11/01ARZ3NDEKTSV4RRFFQ69G5FAV.ndjson.gz"
	if got != want {
		t.Errorf("buildEventPath() = %q, want %q", got, want)
	}
}

func TestBuildEventPath_InvalidTime(t *testing.T) {
	evt := &EventPayload{
		EventULID: "01ARZ3NDEKTSV4RRFFQ69G5FAV",
		TimeUTC:   "not-a-time",
	}

	got, err := buildEventPath(evt)
	if err == nil {
		t.Fatalf("buildEventPath() error = nil, want non-nil; got path %q", got)
	}
}

