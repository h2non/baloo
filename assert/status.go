package assert

import (
	"fmt"

	"gopkg.in/h2non/gentleman.v1"
	"gopkg.in/h2non/gentleman.v1/context"
)

// Func represents the required interface for assertion functions.
type Func func(*gentleman.Response, *context.Context) error

// StatusEqual matches the status code
func StatusEqual(code int) Func {
	return func(res *gentleman.Response, ctx *context.Context) error {
		if res.StatusCode != code {
			return fmt.Errorf("Unexpected status code: %d != %d", res.StatusCode, code)
		}
		return nil
	}
}
