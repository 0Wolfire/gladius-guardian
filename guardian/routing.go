package guardian

import (
	"net/http"
)

func GetServicesHandler(gg *GladiusGuardian) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ResponseHandler(w, r, "Got service status", true, nil, gg.GetServicesStatus())
	}
}
