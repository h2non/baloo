package assert

import (
	"fmt"
	"net/http"
	"regexp"

	"gopkg.in/h2non/gentleman.v2/plugins/bodytype"
)

// Type asserts the response Content-Type header.
// You can pass an alias as type.
// Supported alias are: json, html, xml, text, form, urlencoded.
func Type(kind string) Func {
	return func(res *http.Response, req *http.Request) error {
		value, ok := bodytype.Types[kind]
		if !ok {
			value = kind
		}

		header := res.Header.Get("Content-Type")
		if match, _ := regexp.MatchString(value, header); !match {
			return fmt.Errorf("Unexpected content type: '%s' should match '%s'", kind, header)
		}
		return nil
	}
}
