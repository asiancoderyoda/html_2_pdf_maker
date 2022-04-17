package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/joho/godotenv"
)

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
	awsS3Sess   *session.Session
}

func main() {
	loadEnv()

	var cfg Config

	port, err := strconv.Atoi(GetEnvFromKey("PORT"))

	if err != nil {
		fmt.Println("Error while converting PORT to int: ", err)
		os.Exit(1)
	}

	flag.IntVar(&cfg.port, "port", port, "port to listen on")
	flag.StringVar(&cfg.env, "env", GetEnvFromKey("ENV"), "environment")
	flag.StringVar(&cfg.tempDir, "tempDir", GetEnvFromKey("TEMPDIR"), "temporary directory")
	flag.StringVar(&cfg.templateDir, "templateDir", GetEnvFromKey("TEMPLATES"), "template directory")
	flag.StringVar(&cfg.htmlExtension, "htmlExtension", GetEnvFromKey("HTML"), "html extension")
	flag.StringVar(&cfg.pdfExtension, "pdfExtension", GetEnvFromKey("PDF"), "pdf extension")
	flag.Parse()

	awsS3Sess := GetAwsSession()

	app := Application{
		config:      cfg,
		wkhtmltopdf: &PDFGenerator{},
		awsS3Sess:   awsS3Sess,
	}

	err = PrepareServer()

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

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
		os.Exit(1)
	}
}

func PrepareServer() error {
	// Create temporary direcotry if it doesn't exist
	if _, err := os.Stat(GetEnvFromKey("TEMPDIR")); os.IsNotExist(err) {
		errDir := os.Mkdir(GetEnvFromKey("TEMPDIR"), 0777)
		if errDir != nil {
			fmt.Println("Error while creating directory: ", errDir)
			return errDir
		}
	}

	// Create output directory if it doesn't exist
	if _, err := os.Stat(GetEnvFromKey("OUTPUTDIR")); os.IsNotExist(err) {
		errDir := os.Mkdir(GetEnvFromKey("OUTPUTDIR"), 0777)
		if errDir != nil {
			fmt.Println("Error while creating directory: ", errDir)
			return errDir
		}
	}

	return nil
}
