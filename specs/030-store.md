# Storing incoming events

- Incoming telemetry events should be stored
- Add an EventStorer to the Dependencies structure with the signature defined
  below
- Have the Application POST handler store the incoming event using the
  EventStorer passed in to the Application constructor
- For the time being only a fake, test-double implementation of a storer is
  necessary.
- The event should be stored in the form of gzip compressed NDJSON
- The test for if this event already exists is the generated path
  already being present in the store


## EventStorer interface

```
StoreIfNotExist(evt *EventPayload, path string) error
```

## Building the path parameter for StoreIfNotExist

- The path is built deterministically based on the EventPayload contents.
- It encodes time buckets
- It is derived from the time field inside the incoming Event.

The path should like like this:

```
events/y=2025/m=11/d=01/hour=14/<eventULID>.ndjson.gz
```
