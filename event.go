package function

import (
	"fmt"
	"path/filepath"
	"time"
)

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

