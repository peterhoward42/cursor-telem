package function

// Event names used for analysis.
const (
	EventLaunched               = "launched"
	EventLoadedExample         = "loaded-example"
	EventSignInStarted         = "sign-in-started"
	EventSignInSuccess         = "sign-in-success"
	EventCreatedNewDrawing     = "created-new-drawing"
	EventRetreivedSaveDrawing  = "retreived-save-drawing"
	EventRecoverableJSError    = "recoverable-javascript-error"
	EventFatalJSError          = "fatal-javascript-error"
)

// Analyser performs analysis of event payloads and produces a Report.
type Analyser struct{}

// NewAnalyser constructs an Analyser.
func NewAnalyser() *Analyser {
	return &Analyser{}
}

// Analyse consumes a set of event payloads and returns a Report.
// Malformed events (nil, or missing ProxyUserID or Event) are skipped.
func (a *Analyser) Analyse(events []*EventPayload) (*Report, error) {
	launched := make(map[string]struct{})
	loadedExample := make(map[string]struct{})
	triedSignIn := make(map[string]struct{})
	succeededSignIn := make(map[string]struct{})
	createdDrawing := make(map[string]struct{})
	retreivedDrawing := make(map[string]struct{})
	var recoverableErrors, fatalErrors int

	for _, evt := range events {
		if evt == nil || evt.ProxyUserID == "" || evt.Event == "" {
			continue
		}
		uid := evt.ProxyUserID
		switch evt.Event {
		case EventLaunched:
			launched[uid] = struct{}{}
		case EventLoadedExample:
			loadedExample[uid] = struct{}{}
		case EventSignInStarted:
			triedSignIn[uid] = struct{}{}
		case EventSignInSuccess:
			succeededSignIn[uid] = struct{}{}
		case EventCreatedNewDrawing:
			createdDrawing[uid] = struct{}{}
		case EventRetreivedSaveDrawing:
			retreivedDrawing[uid] = struct{}{}
		case EventRecoverableJSError:
			recoverableErrors++
		case EventFatalJSError:
			fatalErrors++
		default:
			// other event types are ignored for this report
		}
	}

	return &Report{
		HowManyPeopleHave: HowManyPeopleHave{
			Launched:                   len(launched),
			LoadedAnExample:            len(loadedExample),
			TriedToSignIn:              len(triedSignIn),
			SucceededSigningIn:         len(succeededSignIn),
			CreatedTheirOwnDrawing:     len(createdDrawing),
			RetreivedTheirASavedDrawing: len(retreivedDrawing),
		},
		TotalRecoverableErrors: recoverableErrors,
		TotalFatalErrors:       fatalErrors,
	}, nil
}
