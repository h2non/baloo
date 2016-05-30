package assert

import (
	"net/http"
	"testing"

	"github.com/nbio/st"
)

func TestStatusEqual(t *testing.T) {
	st.Expect(t, assert(StatusEqual(200), 200), nil)
	st.Reject(t, assert(StatusEqual(400), 200), nil)
}

func TestStatusRange(t *testing.T) {
	st.Expect(t, assert(StatusRange(200, 300), 204), nil)
	st.Reject(t, assert(StatusRange(205, 300), 200), nil)
}

func TestStatusOk(t *testing.T) {
	st.Expect(t, assert(StatusOk(), 200), nil)
	st.Expect(t, assert(StatusOk(), 300), nil)
	st.Reject(t, assert(StatusOk(), 400), nil)
}

func TestStatusError(t *testing.T) {
	st.Expect(t, assert(StatusError(), 500), nil)
	st.Expect(t, assert(StatusError(), 502), nil)
	st.Expect(t, assert(StatusError(), 400), nil)
	st.Expect(t, assert(StatusError(), 405), nil)
	st.Reject(t, assert(StatusError(), 200), nil)
	st.Reject(t, assert(StatusError(), 0), nil)
}

func TestStatusServerError(t *testing.T) {
	st.Expect(t, assert(StatusServerError(), 500), nil)
	st.Expect(t, assert(StatusServerError(), 502), nil)
	st.Reject(t, assert(StatusServerError(), 400), nil)
	st.Reject(t, assert(StatusServerError(), 200), nil)
	st.Reject(t, assert(StatusServerError(), 0), nil)
}

func TestStatusClientError(t *testing.T) {
	st.Expect(t, assert(StatusClientError(), 400), nil)
	st.Expect(t, assert(StatusClientError(), 415), nil)
	st.Expect(t, assert(StatusClientError(), 403), nil)
	st.Reject(t, assert(StatusClientError(), 200), nil)
	st.Reject(t, assert(StatusClientError(), 500), nil)
	st.Reject(t, assert(StatusClientError(), 302), nil)
}

func assert(fn Func, status int) error {
	res := &http.Response{StatusCode: status}
	return fn(res, &http.Request{})
}
