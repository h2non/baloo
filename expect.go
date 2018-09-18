package baloo

import (
	"fmt"
	"net/http"

	"gopkg.in/h2non/baloo.v3/assert"
	"gopkg.in/h2non/gentleman.v2"
)

// Assertions stores global assertion functions by alias name.
// Use Expect.Assert('<assertion name>') to use
// new assertion at request expectation level.
var Assertions = make(map[string]assert.Func)

// AddAssertFunc adds a new assertion function at global level by alias name.
// Then you can trigger the assertion function in any expectation test.
func AddAssertFunc(name string, fn assert.Func) {
	Assertions[name] = fn
}

// FlushAssertFuncs flushes registered assertion functions.
func FlushAssertFuncs() {
	Assertions = make(map[string]assert.Func)
}

// TestingT implements part of the same interface as testing.T
type TestingT interface {
	Error(args ...interface{})
}

// Expect represents the HTTP expectation suite who is
// able to define multiple assertion functions to match the response.
type Expect struct {
	test       TestingT
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
func (e *Expect) BindTest(t TestingT) *Expect {
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

// HeaderEquals asserts a response header field value
// is equal to the given value.
func (e *Expect) HeaderEquals(key, value string) *Expect {
	e.AssertFunc(assert.HeaderEquals(key, value))
	return e
}

// HeaderNotEquals asserts that a response header field value
// is not equal to the given one.
func (e *Expect) HeaderNotEquals(key, value string) *Expect {
	e.AssertFunc(assert.HeaderNotEquals(key, value))
	return e
}

// HeaderPresent asserts if a header field is present
// in the response.
func (e *Expect) HeaderPresent(key string) *Expect {
	e.AssertFunc(assert.HeaderPresent(key))
	return e
}

// HeaderNotPresent asserts if a header field is not
// present in the response.
func (e *Expect) HeaderNotPresent(key string) *Expect {
	e.AssertFunc(assert.HeaderNotPresent(key))
	return e
}

// RedirectTo asserts the server response redirects
// to the given URL pattern.
// Regular expressions are supported.
func (e *Expect) RedirectTo(uri string) *Expect {
	e.AssertFunc(assert.RedirectTo(uri))
	return e
}

// BodyEquals asserts as strict equality comparison the
// response body with a given string string.
func (e *Expect) BodyEquals(pattern string) *Expect {
	e.AssertFunc(assert.BodyEquals(pattern))
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

// JSON asserts the response body with the given JSON struct.
func (e *Expect) JSON(data interface{}) *Expect {
	e.AssertFunc(assert.JSON(data))
	return e
}

// JSONSchema asserts the response body with the given
// JSON schema definition.
func (e *Expect) JSONSchema(schema string) *Expect {
	e.AssertFunc(assert.JSONSchema(schema))
	return e
}

// Assert adds a new assertion function by alias name.
// Assertion function must be previosly registered
// via baloo.AddAssertFunc("alias", function).
func (e *Expect) Assert(assertions ...string) *Expect {
	for _, alias := range assertions {
		fn, ok := Assertions[alias]
		if !ok {
			panic("No assertion function registered by alias: " + alias)
		}
		e.assertions = append(e.assertions, fn)
	}
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

// End is an alias to `Done()`.
func (e *Expect) End() error {
	return e.Done()
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

// Send does the same as `Done()`, but it also returns the `*http.Response` along with the `error`.
func (e *Expect) Send() (*gentleman.Response, error) {
	// Perform the HTTP request
	res, err := e.request.Send()
	if err != nil {
		err = fmt.Errorf("request error: %s", err)
		e.test.Error(err)
		return res, err
	}

	// Run assertions
	err = e.run(res.RawResponse, res.RawRequest)
	if err != nil {
		e.test.Error(err)
	}

	return res, err
}
