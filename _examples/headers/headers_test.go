package simple

import (
	"testing"

	baloo "gopkg.in/h2non/baloo.v3"
)

// test stores the HTTP testing client preconfigured
var test = baloo.New("http://httpbin.org")

func TestBalooHeaders(t *testing.T) {
	test.Get("/get").
		SetHeader("Foo", "Bar").
		Expect(t).
		Status(200).
		Header("Server", "meinheld").
		Header("Content-Type", "json").
		Type("json").
		Done()
}
