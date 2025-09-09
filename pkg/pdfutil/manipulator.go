
package pdfutil

import (
	"fmt"
	"image"
	"os"

	"github.com/unidoc/unipdf/v3/common/license"
	"github.com/unidoc/unipdf/v3/extractor"
	"github.com/unidoc/unipdf/v3/model"
)

func init() {
	// UniPDF requires a license key to operate. You can get a free one for non-commercial use
	// from their website. It's best to load this from an environment variable.
	// For this example, I will assume it's set.
	licenseKey := os.Getenv("UNIDOC_LICENSE_KEY")
	if licenseKey == "" {
		fmt.Println("Warning: UNIDOC_LICENSE_KEY is not set. UniPDF features will be limited.")
	}
	err := license.SetMeteredKey(licenseKey)
	if err != nil {
		// This can happen if the key is invalid or expired
		fmt.Printf("ERROR: Failed to set UniPDF license key: %v\n", err)
	}
}

// --- File Advanced Functions ---

// openPDFHelper opens a PDF file and returns a UniPDF reader object.
func openPDFHelper(inputFile string) (*model.PdfReader, error) {
	f, err := os.Open(inputFile)
	if err != nil {
		return nil, err
	}

	pdfReader, err := model.NewPdfReader(f)
	if err != nil {
		return nil, err
	}

	return pdfReader, nil
}

// MergePDFs combines multiple PDF files into a single output file.
func MergePDFs(inputFiles []string, outputFile string) error {
	pdfWriter := model.NewPdfWriter()

	for _, inputFile := range inputFiles {
		pdfReader, err := openPDFHelper(inputFile)
		if err != nil {
			return err
		}

		numPages, err := pdfReader.GetNumPages()
		if err != nil {
			return err
		}

		for i := 0; i < numPages; i++ {
			pageNum := i + 1
			page, err := pdfReader.GetPage(pageNum)
			if err != nil {
				return err
			}

			if err := pdfWriter.AddPage(page); err != nil {
				return err
			}
		}
	}

	fWrite, err := os.Create(outputFile)
	if err != nil {
		return err
	}
	defer fWrite.Close()

	return pdfWriter.Write(fWrite)
}

// SplitPDF splits a PDF into multiple files based on page ranges.
// pageRanges: e.g., ["1-3", "4", "5-7"]
func SplitPDF(inputFile string, pageRanges []string, outputDir string) error {
	pdfReader, err := openPDFHelper(inputFile)
	if err != nil {
		return err
	}

	// Create the output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, os.ModePerm); err != nil {
		return err
	}

	for i, prange := range pageRanges {
		pages, err := model.NewPageRange(prange)
		if err != nil {
			return err
		}

		pdfWriter := model.NewPdfWriter()
		if err = pdfWriter.AddPages(pdfReader, pages); err != nil {
			return err
		}

		outputFile := fmt.Sprintf("%s/split_%d.pdf", outputDir, i+1)
		fWrite, err := os.Create(outputFile)
		if err != nil {
			return err
		}
		defer fWrite.Close()

		if err := pdfWriter.Write(fWrite); err != nil {
			return err
		}
	}

	return nil
}

// --- Content Extraction ---

// ExtractText extracts all text from a specific page of a PDF.
func ExtractText(inputFile string, pageNum int) (string, error) {
	pdfReader, err := openPDFHelper(inputFile)
	if err != nil {
		return "", err
	}

	page, err := pdfReader.GetPage(pageNum)
	if err != nil {
		return "", err
	}

	ex, err := extractor.New(page)
	if err != nil {
		return "", err
	}

	return ex.ExtractText()
}

// ExtractImages extracts all images from a specific page of a PDF.
func ExtractImages(inputFile string, pageNum int) ([]image.Image, error) {
	pdfReader, err := openPDFHelper(inputFile)
	if err != nil {
		return nil, err
	}

	page, err := pdfReader.GetPage(pageNum)
	if err != nil {
		return nil, err
	}

	pimages, err := page.GetPageImages(nil)
	if err != nil {
		return nil, err
	}

	var images []image.Image
	for _, pimage := range pimages {
		img, err := pimage.ToImage()
		if err != nil {
			return nil, err
		}
		images = append(images, img)
	}

	return images, nil
}

// --- Security & Optimization (Conceptual) ---
// The following functions require a more advanced setup or a specific UniPDF license.
// I am providing their signatures and a brief explanation.

// EncryptPDF protects a PDF with a password.
func EncryptPDF(inputFile, outputFile, userPassword string) error {
	// This requires creating a `PdfEncrypt` object and applying it with the writer.
	// A full implementation is more involved.
	return fmt.Errorf("encryption not fully implemented yet")
}

// DecryptPDF removes password protection from a PDF.
func DecryptPDF(inputFile, outputFile, password string) error {
	// This involves attempting to read the PDF with the given password.
	return fmt.Errorf("decryption not fully implemented yet")
}

// CompressPDF reduces the file size of a PDF.
func CompressPDF(inputFile, outputFile string) error {
	// True compression often involves re-encoding images, removing unused objects, etc.
	// This is a feature of UniPDF's optimizer package.
	return fmt.Errorf("compression not fully implemented yet")
}
