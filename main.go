package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gladiusio/gladius-guardian/guardian"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func main() {
	r := mux.NewRouter()
	gg := guardian.New()

	// Register our two daemons
	gg.RegisterService(
		"networkd",
		viper.GetString("networkdExecutable"),
		viper.GetStringSlice("networkdEnvironment"),
	)
	gg.RegisterService(
		"cotnrold",
		viper.GetString("cotnroldExecutable"),
		viper.GetStringSlice("cotnroldEnvironment"),
	)

	// Setup routes

	srv := &http.Server{
		Addr: "0.0.0.0:7791",
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
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

	gg.StopAll()
	stopHTTPServer(srv)
}

func stopHTTPServer(srv *http.Server) {
	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()

	srv.Shutdown(ctx)
	os.Exit(0)
}
