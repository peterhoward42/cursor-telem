# Adding Google Cloud Storage implementation of the EventGetter for production

- Create a new additional implementation of EventGetter that uses 
  Google Cloud Storage
- Use the google.com/go/storage Go package 
- Change the event getter used in the dependencies for the cloud 
  function to the new one. 
- Do not change the event getter used by the tests.
