package pdf

import (
	"github.com/go-pdf/fpdf"
	"github.com/huyungtang/go-lib/strings"
	"github.com/phpdave11/gofpdi"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// context ********************************************************************************************************************************
type context struct {
	*fpdf.Fpdf
	*gofpdi.Importer

	opts []Option
}

// AddPage
// ****************************************************************************************************************************************
func (o *context) AddPage(opts ...Option) PDF {
	cfg := applyOption(append(o.opts, opts...)...)
	o.Fpdf.SetMargins(cfg.pageLeft, cfg.pageTop, cfg.pageRight)
	o.Fpdf.AddPageFormat(cfg.orientation, cfg.pageSize.size)

	if cfg.template > -1 && o.Importer != nil && cfg.template < o.Importer.GetNumPages() {
		w, h := o.Fpdf.GetPageSize()
		tn, sX, sY, tX, tY := o.Importer.UseTemplate(0, 0, 0, w, h)
		o.Fpdf.UseImportedTemplate(tn, sX, sY, tX, tY)
	}

	return o
}

// AddBarcode128
// ****************************************************************************************************************************************
func (o *context) AddBarcode128(txt string, opts ...Option) PDF {
	opts = append(opts,
		FontSizeOption(20),
		CellWidthOption(50),
		CellHeightOption(4),
		CellMaringOption(0),
		AlignOption(AlignMC),
		PositionOption(PositionBottom),
	)
	txt = strings.Code128A(txt)

	return o.AddCell(txt, opts...).AddCell(txt, opts...)
}

// AddCell
// ****************************************************************************************************************************************
func (o *context) AddCell(txt string, opts ...Option) PDF {
	cfg := applyOption(append(o.opts, opts...)...)
	o.Fpdf.SetFont(cfg.fontFamily, "", cfg.fontSize)
	o.Fpdf.SetDrawColor(cfg.borderColor[0], cfg.borderColor[1], cfg.borderColor[2])
	o.Fpdf.SetTextColor(cfg.textColor[0], cfg.textColor[1], cfg.textColor[2])
	o.Fpdf.SetCellMargin(cfg.cellMargin)

	cfg.cellWidth = o.getCellWidth(cfg.cellWidth)
	strs := o.getCellText(txt, cfg.cellWidth, cfg.cellMargin)
	cfg.cellHeight = o.getCellHeight(cfg.cellHeight, len(strs), cfg)

	if l := len(strs); l == 1 {
		o.Fpdf.CellFormat(cfg.cellWidth, cfg.cellHeight, strs[0], cfg.cellBorder, cfg.position, cfg.cellAlign, false, 0, "")
	} else if l > 1 {
		cx, cy := o.GetXY()
		o.Fpdf.MultiCell(cfg.cellWidth, cfg.cellHeight, strings.Join(strs, "\n"), cfg.cellBorder, cfg.cellAlign, false)

		nx, ny := o.GetXY()
		switch cfg.position {
		case 0:
			o.Fpdf.SetXY(nx, cy)
		case 2:
			o.Fpdf.SetXY(cx, ny)
		}
	}

	return o
}

// GetXY()
// ****************************************************************************************************************************************
func (o *context) GetXY() (x, y float64) {
	if o.Fpdf != nil {
		x, y = o.Fpdf.GetXY()
	}

	return
}

// SetXY
// ****************************************************************************************************************************************
func (o *context) SetXY(x, y float64) PDF {
	if o.Fpdf != nil {
		o.Fpdf.SetXY(x, y)
	}

	return o
}

// Close
// ****************************************************************************************************************************************
func (o *context) Close() (err error) {
	if o.Fpdf != nil {
		o.Fpdf.Close()
	}

	return
}

// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// getCellText ****************************************************************************************************************************
func (o *context) getCellText(txt string, wd, mg float64) (strs []string) {
	strs = make([]string, 0)
	tstr := make([]rune, 0)
	wd = wd - (mg * 2)
	for _, s := range txt {
		switch s {
		case 13:
		case 10:
			strs = append(strs, string(tstr))
			tstr = make([]rune, 0)
		default:
			tstr = append(tstr, s)
			if o.Fpdf.GetStringWidth(string(tstr)) > wd {
				strs = append(strs, string(tstr[0:len(tstr)-1]))
				tstr = tstr[len(tstr)-1:]
			}
		}
	}
	strs = append(strs, string(tstr))

	return
}

// getCellWidth ***************************************************************************************************************************
func (o *context) getCellWidth(wd float64) float64 {
	if wd <= 1 {
		w, _, _ := o.Fpdf.PageSize(o.Fpdf.PageNo())
		l, _, r, _ := o.Fpdf.GetMargins()
		return (w - l - r) * wd
	}

	return wd
}

// getCellHeight **************************************************************************************************************************
func (o *context) getCellHeight(ht float64, ln int, cfg *option) float64 {
	if ht == 0 {
		_, h := o.Fpdf.GetFontSize()
		if ln <= 1 {
			return h + (cfg.cellMargin * 2)
		} else {
			return h + cfg.cellMargin
		}
	}

	return ht
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
