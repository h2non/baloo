package baloo

import (
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/nbio/st"
	"gopkg.in/h2non/gentleman.v1"
)

func TestExpect(t *testing.T) {
	req := &Request{Request: gentleman.NewRequest()}
	res := &http.Response{StatusCode: 200}
	exp := NewExpect(req)
	exp.Status(200)
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
