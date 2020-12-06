package vlcServer

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/adanalife/tripbot/pkg/config"
	terrors "github.com/adanalife/tripbot/pkg/errors"
	"github.com/adanalife/tripbot/pkg/helpers"
	sentrynegroni "github.com/getsentry/sentry-go/negroni"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	metrics "github.com/slok/go-http-metrics/metrics/prometheus"
	"github.com/slok/go-http-metrics/middleware"
	negronimiddleware "github.com/slok/go-http-metrics/middleware/negroni"
	"github.com/urfave/negroni"
)

// Start starts the web server
func Start() {
	log.Println("Starting VLC web server on host", config.VlcServerHost)

	r := mux.NewRouter()

	// healthcheck endpoints
	hp := r.PathPrefix("/health").Methods("GET").Subrouter()
	hp.HandleFunc("/live", healthHandler)
	hp.HandleFunc("/ready", healthHandler)

	// vlc endpoints
	vlc := r.PathPrefix("/vlc").Methods("GET").Subrouter()
	vlc.HandleFunc("/current", vlcCurrentHandler)
	vlc.HandleFunc("/play", vlcPlayHandler)
	vlc.HandleFunc("/back", vlcBackHandler)
	vlc.HandleFunc("/skip", vlcSkipHandler)
	vlc.HandleFunc("/random", vlcRandomHandler)

	// onscreen endpoints
	osc := r.PathPrefix("/onscreens").Methods("GET").Subrouter()
	osc.HandleFunc("/flag/show", onscreensFlagShowHandler)
	osc.HandleFunc("/gps/hide", onscreensGpsHideHandler)
	osc.HandleFunc("/gps/show", onscreensGpsShowHandler)
	osc.HandleFunc("/timewarp/show", onscreensTimewarpShowHandler)
	osc.HandleFunc("/leaderboard/show", onscreensLeaderboardShowHandler)
	osc.HandleFunc("/middle/hide", onscreensMiddleHideHandler)
	osc.HandleFunc("/middle/show", onscreensMiddleShowHandler)

	// prometheus metrics endpoint
	r.Path("/metrics").Handler(promhttp.Handler())

	// static assets
	r.HandleFunc("/favicon.ico", faviconHandler).Methods("GET")

	// catch everything else
	r.HandleFunc("/", catchAllHandler)

	helpers.PrintAllRoutes(r)

	// negroni classic adds panic recovery, logger, and static file middlewares
	// c.p. https://github.com/urfave/negroni
	//TODO: consider adding HTMLPanicFormatter
	app := negroni.Classic()

	// attach http-metrics (prometheus) middleware
	mdlw := middleware.New(middleware.Config{
		Recorder: metrics.NewRecorder(metrics.Config{
			Prefix: config.ServerType,
		}),
	})
	app.Use(negronimiddleware.Handler("", mdlw))

	// attach sentry middleware
	app.Use(sentrynegroni.New(sentrynegroni.Options{}))

	// attaching routes to handler happens last
	app.UseHandler(r)

	//TODO: error if there's no colon to split on
	port := strings.Split(config.VlcServerHost, ":")[1]

	srv := &http.Server{
		Addr: fmt.Sprintf("0.0.0.0:%s", port),
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout:   time.Second * 15,
		ReadTimeout:    time.Second * 15,
		IdleTimeout:    time.Second * 60,
		MaxHeaderBytes: 1 << 20, // 1 MB
		Handler:        app,     // Pass our instance of negroni in
	}

	//TODO: add graceful shutdown
	if err := srv.ListenAndServe(); err != nil {
		terrors.Fatal(err, "couldn't start server")
	}
}
