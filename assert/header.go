package assert

import (
	"fmt"
	"net/http"
	"regexp"
)

// Header asserts a response header field value matches.
// Regular expressions can be used as value to perform the specific assertions.
func Header(key, value string) Func {
	return func(res *http.Response, req *http.Request) error {
		header := res.Header.Get(key)
		if match, _ := regexp.MatchString(value, header); !match {
			return fmt.Errorf("Header mismatch: '%s' should match '%s'", value, header)
		}
		return nil
	}
}

// HeaderEquals asserts a response header field
// with the given value.
func HeaderEquals(key, value string) Func {
	return func(res *http.Response, req *http.Request) error {
		if header := res.Header.Get(key); header != value {
			return fmt.Errorf("Header mismatch: '%s' == '%s'", value, header)
		}
		return nil
	}
}

// HeaderPresent asserts if a header field is present in the response.
func HeaderPresent(key string) Func {
	return func(res *http.Response, req *http.Request) error {
		if header := res.Header.Get(key); header == "" {
			return fmt.Errorf("Header is not present: %s", key)
		}
		return nil
	}
}
