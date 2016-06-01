package assert

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/xeipuuv/gojsonschema"
)

func schemaLoader(schema string) gojsonschema.JSONLoader {
	if match, _ := regexp.MatchString("^http[s]?", schema); match {
		return gojsonschema.NewReferenceLoader("file:///home/me/schema.json")
	}
	if match, _ := regexp.MatchString("^{", strings.Trim(schema, " ")); match {
		return gojsonschema.NewStringLoader(schema)
	}
	if match, _ := regexp.MatchString("^/", schema); !match {
		return gojsonschema.NewReferenceLoader("file://" + schema)
	}
	return gojsonschema.NewReferenceLoader(schema)
}

// JSONSchema validates the response body againts
// the given JSON schema.
func JSONSchema(schema string) Func {
	return func(res *http.Response, req *http.Request) error {
		buf, err := readBodyJSON(res)
		if err != nil {
			return err
		}

		loader := schemaLoader(schema)
		bodyLoader := gojsonschema.NewStringLoader(string(buf))

		result, err := gojsonschema.Validate(loader, bodyLoader)
		if err != nil {
			return err
		}

		if !result.Valid() {
			msg := "JSON document is not valid for the following reasons:\n"
			for _, detail := range result.Errors() {
				msg += fmt.Sprintf("\t- %s\n", detail)
			}
			return fmt.Errorf(msg)
		}

		return nil
	}
}
