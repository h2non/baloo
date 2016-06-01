package assert_json

import (
	"testing"

	"gopkg.in/h2non/baloo.v0"
)

// test stores the HTTP testing client preconfigured
var test = baloo.New().URL("http://httpbin.org")

func TestBalooJSONAssertion(t *testing.T) {
	test.Get("/user-agent").
		SetHeader("Foo", "Bar").
		Expect(t).
		Status(200).
		Type("json").
		JSON(`{"user-agent":"baloo/0.1.0"}`).
		Done()
}
