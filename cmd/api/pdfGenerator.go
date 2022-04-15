package main

import (
	"fmt"
	"os"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type wkhtmltopdfInterface interface {
	createPdf(string) (bool, error)
}

type PDFGenerator struct{}

func (pdfGen *PDFGenerator) createPdf(pathToFile string) (bool, error) {
	// Create temporary direcotry if it doesn't exist
	if _, err := os.Stat(TEMPDIR); os.IsNotExist(err) {
		errDir := os.Mkdir(TEMPDIR, 0777)
		if errDir != nil {
			fmt.Println("Error while creating directory: ", errDir)
			return false, errDir
		}
	}

	// Create new PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return false, err
	}

	// Set global options
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.Cover.Zoom.Set(0.75)
	pdfg.Dpi.Set(300)
	page := wkhtmltopdf.NewPage(pathToFile)
	pdfg.AddPage(page)

	err = pdfg.Create()

	if err != nil {
		return false, err
	}

	pdfPath := fmt.Sprintf("%s%s-%d%s", OUTPUTDIR, "INV-GEN", int32(time.Now().UnixNano()), PDF)
	err = pdfg.WriteFile(pdfPath)

	if err != nil {
		return false, err
	}

	// dir, err := os.Getwd()
	// if err != nil {
	// 	panic(err)
	// }

	defer os.RemoveAll(TEMPDIR)

	return true, nil
}
