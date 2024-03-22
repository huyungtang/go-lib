package pdf

import (
	"testing"

	"github.com/huyungtang/go-lib/strings"
	"github.com/signintech/gopdf"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// TestMain
// ****************************************************************************************************************************************
func TestMain(t *testing.T) {
	pdf := gopdf.GoPdf{}
	defer pdf.Close()

	pdf.Start(gopdf.Config{PageSize: *gopdf.PageSizeA4, Unit: gopdf.UnitPT})
	err := pdf.AddTTFFont("default", "/Users/huyungtang/Projects/golang.batches/fonts/TaipeiSansTC/TaipeiSansTCBeta-Regular.ttf")
	if err != nil {
		t.Error(err)
	}
	err = pdf.AddTTFFont("barcode", "/Users/huyungtang/Projects/golang.batches/fonts/LibreBarcode128/LibreBarcode128-Regular.ttf")
	if err != nil {
		t.Error(err)
	}
	err = pdf.SetFont("default", "", 12)
	if err != nil {
		t.Error(err)
	}

	tpl := pdf.ImportPage("/Users/huyungtang/Downloads/文件範本.pdf", 1, "/MediaBox")
	pdf.AddPage()
	pdf.UseImportedTemplate(tpl, 0, 0, gopdf.PageSizeA4.W, gopdf.PageSizeA4.H)

	pdf.SetFont("barcode", "", 28)
	t.Log(pdf.GetY())
	pdf.SetXY(gopdf.UnitsToPoints(gopdf.UnitCM, 14.8), gopdf.UnitsToPoints(gopdf.UnitCM, 1.55))
	pdf.Cell(nil, strings.Code128A("E20230322001"))
	t.Log(pdf.GetY(), gopdf.UnitsToPoints(gopdf.UnitCM, 1.55), gopdf.UnitsToPoints(gopdf.UnitCM, 2.1))

	pdf.SetXY(gopdf.UnitsToPoints(gopdf.UnitCM, 14.8), gopdf.UnitsToPoints(gopdf.UnitCM, 2.1))
	pdf.Cell(nil, strings.Code128A("E20230322001"))

	pdf.SetFont("default", "", 12)
	pdf.SetXY(gopdf.UnitsToPoints(gopdf.UnitCM, 1), gopdf.UnitsToPoints(gopdf.UnitCM, 2.3))
	pdf.Cell(nil, "單據編號：E20230322001")

	pdf.SetXY(gopdf.UnitsToPoints(gopdf.UnitCM, 3.5), gopdf.UnitsToPoints(gopdf.UnitCM, 3.15))
	pdf.Cell(nil, "某某某")
	pdf.SetXY(gopdf.UnitsToPoints(gopdf.UnitCM, 9.7), gopdf.UnitsToPoints(gopdf.UnitCM, 3.15))
	pdf.Cell(nil, "YYYY 部門")
	pdf.SetXY(gopdf.UnitsToPoints(gopdf.UnitCM, 16.1), gopdf.UnitsToPoints(gopdf.UnitCM, 3.15))
	pdf.Cell(nil, "2024/03/22 18:33")

	if err = pdf.WritePdf("/Users/huyungtang/Downloads/pdf.pdf"); err != nil {
		t.Error(err)
	}
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
