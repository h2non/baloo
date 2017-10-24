package assert

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/nbio/st"
)

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
