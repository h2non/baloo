package assert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
)

func unmarshal(buf []byte) (interface{}, error) {
	var data interface{}
	if err := json.Unmarshal(buf, &data); err != nil {
		return data, fmt.Errorf("failed to Unmarshal: %s", err)
	}
	return data, nil
}

func unmarshalBody(res *http.Response) (interface{}, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	if len(body) == 0 {
		return nil, nil
	}
	return unmarshal(body)
}

func marshal(data interface{}) ([]byte, error) {
	return json.Marshal(data)
}

func readBodyJSON(res *http.Response) ([]byte, error) {
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}

	// Re-fill body reader stream after reading it
	res.Body = ioutil.NopCloser(bytes.NewBuffer(body))
	return body, err
}

func compare(body interface{}, data interface{}) error {
	var err error
	var matchBytes []byte

	// taking pointer to json.RawMessage due to regression in go 1.7 (https://github.com/golang/go/issues/14493)
	var matchData interface{}
	switch data := data.(type) {
	case json.Marshaler:
		matchData = data
	case string:
		v := json.RawMessage([]byte(data))
		matchData = &v
	case []byte:
		v := json.RawMessage(data)
		matchData = &v
	default:
		matchData = data
	}
	matchBytes, err = marshal(matchData)

	if err != nil {
		return err
	}

	// Assert via string
	bodyBytes, err := marshal(body)
	if err != nil {
		return err
	}

	bodyValue, err := unmarshal(bodyBytes)
	if err != nil {
		return err
	}
	matchValue, err := unmarshal(matchBytes)
	if err != nil {
		return err
	}

	// Compare values so order of keys in maps does not influence the result
	if !reflect.DeepEqual(bodyValue, matchValue) {
		return fmt.Errorf("failed due to JSON mismatch:\n\thave: %#v\n\twant: %#v", string(bodyBytes), string(matchBytes))
	}

	return nil
}

// JSON deeply and strictly compares the JSON
// response body with the given JSON structure.
func JSON(data interface{}) Func {
	return func(res *http.Response, req *http.Request) error {
		// Read and unmarshal response body as JSON
		body, err := unmarshalBody(res)
		if err != nil {
			return err
		}

		return compare(body, data)
	}
}
