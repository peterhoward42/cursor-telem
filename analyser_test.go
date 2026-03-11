package function

import (
	"testing"
)

func TestAnalyser_Analyse_EmptyInput_ReturnsAllZeros(t *testing.T) {
	a := NewAnalyser()

	got, err := a.Analyse(nil)
	if err != nil {
		t.Fatalf("Analyse() error = %v, want nil", err)
	}

	want := &Report{
		HowManyPeopleHave:     HowManyPeopleHave{},
		TotalRecoverableErrors: 0,
		TotalFatalErrors:       0,
	}
	assertReportEqual(t, got, want)
}

func TestAnalyser_Analyse_EmptySlice_ReturnsAllZeros(t *testing.T) {
	a := NewAnalyser()

	got, err := a.Analyse([]*EventPayload{})
	if err != nil {
		t.Fatalf("Analyse() error = %v, want nil", err)
	}

	want := &Report{
		HowManyPeopleHave:     HowManyPeopleHave{},
		TotalRecoverableErrors: 0,
		TotalFatalErrors:       0,
	}
	assertReportEqual(t, got, want)
}

func TestAnalyser_Analyse_OneEventPerType_CountsOneInEachBucket(t *testing.T) {
	events := []*EventPayload{
		{ProxyUserID: "u1", Event: EventLaunched},
		{ProxyUserID: "u2", Event: EventLoadedExample},
		{ProxyUserID: "u3", Event: EventSignInStarted},
		{ProxyUserID: "u4", Event: EventSignInSuccess},
		{ProxyUserID: "u5", Event: EventCreatedNewDrawing},
		{ProxyUserID: "u6", Event: EventRetreivedSaveDrawing},
	}
	a := NewAnalyser()

	got, err := a.Analyse(events)
	if err != nil {
		t.Fatalf("Analyse() error = %v, want nil", err)
	}

	want := &Report{
		HowManyPeopleHave: HowManyPeopleHave{
			Launched:                   1,
			LoadedAnExample:            1,
			TriedToSignIn:              1,
			SucceededSigningIn:         1,
			CreatedTheirOwnDrawing:     1,
			RetreivedTheirASavedDrawing: 1,
		},
		TotalRecoverableErrors: 0,
		TotalFatalErrors:       0,
	}
	assertReportEqual(t, got, want)
}

func TestAnalyser_Analyse_MalformedEventsSkipped(t *testing.T) {
	events := []*EventPayload{
		nil,
		{ProxyUserID: "u1", Event: EventLaunched},
		{ProxyUserID: "", Event: EventLaunched},
		{ProxyUserID: "u2", Event: ""},
		{ProxyUserID: "u3", Event: EventLaunched},
	}
	a := NewAnalyser()

	got, err := a.Analyse(events)
	if err != nil {
		t.Fatalf("Analyse() error = %v, want nil", err)
	}

	// Only u1 and u3 have both ProxyUserID and Event; both launched → 2 distinct.
	if got.HowManyPeopleHave.Launched != 2 {
		t.Errorf("HowManyPeopleHave.Launched = %d, want 2 (malformed skipped)", got.HowManyPeopleHave.Launched)
	}
}

func TestAnalyser_Analyse_DistinctUsersSameAction_CountsDistinct(t *testing.T) {
	events := []*EventPayload{
		{ProxyUserID: "user-a", Event: EventLaunched},
		{ProxyUserID: "user-b", Event: EventLaunched},
		{ProxyUserID: "user-c", Event: EventLaunched},
	}
	a := NewAnalyser()

	got, err := a.Analyse(events)
	if err != nil {
		t.Fatalf("Analyse() error = %v, want nil", err)
	}

	if got.HowManyPeopleHave.Launched != 3 {
		t.Errorf("HowManyPeopleHave.Launched = %d, want 3 (distinct ProxyUserID)", got.HowManyPeopleHave.Launched)
	}
}

func TestAnalyser_Analyse_SameUserMultipleSameAction_CountsOnce(t *testing.T) {
	events := []*EventPayload{
		{ProxyUserID: "user-a", Event: EventLaunched},
		{ProxyUserID: "user-a", Event: EventLaunched},
		{ProxyUserID: "user-a", Event: EventLaunched},
	}
	a := NewAnalyser()

	got, err := a.Analyse(events)
	if err != nil {
		t.Fatalf("Analyse() error = %v, want nil", err)
	}

	if got.HowManyPeopleHave.Launched != 1 {
		t.Errorf("HowManyPeopleHave.Launched = %d, want 1 (same user counted once)", got.HowManyPeopleHave.Launched)
	}
}

func TestAnalyser_Analyse_ErrorEvents_CountTotalOccurrences(t *testing.T) {
	events := []*EventPayload{
		{ProxyUserID: "u1", Event: EventRecoverableJSError},
		{ProxyUserID: "u1", Event: EventRecoverableJSError},
		{ProxyUserID: "u2", Event: EventFatalJSError},
		{ProxyUserID: "u2", Event: EventFatalJSError},
		{ProxyUserID: "u2", Event: EventFatalJSError},
	}
	a := NewAnalyser()

	got, err := a.Analyse(events)
	if err != nil {
		t.Fatalf("Analyse() error = %v, want nil", err)
	}

	if got.TotalRecoverableErrors != 2 {
		t.Errorf("TotalRecoverableErrors = %d, want 2", got.TotalRecoverableErrors)
	}
	if got.TotalFatalErrors != 3 {
		t.Errorf("TotalFatalErrors = %d, want 3", got.TotalFatalErrors)
	}
}

func TestAnalyser_Analyse_UnknownEventType_Ignored(t *testing.T) {
	events := []*EventPayload{
		{ProxyUserID: "u1", Event: "unknown-event"},
		{ProxyUserID: "u2", Event: "view"},
	}
	a := NewAnalyser()

	got, err := a.Analyse(events)
	if err != nil {
		t.Fatalf("Analyse() error = %v, want nil", err)
	}

	// No engagement or error events; all counts stay zero.
	want := &Report{
		HowManyPeopleHave:     HowManyPeopleHave{},
		TotalRecoverableErrors: 0,
		TotalFatalErrors:       0,
	}
	assertReportEqual(t, got, want)
}

func assertReportEqual(t *testing.T, got, want *Report) {
	t.Helper()
	if got == nil || want == nil {
		if got != want {
			t.Fatalf("Report: got %v, want %v", got, want)
		}
		return
	}
	h := got.HowManyPeopleHave
	w := want.HowManyPeopleHave
	if h.Launched != w.Launched || h.LoadedAnExample != w.LoadedAnExample ||
		h.TriedToSignIn != w.TriedToSignIn || h.SucceededSigningIn != w.SucceededSigningIn ||
		h.CreatedTheirOwnDrawing != w.CreatedTheirOwnDrawing ||
		h.RetreivedTheirASavedDrawing != w.RetreivedTheirASavedDrawing {
		t.Errorf("HowManyPeopleHave: got %+v, want %+v", h, w)
	}
	if got.TotalRecoverableErrors != want.TotalRecoverableErrors {
		t.Errorf("TotalRecoverableErrors = %d, want %d", got.TotalRecoverableErrors, want.TotalRecoverableErrors)
	}
	if got.TotalFatalErrors != want.TotalFatalErrors {
		t.Errorf("TotalFatalErrors = %d, want %d", got.TotalFatalErrors, want.TotalFatalErrors)
	}
}
