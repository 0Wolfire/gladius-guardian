package guardian

import (
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	ResponseHandler(w, r, "There's nothing here, check our API docs at https://github.com/gladiusio/gladius-guardian", true, nil, "")
}

func GetServicesHandler(gg *GladiusGuardian) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		ResponseHandler(w, r, "Got service status", true, nil, gg.GetServicesStatus())
	}
}

func StartServiceHandler(gg *GladiusGuardian) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get service name, optionally location and environment variables
		ResponseHandler(w, r, "Got service status", true, nil, gg.GetServicesStatus())
	}
}

func StopServiceHandler(gg *GladiusGuardian) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get service name to stop it
		ResponseHandler(w, r, "Got service status", true, nil, gg.GetServicesStatus())
	}
}

func StopAllServiceHandler(gg *GladiusGuardian) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get service name to stop it
		ResponseHandler(w, r, "Got service status", true, nil, gg.GetServicesStatus())
	}
}

func SetStartTimeoutHandler(gg *GladiusGuardian) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get service name to stop it
		ResponseHandler(w, r, "Got service status", true, nil, gg.GetServicesStatus())
	}
}

func GetLogsHandler(gg *GladiusGuardian) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get service name to stop it
		ResponseHandler(w, r, "Not implemented", true, nil, gg.GetServicesStatus())
	}
}
