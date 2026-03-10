package function

import (
	"fmt"
	"net/http"

	"github.com/GoogleCloudPlatform/functions-framework-go/functions"
)

// Dependencies aggregates all external systems required by the application.
// Start empty and add fields only when needed.
type Dependencies struct {
}

// Application owns all request handling logic.
type Application struct {
	deps Dependencies
}

// NewApplication constructs an Application with its explicit dependencies.
func NewApplication(deps Dependencies) *Application {
	return &Application{
		deps: deps,
	}
}

// Hello is the minimal HTTP handler method implementing hello-world behaviour.
func (a *Application) Hello(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello, world")
}

// init is the composition root for the Cloud Function.
// It constructs production Dependencies and Application, then registers the HTTP function.
func init() {
	deps := Dependencies{}
	app := NewApplication(deps)

	functions.HTTP("Hello", app.Hello)
}

