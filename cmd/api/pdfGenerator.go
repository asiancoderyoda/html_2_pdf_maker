package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/SebastiaanKlippert/go-wkhtmltopdf"
)

type wkhtmltopdfInterface interface {
	createPdf(string) (string, error)
}

type PDFGenerator struct{}

func (pdfGen *PDFGenerator) createPdf(pathToFile string) (string, error) {
	// Create new PDF generator
	pdfg, err := wkhtmltopdf.NewPDFGenerator()
	if err != nil {
		return "", err
	}

	// Set global options
	pdfg.PageSize.Set(wkhtmltopdf.PageSizeA4)
	pdfg.Cover.Zoom.Set(0.75)
	pdfg.Dpi.Set(300)
	page := wkhtmltopdf.NewPage(pathToFile)
	pdfg.AddPage(page)

	err = pdfg.Create()

	if err != nil {
		return "", err
	}

	pdfPath := fmt.Sprintf("%s%s-%d%s", OUTPUTDIR, "INV-GEN", int32(time.Now().UnixNano()), PDF)
	err = pdfg.WriteFile(pdfPath)

	if err != nil {
		return "", err
	}

	err = RemoveContents(TEMPDIR)

	if err != nil {
		return "", err
	}

	return pdfPath, nil
}

func RemoveContents(dir string) error {
	d, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer d.Close()
	names, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}
	for _, name := range names {
		err = os.RemoveAll(filepath.Join(dir, name))
		if err != nil {
			return err
		}
	}
	return nil
}
