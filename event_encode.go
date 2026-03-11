package function

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
)

// encodeEventToGzippedNDJSON encodes a single event as gzip-compressed NDJSON.
func encodeEventToGzippedNDJSON(evt *EventPayload) ([]byte, error) {
	var buf bytes.Buffer

	gw := gzip.NewWriter(&buf)
	encoder := json.NewEncoder(gw)

	if err := encoder.Encode(evt); err != nil {
		gw.Close()
		return nil, err
	}

	if err := gw.Close(); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
