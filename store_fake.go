package function

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
)

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

