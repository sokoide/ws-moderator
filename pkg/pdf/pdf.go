package pdf

import "github.com/jung-kurt/gofpdf"

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

const fontFamily = "NotSansJP"
const regularFont = "NotoSansJP-Regular.ttf"
const boldFont = "NotoSansJP-Bold.ttf"

func GeneratePdf(title string, header string, imagePath string, user string, longText string, outputFile string) {
	var y float64 = 10

	pdfObj := gofpdf.New("P", "mm", "A4", "")
	pdfObj.AddUTF8Font(fontFamily, "", regularFont)
	pdfObj.AddUTF8Font(fontFamily, "B", boldFont)

	// Set margins: left, top, and right (in millimeters)
	pdfObj.SetMargins(20, 20, 20)

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
	pdfObj.MultiCell(0, y, longText, "", "", false)

	// --- Save the file
	err := pdfObj.OutputFileAndClose(outputFile)

	if err != nil {
		panic(err)
	}
}
