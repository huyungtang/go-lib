package pdf

import (
	"fmt"
	"io"
	bases "strings"

	"github.com/go-pdf/fpdf"
	"github.com/huyungtang/go-lib/barcode"
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
	o.applyOptions(o.cellOptions, opts...)

	if br, err := barcode.Code128(text, 280, 280); err == nil {
		k := strings.Random(32)
		o.fpdf.RegisterImageReader(k, "png", br)
		x, y := o.fpdf.GetXY()
		o.fpdf.Image(k, x, y, o.cellWidth, o.cellHeight, false, "png", 0, "")
	}

	return o
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

// AddRow
// ****************************************************************************************************************************************
func (o *context) AddRow(dtos []*PDFRowDTO) PDF {
	x, y := o.GetXY()
	my := y
	for _, dto := range dtos {
		dto.Opts = append(dto.Opts, Location(x, y), Below())
		if _, cy := o.AddCell(dto.Text, dto.Opts...).GetXY(); cy > my {
			my = cy
		}
		x += o.cellWidth
	}

	x, _, _, _ = o.fpdf.GetMargins()
	o.SetXY(x, my)

	return o
}

// AddLink
// ****************************************************************************************************************************************
func (o *context) AddLink(txt, url string, opts ...option) PDF {
	if url != "" {
		o.applyOptions(o.cellOptions, opts...)

		if txt == "" {
			txt = url
		}

		x, y := o.GetXY()
		_, ht := o.fpdf.GetFontSize()
		ht = ht + (o.fpdf.GetCellMargin() * 2)

		if wd := o.fpdf.GetStringWidth(txt) + (o.fpdf.GetCellMargin() * 2); o.cellWidth < wd {
			rs := []rune(txt)
			tail := string(rs[len(rs)-3:])

			var bs bases.Builder
			bs.Grow(32)
			for i := len(rs) - 4; i > 0; i-- {
				fmt.Fprintf(&bs, "%s…%s", string(rs[0:i]), tail)
				if wd = o.fpdf.GetStringWidth(bs.String()) + (o.fpdf.GetCellMargin() * 2); wd <= o.cellWidth {
					txt = bs.String()
					break
				}
				bs.Reset()
			}
		}

		writer := o.fpdf.HTMLBasicNew()
		writer.Write(ht, strings.Format(`<a href="%s" target="_blank">%s</a>`, url, txt))

		o.fpdf.SetXY(x, y)
		o.fpdf.CellFormat(o.cellWidth, ht, "", o.cellBorder, o.position, o.cellAlign, false, 0, "")
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

// SetXY
// ****************************************************************************************************************************************
func (o *context) SetXY(x, y float64) PDF {
	o.fpdf.SetXY(x, y)

	return o
}

// GetContainerSize
// ****************************************************************************************************************************************
func (o *context) GetContainerSize() (wd, ht float64) {
	wd, ht = o.fpdf.GetPageSize()
	l, t, r, b := o.fpdf.GetMargins()
	wd = wd - l - r
	ht = ht - t - b

	return
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

	if len(tstr) > 0 {
		strs = append(strs, string(tstr))
	}

	if len(strs) == 0 {
		strs = append(strs, "")
	}

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
