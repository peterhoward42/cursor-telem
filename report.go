package function

// HowManyPeopleHave holds counts of distinct ProxyUserID per engagement action.
type HowManyPeopleHave struct {
	Launched                   int `json:"Launched"`
	LoadedAnExample            int `json:"LoadedAnExample"`
	TriedToSignIn              int `json:"TriedToSignIn"`
	SucceededSigningIn         int `json:"SucceededSigningIn"`
	CreatedTheirOwnDrawing     int `json:"CreatedTheirOwnDrawing"`
	RetreivedTheirASavedDrawing int `json:"RetreivedTheirASavedDrawing"`
}

// Report is the result of analysing a set of event payloads.
type Report struct {
	HowManyPeopleHave     HowManyPeopleHave `json:"HowManyPeopleHave"`
	TotalRecoverableErrors int               `json:"TotalRecoverableErrors"`
	TotalFatalErrors       int               `json:"TotalFatalErrors"`
}
