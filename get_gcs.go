package function

import (
	"context"
	"encoding/json"
	"io"

	"cloud.google.com/go/storage"
	"google.golang.org/api/iterator"
)

// GCSEventGetter retrieves stored events from Google Cloud Storage.
type GCSEventGetter struct {
	client     *storage.Client
	bucketName string
}

// NewGCSEventGetter constructs a GCSEventGetter that reads from the given bucket.
func NewGCSEventGetter(client *storage.Client, bucketName string) *GCSEventGetter {
	return &GCSEventGetter{
		client:     client,
		bucketName: bucketName,
	}
}

// GetStoredEvents lists all objects under the events/ prefix, reads each
// gzip-compressed NDJSON object, and returns the decoded events.
func (g *GCSEventGetter) GetStoredEvents() ([]*EventPayload, error) {
	ctx := context.Background()
	bucket := g.client.Bucket(g.bucketName)
	it := bucket.Objects(ctx, &storage.Query{Prefix: "events/"})

	var events []*EventPayload
	for {
		attrs, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}

		evt, err := g.readEvent(ctx, bucket.Object(attrs.Name))
		if err != nil {
			return nil, err
		}
		events = append(events, evt)
	}

	return events, nil
}

// readEvent opens the object (GCS decompresses when ContentEncoding is gzip),
// then decodes the single NDJSON line into an EventPayload.
func (g *GCSEventGetter) readEvent(ctx context.Context, obj *storage.ObjectHandle) (*EventPayload, error) {
	r, err := obj.NewReader(ctx)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var payload EventPayload
	if err := json.NewDecoder(r).Decode(&payload); err != nil && err != io.EOF {
		return nil, err
	}
	return &payload, nil
}
