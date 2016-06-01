package assert

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/nbio/st"
)

func TestJSONSchemaSimple(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewBufferString(`{"foo":1}`))
	res := &http.Response{Body: body}
	match := `{"properties":{"bar":{},"foo":{}},"required":["foo"]}`
	st.Expect(t, JSONSchema(match)(res, nil), nil)
}

func TestJSONSchemaComplex(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewBufferString(`{"foo":"baz","bar":"baz", "age": 1}`))
	res := &http.Response{Body: body}
	match := `{
  "title": "Example Schema",
  "type": "object",
  "properties": {
    "foo": {
      "type": "string"
    },
    "bar": {
      "type": "string"
    },
    "age": {
      "description": "Age in years",
      "type": "integer",
      "minimum": 0
    }
  },
  "required": ["foo", "bar"]
}`
	st.Expect(t, JSONSchema(match)(res, nil), nil)
}

func TestJSONSchemaInvalid(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewBufferString(`{"foo":"baz","age": -1}`))
	res := &http.Response{Body: body}
	match := `{
  "title": "Example Schema",
  "type": "object",
  "properties": {
    "foo": {
      "type": "string"
    },
    "bar": {
      "type": "string"
    },
    "age": {
      "description": "Age in years",
      "type": "integer",
      "minimum": 0
    }
  },
  "required": ["foo", "bar"]
}`
	err := JSONSchema(match)(res, nil)
	st.Reject(t, err, nil)
	st.Expect(t, strings.Contains(err.Error(), "bar: bar is required"), true)
	st.Expect(t, strings.Contains(err.Error(), "age: Must be greater than or equal to 0"), true)
}
