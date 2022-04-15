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
const TEMPDIR = "temp/"
const OUTPUTDIR = "output/"
const TEMPLATES = "templates/"

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

	flag.IntVar(&cfg.port, "port", 8090, "port to listen on")
	flag.StringVar(&cfg.env, "env", "development", "environment")
	flag.StringVar(&cfg.tempDir, "tempDir", TEMPDIR, "temporary directory")
	flag.StringVar(&cfg.templateDir, "templateDir", TEMPLATES, "template directory")
	flag.StringVar(&cfg.htmlExtension, "htmlExtension", HTML, "html extension")
	flag.StringVar(&cfg.pdfExtension, "pdfExtension", PDF, "pdf extension")
	flag.Parse()

	app := Application{
		config:      cfg,
		wkhtmltopdf: &PDFGenerator{},
	}

	err := prepareServer()

	if err != nil {
		fmt.Println("Error while preparing server: ", err)
		return
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
		return
	}
}

func prepareServer() error {
	// Create temporary direcotry if it doesn't exist
	if _, err := os.Stat(TEMPDIR); os.IsNotExist(err) {
		errDir := os.Mkdir(TEMPDIR, 0777)
		if errDir != nil {
			fmt.Println("Error while creating directory: ", errDir)
			return errDir
		}
	}

	// Create output directory if it doesn't exist
	if _, err := os.Stat(OUTPUTDIR); os.IsNotExist(err) {
		errDir := os.Mkdir(OUTPUTDIR, 0777)
		if errDir != nil {
			fmt.Println("Error while creating directory: ", errDir)
			return errDir
		}
	}

	return nil
}
