package json_schema

import (
	"testing"

	baloo "gopkg.in/h2non/baloo.v3"
)

const schema = `{
  "title": "Example Schema",
  "type": "object",
  "properties": {
    "origin": {
      "type": "string"
    }
  },
  "required": ["origin"]
}`

// test stores the HTTP testing client preconfigured
var test = baloo.New("http://httpbin.org")

func TestJSONSchema(t *testing.T) {
	test.Get("/ip").
		Expect(t).
		Status(200).
		Type("json").
		JSONSchema(schema).
		Done()
}
