package assert

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/mitchellh/mapstructure"

	"github.com/nbio/st"
)

type UserAgent struct {
	Value string `mapstructure:"user-agent"`
}

type items map[string]string

func TestJSONString(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewBufferString(`{"foo":"bar"}`))
	res := &http.Response{Body: body}
	match := `{"foo": "bar"}`
	st.Expect(t, JSON(match)(res, nil), nil)
}

func TestJSONBuffer(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewBufferString(`{"foo":"bar"}`))
	res := &http.Response{Body: body}
	match := []byte(`{"foo": "bar"}`)
	st.Expect(t, JSON(match)(res, nil), nil)
}

type literalJSON string

func (v literalJSON) MarshalJSON() ([]byte, error) {
	return []byte(v), nil
}

func TestJSON(t *testing.T) {
	type genMap map[string]interface{}
	type strMap map[string]string
	testcases := []struct {
		name  string
		body  string
		match interface{}
	}{
		{
			name:  "generic map",
			body:  `{"foo": [{"bar":"baz"}]}`,
			match: genMap{"foo": []items{{"bar": "baz"}}},
		},
		{
			name:  "string to sring map",
			body:  `{"foo": "bar"}`,
			match: strMap{"foo": "bar"},
		},
		{
			name:  "array",
			body:  `["foo", 1.0]`,
			match: []interface{}{"foo", 1.0},
		},
		{
			name:  "custom type with json marshaller",
			body:  `"foo"`,
			match: literalJSON(`"foo"`),
		},
		{
			name: "keys order in map does not matter",
			body: `
        {
          "args": {},
          "headers": {
          	"Accept-Encoding": "gzip",
          	"Connection": "close",
          	"Foo": "Bar",
          	"Host": "httpbin.org",
          	"User-Agent": "baloo/2.0.0"
          },
          "origin": "0.0.0.0",
          "url": "http://httpbin.org/get"
        }`,
			match: literalJSON(`
				{
          "args": {},
          "headers": {
          	"Connection": "close",
          	"Accept-Encoding": "gzip",
          	"Foo": "Bar",
          	"Host": "httpbin.org",
          	"User-Agent": "baloo/2.0.0"
          },
          "origin": "0.0.0.0",
          "url": "http://httpbin.org/get"
        }`),
		},
	}
	for i, tc := range testcases {
		t.Run(tc.name, func(t *testing.T) {
			body := ioutil.NopCloser(bytes.NewBufferString(tc.body))
			res := &http.Response{Body: body}
			st.Expect(t, JSON(tc.match)(res, nil), nil, i)
		})
	}
}

func TestCompare(t *testing.T) {
	st.Expect(t, compare(map[string]interface{}{"a": "b"}, `{"a":"b"}`), nil)
	st.Expect(t, compare(map[string]interface{}{"a": 5.5}, `{"a":5.5}`), nil)
	st.Expect(t, compare(map[string]interface{}{"a": true}, `{"a": true}`), nil)
	st.Expect(t, compare(map[string]interface{}{"a": map[string]interface{}{"b": "c"}}, `{"a": {"b":"c"}}`), nil)
	st.Expect(t, compare(map[string]interface{}{"a": []float64{1.1, 2.2}}, `{"a":[1.1, 2.2]}`), nil)
	st.Expect(t, compare(map[string]interface{}{"a": 5.1}, `{"a": 5.1}`), nil)
}

func TestOnJSONCustomAssertion(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewBufferString(`{"user-agent":"baloo/v3"}`))
	res := &http.Response{Body: body}
	st.Expect(t, OnJSON(func(data interface{}) error {
		var result UserAgent
		err := mapstructure.Decode(data, &result)
		if err != nil {
			return err
		}
		if !strings.Contains(result.Value, "baloo") {
			return fmt.Errorf("bad user-agent: %s, %s", result.Value, data)
		}
		return nil
	})(res, nil), nil)
}

func TestOnJSONCustomAssertionBadResponse(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewBufferString(`{"user-agent":"toto"}`))
	res := &http.Response{Body: body}
	st.Expect(t, OnJSON(func(data interface{}) error {
		var result UserAgent
		err := mapstructure.Decode(data, &result)
		if err != nil {
			return err
		}
		if !strings.Contains(result.Value, "baloo") {
			return fmt.Errorf("bad user-agent: %s, %s", result.Value, data)
		}
		return nil
	})(res, nil), fmt.Errorf("bad user-agent: toto, map[user-agent:toto]"))
}
