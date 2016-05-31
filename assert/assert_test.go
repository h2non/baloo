package assert

import (
	"errors"
	"net/http"
	"testing"

	"github.com/nbio/st"
)

func TestFuncInterface(t *testing.T) {
	fn := func(res *http.Response, req *http.Request) error {
		return errors.New("foo error")
	}

	receiver := func(fn Func) bool {
		return true
	}

	st.Expect(t, receiver(fn), true)
}
