# Implement a POST request handler

	- JSON payload, representing a single instance of the EventPayload go 
	  struct defined below
	- Apply validation to the payload

```
type EventPayload struct {
	SchemaVersion int `validate:"max=1,min=0"`
	// The EventULID not only makes an event global unique, but also facilitates
	// sorting events into a time ordered sequence.
	EventULID string `validate:"len=26"`
	// The ProxyUserID is our anonymised DrawExact user-key.
	ProxyUserID string `validate:"uuid4"`
	TimeUTC     string `validate:"datetime=2006-01-02T15:04:05Z07:00"` // RFC3339
	// The Visit (count) tells us how many cumulative DrawExact sessions the user had launched at the time of this event.
	Visit int `validate:"max=100000,min=1"`
	// The Event tells us what event occurred. E.g. the user completely training step 3.
	Event string `validate:"max=40,min=4"`
	// Some events require additional parameters to adequately describe them.
	// Each event type conveys those parameters in any way it thinks fit encoded into this string.
	Parameters string `validate:"max=80"`
}
```


