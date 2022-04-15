package main

import (
	"flag"
	"fmt"
	"net/http"
	"time"
)

const version = "1.0.0"

type Config struct {
	port int
	env  string
}

type ServerStatus struct {
	Version string `json:"version"`
	Status  string `json:"status"`
	Env     string `json:"env"`
}

type Application struct {
	config Config
}

func main() {
	var cfg Config
	flag.IntVar(&cfg.port, "port", 8090, "port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "environment")
	flag.Parse()

	app := Application{
		config: cfg,
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}

	fmt.Printf("Starting server on port %d in %s mode\n", cfg.port, cfg.env)

	err := srv.ListenAndServe()

	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
