# Adding Google Cloud Storage for production

- Create a new additional implementation of EventStorer that uses 
  Google Cloud Storage
- Use the google.com/go/storage Go package 
- Hard code the bucket name `drawexact-telemetry'
- Change the event storer used in the dependencies for the cloud 
  function to the new one. 
- Do not change the event storer used by the tests.
