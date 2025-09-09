
package pdfutil

import (
	"github.com/jung-kurt/gofpdf"
)

// --- File Core Functions ---

// CreatePDF initializes a new PDF document.
// orientation: "P" (Portrait) or "L" (Landscape)
// unit: "mm", "pt", "cm", or "in"
// size: "A3", "A4", "A5", "Letter", "Legal"
func CreatePDF(orientation, unit, size string) (*gofpdf.Fpdf, error) {
	pdf := gofpdf.New(orientation, unit, size, "")
	pdf.AddPage()
	return pdf, nil
}

// SavePDF saves the generated PDF to a file.
func SavePDF(pdf *gofpdf.Fpdf, path string) error {
	return pdf.OutputFileAndClose(path)
}

// --- Page Core Functions ---

// AddPage adds a new page to the PDF.
func AddPage(pdf *gofpdf.Fpdf) {
	pdf.AddPage()
}

// SetPageSize sets the size of the current page. This is not a standard gofpdf function,
// as page size is typically set on creation. This is a helper to illustrate the concept.
func SetPageSize(pdf *gofpdf.Fpdf, width, height float64) {
	// gofpdf sets page size on a per-page basis during creation or with AddPageFormat
	// This function is a conceptual placeholder.
	// To do this properly, one would call: pdf.AddPageFormat("P", gofpdf.SizeType{Wd: width, Ht: height})
}

// SetMargins sets the page margins.
func SetMargins(pdf *gofpdf.Fpdf, left, top, right float64) {
	pdf.SetMargins(left, top, right)
}

// --- Style Core Functions ---

// SetFont sets the font for subsequent text operations.
func SetFont(pdf *gofpdf.Fpdf, family, style string, size float64) {
	pdf.SetFont(family, style, size)
}

// SetTextColor sets the color for text.
func SetTextColor(pdf *gofpdf.Fpdf, r, g, b int) {
	pdf.SetTextColor(r, g, b)
}

// SetFillColor sets the color for filling shapes (like rectangles).
func SetFillColor(pdf *gofpdf.Fpdf, r, g, b int) {
	pdf.SetFillColor(r, g, b)
}

// DrawLine draws a line from (x1, y1) to (x2, y2).
func DrawLine(pdf *gofpdf.Fpdf, x1, y1, x2, y2 float64) {
	pdf.Line(x1, y1, x2, y2)
}

// DrawRect draws a rectangle.
// style: "D" (draw), "F" (fill), or "DF" (draw and fill)
func DrawRect(pdf *gofpdf.Fpdf, x, y, w, h float64, style string) {
	pdf.Rect(x, y, w, h, style)
}

// --- Text Core Functions ---

// AddText adds a single line of text at a specific position.
func AddText(pdf *gofpdf.Fpdf, x, y float64, text string, font string, size float64) {
	pdf.SetFont(font, "", size)
	pdf.Text(x, y, text)
}

// AddParagraph adds a block of text that wraps automatically.
func AddParagraph(pdf *gofpdf.Fpdf, text string, align string) {
	// A=align: "L", "C", "R", "J" (Justify)
	pdf.MultiCell(0, 5, text, "", align, false)
}

// AddHeader adds a header to each page.
func AddHeader(pdf *gofpdf.Fpdf, text string) {
	pdf.SetHeaderFunc(func() {
		pdf.SetY(5)
		pdf.SetFont("Arial", "I", 8)
		pdf.CellFormat(0, 10, text, "", 0, "C", false, 0, "")
	})
}

// AddFooter adds a footer to each page (often with page numbers).
func AddFooter(pdf *gofpdf.Fpdf, text string) {
	pdf.SetFooterFunc(func() {
		pdf.SetY(-15)
		pdf.SetFont("Arial", "I", 8)
		footerText := fmt.Sprintf("%s - Page %d", text, pdf.PageNo())
		pdf.CellFormat(0, 10, footerText, "", 0, "C", false, 0, "")
	})
}

// --- Table Core Functions ---

// AddTable creates a simple table.
func AddTable(pdf *gofpdf.Fpdf, headers []string, data [][]string) {
	// Header
	pdf.SetFont("Arial", "B", 12)
	for _, header := range headers {
		pdf.CellFormat(40, 7, header, "1", 0, "C", false, 0, "")
	}
	pdf.Ln(-1)

	// Data
	pdf.SetFont("Arial", "", 12)
	for _, row := range data {
		AddRow(pdf, row)
	}
}

// AddRow adds a single row to a table.
func AddRow(pdf *gofpdf.Fpdf, row []string) {
	for _, cell := range row {
		pdf.CellFormat(40, 7, cell, "1", 0, "", false, 0, "")
	}
	pdf.Ln(-1)
}

// --- Image Core Functions ---

// AddImage adds an image to the PDF.
func AddImage(pdf *gofpdf.Fpdf, path string, x, y, width, height float64) {
	// The last boolean parameter indicates whether to re-encode the image.
	// False is usually fine and faster.
	pdf.Image(path, x, y, width, height, false, "", 0, "")
}

