package util

import "github.com/jung-kurt/gofpdf"

type cellType struct {
	list [][]byte
	ht   float64
	wd   float64
	bld  bool
}

type PDFTableHeader struct {
	IsBold      bool
	Name        string
	ColumnWidth float64
}

func AddTable(pdf *gofpdf.Fpdf, header []PDFTableHeader, strList [][]string, yStart float64, border bool, isShowHeader bool, cellGap float64, lineHt float64, marginH float64) float64 {
	autoBreak, _ := pdf.GetAutoPageBreak()
	if autoBreak {
		pdf.SetAutoPageBreak(false, 0)
	}

	var cell cellType
	var cellList [][]cellType
	cellList = make([][]cellType, len(strList)+1)

	// Rows
	y := yStart
	count := 0
	var maxHt = lineHt
	var headerMaxHt = lineHt

	startFrom := 1
	if isShowHeader {
		startFrom = 0
	}

	for i := 0; i < len(header); i++ {
		cell.list = pdf.SplitLines([]byte(header[i].Name), header[i].ColumnWidth-cellGap-cellGap)
		cell.ht = float64(len(cell.list)) * lineHt
		if cell.ht > headerMaxHt {
			headerMaxHt = cell.ht
		}
		cellList[0] = append(cellList[0], cell)
	}

	for i := 0; i < len(strList); i++ {
		count++
		for j := 0; j < len(strList[i]); j++ {
			cell.list = pdf.SplitLines([]byte(strList[i][j]), header[j].ColumnWidth-cellGap-cellGap)
			cell.ht = float64(len(cell.list)) * lineHt
			if cell.ht > maxHt {
				maxHt = cell.ht
			}
			cellList[i+1] = append(cellList[i+1], cell)
		}
	}

	_, height := pdf.GetPageSize()
	_, top, _, bottom := pdf.GetMargins()
	for i := startFrom; i < len(cellList); i++ {
		usedMax := maxHt
		if i == 0 {
			style := ""
			if header[i].IsBold {
				style = "B"
			}
			usedMax = headerMaxHt
			pdf.SetFontStyle(style)
		}
		if y+marginH+top+bottom > height {
			pdf.AddPage()
			y = pdf.GetY()
		}
		x := marginH
		for j := 0; j < len(cellList[i]); j++ {
			align := ""
			if i == 0 {
				align = "C"
			}
			pdf.ClipRect(x, y, header[j].ColumnWidth, usedMax+cellGap+cellGap, border)
			cell = cellList[i][j]
			cellY := y + cellGap + (usedMax-cell.ht)/2
			for splitJ := 0; splitJ < len(cell.list); splitJ++ {
				pdf.SetXY(x+cellGap, cellY)
				pdf.MultiCell(header[j].ColumnWidth-cellGap-cellGap, lineHt, string(cell.list[splitJ]), "", align, false)
				cellY += lineHt
			}
			pdf.ClipEnd()
			x += header[j].ColumnWidth
		}
		if i == 0 {
			pdf.SetFontStyle("")
		}
		y += usedMax + cellGap + cellGap
	}
	return y
}
