package assert

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/nbio/st"
)

func testJSONSchema(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewBufferString(`{"foo":"bar"}`))
	res := &http.Response{Body: body}
	match := `{"foo": "bar"}`
	st.Expect(t, JSONEquals(match)(res, nil), nil)
}
