package assert

import (
	"testing"

	"github.com/nbio/st"
	"gopkg.in/h2non/gentleman.v1"
)

func TestStatusEqual(t *testing.T) {
	ctx := gentleman.NewContext()
	res := &gentleman.Response{StatusCode: 200}

	st.Expect(t, StatusEqual(200)(res, ctx), nil)
}
