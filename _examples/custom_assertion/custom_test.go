package custom_assertion

import (
	"fmt"
	"testing"

	"gopkg.in/h2non/baloo.v1"
	"gopkg.in/h2non/gentleman.v1"
)

// test stores the HTTP testing client preconfigured
var test = baloo.New().URL("http://httpbin.org")

func assert(r *gentleman.Response, req *http.Request) error {
	if r.StatusCode >= 400 {
		return errors.New("Invalid server response (> 400)")
	}
	return nil
}

func TestBalooClient(t *testing.T) {
	test.Get("/foo").
		SetHeader("Foo", "Bar").
		JSON(map[string]string{"foo": "bar"}).
		Expect(t).
		Status(200).
		Type("json").
		AssertFunc(assert).
		Done()
}
