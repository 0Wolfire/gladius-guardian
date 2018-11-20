package guardian

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
	"github.com/gladiusio/gladius-guardian/updater"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	ResponseHandler(w, r, "There's nothing here, check our API docs at https://github.com/gladiusio/gladius-guardian", true, nil, "")
}

func GetServicesHandler(gg *GladiusGuardian) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sn := vars["service_name"]

		ResponseHandler(w, r, "Got service status", true, nil, gg.GetServicesStatus(sn))
	}
}

func ServiceStateHandler(gg *GladiusGuardian) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		// Get desired run state, optionally environment variables
		vals, err := getJSONFields(w, r, "running", "environment_vars")
		if err != nil {
			ErrorHandler(w, r, "Couldn't parse body", err, http.StatusBadRequest)
			return
		}
		if _, ok := vals["running"]; !ok {
			ErrorHandler(w, r, "Need 'running' in request", err, http.StatusBadRequest)
			return
		}

		// Get the service name from the URL
		vars := mux.Vars(r)
		sn := vars["service_name"]

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

		// Parse the run state they want
		setRunning, err := strconv.ParseBool(string(vals["running"]))
		if err != nil {
			ErrorHandler(w, r, "Could not parse running as type bool", err, http.StatusBadRequest)
			return
		}

		// Start or stop the service
		if setRunning {
			err = gg.StartService(sn, environmentVars)
			if err != nil {
				ErrorHandler(w, r, "Error starting service", err, http.StatusBadRequest)
				return
			}
			ResponseHandler(w, r, "Attempted to start service, check logs to make sure it didn't fail after timeout", true, nil, gg.GetServicesStatus(sn))
		} else {
			err = gg.StopService(sn)
			if err != nil {
				ErrorHandler(w, r, "Error stoping service", err, http.StatusBadRequest)
				return
			}
			time.Sleep(200 * time.Millisecond)
			ResponseHandler(w, r, "Stopped Service", true, nil, gg.GetServicesStatus(sn))
		}

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
		response := make(map[string]int)
		response["timeout"] = t
		ResponseHandler(w, r, "Set timeout", true, nil, response)
	}
}

func GetOldLogsHandler(gg *GladiusGuardian) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		toReturn := make(map[string]([]string))
		for name, fsl := range gg.serviceLogs {
			toReturn[name] = fsl.LogLines()
		}
		ResponseHandler(w, r, "Got logs", true, nil, toReturn)
	}
}

func GetNewLogsWebSocketHandler(gg *GladiusGuardian) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		sn := vars["service_name"]
		if sn != "" {
			gg.AddLogClient(sn, w, r)
		}
	}
}

// VersionHandler - Sends back the version information on whether you need to update or not
func VersionHandler() func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		service := vars["service_name"]

		officialVersions, err := updater.GetOfficialVersions()
		if err != nil {
			ErrorHandler(w, r, "Couldn't get version", err, http.StatusBadRequest)
		}

		modules := [3]string{"guardian", "edged", "network-gateway"}
		response := make(map[string]int)
		var version int

		for i := 0; i < 3; i++ {
			realVersion, err := updater.GetVersion("guardian")
			if err != nil {
				ErrorHandler(w, r, "Couldn't get version", err, http.StatusBadRequest)
				return
			}
			// realVersion, err := updater.GetVersion(modules[i])
			version, err = updater.CompareVersion(realVersion, officialVersions[fmt.Sprintf("gladius-%s", modules[i])])
			if err != nil {
				ErrorHandler(w, r, "Couldn't get version", err, http.StatusBadRequest)
				return
			}
			response[fmt.Sprintf("gladius-%s", modules[i])] = version
		}

		if service == "all" {
			ResponseHandler(w, r, "Got all versions", true, nil, response)
			return
		}

		ResponseHandler(w, r, fmt.Sprintf("Got version for %s", service), true, nil, response["service"])
	}
}
