package pdf

import (
	"os"
	"testing"
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
		FontPathOption("/Users/huyungtang/Projects/golang.batches/fonts/"),
		FontFamilyOption("TaipeiSansTCBeta-Regular"),
		TemplateFileOption("/Users/huyungtang/Downloads/文件範本.pdf"),
		FontSizeOption(9),
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
		SetXY(10, 50).
		AddCell("永盛車電股份有限公司", AlignOption(AlignMC), FontSizeOption(15), CellHeightOption(6)).
		AddCell("請購單", AlignOption(AlignMC), FontSizeOption(13), CellHeightOption(7)).
		AddCell("申請人", CellWidthOption(.13), AlignOption(AlignMC), BorderOption(BorderFull), PositionOption(PositionTail)).
		AddCell("張三", CellWidthOption(.2), BorderOption(BorderFull), PositionOption(PositionTail)).
		AddCell("部門", CellWidthOption(.13), AlignOption(AlignMC), BorderOption(BorderFull), PositionOption(PositionTail)).
		AddCell("XX 部", CellWidthOption(.2), BorderOption(BorderFull), PositionOption(PositionTail)).
		AddCell("申請日期", CellWidthOption(.13), AlignOption(AlignMC), BorderOption(BorderFull), PositionOption(PositionTail)).
		AddCell("2024/03/26 11:22", CellWidthOption(.21), BorderOption(BorderFull)).
		SetXY(150, 20).
		AddBarcode128("E12345678901111", FontFamilyOption("LibreBarcode128-Regular")).
		SetXY(33, 29.5).AddCell("張三", CellWidthOption(40), CellHeightOption(9)).
		SetXY(96, 29.5).AddCell("XX 部", CellWidthOption(40), CellHeightOption(9)).
		SetXY(159, 29.5).AddCell("2024/03/26 11:22", CellWidthOption(40), CellHeightOption(9)).
		Output(f); err != nil {
		t.Error(err)
	}
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************