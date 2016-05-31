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

// HeaderEquals asserts a response header field value
// is equal to the given value.
func HeaderEquals(key, value string) Func {
	return func(res *http.Response, req *http.Request) error {
		if header := res.Header.Get(key); header != value {
			return fmt.Errorf("Header mismatch: '%s' == '%s'", value, header)
		}
		return nil
	}
}

// HeaderNotEquals asserts a response header field
// with the given value.
func HeaderNotEquals(key, value string) Func {
	return func(res *http.Response, req *http.Request) error {
		if err := HeaderEquals(key, value)(res, req); err == nil {
			header := res.Header.Get(key)
			return fmt.Errorf("Header should not be equal: '%s' != '%s'", value, header)
		}
		return nil
	}
}

// HeaderPresent asserts if a header field is present
// in the response.
func HeaderPresent(key string) Func {
	return func(res *http.Response, req *http.Request) error {
		if header := res.Header.Get(key); header == "" {
			return fmt.Errorf("Header is not present: %s", key)
		}
		return nil
	}
}

// HeaderNotPresent asserts if a header field is not
// present in the response.
func HeaderNotPresent(key string) Func {
	return func(res *http.Response, req *http.Request) error {
		if err := HeaderPresent(key)(res, req); err == nil {
			return fmt.Errorf("Header should not be present: %s", err)
		}
		return nil
	}
}

// RedirectTo asserts the server response redirects
// to the given URL pattern.
// Regular expressions are supported.
func RedirectTo(uri string) Func {
	return func(res *http.Response, req *http.Request) error {
		header := res.Header.Get("Location")
		if header == "" {
			return fmt.Errorf("Location headear is not present")
		}
		if match, _ := regexp.MatchString(uri, header); !match {
			return fmt.Errorf("Invalid location header: '%s' should match '%s'", uri, header)
		}
		return nil
	}
}
