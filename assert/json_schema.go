package assert

import (
	"fmt"
	"net/http"

	"github.com/xeipuuv/gojsonschema"
)

// JSONSchema validates the response body againts
// the given JSON schema.
func JSONSchema(schema string) Func {
	return func(res *http.Response, req *http.Request) error {
		buf, err := readBodyJSON(res)
		if err != nil {
			return err
		}

		loader := gojsonschema.NewStringLoader(schema)
		bodyLoader := gojsonschema.NewStringLoader(string(buf))

		result, err := gojsonschema.Validate(loader, bodyLoader)
		if err != nil {
			return err
		}

		if !result.Valid() {
			msg := "JSON document is not valid:\n"
			for _, detail := range result.Errors() {
				msg += fmt.Sprintf("- %s\n", detail)
			}
			return fmt.Errorf(msg)
		}

		return nil
	}
}
