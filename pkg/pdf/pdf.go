package pdf

import (
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/jung-kurt/gofpdf/v2"
	"github.com/signintech/gopdf"
)

const fontFamily = "NotoSansJP"
const regularFont = "NotoSansJP-Regular.ttf"
const boldFont = "NotoSansJP-Bold.ttf"

// var re = regexp.MustCompile(`^[A-Za-z\s\-–’‘“”()[]{}.,!?;:'"]*$`)

func isAllEnglish(text string) bool {
	// check the first 100 chars
	len := len(text)
	if len > 100 {
		len = 100
	}
	target := text[:len]
	nonEnglishChars := 0

	for _, char := range target {
		if char > unicode.MaxASCII && char != '–' {
			nonEnglishChars++
		}
	}

	// if <10% is non English chars, assume they are symbols and return true
	if nonEnglishChars < len/10 {
		return true
	}
	return false
}

func DrawLongText(isEnglish bool, pdfObj *gopdf.GoPdf, line string, x float64, y float64, lineHeight float64,
	pageWidth float64, pageHeight float64,
	marginX float64, marginY float64) (float64, float64) {
	if isEnglish {
		return DrawLongTextEnglish(pdfObj, line, x, y, lineHeight, pageWidth, pageHeight, marginX, marginY)
	} else {
		return DrawLongTextNonEnglish(pdfObj, line, x, y, lineHeight, pageWidth, pageHeight, marginX, marginY)
	}
}

func DrawLongTextEnglish(pdfObj *gopdf.GoPdf, line string, x float64, y float64, lineHeight float64,
	pageWidth float64, pageHeight float64,
	marginX float64, marginY float64) (float64, float64) {
	words := strings.Split(line, " ")
	for _, word := range words {
		wordWidth, err := pdfObj.MeasureTextWidth(word + " ")
		if err != nil {
			panic(err)
		}

		if x+wordWidth > pageWidth-marginX {
			x = marginX
			y += lineHeight
		}

		if y > pageHeight-marginY {
			pdfObj.AddPage()
			y = marginY
		}

		pdfObj.SetXY(x, y)
		pdfObj.Text(word + " ")
		x += wordWidth
	}
	// Move to next line after finishing current line
	x = marginX
	y += lineHeight

	return x, y
}

func DrawLongTextNonEnglish(pdfObj *gopdf.GoPdf, line string, x float64, y float64, lineHeight float64,
	pageWidth float64, pageHeight float64,
	marginX float64, marginY float64) (float64, float64) {
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
	return x, y
}

func SplitString(s string) (string, string) {
	// Get the length of the string in runes (characters), not bytes
	length := utf8.RuneCountInString(s)

	// Find the split point
	splitPoint := length / 2

	// Iterate to find the exact rune boundary for the split point
	var currentRuneIndex int
	var splitIndex int

	for i, _ := range s {
		if currentRuneIndex == splitPoint {
			splitIndex = i
			break
		}
		currentRuneIndex++
	}

	// Split the string into two parts
	firstPart := s[:splitIndex]
	secondPart := s[splitIndex:]

	return firstPart, secondPart
}

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

	if textWidth > pageWidth-marginX {
		// make it 2 lines
		title1, title2 := SplitString(title)

		textWidth, _ = pdfObj.MeasureTextWidth(title1)
		x = (pageWidth - textWidth) / 2
		pdfObj.SetXY(x, y)
		pdfObj.Cell(nil, title1)
		y += 40

		textWidth, _ = pdfObj.MeasureTextWidth(title2)
		x = (pageWidth - textWidth) / 2
		pdfObj.SetXY(x, y)
		pdfObj.Cell(nil, title2)
		y += 40
	} else {
		x = (pageWidth - textWidth) / 2
		pdfObj.SetXY(x, y)
		pdfObj.Cell(nil, title)
		y += 40
	}
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

	isEnglish := isAllEnglish(longText)
	// fmt.Printf("isEnglish: %v\n", isEnglish)
	lines := strings.Split(longText, "\n")
	lineHeight := 20.0

	for _, line := range lines {
		x, y = DrawLongText(isEnglish, &pdfObj, line, x, y, lineHeight, pageWidth, pageHeight, marginX, marginY)
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
