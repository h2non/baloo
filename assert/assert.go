package assert

import (
	"net/http"
)

// Func represents the required interface for assertion functions.
type Func func(*http.Response, *http.Request) error
