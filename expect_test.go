package baloo

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nbio/st"
	"gopkg.in/h2non/gentleman.v2"
)

func TestExpect(t *testing.T) {
	req := &Request{Request: gentleman.NewRequest()}
	res := &http.Response{StatusCode: 200}
	exp := NewExpect(req)
	exp.Status(200)
	st.Expect(t, exp.run(res, nil), nil)
}

func TestExpectBindTest(t *testing.T) {
	req := &Request{Request: gentleman.NewRequest()}
	exp := NewExpect(req)
	exp.BindTest(t)
	st.Expect(t, exp.test, t)
}

func TestExpectStatusOk(t *testing.T) {
	req := &Request{Request: gentleman.NewRequest()}
	res := &http.Response{StatusCode: 200}
	exp := NewExpect(req)
	exp.StatusOk()
	st.Expect(t, exp.run(res, nil), nil)
}

func TestExpectStatusError(t *testing.T) {
	req := &Request{Request: gentleman.NewRequest()}
	res := &http.Response{StatusCode: 400}
	exp := NewExpect(req)
	exp.StatusError()
	st.Expect(t, exp.run(res, nil), nil)
}

func TestExpectServerError(t *testing.T) {
	req := &Request{Request: gentleman.NewRequest()}
	res := &http.Response{StatusCode: 500}
	exp := NewExpect(req)
	exp.StatusServerError()
	st.Expect(t, exp.run(res, nil), nil)
}

func TestExpectStatusClientError(t *testing.T) {
	req := &Request{Request: gentleman.NewRequest()}
	res := &http.Response{StatusCode: 404}
	exp := NewExpect(req)
	exp.StatusClientError()
	st.Expect(t, exp.run(res, nil), nil)
}

func TestExpectType(t *testing.T) {
	req := &Request{Request: gentleman.NewRequest()}
	headers := http.Header{"Content-Type": []string{"application/json"}}
	res := &http.Response{StatusCode: 404, Header: headers}
	exp := NewExpect(req)
	exp.Type("json")
	st.Expect(t, exp.run(res, nil), nil)
}

func TestExpectHeader(t *testing.T) {
	req := &Request{Request: gentleman.NewRequest()}
	headers := http.Header{"Foo": []string{"bar"}}
	res := &http.Response{StatusCode: 404, Header: headers}
	exp := NewExpect(req)
	exp.Header("foo", "bar")
	st.Expect(t, exp.run(res, nil), nil)
}

func TestExpectHeaderEquals(t *testing.T) {
	req := &Request{Request: gentleman.NewRequest()}
	headers := http.Header{"Foo": []string{"bar"}}
	res := &http.Response{StatusCode: 404, Header: headers}
	exp := NewExpect(req)
	exp.HeaderEquals("foo", "bar")
	st.Expect(t, exp.run(res, nil), nil)
}

func TestExpectHeaderNotEquals(t *testing.T) {
	req := &Request{Request: gentleman.NewRequest()}
	headers := http.Header{"Foo": []string{"bar"}}
	res := &http.Response{StatusCode: 404, Header: headers}
	exp := NewExpect(req)
	exp.HeaderNotEquals("foo", "foo")
	st.Expect(t, exp.run(res, nil), nil)
}

func TestExpectHeaderPresent(t *testing.T) {
	req := &Request{Request: gentleman.NewRequest()}
	headers := http.Header{"Foo": []string{"bar"}}
	res := &http.Response{StatusCode: 404, Header: headers}
	exp := NewExpect(req)
	exp.HeaderPresent("foo")
	st.Expect(t, exp.run(res, nil), nil)
}

func TestExpectHeaderNotPresent(t *testing.T) {
	req := &Request{Request: gentleman.NewRequest()}
	headers := http.Header{"Foo": []string{"bar"}}
	res := &http.Response{StatusCode: 404, Header: headers}
	exp := NewExpect(req)
	exp.HeaderNotPresent("bar")
	st.Expect(t, exp.run(res, nil), nil)
}

func assertStatus(res *http.Response, req *http.Request) error {
	if res.StatusCode >= 400 {
		return errors.New("Invalid server response (> 400)")
	}
	return nil
}

func TestRealRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; encoding=utf8")
		fmt.Fprintln(w, "Hello, "+r.Header.Get("Foo"))
	}))
	defer ts.Close()

	cli := New(ts.URL)
	cli.Get("/foo").
		SetHeader("Foo", "Bar").
		Expect(t).
		Status(200).
		Type("json").
		AssertFunc(assertStatus).
		Done()
}

func TestGlobalAssertFunc(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json; encoding=utf8")
		fmt.Fprintln(w, "Hello, "+r.Header.Get("Foo"))
	}))
	defer ts.Close()

	AddAssertFunc("foo", assertStatus)
	defer FlushAssertFuncs()

	cli := New(ts.URL)
	cli.Get("/foo").
		SetHeader("Foo", "Bar").
		Expect(t).
		Status(200).
		Type("json").
		Assert("foo").
		Done()
}
