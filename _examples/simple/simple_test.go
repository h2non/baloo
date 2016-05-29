package simple

import (
	"fmt"
	"testing"

	"gopkg.in/h2non/baloo.v1"
)

// test stores the HTTP testing client preconfigured
var test = baloo.New().URL("http://httpbin.org")

func TestBalooSimple(t *testing.T) {
	test.Get("/foo").
		SetHeader("Foo", "Bar").
		JSON(map[string]string{"foo": "bar"}).
		Expect(t).
		Status(200).
		Header("Server", "apache").
		Type("json").
		JSON(map[string]string{"bar": "foo"}).
		Done()
}
