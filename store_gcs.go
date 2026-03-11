package function

import (
	"context"

	"cloud.google.com/go/storage"
)

// GCSEventStorer persists events to Google Cloud Storage.
type GCSEventStorer struct {
	client     *storage.Client
	bucketName string
}

// NewGCSEventStorer constructs a GCSEventStorer that writes to the given bucket.
func NewGCSEventStorer(client *storage.Client, bucketName string) *GCSEventStorer {
	return &GCSEventStorer{
		client:     client,
		bucketName: bucketName,
	}
}

// StoreIfNotExist stores the event as gzip-compressed NDJSON at the given path
// if that object does not already exist.
func (s *GCSEventStorer) StoreIfNotExist(evt *EventPayload, path string) error {
	ctx := context.Background()

	obj := s.client.Bucket(s.bucketName).Object(path)

	// Check if the object already exists.
	if _, err := obj.Attrs(ctx); err == nil {
		return nil
	} else if err != storage.ErrObjectNotExist {
		return err
	}

	data, err := encodeEventToGzippedNDJSON(evt)
	if err != nil {
		return err
	}

	w := obj.NewWriter(ctx)
	w.ContentType = "application/x-ndjson"
	w.ContentEncoding = "gzip"

	if _, err := w.Write(data); err != nil {
		_ = w.Close()
		return err
	}

	if err := w.Close(); err != nil {
		return err
	}

	return nil
}

