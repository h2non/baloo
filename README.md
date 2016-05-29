# baloo [![Build Status](https://travis-ci.org/h2non/baloo.png)](https://travis-ci.org/h2non/baloo) [![GitHub release](https://img.shields.io/badge/version-1.0.0-orange.svg?style=flat)](https://github.com/h2non/baloo/releases) [![GoDoc](https://godoc.org/github.com/h2non/baloo?status.svg)](https://godoc.org/github.com/h2non/baloo) [![Coverage Status](https://coveralls.io/repos/github/h2non/baloo/badge.svg?branch=master)](https://coveralls.io/github/h2non/baloo?branch=master) [![Go Report Card](https://goreportcard.com/badge/github.com/h2non/baloo)](https://goreportcard.com/report/github.com/h2non/baloo)

Expressive and versatile programmatic HTTP API testing made easy in [Go](http://golang.org).
Built-on-top of [gentleman](https://github.com/h2non/gentleman) HTTP client toolkit.

Take a look to the [examples](#examples) to get started.

## Features

- Full-featured HTTP expectation and built-in assertions.
- Declarative, expressive, fluent API.
- Full-featured HTTP client built on top of [gentleman](https://github.com/h2non/gentleman).
- Easy to configure and use.
- Convenient helpers and abstractions over Go's HTTP primitives.
- Built-in JSON, XML and multipart bodies matching.
- Extensible custom assertion functions.
- Hackable API.

## Installation

```bash
go get -u gopkg.in/h2non/baloo.v1
```

## API

See [godoc reference](https://godoc.org/github.com/h2non/baloo) for detailed API documentation.

### HTTP assertions

#### Status(code int)

#### StatusRange(start, end int)

#### AssertFunc(func (*gentleman.Response, *context.Context) error)

## Examples

See [examples](https://github.com/h2non/baloo/blob/master/_examples) directory for featured examples.

#### Simple request expectation

```go
package foo

import (
  "fmt"
  "testing"
  "gopkg.in/h2non/baloo.v1"
)

// test stores the HTTP testing client preconfigured
var test = baloo.New().URL("http://httpbin.org")

func TestBalooClient(t *testing.T) {
  test.Get("/foo").
    SetHeader("Foo", "Bar").
    JSON(map[string]string{"foo":"bar"}).
    Expect(t).
    Status(200).
    Header("Server", "apache").
    Type("json").
    JSON(map[string]string{"bar":"foo"}).
    Done()
}
```

#### Custom assertion functions

```go
package foo

import (
  "fmt"
  "testing"
  "gopkg.in/h2non/baloo.v1"
  "gopkg.in/h2non/gentleman.v1"
)

// test stores the HTTP testing client preconfigured
var test = baloo.New().URL("http://httpbin.org")

func TestBalooClient(t *testing.T) {
  test.Get("/foo").
    SetHeader("Foo", "Bar").
    JSON(map[string]string{"foo":"bar"}).
    Expect(t).
    Status(200).
    Type("json").
    AssertFunc(func (r *gentleman.Response, req *http.Request) error {
      if r.StatusCode >= 400 {
        return errors.New("Invalid server response (> 400)")
      }
      return nil
    }).
    Done()
}
```

## License 

MIT - Tomas Aparicio
