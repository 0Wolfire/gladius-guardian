package guardian

import (
	"net/http"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
	"github.com/spf13/viper"
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
		// Get service name, optionally environment variables
		vals, err := getJSONFields(w, r, "service_name", "environment_vars")
		if err != nil {
			ErrorHandler(w, r, "Couldn't parse body", err, http.StatusBadRequest)
			return
		}
		if _, ok := vals["service_name"]; !ok {
			ErrorHandler(w, r, "Need 'service_name' in request", err, http.StatusBadRequest)
			return
		}

		environmentVars := make([]string, 0)

		// Defaults will be used if empty, so only specify if we have some to add
		if envBytes, ok := vals["environment_vars"]; ok {
			// Add our defaults just in case, they can be overriden if they are redefined
			environmentVars = append(environmentVars, viper.GetStringSlice("DefaultEnvironment")...)

			// Add the desired environment vars
			jsonparser.ArrayEach(envBytes, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				environmentVars = append(environmentVars, string(value))
			})
		}

		err = gg.StartService(string(vals["service_name"]), environmentVars)
		if err != nil {
			ErrorHandler(w, r, "Error starting service", err, http.StatusBadRequest)
			return
		}

		ResponseHandler(w, r, "Attempted to start service", true, nil, gg.GetServicesStatus())
	}
}

func StopServiceHandler(gg *GladiusGuardian) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get service name to stop it
		vals, err := getJSONFields(w, r, "service_name")
		if err != nil {
			ErrorHandler(w, r, "Couldn't parse body", err, http.StatusBadRequest)
			return
		}
		if _, ok := vals["service_name"]; !ok {
			ErrorHandler(w, r, "Need 'service_name' in request", err, http.StatusBadRequest)
			return
		}

		err = gg.StopService(string(vals["service_name"]))
		if err != nil {
			ErrorHandler(w, r, "Error stoping service", err, http.StatusBadRequest)
			return
		}

		ResponseHandler(w, r, "Stopped Service", true, nil, gg.GetServicesStatus())
	}
}

func StopAllServiceHandler(gg *GladiusGuardian) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		err := gg.StopAll()
		if err != nil {
			ErrorHandler(w, r, "Error stoping service", err, http.StatusBadRequest)
			return
		}

		ResponseHandler(w, r, "Stopped all services", true, nil, gg.GetServicesStatus())
	}
}

func SetStartTimeoutHandler(gg *GladiusGuardian) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vals, err := getJSONFields(w, r, "timeout")
		if err != nil {
			ErrorHandler(w, r, "Couldn't parse body", err, http.StatusBadRequest)
			return
		}
		if _, ok := vals["timeout"]; !ok {
			ErrorHandler(w, r, "Need 'timeout' in request", err, http.StatusBadRequest)
			return
		}
		t, err := strconv.Atoi(string(vals["timeout"]))
		if err != nil {
			ErrorHandler(w, r, "Couldn't parse timeout, must be in seconds", err, http.StatusBadRequest)
			return
		}
		dur := time.Duration(t) * time.Second
		gg.SetTimeout(&dur)
		ResponseHandler(w, r, "Set timeout", true, nil, gg.GetServicesStatus())
	}
}

func GetLogsHandler(gg *GladiusGuardian) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get service name to stop it
		ResponseHandler(w, r, "Not implemented", true, nil, gg.GetServicesStatus())
	}
}
