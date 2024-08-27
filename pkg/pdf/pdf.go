package pdf

import (
	"strings"

	"github.com/jung-kurt/gofpdf/v2"
	"github.com/signintech/gopdf"
)

const fontFamily = "NotoSansJP"
const regularFont = "NotoSansJP-Regular.ttf"
const boldFont = "NotoSansJP-Bold.ttf"

func GeneratePdf2(title string, header string, imagePath string, user string, longText string, outputFile string) error {
	var marginX float64 = 50.0
	var marginY float64 = 80.0
	var x float64 = marginX
	var y float64 = marginY

	pdfObj := gopdf.GoPdf{}
	pdfObj.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4})

	// Get page dimensions
	pageWidth, pageHeight := gopdf.PageSizeA4.W, gopdf.PageSizeA4.H

	err := pdfObj.AddTTFFont(fontFamily, regularFont)
	if err != nil {
		return err
	}

	pdfObj.SetMargins(marginX, marginY, marginX, marginY)
	pdfObj.AddPage()

	// Header
	pdfObj.SetFont(fontFamily, "", 12)
	pdfObj.SetXY(x, y)
	pdfObj.Cell(nil, header)
	y += 40

	// --- Add title
	pdfObj.SetFont(fontFamily, "", 24)

	// Get page dimensions
	textWidth, _ := pdfObj.MeasureTextWidth(title)
	x = (pageWidth - textWidth) / 2
	pdfObj.SetXY(x, y)

	// TODO: make it 2 lines if it's long
	pdfObj.SetXY(x, y)
	pdfObj.Cell(nil, title)
	y += 40

	x = (pageWidth - textWidth) / 2
	pdfObj.SetXY(x, y)
	pdfObj.Cell(nil, title)
	y += 40
	x = marginX

	// --- Draw a horizontal line
	// Coordinates (x1, y1) to (x2, y2)
	x1 := x             // X-coordinate of start point
	y1 := y             // Y-coordinate (vertical position)
	x2 := pageWidth - x // X-coordinate of end point
	y2 := y1            // Y-coordinate (same as y1 to keep it horizontal)

	// Set line width
	pdfObj.SetLineWidth(0.3)
	// Set line color to light gray (192, 192, 192)
	pdfObj.SetStrokeColor(192, 192, 192)

	// Draw line
	pdfObj.Line(x1, y1, x2, y2)
	y += 10

	// --- User
	// Set font for Header
	pdfObj.SetFont(fontFamily, "", 16)
	textWidth, _ = pdfObj.MeasureTextWidth(user)
	x = (pageWidth - textWidth) / 2
	pdfObj.SetXY(x, y)
	pdfObj.Cell(nil, user)
	y += 40

	// --- Image
	x = marginX
	pdfObj.Image(imagePath, x, y, &gopdf.Rect{W: pageWidth - marginX*2, H: 512})

	// --- Add longText
	pdfObj.AddPage()
	x = marginX
	y = marginY
	pdfObj.SetFont(fontFamily, "", 14)
	pdfObj.SetXY(x, y)

	lines := strings.Split(longText, "\n")
	lineHeight := 20.0

	for _, line := range lines {
		if line == "" {
			// Empty line, just move to next line
			y += lineHeight
		} else {
			for _, r := range line {
				char := string(r)
				charWidth, err := pdfObj.MeasureTextWidth(char)
				if err != nil {
					panic(err)
				}

				if x+charWidth > pageWidth-marginX {
					x = marginX
					y += lineHeight
				}

				if y > pageHeight-marginY {
					pdfObj.AddPage()
					y = marginY
				}

				pdfObj.SetXY(x, y)
				pdfObj.Text(char)
				x += charWidth
			}
			x = marginX
			y += lineHeight
		}
	}

	pdfObj.WritePdf(outputFile)
	return nil
}

func drawCenter(pdfObj *gofpdf.Fpdf, str string) {
	// Get the width of the page
	pageWidth, _ := pdfObj.GetPageSize()

	// Calculate the width of the text
	textWidth := pdfObj.GetStringWidth(str)

	// Set X position to center the text
	pdfObj.SetX((pageWidth - textWidth) / 2)

	// Print the text centered
	pdfObj.Cell(textWidth, 10, str)
}

func GeneratePdf(title string, header string, imagePath string, user string, longText string, outputFile string) error {
	var y float64 = 10

	pdfObj := gofpdf.New("P", "mm", "A4", "")
	pdfObj.AddUTF8Font(fontFamily, "", regularFont)
	pdfObj.AddUTF8Font(fontFamily, "B", boldFont)

	// Set margins: left, top, and right (in millimeters)
	pdfObj.SetMargins(10, 20, 10)

	// Add a new page
	pdfObj.AddPage()

	// Header
	pdfObj.SetFont("Arial", "B", 12)
	pdfObj.Cell(0, 10, header)

	y += 30
	pdfObj.SetXY(10, y)

	// Set font
	pdfObj.SetFont(fontFamily, "B", 24)

	// --- Add title
	drawCenter(pdfObj, title)

	// Add line break
	pdfObj.Ln(12)

	// Draw a horizontal line
	// Coordinates (x1, y1) to (x2, y2)
	x1 := 10.0  // X-coordinate of start point
	y1 := 60.0  // Y-coordinate (vertical position)
	x2 := 200.0 // X-coordinate of end point
	y2 := y1    // Y-coordinate (same as y1 to keep it horizontal)

	// Set line width
	pdfObj.SetLineWidth(0.3)
	// Set line color to light gray (192, 192, 192)
	pdfObj.SetDrawColor(192, 192, 192)

	// Draw line
	pdfObj.Line(x1, y1, x2, y2)

	// --- User
	// Set font for Header
	pdfObj.SetFont(fontFamily, "B", 16)
	pdfObj.SetXY(10, y1+10)
	drawCenter(pdfObj, user)

	// Add line break
	pdfObj.Ln(10)

	// --- Add an image (make sure the path is correct)
	pdfObj.Image(imagePath, 10, 100, 190, 0, false, "", 0, "")

	// --- Story
	pdfObj.AddPage()
	y = 10
	pdfObj.SetFont(fontFamily, "", 14)
	//pdfObj.MultiCell(0, y, longText, "", "L", false)
	pdfObj.MultiCell(0, y, longText, "", "J", false)

	// --- Save the file
	err := pdfObj.OutputFileAndClose(outputFile)

	return err
}
