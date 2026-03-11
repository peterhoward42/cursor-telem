package function

// FakeEventGetter is an in-memory, deterministic implementation of EventGetter.
// It always returns the same set of events provided at construction time.
type FakeEventGetter struct {
	events []*EventPayload
}

// NewFakeEventGetter constructs a FakeEventGetter that will return the given events.
func NewFakeEventGetter(events []*EventPayload) *FakeEventGetter {
	return &FakeEventGetter{
		events: events,
	}
}

// GetStoredEvents returns the stored events without modification.
func (g *FakeEventGetter) GetStoredEvents() ([]*EventPayload, error) {
	return g.events, nil
}

