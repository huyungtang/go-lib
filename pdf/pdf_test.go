package pdf

import (
	"os"
	"testing"

	"github.com/go-pdf/fpdf"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// TestInit
// ****************************************************************************************************************************************
func TestInit(t *testing.T) {
	pdf, err := Init(
		[]PageOption{
			FontPathOption("/Users/huyungtang/Projects/golang.batches/fonts/"),
		},
		[]CellOption{
			FontFamilyOption("TaipeiSansTCBeta-Regular"),
		},
	)
	if err != nil {
		t.Error(err)
	}
	defer pdf.Close()

	var f *os.File
	if f, err = os.OpenFile("/Users/huyungtang/Downloads/fpdf.pdf", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777); err != nil {
		t.Error(err)
	}
	defer f.Close()

	if err = pdf.
		AddPage(
			TemplateOption(0),
		).
		AddBarcode128("E12345678901111", LocationOption(150, 10), WidthOption(50), HeightOption(12)).
		AddPage(
			PageSizeA4Option(true),
		).
		AddCell("永盛車電股份有限公司", AlignOption(AlignMC), FontSizeOption(15), BorderOption("1"), HeightOption(6)).
		AddCell("永盛車電股份有限公司", AlignOption(AlignMC), FontSizeOption(15), BorderOption("1")).
		AddCell("一二三四五六七八九十",
			LocationOption(10, 50), WidthOption(.2), PositionTailOption(), BorderOption("1")).
		AddCell("一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十一二三四五六七八九十",
			WidthOption(.2), PositionTailOption(), BorderOption("1")).
		AddCell("XXXXXXXXXXXXXXXXXXXX",
			WidthOption(.2), PositionTailOption(), BorderOption("1")).
		// 	AddCell("請購單", AlignOption(AlignMC), FontSizeOption(13), CellHeightOption(7)).
		// 	AddCell("申請人", CellWidthOption(.13), AlignOption(AlignMC), BorderOption(BorderFull), PositionOption(PositionTail)).
		// 	AddCell("張三", CellWidthOption(.2), BorderOption(BorderFull), PositionOption(PositionTail)).
		// 	AddCell("部門", CellWidthOption(.13), AlignOption(AlignMC), BorderOption(BorderFull), PositionOption(PositionTail)).
		// 	AddCell("XX 部", CellWidthOption(.2), BorderOption(BorderFull), PositionOption(PositionTail)).
		// 	AddCell("申請日期", CellWidthOption(.13), AlignOption(AlignMC), BorderOption(BorderFull), PositionOption(PositionTail)).
		// 	AddCell("2024/03/26 11:22", CellWidthOption(.21), BorderOption(BorderFull)).
		// 	SetXY(150, 20).
		// 	SetXY(33, 29.5).AddCell("張三", CellWidthOption(40), CellHeightOption(9)).
		// 	SetXY(96, 29.5).AddCell("XX 部", CellWidthOption(40), CellHeightOption(9)).
		// 	SetXY(159, 29.5).AddCell("2024/03/26 11:22", CellWidthOption(40), CellHeightOption(9)).
		Output(f); err != nil {
		t.Error(err)
	}
}

// TestSDK
// ****************************************************************************************************************************************
func TestSDK(t *testing.T) {
	out := fpdf.New("P", "mm", "A4", "/")
	out.AddPage()
	out.AddUTF8Font("default", "", "/Users/huyungtang/Projects/golang.batches/fonts/TaipeiSansTC/TaipeiSansTCBeta-Regular.ttf")
	out.SetFont("default", "", 12)
	_, ht := out.GetFontSize()

	out.Cell(out.GetStringWidth("googlegg"), ht, "googlegg")
	// html := `<a href="https://www.google.com" title="google" target="_blank">google search</a><br>` +
	// 	`<a href="https://www.google.com" target="_blank">google search</a><br>` +
	// 	`<a href="https://www.google.com" target="_blank">google search</a><br>`
	// hh := out.HTMLBasicNew()
	// hh.Write(ht, html)

	x, y := out.GetXY()
	out.AddPage()
	out.SetFont("default", "", 12)

	out.SetXY(10, 10)
	out.Cell(100, 10, "abc")
	out.CellFormat(100, ht, "aaaa", "", 1, "", false, 0, "abb")
	link := out.AddLink()
	out.Link(10, 10, 10, 10, link)
	out.SetPage(1)
	out.SetXY(x, y+10)
	out.CellFormat(100, ht, "test", "", 1, "", false, link, "link")
	out.SetPage(2)

	if err := out.OutputFileAndClose("/Users/huyungtang/Downloads/fpdf.pdf"); err != nil {
		t.Error(err)
	}

}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
