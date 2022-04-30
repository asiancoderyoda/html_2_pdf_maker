package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

var AwsRegion string
var SecretAccessKey string
var AccessKeyId string

func GetEnvFromKey(key string) string {
	return os.Getenv(key)
}

func (app *Application) writeJSON(w http.ResponseWriter, statusCode int, payload interface{}, wrap string) error {
	wrapper := make(map[string]interface{})
	wrapper[wrap] = payload
	js, err := json.Marshal(wrapper)
	if err != nil {
		return err
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	w.Write(js)

	return nil
}

func (app *Application) writeError(w http.ResponseWriter, err error) {
	type errorResponse struct {
		Error string `json:"error"`
	}

	httpError := errorResponse{
		Error: err.Error(),
	}

	app.writeJSON(w, http.StatusUnprocessableEntity, httpError, "error")

}

func (app *Application) parseTemplate(templateType string, data TemplateInterface) (string, error) {
	var document bytes.Buffer // buffer to hold the final document

	// Load the HTML template
	templatePath := fmt.Sprintf("%s%s%s", app.config.templateDir, templateType, app.config.htmlExtension)
	tmpl, err := template.ParseFiles(templatePath)

	if err != nil {
		return "", err
	}

	// Execute the template
	err = tmpl.Execute(&document, data)

	if err != nil {
		return "", err
	}

	// Create populated HTML template
	populatedTemplate := fmt.Sprintf("%s%d-%d%s", app.config.tempDir, data.GetID(), int32(time.Now().UnixNano()), app.config.htmlExtension)
	file, _ := os.Create(populatedTemplate)
	defer file.Close()

	// Write the populated HTML template to file
	file.Write(document.Bytes())

	return populatedTemplate, nil
}

func (app *Application) fetchTemplate(templateType string, pdfData []byte) (TemplateInterface, error) {
	var data TemplateInterface
	switch templateType {
	case "invoice":
		data = &Invoice{}
		err := json.NewDecoder(bytes.NewReader(pdfData)).Decode(data)
		if err != nil {
			return nil, err
		}

	default:
		err := fmt.Errorf("%s%s", "Unknown template type: ", templateType)
		return nil, err
	}

	return data, nil
}

/*
 *Config setup for S3 bucket
 */
func GetAwsSession() *session.Session {
	AccessKeyId = GetEnvFromKey("AWS_ACCESS_KEY_ID")
	SecretAccessKey = GetEnvFromKey("AWS_SECRET_ACCESS_KEY")
	AwsRegion = GetEnvFromKey("AWS_REGION")

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(AwsRegion),
		Credentials: credentials.NewStaticCredentials(
			AccessKeyId,
			SecretAccessKey,
			"",
		),
	})

	if err != nil {
		fmt.Println("Error creating aws session: ", err)
		panic(err)
	}

	fmt.Println("Successfully created aws session")

	return sess
}

/*
 * Utility to upload file to S3 bucket
 */
func UploadFileToS3(templateType string, filePath string, sess *session.Session) (string, error) {
	// Open the file for use
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("Error opening file: ", err)
		return "", err
	}
	defer file.Close()

	// Get file size and read the file content into a buffer
	fileInfo, _ := file.Stat()
	var size int64 = fileInfo.Size()
	var fileName string = fileInfo.Name()
	buffer := make([]byte, size)
	file.Read(buffer)

	// NewUploader creates a new Uploader instance to upload objects to S3
	uploader := s3manager.NewUploader(sess)
	AwsBucket := GetEnvFromKey("AWS_S3_BUCKET")
	path_to_directory := templateType

	/*
	 * Config settings: this is where you choose the bucket, filename, content-type etc.
	 * of the file you're uploading.
	 */
	uploadOutput, err := uploader.Upload(&s3manager.UploadInput{
		Bucket:               aws.String(AwsBucket),
		Key:                  aws.String(path.Join(path_to_directory, fileName)), // dir/filename.ext
		Body:                 bytes.NewReader(buffer),
		ContentType:          aws.String(http.DetectContentType(buffer)),
		ContentDisposition:   aws.String("attachment"),
		ServerSideEncryption: aws.String("AES256"),
		// ACL:               	  aws.String("private"),
		// ContentLength:     	  aws.Int64(size),
	})

	if err != nil {
		fmt.Println("Error uploading file to S3: ", err)
		return "", err
	}

	fmt.Printf("Successfully uploaded file to S3: %s\n", uploadOutput.Location)
	fmt.Println(uploadOutput.ETag, uploadOutput.VersionID, uploadOutput.UploadID)

	return uploadOutput.Location, nil
}

/*
 * Utility to download file to S3 bucket
 */
func DownloadFileFromS3(docType string, key string, sess *session.Session) error {
	AwsBucket := GetEnvFromKey("AWS_S3_BUCKET")
	downloader := s3manager.NewDownloader(sess)
	path_to_directory := docType

	// Create a file to write the S3 Object contents to.
	file, err := os.Create(path.Join(GetEnvFromKey("TEMPDIR"), key))
	if err != nil {
		fmt.Println("Error creating file: ", err)
		return err
	}
	defer file.Close()

	// Write the contents of S3 Object to the file
	n, err := downloader.Download(file,
		&s3.GetObjectInput{
			Bucket: aws.String(AwsBucket),
			Key:    aws.String(path.Join(path_to_directory, key+".pdf")),
		})
	if err != nil {
		fmt.Println("Error downloading file: ", err)
		return err
	}
	fmt.Println("Successfully downloaded file: ", n)

	return nil
}
