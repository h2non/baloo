package baloo

import (
	"fmt"
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

// Status asserts the status code.
func (e *Expect) Status(code int) *Expect {
	e.AssertFunc(assert.StatusEqual(code))
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
		err = fmt.Errorf("[baloo] request error: %s", err)
		e.test.Error(err)
		return err
	}

	for _, assertion := range e.assertions {
		err = assertion(res, e.request.Request.Context)
		if err != nil {
			e.test.Error(err)
			break
		}
	}

	return err
}
