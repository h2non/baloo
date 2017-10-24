package assert_body

import (
	"testing"

	"gopkg.in/h2non/baloo.v3"
)

// test stores the HTTP testing client preconfigured
var test = baloo.New("http://httpbin.org")

func TestBalooBodyAssertion(t *testing.T) {
	test.Get("/headers").
		SetHeader("Foo", "Bar").
		Expect(t).
		Status(200).
		Type("json").
		BodyMatchString(`"Foo"`).
		BodyMatchString(`"Bar"`).
		Done()
}
