package pdf

import (
	"io"

	"github.com/go-pdf/fpdf"
	"github.com/huyungtang/go-lib/strings"
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
		TextColorOption("#383838"),
		BorderColorOption("#383838"),
		FontSizeOption(6.5),
		FontRemOption(1),
		FontStyleOpiton("default"),
		PageSizeOption(PageSizeA4),
		OrientationOption(Portrait),
		UnitOption(Millimeter),
		PageMarginOpiton(10, 10, 10),
		AutoPageBreak(true, 10),
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
	SetXY(float64, float64) PDF
	Close()
}

// AddPage
// ****************************************************************************************************************************************
func (o *context) AddPage(args ...Option) {
	opt, _ := getOptions(o.options, args...)
	o.Fpdf.SetMargins(opt.marginLeft, opt.marginTop, opt.marginRight)
	o.Fpdf.SetAutoPageBreak(opt.autoPaged, opt.marginBottom)
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
	if w+o.Fpdf.GetX() > mw+1 {
		o.NewLine(args...)
	} else if w+o.Fpdf.GetX() == mw {
		ln = 1
	}

	r, g, b := opt.getTextColor()
	o.Fpdf.SetTextColor(r, g, b)
	o.Fpdf.SetFont(opt.font, "", opt.getFontSize())

	r, g, b = opt.getBorderColor()
	o.Fpdf.SetDrawColor(r, g, b)

	if opt.wrap {
		strs := make([]string, 0)
		rs := make([]rune, 0)
		for _, j := range txt {
			if j == 10 {
				strs = append(strs, string(rs))
				rs = make([]rune, 0)
				continue
			}

			if rs = append(rs, j); o.Fpdf.GetStringWidth(string(rs)) >= w-6 {
				strs = append(strs, string(rs))
				rs = make([]rune, 0)
			}
		}
		if len(rs) > 0 {
			strs = append(strs, string(rs))
		}
		o.Fpdf.MultiCell(w, opt.getLineHeight(), strings.Join(strs, "\n"), opt.getBorder(), opt.getAlign(), false)
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
func (o *context) SetXY(x, y float64) PDF {
	o.Fpdf.SetXY(x, y)

	return o
}

// Close
// ****************************************************************************************************************************************
func (o *context) Close() {
	o.Fpdf.Close()
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
