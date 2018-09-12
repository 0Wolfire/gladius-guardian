package guardian

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Message  string      `json:"message"`
	Success  bool        `json:"success"`
	Error    string      `json:"error"`
	Response interface{} `json:"response"`
	Endpoint string      `json:"endpoint"`
}

// ErrorHandler - Default Error Handler
func ErrorHandler(w http.ResponseWriter, r *http.Request, m string, e error, statusCode int) {
	w.WriteHeader(statusCode)

	ResponseHandler(w, r, m, false, e, nil)
}

func ResponseHandler(w http.ResponseWriter, r *http.Request, m string, success bool, err error, res interface{}) {
	errorString := ""

	if err != nil {
		errorString = err.Error()
	}

	responseStruct := Response{
		Message:  m,
		Success:  success,
		Error:    errorString,
		Response: res,
		Endpoint: r.URL.String(),
	}

	enc := json.NewEncoder(w)
	enc.SetEscapeHTML(false) // So we can have an & come through in our URL's
	parseErr := enc.Encode(responseStruct)

	if parseErr != nil {
		ErrorHandler(w, r, "Could not parse response JSON", parseErr, http.StatusInternalServerError)
	}
}
