package pdf

import (
	"io"

	"github.com/go-pdf/fpdf"
	"github.com/huyungtang/go-lib/strings"
	"github.com/phpdave11/gofpdi"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

const (
	Baseline string = "A"
	Bottom   string = "B"
	Center   string = "C"
	Left     string = "L"
	Middle   string = "M"
	Right    string = "R"
	Top      string = "T"

	BorderTRB  border = "TRB"
	BorderLRB  border = "LRB"
	BorderRB   border = "RB"
	BorderFull border = "1"

	AlignTL align = "TL"
	AlignTC align = "TC"
	AlignTR align = "TR"
	AlignML align = "ML"
	AlignMC align = "MC"
	AlignMR align = "MR"
	AlignBL align = "BL"
	AlignBC align = "BC"
	AlignBR align = "BR"
)

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// context ********************************************************************************************************************************
type context struct {
	fpdf                     *fpdf.Fpdf
	importer                 *gofpdi.Importer
	cellOptions, pageOptions []option
	pageIndex                int
	cellWidth, cellHeight    float64
	cellBorder, cellAlign    border
	position                 int
	fontFamily               string
	fontSize                 float64
}

// AddPage
// ****************************************************************************************************************************************
func (o *context) AddPage(opts ...option) PDF {
	o.pageIndex = o.fpdf.PageNo()
	o.applyOptions(o.pageOptions, opts...)
	if o.pageIndex == o.fpdf.PageNo() {
		o.fpdf.AddPage()
		o.applyOptions(o.pageOptions, opts...)
	}

	return o
}

// AddBarcode128
// ****************************************************************************************************************************************
func (o *context) AddBarcode128(text string, opts ...option) PDF {
	text = strings.Code128A(text)

	cpts := append(opts,
		FontSizeOption(20),
		CellWidthOption(50),
		CellHeightOption(4),
		CellMarginsOption(0),
		CellAlignOption(AlignMC),
		PositionBottomOption())
	x, y := o.AddCell(text, cpts...).GetXY()

	cpts = append(cpts, LocationOption(x, y), PositionNewLineOption())

	return o.AddCell(text, cpts...)
}

// AddCell
// ****************************************************************************************************************************************
func (o *context) AddCell(text string, opts ...option) PDF {
	o.applyOptions(o.cellOptions, opts...)

	strs := o.getCellText(text)

	ht, mg := o.cellHeight, o.fpdf.GetCellMargin()
	if ht == 0 {
		_, ht = o.fpdf.GetFontSize()
	}

	if l := len(strs); l == 1 {
		o.fpdf.CellFormat(o.cellWidth, ht+(mg*2), strs[0], o.cellBorder, o.position, o.cellAlign, false, 0, "")
	} else if l > 1 {
		cx, cy := o.GetXY()
		o.fpdf.MultiCell(o.cellWidth, ht+(mg*2), strings.Join(strs, "\n"), o.cellBorder, o.cellAlign, false)

		_, ny := o.GetXY()
		switch o.position {
		case 0:
			o.fpdf.SetXY(cx+o.cellWidth, cy)
		case 2:
			o.fpdf.SetXY(cx, ny)
		}
	}

	return o
}

// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// GetXY
// ****************************************************************************************************************************************
func (o *context) GetXY() (x float64, y float64) {

	return o.fpdf.GetXY()
}

// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Output
// ****************************************************************************************************************************************
func (o *context) Output(w io.Writer) (err error) {

	return o.fpdf.Output(w)
}

// Close
// ****************************************************************************************************************************************
func (o *context) Close() (err error) {
	o.fpdf.Close()

	return
}

// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// applyOptions ***************************************************************************************************************************
func (o *context) applyOptions(defa []option, opts ...option) {
	for _, opt := range defa {
		opt(o)
	}

	for _, opt := range opts {
		opt(o)
	}
}

// getCellText ****************************************************************************************************************************
func (o *context) getCellText(text string) (strs []string) {
	wd, mg := o.cellWidth, o.fpdf.GetCellMargin()

	strs = make([]string, 0)
	tstr := make([]rune, 0)
	wd = wd - (mg * 2)
	for _, s := range text {
		switch s {
		case 13:
		case 10:
			strs = append(strs, string(tstr))
			tstr = make([]rune, 0)
		default:
			tstr = append(tstr, s)
			if o.fpdf.GetStringWidth(string(tstr)) > wd {
				strs = append(strs, string(tstr[0:len(tstr)-1]))
				tstr = tstr[len(tstr)-1:]
			}
		}
	}
	strs = append(strs, string(tstr))

	return
}

// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// align **********************************************************************************************************************************
type align = string

// border *********************************************************************************************************************************
type border = string

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
