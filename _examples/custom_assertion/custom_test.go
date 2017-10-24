package custom_assertion

import (
	"errors"
	"net/http"
	"testing"

	baloo "gopkg.in/h2non/baloo.v3"
)

// test stores the HTTP testing client preconfigured
var test = baloo.New("http://httpbin.org")

func assert(res *http.Response, req *http.Request) error {
	if res.StatusCode >= 400 {
		return errors.New("Invalid server response (> 400)")
	}
	return nil
}

func TestBalooCustomAssertion(t *testing.T) {
	test.Post("/post").
		SetHeader("Foo", "Bar").
		JSON(map[string]string{"foo": "bar"}).
		Expect(t).
		Status(200).
		Type("json").
		AssertFunc(assert).
		Done()
}
