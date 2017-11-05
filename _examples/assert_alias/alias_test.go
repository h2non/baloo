package alias_assertion

import (
	"errors"
	"net/http"
	"testing"

	"gopkg.in/h2non/baloo.v3"
)

// test stores the HTTP testing client preconfigured
var test = baloo.New("http://httpbin.org")

func assert(res *http.Response, req *http.Request) error {
	if res.StatusCode >= 400 {
		return errors.New("Invalid server response (> 400)")
	}
	return nil
}

func init() {
	// Register assertion function at global level
	baloo.AddAssertFunc("test", assert)
}

func TestBalooClient(t *testing.T) {
	test.Post("/post").
		SetHeader("Foo", "Bar").
		JSON(map[string]string{"foo": "bar"}).
		Expect(t).
		Status(200).
		Type("json").
		Assert("test").
		Done()
}
