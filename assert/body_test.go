package assert

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/nbio/st"
)

func TestBodyMatchString(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewBufferString("hello world"))
	res := &http.Response{Body: body}
	st.Expect(t, BodyMatchString("hello")(res, nil), nil)
	st.Expect(t, BodyMatchString("^hello world$")(res, nil), nil)
	st.Expect(t, BodyMatchString("world$")(res, nil), nil)
	st.Expect(t, BodyMatchString("he[a-z]+")(res, nil), nil)
	st.Reject(t, BodyMatchString("foo")(res, nil), nil)
	st.Reject(t, BodyMatchString("bar")(res, nil), nil)
}

func TestBodyLength(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewBufferString("hello world"))
	res := &http.Response{Body: body}
	st.Expect(t, BodyLength(11)(res, nil), nil)
	st.Reject(t, BodyLength(10)(res, nil), nil)
	st.Reject(t, BodyLength(0)(res, nil), nil)
}

func TestBodyEquals(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewBufferString("hello world"))
	res := &http.Response{Body: body}
	st.Expect(t, BodyEquals("hello world")(res, nil), nil)
	st.Reject(t, BodyEquals("hello")(res, nil), nil)
	st.Reject(t, BodyEquals("world")(res, nil), nil)
	st.Reject(t, BodyEquals("foo")(res, nil), nil)
	st.Reject(t, BodyEquals("")(res, nil), nil)

	body = ioutil.NopCloser(bytes.NewBufferString("hello world\n"))
	res = &http.Response{Body: body}
	st.Expect(t, BodyEquals("hello world")(res, nil), nil)
	st.Reject(t, BodyEquals("hello world\n")(res, nil), nil)
	st.Reject(t, BodyEquals("hello")(res, nil), nil)
	st.Reject(t, BodyEquals("foo")(res, nil), nil)
	st.Reject(t, BodyEquals("")(res, nil), nil)
}
