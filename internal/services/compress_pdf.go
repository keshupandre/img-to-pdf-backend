package services

import (
	"fmt"
	"os/exec"
)

func CompressPDF(tempPath string, outputPath string, quality int) error {
	// Map quality to Ghostscript PDFSETTINGS
	// /screen (low), /ebook (medium), /printer (high), /prepress (very high)
	var setting string
	switch {
	case quality <= 40:
		setting = "/screen"
	case quality <= 70:
		setting = "/ebook"
	case quality <= 90:
		setting = "/printer"
	default:
		setting = "/prepress"
	}

	cmd := exec.Command(
		"gs",
		"-sDEVICE=pdfwrite",
		"-dCompatibilityLevel=1.4",
		"-dPDFSETTINGS="+setting,
		"-dNOPAUSE",
		"-dQUIET",
		"-dBATCH",
		"-sOutputFile="+outputPath,
		tempPath,
	)

	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("ghostscript error: %v, output: %s", err, string(out))
	}
	return nil

}
