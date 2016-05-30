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
			return fmt.Errorf("Body mismatch: cannot match pattern '%s'", pattern)
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
			return fmt.Errorf("Body length mismatch: '%d' should be equal '%d'", cl, length)
		}
		return nil
	}
}
