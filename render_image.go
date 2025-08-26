package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"image/png"

	_ "embed" // Required for embedding font data

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"golang.org/x/image/math/fixed"
)

//go:embed fonts/MPLUS1p-Regular.ttf
var fontData []byte

// renderAsPNG formats the table as a PNG image with grid lines and alternating row colors.
func renderAsPNG(table [][]string, title string, fontSize float64) ([]byte, error) {
	if len(table) == 0 {
		return nil, errors.New("cannot generate image from empty data")
	}

	// --- Colors ---
	bgColorHeader := color.RGBA{R: 238, G: 242, B: 249, A: 255} // Light blue-gray
	bgColorEven := color.RGBA{R: 248, G: 249, B: 250, A: 255} // Very light gray
	bgColorOdd := color.White
	lineColor := color.RGBA{R: 222, G: 226, B: 230, A: 255} // Light gray

	// --- Font and Metrics ---
	parsedFont, err := opentype.Parse(fontData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %w", err)
	}

	face, err := opentype.NewFace(parsedFont, &opentype.FaceOptions{
		Size:    fontSize,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create font face: %w", err)
	}
	defer face.Close()

	// --- Layout Calculation ---
	padding := int(fontSize)
	cellPadding := padding / 2
	lineHeight := face.Metrics().Height.Ceil() + cellPadding*2

	titleHeight := 0
	if title != "" {
		titleHeight = lineHeight + padding
	}

	colWidths := make([]int, len(table[0]))
	for _, row := range table {
		for i, cell := range row {
			width := font.MeasureString(face, cell).Ceil()
			if width > colWidths[i] {
				colWidths[i] = width
			}
		}
	}

	totalWidth := 0
	for _, w := range colWidths {
		totalWidth += w + padding
	}
	totalHeight := titleHeight + len(table)*lineHeight + padding

	// --- Image Drawing ---
	img := image.NewRGBA(image.Rect(0, 0, totalWidth, totalHeight))
	draw.Draw(img, img.Bounds(), image.White, image.Point{}, draw.Src)

	drawer := &font.Drawer{
		Dst:  img,
		Src:  image.Black, // Use image.Black which is an image.Image
		Face: face,
	}

	// --- Draw Backgrounds and Text ---
	y := titleHeight + padding/2
	if title != "" {
		titleX := (totalWidth - font.MeasureString(face, title).Ceil()) / 2
		drawer.Dot = fixed.P(titleX, titleHeight-padding/2)
		drawer.DrawString(title)
	}

	for i, row := range table {
		rowY := y + i*lineHeight
		
		var bgColor color.Color
		if i == 0 {
			bgColor = bgColorHeader
		} else if (i-1)%2 == 0 {
			bgColor = bgColorOdd
		} else {
			bgColor = bgColorEven
		}
		draw.Draw(img, image.Rect(0, rowY, totalWidth, rowY+lineHeight), &image.Uniform{C: bgColor}, image.Point{}, draw.Src)

		x := 0
		for j, cell := range row {
			textX := x + cellPadding
			textY := rowY + (lineHeight-face.Metrics().Height.Ceil())/2 + face.Metrics().Ascent.Ceil()
			drawer.Dot = fixed.P(textX, textY)
			drawer.DrawString(cell)
			x += colWidths[j] + padding
		}
	}

	// --- Draw Grid Lines ---
	tableTop := titleHeight + padding/2
	tableBottom := tableTop + len(table)*lineHeight
	// Horizontal lines
	for i := 0; i <= len(table); i++ {
		yLine := tableTop + i*lineHeight
		for x := 0; x < totalWidth; x++ {
			img.Set(x, yLine, lineColor)
		}
	}
	// Vertical lines
	x := 0
	for i := 0; i < len(colWidths); i++ {
		for y := tableTop; y < tableBottom; y++ {
			img.Set(x, y, lineColor)
		}
		x += colWidths[i] + padding
	}
	// Last vertical line
	for y := tableTop; y < tableBottom; y++ {
		img.Set(totalWidth-1, y, lineColor)
	}

	// --- Encoding ---
	var buf bytes.Buffer
	if err := png.Encode(&buf, img); err != nil {
		return nil, fmt.Errorf("failed to encode png: %w", err)
	}
	return buf.Bytes(), nil
}
