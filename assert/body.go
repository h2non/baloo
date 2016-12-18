package assert

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
)

func readBody(res *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	// Re-fill body reader stream after reading it
	res.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return body, err
}

// BodyMatchString asserts a response body matching a string expression.
// Regular expressions can be used as value to perform the specific assertions.
func BodyMatchString(pattern string) Func {
	return func(res *http.Response, req *http.Request) error {
		body, err := readBody(res)
		if err != nil {
			return err
		}
		if match, _ := regexp.MatchString(pattern, string(body)); !match {
			return fmt.Errorf("body mismatch: pattern not found '%s'", pattern)
		}
		return nil
	}
}

// BodyEquals asserts as strict equality comparison the
// response body with a given string string.
func BodyEquals(value string) Func {
	return func(res *http.Response, req *http.Request) error {
		body, err := readBody(res)
		if err != nil {
			return err
		}

		bodyStr := string(body)
		err = fmt.Errorf("bodies mismatch:\n\thave: %#v\n\twant: %#v", bodyStr, value)

		// Remove line feed sequence
		if len(bodyStr) > 0 && bodyStr[len(bodyStr)-1] == '\n' {
			bodyStr = bodyStr[:len(bodyStr)-1]
		}

		// Perform the comparison
		if len(bodyStr) != len(value) || value != bodyStr {
			return err
		}

		return nil
	}
}

// BodyLength asserts a response body length.
func BodyLength(length int) Func {
	return func(res *http.Response, req *http.Request) error {
		cl, err := strconv.Atoi(res.Header.Get("Content-Length"))
		// Infer length from body buffer
		if err != nil || cl == 0 {
			body, err := readBody(res)
			if err != nil {
				return err
			}
			cl = len(body)
		}
		// Match body length
		if cl != length {
			return fmt.Errorf("body length mismatch: '%d' should be equal to '%d'", cl, length)
		}
		return nil
	}
}
