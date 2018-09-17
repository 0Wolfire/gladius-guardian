package guardian

import (
	"errors"
	"io/ioutil"
	"net/http"

	"github.com/buger/jsonparser"
)

func getJSONFields(w http.ResponseWriter, r *http.Request, fields ...string) (map[string][]byte, error) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		ErrorHandler(w, r, "", err, http.StatusBadRequest)
		return nil, errors.New("error decoding body")
	}

	vals := make(map[string][]byte)

	for _, field := range fields {
		valBytes, _, _, err := jsonparser.Get(body, field)
		if err == nil {
			vals[field] = valBytes
		}
	}

	return vals, nil
}
