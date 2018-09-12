package guardian

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

func GetServicesHandler(gg *GladiusGuardian) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "Category: %v\n", vars["category"])
	}
}
