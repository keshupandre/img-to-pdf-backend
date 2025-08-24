package services

import (
	"github.com/jung-kurt/gofpdf"
)

type Margins struct {
	Top, Right, Bottom, Left float64
}

func ImagesToPDF(images []string, output string, fitSmall bool, position string, orientation string) error {
	if orientation != "L" {
		orientation = "P"
	}

	pdf := gofpdf.New(orientation, "mm", "A4", "")

	var pageW, pageH float64
	if orientation == "L" {
		pageW, pageH = 297.0, 210.0
	} else {
		pageW, pageH = 210.0, 297.0
	}

	margins := Margins{
		Top:    10,
		Right:  10,
		Bottom: 10,
		Left:   10,
	}

	usableW := pageW - (margins.Left + margins.Right)
	usableH := pageH - (margins.Top + margins.Bottom)

	for _, img := range images {
		pdf.AddPage()

		info := pdf.RegisterImage(img, "")
		imgW, imgH := info.Extent()

		newW, newH := imgW, imgH

		if imgW > usableW || imgH > usableH {
			ratio := min(usableW/imgW, usableH/imgH)
			newW = imgW * ratio
			newH = imgH * ratio
		} else if fitSmall {
			ratio := min(usableW/imgW, usableH/imgH)
			newW = imgW * ratio
			newH = imgH * ratio
		}

		x, y := calculatePosition(position, usableW, usableH, newW, newH)
		x += margins.Left
		y += margins.Top

		pdf.ImageOptions(img, x, y, newW, newH, false, gofpdf.ImageOptions{}, 0, "")
	}

	return pdf.OutputFileAndClose(output)
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func calculatePosition(pos string, areaW, areaH, imgW, imgH float64) (float64, float64) {
	var x, y float64

	switch pos {
	case "top-left":
		x, y = 0, 0
	case "top-center":
		x, y = (areaW-imgW)/2, 0
	case "top-right":
		x, y = areaW-imgW, 0
	case "center-left":
		x, y = 0, (areaH-imgH)/2
	case "center", "center-center":
		x, y = (areaW-imgW)/2, (areaH-imgH)/2
	case "center-right":
		x, y = areaW-imgW, (areaH-imgH)/2
	case "bottom-left":
		x, y = 0, areaH-imgH
	case "bottom-center":
		x, y = (areaW-imgW)/2, areaH-imgH
	case "bottom-right":
		x, y = areaW-imgW, areaH-imgH
	default:
		x, y = (areaW-imgW)/2, (areaH-imgH)/2
	}
	return x, y
}
