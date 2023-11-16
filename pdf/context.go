package pdf

import (
	"io"
	"math"

	"github.com/go-pdf/fpdf"
	"github.com/huyungtang/go-lib/slices"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Init
// ****************************************************************************************************************************************
func Init(args ...Option) PDF {
	opt, opts := getOptions([]Option{
		ColorOption("#383838"),
		FontSizeOption(6.5),
		FontRemOption(1),
		FontStyleOpiton("default"),
		PageSizeOption(PageSizeA4),
		OrientationOption(Portrait),
		UnitOption(Millimeter),
		PageMarginOpiton(10, 10, 10),
		WidthOption(1),
	}, args...)

	ctx := &context{
		fpdf.New(opt.orientation, opt.unit, opt.pageSize, ""),
		opts,
	}
	ctx.Fpdf.SetFont("Arial", "B", opt.getFontSize())
	for i, j := range opt.fonts {
		ctx.Fpdf.AddUTF8FontFromBytes(i, "", j)
		if opt.font == "" {
			opt.font = i
			ctx.Fpdf.SetFont(opt.font, "", opt.getFontSize())
		}
	}
	ctx.AddPage()

	return ctx
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// context ********************************************************************************************************************************
type context struct {
	*fpdf.Fpdf
	options []Option
}

// PDF
// ****************************************************************************************************************************************
type PDF interface {
	AddPage(...Option)
	NewLine(...Option)
	Output(io.Writer) error
	Text(string, ...Option) PDF
	GetXY() []float64
	SetXY(float64, float64)
	Close()
}

// AddPage
// ****************************************************************************************************************************************
func (o *context) AddPage(args ...Option) {
	opt, _ := getOptions(o.options, args...)
	o.Fpdf.SetMargins(opt.marginLeft, opt.marginTop, opt.marginRight)
	o.Fpdf.AddPageFormat(opt.orientation, pageSizes[opt.pageSize][opt.unit])
}

// NewLine
// ****************************************************************************************************************************************
func (o *context) NewLine(args ...Option) {
	opt, _ := getOptions(o.options, args...)
	o.Fpdf.Ln(opt.getLineHeight())
}

// Output
// ****************************************************************************************************************************************
func (o *context) Output(writer io.Writer) (err error) {
	return o.Fpdf.Output(writer)
}

// Text
// ****************************************************************************************************************************************
func (o *context) Text(txt string, args ...Option) PDF {
	opt, _ := getOptions(o.options, args...)

	ln := 0
	w, mw := opt.getCellWidth(o.Fpdf)
	if w+o.Fpdf.GetX() > mw {
		o.NewLine(args...)
	} else if w+o.Fpdf.GetX() == mw {
		ln = 1
	}

	r, g, b := opt.getTextColor()
	o.Fpdf.SetTextColor(r, g, b)
	o.Fpdf.SetFont(opt.font, "", opt.getFontSize())

	if opt.wrap {
		_, uu := o.Fpdf.GetFontSize()
		ml := int(math.Round(w * .96 / uu))
		rs := []rune(txt)
		rr := make([]rune, 0)
		for l := len(rs); l > 0; {
			x, isMatched := slices.IndexOf(l, func(i int) bool { return rs[i] == 10 })
			if !isMatched {
				x = l
			}
			if x > ml {
				x = ml
			}
			rr = append(rr, rs[0:x]...)
			rr = append(rr, 10)
			if isMatched {
				x++
			}
			rs = rs[x:]
			l = len(rs)
		}

		o.Fpdf.MultiCell(w, opt.getLineHeight(), string(rr), opt.getBorder(), opt.getAlign(), false)
	} else {
		o.Fpdf.CellFormat(w, opt.getLineHeight(), txt, opt.getBorder(), ln, opt.getAlign(), false, 0, "")
	}

	return o
}

// GetXY
// ****************************************************************************************************************************************
func (o *context) GetXY() []float64 {
	x, y := o.Fpdf.GetXY()
	return []float64{x, y}
}

// SetXY
// ****************************************************************************************************************************************
func (o *context) SetXY(x, y float64) {
	o.Fpdf.SetXY(x, y)
}

// Close
// ****************************************************************************************************************************************
func (o *context) Close() {
	o.Fpdf.Close()
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
