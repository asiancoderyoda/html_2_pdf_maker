package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"time"
)

const version = "1.0.0"
const HTML = ".html"
const PDF = ".pdf"
const TEMPDIR = "/tmp/"
const OUTPUTDIR = "/output/"
const TEMPLATES = "/templates/"

type Config struct {
	port          int
	env           string
	tempDir       string
	templateDir   string
	htmlExtension string
	pdfExtension  string
}

type ServerStatus struct {
	Version string `json:"version"`
	Status  string `json:"status"`
	Env     string `json:"env"`
}

type Application struct {
	config      Config
	wkhtmltopdf wkhtmltopdfInterface
}

func main() {
	var cfg Config

	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	flag.IntVar(&cfg.port, "port", 8090, "port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "environment")
	flag.StringVar(&cfg.tempDir, "tempDir", dir+TEMPDIR, "temporary directory")
	flag.StringVar(&cfg.templateDir, "templateDir", dir+TEMPLATES, "template directory")
	flag.StringVar(&cfg.htmlExtension, "htmlExtension", HTML, "html extension")
	flag.StringVar(&cfg.pdfExtension, "pdfExtension", PDF, "pdf extension")
	flag.Parse()

	app := Application{
		config:      cfg,
		wkhtmltopdf: &PDFGenerator{},
	}

	srv := &http.Server{
		Addr:         fmt.Sprintf(":%d", cfg.port),
		Handler:      app.routes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  time.Second * 10,
		WriteTimeout: time.Second * 30,
	}

	fmt.Printf("Starting server on port %d in %s mode\n", cfg.port, cfg.env)

	err = srv.ListenAndServe()

	if err != nil {
		fmt.Printf("Error starting server: %s\n", err)
	}
}
