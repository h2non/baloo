package simple

import (
	"testing"

	"gopkg.in/h2non/baloo.v0"
)

// test stores the HTTP testing client preconfigured
var test = baloo.New("http://httpbin.org")

func TestBalooHeaders(t *testing.T) {
	test.Get("/get").
		SetHeader("Foo", "Bar").
		Expect(t).
		Status(200).
		Header("Server", "nginx").
		Header("Content-Type", "json").
		Type("json").
		Done()
}
