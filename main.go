package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gladiusio/gladius-common/pkg/routing"
	"github.com/gladiusio/gladius-common/pkg/utils"
	"github.com/gladiusio/gladius-guardian/config"
	"github.com/gladiusio/gladius-guardian/guardian"
	"github.com/gladiusio/gladius-guardian/service"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	service.SetupService(run)
}

func run() {
	base, err := utils.GetGladiusBase()
	if err != nil {
		log.WithFields(log.Fields{
			"err": err,
		}).Fatal("Couldn't get Gladius base")
	}
	config.SetupConfig(base)

	r := mux.NewRouter()
	gg := guardian.New()

	// Register our two daemons
	gg.RegisterService(
		"edged",
		viper.GetString("NetworkdExecutable"),
		viper.GetStringSlice("DefaultEnvironment"),
	)
	gg.RegisterService(
		"network-gateway",
		viper.GetString("ControldExecutable"),
		viper.GetStringSlice("DefaultEnvironment"),
	)

	// Handle the index
	r.HandleFunc("/", guardian.IndexHandler)

	// Guardian related endpoints
	r.HandleFunc("/service/stats/{service_name}", guardian.GetServicesHandler(gg)).Methods("GET")
	r.HandleFunc("/service/set_state/{service_name}", guardian.ServiceStateHandler(gg)).Methods("PUT")
	r.HandleFunc("/service/set_timeout", guardian.SetStartTimeoutHandler(gg)).Methods("POST")
	r.HandleFunc("/service/logs", guardian.GetOldLogsHandler(gg)).Methods("GET")
	r.HandleFunc("/service/ws/logs/{service_name}", guardian.GetNewLogsWebSocketHandler(gg))

	// Version
	r.HandleFunc("/service/version/{service_name}", guardian.VersionHandler()).Methods("GET")

	routing.AppendVersionEndpoints(r)

	// Setup a custom server so we can gracefully stop later
	srv := &http.Server{
		Addr:         "0.0.0.0:7791",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r,
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Println(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)

	<-c // Block until we receive our signal.

	gg.StopService("all")
	stopHTTPServer(srv)
}

func stopHTTPServer(srv *http.Server) {
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	srv.Shutdown(ctx)
	os.Exit(0)
}
