package assert_json

import (
	"fmt"
	"strings"
	"testing"

	baloo "gopkg.in/h2non/baloo.v3"
	"github.com/mitchellh/mapstructure")

type UserAgent struct {
	Value string `mapstructure:"user-agent"`
}

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
			var result UserAgent
			err := mapstructure.Decode(data, &result)
			if err != nil {
				return err
			}
			if !strings.Contains(result.Value, "baloo") {
				return fmt.Errorf("bad user-agent: %s, %s", result.Value, data)
			}
			return nil
		}).
		Done()
}
