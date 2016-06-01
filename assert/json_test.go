package assert

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/nbio/st"
)

type items map[string]string

func TestJSONString(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewBufferString(`{"foo":"bar"}`))
	res := &http.Response{Body: body}
	match := `{"foo": "bar"}`
	st.Expect(t, JSON(match)(res, nil), nil)
}

func TestJSONBuffer(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewBufferString(`{"foo":"bar"}`))
	res := &http.Response{Body: body}
	match := []byte(`{"foo": "bar"}`)
	st.Expect(t, JSON(match)(res, nil), nil)
}

func TestJSONMap(t *testing.T) {
	body := ioutil.NopCloser(bytes.NewBufferString(`{"foo": [{"bar":"baz"}]}`))
	res := &http.Response{Body: body}
	list := []items{{"bar": "baz"}}
	match := map[string]interface{}{"foo": list}
	st.Expect(t, JSON(match)(res, nil), nil)
}

func TestCompare(t *testing.T) {
	st.Expect(t, compare(map[string]interface{}{"a": "b"}, `{"a":"b"}`), nil)
	st.Expect(t, compare(map[string]interface{}{"a": 5.5}, `{"a":5.5}`), nil)
	st.Expect(t, compare(map[string]interface{}{"a": true}, `{"a": true}`), nil)
	st.Expect(t, compare(map[string]interface{}{"a": map[string]interface{}{"b": "c"}}, `{"a": {"b":"c"}}`), nil)
	st.Expect(t, compare(map[string]interface{}{"a": []float64{1.1, 2.2}}, `{"a":[1.1, 2.2]}`), nil)
	st.Expect(t, compare(map[string]interface{}{"a": 5.1}, `{"a": 5.1}`), nil)
}
