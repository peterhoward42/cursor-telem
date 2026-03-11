# Create a new dependency interface

## EventGetter interface

```
GetStoredEvents() ([]*EventPayload, error)
```

## New dependency field

- Update the Dependencies struct to have an EventGetter field
- Initially use a deterministic test double implementation to implement
  the  interface in all cases
