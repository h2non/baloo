package assert

import (
	"net/http"
	"testing"

	"github.com/nbio/st"
)

func TestType(t *testing.T) {
	headers := http.Header{"Content-Type": []string{"application/json; encoding=utf8"}}
	res := &http.Response{Header: headers}

	st.Expect(t, Type("json")(res, nil), nil)
	st.Expect(t, Type("application/json")(res, nil), nil)
	st.Reject(t, Type("xml")(res, nil), nil)
	st.Reject(t, Type("foo")(res, nil), nil)
}
