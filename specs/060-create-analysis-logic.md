# Creating the analysis logic

- Create an encapsulated "Analyser" that performs an analysis of a set of
  given  *EventPayload objects
- The analysis should yield the Report data structure defined below using the logic also
  defined below, or an error.
- Write unit tests for the analyser

## Report data structure

```
{
  "HowManyPeopleHave": {
    "Launched": 0,
    "LoadedAnExample": 0,
    "TriedToSignIn": 0,
    "SucceededSigningIn": 0,
    "CreatedTheirOwnDrawing": 0,
    "RetreivedTheirASavedDrawing": 0
  },
  "TotalRecoverableErrors": 0,
  "TotalFatalErrors": 0
}
```

## Analysis logic

HowManyPeopleHave.Launched: count of distinct ProxyUserID that have at least one launched event.

HowManyPeopleHave.LoadedAnExample: distinct ProxyUserID with at least one loaded-example event.

HowManyPeopleHave.TriedToSignIn: distinct ProxyUserID with at least one sign-in-started event.

HowManyPeopleHave.SucceededSigningIn: distinct ProxyUserID with at least one sign-in-success event.

HowManyPeopleHave.CreatedTheirOwnDrawing: distinct ProxyUserID with at least one created-new-drawing event.

HowManyPeopleHave.RetreivedTheirASavedDrawing: distinct ProxyUserID with at least one retreived-save-drawing event.

TotalRecoverableErrors: total count of events with Event == recoverable-javascript-error.

TotalFatalErrors: total count of events with Event == fatal-javascript-error.

Malformed stored events MUST be skipped.
