package assert

import (
	"net/http"
	"testing"

	"github.com/nbio/st"
)

func TestHeader(t *testing.T) {
	headers := http.Header{
		"Content-Type": []string{"application/json; encoding=utf8"},
		"Server":       []string{"nginx"},
	}
	res := &http.Response{Header: headers}

	st.Expect(t, Header("Content-Type", "json")(res, nil), nil)
	st.Expect(t, Header("Content-Type", "application/json")(res, nil), nil)
	st.Expect(t, Header("Content-Type", "utf8")(res, nil), nil)
	st.Expect(t, Header("Content-Type", "^application/json")(res, nil), nil)
	st.Expect(t, Header("Content-Type", "^application/json; encoding=utf8$")(res, nil), nil)

	st.Reject(t, Header("Content-Type", "xml")(res, nil), nil)
	st.Reject(t, Header("Content-Type", "foo")(res, nil), nil)
}

func TestHeaderEquals(t *testing.T) {
	headers := http.Header{
		"Content-Type": []string{"application/json"},
		"Server":       []string{"nginx"},
	}
	res := &http.Response{Header: headers}

	st.Expect(t, HeaderEquals("Content-Type", "application/json")(res, nil), nil)
	st.Expect(t, HeaderEquals("server", "nginx")(res, nil), nil)
	st.Reject(t, HeaderEquals("Content-Type", "application/foo")(res, nil), nil)
	st.Reject(t, HeaderEquals("Content-Type", "foo")(res, nil), nil)
	st.Reject(t, HeaderEquals("server", "foo")(res, nil), nil)
}

func TestHeaderNotEquals(t *testing.T) {
	headers := http.Header{
		"Content-Type": []string{"application/json"},
		"Server":       []string{"nginx"},
	}
	res := &http.Response{Header: headers}

	st.Reject(t, HeaderNotEquals("Content-Type", "application/json")(res, nil), nil)
	st.Reject(t, HeaderNotEquals("server", "nginx")(res, nil), nil)
	st.Expect(t, HeaderNotEquals("Content-Type", "application/foo")(res, nil), nil)
	st.Expect(t, HeaderNotEquals("Content-Type", "foo")(res, nil), nil)
	st.Expect(t, HeaderNotEquals("server", "foo")(res, nil), nil)
}

func TestHeaderPresent(t *testing.T) {
	headers := http.Header{
		"Content-Type": []string{"application/json"},
		"Server":       []string{"nginx"},
	}
	res := &http.Response{Header: headers}

	st.Expect(t, HeaderPresent("Content-Type")(res, nil), nil)
	st.Expect(t, HeaderPresent("server")(res, nil), nil)
	st.Reject(t, HeaderPresent("Accept")(res, nil), nil)
	st.Reject(t, HeaderPresent("Cookie")(res, nil), nil)
}

func TestHeaderNotPresent(t *testing.T) {
	headers := http.Header{
		"Content-Type": []string{"application/json"},
		"Server":       []string{"nginx"},
	}
	res := &http.Response{Header: headers}

	st.Reject(t, HeaderNotPresent("Content-Type")(res, nil), nil)
	st.Reject(t, HeaderNotPresent("server")(res, nil), nil)
	st.Expect(t, HeaderNotPresent("Accept")(res, nil), nil)
	st.Expect(t, HeaderNotPresent("Cookie")(res, nil), nil)
}

func TestRedirectTo(t *testing.T) {
	headers := http.Header{"Location": []string{"http://foo/bar"}}
	res := &http.Response{Header: headers}

	st.Expect(t, RedirectTo("http://foo/bar")(res, nil), nil)
	st.Expect(t, RedirectTo("foo/bar")(res, nil), nil)
	st.Expect(t, RedirectTo("^http://foo/bar$")(res, nil), nil)
	st.Reject(t, RedirectTo("^http://foo$")(res, nil), nil)
	st.Reject(t, RedirectTo("baz")(res, nil), nil)
}
