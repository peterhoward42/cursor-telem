# Implement a GET request handler wired up to report generation

- Extend Application::Injest to also support a GET request
- The GET request should:
  - retrieve all the EventPayloads from storage using the EventGetter
  - use an Analyser to create a Report
  - return a JSON payload representing the Report
- Write tests for the new behaviour