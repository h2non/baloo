package assert

import (
	"fmt"
	"net/http"
)

// StatusOk asserts the response status code as
// valid response (>= 200 && < 400).
func StatusOk() Func {
	return StatusRange(200, 399)
}

// StatusError asserts the response status code as
// server error response (>= 400 && < 600).
func StatusError() Func {
	return StatusRange(400, 599)
}

// StatusServerError asserts the response status code as
// server error response (>= 500 && < 600).
func StatusServerError() Func {
	return StatusRange(500, 599)
}

// StatusClientError asserts the response status code as
// client error response (>= 400 && < 500).
func StatusClientError() Func {
	return StatusRange(400, 499)
}

// StatusRedirect asserts the response status code as
// server redirect status (= 301 || = 302).
func StatusRedirect() Func {
	return StatusRange(301, 302)
}

// StatusEqual asserts the response status code with the given number.
func StatusEqual(code int) Func {
	return func(res *http.Response, req *http.Request) error {
		if res.StatusCode != code {
			return fmt.Errorf("Unexpected status code: %d != %d", res.StatusCode, code)
		}
		return nil
	}
}

// StatusRange asserts the response status code if
// it is within the given numeric range.
func StatusRange(start, end int) Func {
	return func(res *http.Response, req *http.Request) error {
		if res.StatusCode >= start && res.StatusCode <= end {
			return nil
		}
		return fmt.Errorf("Status code outside the given range: %d >= %d && %d <= %d",
			res.StatusCode, start, res.StatusCode, end)
	}
}
