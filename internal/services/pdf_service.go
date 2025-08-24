package services

import "github.com/jung-kurt/gofpdf"

func ImagesToPDF(images []string, output string) error {
	pdf := gofpdf.New("P", "mm", "A4", "")
	for _, img := range images {
		pdf.AddPage()
		pdf.ImageOptions(img, 10, 10, 190, 0, false, gofpdf.ImageOptions{}, 0, "")
	}
	return pdf.OutputFileAndClose(output)
}
