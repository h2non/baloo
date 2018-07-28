# baloo [![Build Status](https://travis-ci.org/h2non/baloo.svg)](https://travis-ci.org/h2non/baloo) [![GitHub release](https://img.shields.io/badge/version-3.0-orange.svg?style=flat)](https://github.com/h2non/baloo/releases) [![GoDoc](https://godoc.org/github.com/h2non/baloo?status.svg)](https://godoc.org/github.com/h2non/baloo) [![Coverage Status](https://coveralls.io/repos/github/h2non/baloo/badge.svg?branch=master)](https://coveralls.io/github/h2non/baloo?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/h2non/baloo)](https://goreportcard.com/report/github.com/h2non/baloo)

Expressive and versatile end-to-end HTTP API testing made easy in [Go](http://golang.org) (golang), built on top of [gentleman](https://github.com/h2non/gentleman) HTTP client toolkit.

Take a look to the [examples](#examples) to get started.

## Features

- Versatile built-in expectations.
- Extensible custom expectations.
- Declarative, expressive, fluent API.
- Response body matching and strict equality expectations.
- Deep JSON comparison.
- JSON Schema validation.
- Full-featured HTTP client built on top of [gentleman](https://github.com/h2non/gentleman) toolkit.
- Intuitive and semantic HTTP client DSL.
- Easy to configure and use.
- Composable chainable assertions.
- Works with Go's `testing` package (more test engines might be added in the future).
- Convenient helpers and abstractions over Go's HTTP primitives.
- Middleware-oriented via gentleman's [middleware layer](https://github.com/h2non/gentleman#middleware).
- Extensible and hackable API.

## Versions

- [v3](https://github.com/h2non/baloo) - Latest stable version with better JSON assertion. Uses `gentleman@v2`. Recommended.
- [v2](https://github.com/h2non/baloo/tree/v2.0.0) - Stable version. Uses `gentleman@v2`.
- [v1](https://github.com/h2non/baloo/tree/v1) - First version. Stable. Uses `gentleman@v1`. Actively maintained.

## Installation

```bash
go get -u gopkg.in/h2non/baloo.v3
```

## Requirements

- Go 1.7+

## Examples

See [examples](https://github.com/h2non/baloo/blob/master/_examples) directory for featured examples.

#### Simple request expectation

```go
package simple

import (
  "testing"

  "gopkg.in/h2non/baloo.v3"
)

// test stores the HTTP testing client preconfigured
var test = baloo.New("http://httpbin.org")

func TestBalooSimple(t *testing.T) {
  test.Get("/get").
    SetHeader("Foo", "Bar").
    Expect(t).
    Status(200).
    Header("Server", "apache").
    Type("json").
    JSON(map[string]string{"bar": "foo"}).
    Done()
}
```

#### Custom assertion function

```go
package custom_assertion

import (
  "errors"
  "net/http"
  "testing"

  "gopkg.in/h2non/baloo.v3"
)

// test stores the HTTP testing client preconfigured
var test = baloo.New("http://httpbin.org")

// assert implements an assertion function with custom validation logic.
// If the assertion fails it should return an error.
func assert(res *http.Response, req *http.Request) error {
  if res.StatusCode >= 400 {
    return errors.New("Invalid server response (> 400)")
  }
  return nil
}

func TestBalooClient(t *testing.T) {
  test.Post("/post").
    SetHeader("Foo", "Bar").
    JSON(map[string]string{"foo": "bar"}).
    Expect(t).
    Status(200).
    Type("json").
    AssertFunc(assert).
    Done()
}
```

#### JSON Schema assertion

```go
package json_schema

import (
  "testing"

  "gopkg.in/h2non/baloo.v3"
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
```

#### Custom global assertion by alias

```go
package alias_assertion

import (
  "errors"
  "net/http"
  "testing"

  "gopkg.in/h2non/baloo.v3"
)

// test stores the HTTP testing client preconfigured
var test = baloo.New("http://httpbin.org")

func assert(res *http.Response, req *http.Request) error {
  if res.StatusCode >= 400 {
    return errors.New("Invalid server response (> 400)")
  }
  return nil
}

func init() {
  // Register assertion function at global level
  baloo.AddAssertFunc("test", assert)
}

func TestBalooClient(t *testing.T) {
  test.Post("/post").
    SetHeader("Foo", "Bar").
    JSON(map[string]string{"foo": "bar"}).
    Expect(t).
    Status(200).
    Type("json").
    Assert("test").
    Done()
}
```

## API

See [godoc reference](https://godoc.org/github.com/h2non/baloo) for detailed API documentation.

### HTTP assertions

#### Status(code int)

Asserts the response HTTP status code to be equal.

#### StatusRange(start, end int)

Asserts the response HTTP status to be within the given numeric range.

#### StatusOk()

Asserts the response HTTP status to be a valid server response (>= 200 && < 400).

#### StatusError()

Asserts the response HTTP status to be a valid clint/server error response (>= 400 && < 600).

#### StatusServerError()

Asserts the response HTTP status to be a valid server error response (>= 500 && < 600).

#### StatusClientError()

Asserts the response HTTP status to be a valid client error response (>= 400 && < 500).

#### Type(kind string)

Asserts the `Content-Type` header. MIME type aliases can be used as `kind` argument.

Supported aliases: `json`, `xml`, `html`, `form`, `text` and `urlencoded`.

#### Header(key, value string)

Asserts a response header field value matches.

Regular expressions can be used as value to perform the specific assertions.

#### HeaderEquals(key, value string)

Asserts a response header field with the given value.

#### HeaderNotEquals(key, value string)

Asserts that a response header field is not equal to the given value.

#### HeaderPresent(key string)

Asserts if a header field is present in the response.

#### HeaderNotPresent(key string)

Asserts if a header field is not present in the response.

#### BodyEquals(value string)

Asserts a response body as string using strict comparison.

Regular expressions can be used as value to perform the specific assertions.

#### BodyMatchString(pattern string)

Asserts a response body matching a string expression.

Regular expressions can be used as value to perform the specific assertions.

#### BodyLength(length int)

Asserts the response body length.

#### JSON(match interface{})

Asserts the response body with the given JSON struct.

#### JSONSchema(schema string)

Asserts the response body againts the given JSON schema definition.

`data` argument can be a `string` containing the JSON schema, a file path
or an URL pointing to the JSON schema definition.

#### Assert(alias string)

Assert adds a new assertion function by alias name.

Assertion function must be previosly registered
via baloo.AddAssertFunc("alias", function).

See [an example here](#custom-global-assertion-by-alias).

#### AssertFunc(func (*http.Response, *http.Request) error)

Adds a new custom assertion function who should return an
detailed error in case that the assertion fails.

## Development

Clone this repository:
```bash
git clone https://github.com/h2non/baloo.git && cd baloo
```

Install dependencies:
```bash
go get -u ./...
```

Run tests:
```bash
go test ./...
```

Lint code:
```bash
go test ./...
```

Run example:
```bash
go test ./_examples/simple/simple_test.go
```

## License

MIT - Tomas Aparicio
