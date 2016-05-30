# baloo [![Build Status](https://travis-ci.org/h2non/baloo.png)](https://travis-ci.org/h2non/baloo) [![GitHub release](https://img.shields.io/badge/version-1.0.0-orange.svg?style=flat)](https://github.com/h2non/baloo/releases) [![GoDoc](https://godoc.org/github.com/h2non/baloo?status.svg)](https://godoc.org/github.com/h2non/baloo) [![Coverage Status](https://coveralls.io/repos/github/h2non/baloo/badge.svg?branch=master)](https://coveralls.io/github/h2non/baloo?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/h2non/baloo)](https://goreportcard.com/report/github.com/h2non/baloo)

Expressive and versatile End-to-End HTTP API testing made easy in [Go](http://golang.org) (golang).
Built on top of [gentleman](https://github.com/h2non/gentleman) HTTP client toolkit.

Take a look to the [examples](#examples) to get started.

**Note**: still beta.

## Features

- Versatile built-in expectations.
- Extensible custom expectations.
- Declarative, expressive, fluent API.
- Full-featured HTTP client built on top of [gentleman](https://github.com/h2non/gentleman).
- Intuitive and semantic HTTP client DSL.
- Easy to configure and use.
- Works with Go's `testing` package (more test engines might be added in the future).
- Convenient helpers and abstractions over Go's HTTP primitives.
- Built-in JSON, XML and multipart bodies definition and expectation.
- Middleware-oriented via gentleman's [middleware layer](https://github.com/h2non/gentleman#middleware).
- Extensible and hackable API.

## Upcoming features

- [ ] JSON expectations using deep comparison.
- [ ] JSON Schema body expectations.
- [ ] Cookies expectations.
- [ ] Response timing expectations.

## Installation

```bash
go get -u gopkg.in/h2non/baloo.v1
```

## API

See [godoc reference](https://godoc.org/github.com/h2non/baloo) for detailed API documentation.

### HTTP assertions

#### Status(code int)

Asserts the response HTTP status code to be equal.

#### StatusRange(start, end int)

Asserts the response HTTP status to be within the given numeric range.

#### StatusOk(start, end int)

Asserts the response HTTP status to be a valid server response (>= 200 && < 400).

#### StatusError(start, end int)

Asserts the response HTTP status to be a valid clint/server error response (>= 400 && < 600).

#### StatusServerError(start, end int)

Asserts the response HTTP status to be a valid server error response (>= 500 && < 600).

#### StatusClientError(start, end int)

Asserts the response HTTP status to be a valid client error response (>= 400 && < 500).

#### Type(kind string)

Asserts the `Content-Type` header. MIME type aliases can be used as `kind` argument.

Supported aliases: `json`, `xml`, `html`, `form`, `text` and `urlencoded`.

#### Header(key, value string)

Asserts a response header field value matches.

Regular expressions can be used as value to perform the specific assertions.

#### HeaderEquals(key, value string)

Asserts a response header field with the given value.

#### HeaderPresent(key string)

Asserts if a header field is present in the response.

#### BodyMatchString(pattern string)

Asserts a response body matching a string expression.

Regular expressions can be used as value to perform the specific assertions.

#### BodyLength(length int)

Asserts a response body length.

#### AssertFunc(func (*gentleman.Response, *context.Context) error)

Adds a new custom assertion function who should return an 
detailed error in case that the assertion fails.

## Examples

See [examples](https://github.com/h2non/baloo/blob/master/_examples) directory for featured examples.

#### Simple request expectation

```go
package simple

import (
  "testing"

  "gopkg.in/h2non/baloo.v0"
)

// test stores the HTTP testing client preconfigured
var test = baloo.New().URL("http://httpbin.org")

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

  "gopkg.in/h2non/baloo.v0"
)

// test stores the HTTP testing client preconfigured
var test = baloo.New().URL("http://httpbin.org")

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

## License 

MIT - Tomas Aparicio
