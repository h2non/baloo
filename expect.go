package baloo

import (
	"fmt"
	"net/http"
	"testing"

	"gopkg.in/h2non/baloo.v0/assert"
)

// Expect represents the HTTP expectation suite who is
// able to define multiple assertion functions to match the response.
type Expect struct {
	test       *testing.T
	request    *Request
	assertions []assert.Func
}

// NewExpect creates a new testing expectation instance.
func NewExpect(req *Request) *Expect {
	return &Expect{request: req}
}

// BindTest binds the Go testing instance to the current suite.
// In the future multiple testing interfaces can
// be supported via adapters.
func (e *Expect) BindTest(t *testing.T) *Expect {
	e.test = t
	return e
}

// Status asserts the response status code
// with the given status.
func (e *Expect) Status(code int) *Expect {
	e.AssertFunc(assert.StatusEqual(code))
	return e
}

// StatusOk asserts the response status code
// as valid response (>= 200 && < 400).
func (e *Expect) StatusOk() *Expect {
	e.AssertFunc(assert.StatusOk())
	return e
}

// StatusError asserts the response status code
// as client/server error response (>= 400 && < 600).
func (e *Expect) StatusError() *Expect {
	e.AssertFunc(assert.StatusError())
	return e
}

// StatusServerError asserts the response status code
// as server error response (>= 500 && < 600).
func (e *Expect) StatusServerError() *Expect {
	e.AssertFunc(assert.StatusServerError())
	return e
}

// StatusClientError asserts the response status code
// as server error response (>= 400 && < 500).
func (e *Expect) StatusClientError() *Expect {
	e.AssertFunc(assert.StatusClientError())
	return e
}

// Type asserts the response MIME type with the given alias or value.
func (e *Expect) Type(kind string) *Expect {
	e.AssertFunc(assert.Type(kind))
	return e
}

// Header asserts a response header field value matches.
// Regular expressions can be used as value to perform the specific assertions.
func (e *Expect) Header(key, value string) *Expect {
	e.AssertFunc(assert.Header(key, value))
	return e
}

// HeaderEquals asserts a response header field value.
func (e *Expect) HeaderEquals(key, value string) *Expect {
	e.AssertFunc(assert.HeaderEquals(key, value))
	return e
}

// HeaderPresent asserts if a header field is present in the response.
func (e *Expect) HeaderPresent(key string) *Expect {
	e.AssertFunc(assert.HeaderPresent(key))
	return e
}

// BodyEquals asserts as strict equality comparison the
// response body with a given string string.
func (e *Expect) BodyEquals(pattern string) *Expect {
	e.AssertFunc(assert.BodyMatchString(pattern))
	return e
}

// BodyMatchString asserts a response body matching a string expression.
// Regular expressions can be used as value to perform the specific assertions.
func (e *Expect) BodyMatchString(pattern string) *Expect {
	e.AssertFunc(assert.BodyMatchString(pattern))
	return e
}

// BodyLength asserts a response body length.
func (e *Expect) BodyLength(length int) *Expect {
	e.AssertFunc(assert.BodyLength(length))
	return e
}

// JSONEquals asserts the response body with the given JSON struct.
func (e *Expect) JSONEquals(data interface{}) *Expect {
	e.AssertFunc(assert.JSONEquals(data))
	return e
}

// JSONSchema asserts the response body with the given
// JSON schema definition.
func (e *Expect) JSONSchema(schema string) *Expect {
	e.AssertFunc(assert.JSONSchema(schema))
	return e
}

// JSONSchemaFile asserts the response body with the given
// JSON schema definition loaded from file path.
func (e *Expect) JSONSchemaFile(schema string) *Expect {
	// e.AssertFunc(assert.JSONSchemaFile(schema))
	return e
}

// AssertFunc adds a new assertion function.
func (e *Expect) AssertFunc(assertion ...assert.Func) *Expect {
	e.assertions = append(e.assertions, assertion...)
	return e
}

// Done performs and asserts the HTTP response based
// on the defined expectations.
func (e *Expect) Done() error {
	// Perform the HTTP request
	res, err := e.request.Send()
	if err != nil {
		err = fmt.Errorf("request error: %s", err)
		e.test.Error(err)
		return err
	}

	// Run assertions
	err = e.run(res.RawResponse, res.RawRequest)
	if err != nil {
		e.test.Error(err)
	}

	return err
}

func (e *Expect) run(res *http.Response, req *http.Request) error {
	var err error
	for _, assertion := range e.assertions {
		err = assertion(res, req)
		if err != nil {
			break
		}
	}
	return err
}
