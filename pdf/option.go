package pdf

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/go-pdf/fpdf"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

const (
	Top position = 1 << iota
	Middle
	Baseline
	Left
	Center
	Right
	Bottom
	All position = Top | Left | Right | Bottom
)

const (
	PageSizeA4 pageSize = "A4"
	PageSizeA5 pageSize = "A5"

	Landscape orientation = "L"
	Portrait  orientation = "P"

	Centimeter unit = "cm"
	Inch       unit = "in"
	Millimeter unit = "mm"
	Point      unit = "pt"
)

var (
	pageSizes map[string]map[string]fpdf.SizeType = map[string]map[string]fpdf.SizeType{
		PageSizeA4: {
			Centimeter: {Wd: 21.0, Ht: 29.7},
			Inch:       {Wd: 8.27, Ht: 11.69},
			Millimeter: {Wd: 210, Ht: 297},
		},
		PageSizeA5: {
			Centimeter: {Wd: 14.8, Ht: 21.0},
			Inch:       {Wd: 5.83, Ht: 8.27},
			Millimeter: {Wd: 148, Ht: 210},
		},
	}
)

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// AlignOption
// ****************************************************************************************************************************************
func AlignOption(pos position) Option {
	return func(o *option) {
		o.align = pos
	}
}

// BorderOption
// ****************************************************************************************************************************************
func BorderOption(border position) Option {
	return func(o *option) {
		o.border = border
	}
}

// TextColorOption
// ****************************************************************************************************************************************
func TextColorOption(rgb string) Option {
	return func(o *option) {
		o.textColor = rgb
	}
}

// BorderColorOption
// ****************************************************************************************************************************************
func BorderColorOption(rgb string) Option {
	return func(o *option) {
		o.borderColor = rgb
	}
}

// FontOption
// ****************************************************************************************************************************************
func FontOption(name string, font string) Option {
	return func(o *option) {
		if bs, err := ioutil.ReadFile(font); err == nil {
			if o.fonts == nil {
				o.fonts = make(map[string][]byte)
			}
			o.fonts[strings.ToLower(name)] = bs
		}
	}
}

// FontStyleOpiton
// ****************************************************************************************************************************************
func FontStyleOpiton(name string) Option {
	return func(o *option) {
		o.font = name
	}
}

// FontSizeOption
// ****************************************************************************************************************************************
func FontSizeOption(size float64) Option {
	return func(o *option) {
		o.fontSize = size
	}
}

// FontRemOption
// ****************************************************************************************************************************************
func FontRemOption(rem float64) Option {
	return func(o *option) {
		o.fontRem = rem
	}
}

// OrientationOption
// ****************************************************************************************************************************************
func OrientationOption(orient orientation) Option {
	return func(o *option) {
		o.orientation = orient
	}
}

// PageMarginOpiton
// ****************************************************************************************************************************************
func PageMarginOpiton(top, left, right float64) Option {
	return func(o *option) {
		o.marginTop = top
		o.marginLeft = left
		o.marginRight = right
	}
}

// PageSizeOption
// ****************************************************************************************************************************************
func PageSizeOption(size pageSize) Option {
	return func(o *option) {
		o.pageSize = size
	}
}

// UnitOption
// ****************************************************************************************************************************************
func UnitOption(unit unit) Option {
	return func(o *option) {
		o.unit = unit
	}
}

// WidthOption
// ****************************************************************************************************************************************
func WidthOption(wd float64) Option {
	return func(o *option) {
		o.width = wd
	}
}

// HeightOption
// ****************************************************************************************************************************************
func HeightOption(ht float64) Option {
	return func(o *option) {
		o.height = ht
	}
}

// WrapOption
// ****************************************************************************************************************************************
func WrapOption(wrap bool) Option {
	return func(o *option) {
		o.wrap = wrap
	}
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// option *********************************************************************************************************************************
type option struct {
	font        string
	fonts       map[string][]byte
	fontSize    float64
	fontRem     float64
	textColor   string
	borderColor string
	align       position
	border      position
	pageSize    string
	orientation string
	unit        string
	marginTop   float64
	marginLeft  float64
	marginRight float64
	width       float64
	height      float64
	wrap        bool
}

// getAlign *******************************************************************************************************************************
func (o *option) getAlign() string {
	bs := make([]string, 0, 4)
	if o.align&Left == Left {
		bs = append(bs, "L")
	}
	if o.align&Center == Center {
		bs = append(bs, "C")
	}
	if o.align&Right == Right {
		bs = append(bs, "R")
	}
	if o.align&Top == Top {
		bs = append(bs, "T")
	}
	if o.align&Middle == Middle {
		bs = append(bs, "M")
	}
	if o.align&Bottom == Bottom {
		bs = append(bs, "B")
	}
	if o.align&Baseline == Baseline {
		bs = append(bs, "A")
	}

	return strings.Join(bs, "")
}

// getBorder ******************************************************************************************************************************
func (o *option) getBorder() string {
	bs := make([]string, 0, 4)
	if o.border&Top == Top {
		bs = append(bs, "T")
	}
	if o.border&Left == Left {
		bs = append(bs, "L")
	}
	if o.border&Right == Right {
		bs = append(bs, "R")
	}
	if o.border&Bottom == Bottom {
		bs = append(bs, "B")
	}

	return strings.Join(bs, "")
}

// getCellWidth ***************************************************************************************************************************
func (o *option) getCellWidth(pdf *fpdf.Fpdf) (float64, float64) {
	wd, _ := pdf.GetPageSize()
	ml, _, mr, _ := pdf.GetMargins()

	fw := (wd - ml - mr)

	return fw * o.width, wd - mr
}

// getFontSize ****************************************************************************************************************************
func (o *option) getFontSize() float64 {
	return o.fontSize * o.fontRem
}

// getLineHeight **************************************************************************************************************************
func (o *option) getLineHeight() float64 {
	if o.height != 0 {
		return o.fontSize * o.height * 0.6
	}

	return o.fontSize * o.fontRem * 0.6
}

// getTextColor ***************************************************************************************************************************
func (o *option) getTextColor() (r, g, b int) {
	fmt.Sscanf(strings.ToUpper(o.textColor), "#%2X%2X%2X", &r, &g, &b)

	return
}

// getBorderColor *************************************************************************************************************************
func (o *option) getBorderColor() (r, g, b int) {
	fmt.Sscanf(strings.ToUpper(o.borderColor), "#%2X%2X%2X", &r, &g, &b)

	return
}

// position *******************************************************************************************************************************
type position = int

// pageSize *******************************************************************************************************************************
type pageSize = string

// orientation ****************************************************************************************************************************
type orientation = string

// unit ***********************************************************************************************************************************
type unit = string

// Option
// ****************************************************************************************************************************************
type Option func(*option)

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// getOptions *****************************************************************************************************************************
func getOptions(defa []Option, args ...Option) (*option, []Option) {
	opts := make([]Option, len(defa), len(defa)+len(args))
	copy(opts, defa)
	opts = append(opts, args...)

	opt := new(option)
	for _, optFn := range opts {
		optFn(opt)
	}

	return opt, opts
}
