package assert_json

import (
	"testing"

	baloo "gopkg.in/h2non/baloo.v2"
)

// test stores the HTTP testing client preconfigured
var test = baloo.New("http://httpbin.org")

func TestBalooJSONAssertion(t *testing.T) {
	test.Get("/user-agent").
		SetHeader("Foo", "Bar").
		Expect(t).
		Status(200).
		Type("json").
		JSON(`{"user-agent":"baloo/` + baloo.Version + `"}`).
		Done()
}

func TestBalooJSONCustomAssertion(t *testing.T) {
	test.Get("/user-agent").
		SetHeader("Foo", "Bar").
		Expect(t).
		Status(200).
		Type("json").
		JSON(`{"user-agent":"baloo/` + baloo.Version + `"}`).
		VerifyJSON(func(data map[string]interface{}) error {
			// check your json response here
			return nil
		}).
		Done()
}
