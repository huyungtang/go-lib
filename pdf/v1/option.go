package pdf

import (
	"fmt"
	"strings"

	"github.com/go-pdf/fpdf"
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

	PositionTail    position = 0
	PositionNewLine position = 1
	PositionBottom  position = 2
)

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// AlignOption
// ****************************************************************************************************************************************
func AlignOption(align align) Option {
	return func(o *option) {
		o.cellAlign = align
	}
}

// BorderOption
// ****************************************************************************************************************************************
func BorderOption(border border) Option {
	return func(o *option) {
		o.cellBorder = border
	}
}

// BorderColor
// ****************************************************************************************************************************************
func BorderColor(rgb string) Option {
	return func(o *option) {
		var r, g, b int
		fmt.Sscanf(strings.ToUpper(rgb), "#%2X%2X%2X", &r, &g, &b)
		o.borderColor = []int{r, g, b}
	}
}

// TextColor
// ****************************************************************************************************************************************
func TextColor(rgb string) Option {
	return func(o *option) {
		var r, g, b int
		fmt.Sscanf(strings.ToUpper(rgb), "#%2X%2X%2X", &r, &g, &b)
		o.textColor = []int{r, g, b}
	}
}

// cell left & right margins
//	`margin`	default 2
// ****************************************************************************************************************************************
func CellMaringOption(margin float64) Option {
	return func(o *option) {
		o.cellMargin = margin
	}
}

// set cell width
//	`wd`	<=1 percentage of page available width
// ****************************************************************************************************************************************
func CellWidthOption(wd float64) Option {
	return func(o *option) {
		o.cellWidth = wd
	}
}

// set cell height
//	`ht`	0 auto height
// ****************************************************************************************************************************************
func CellHeightOption(ht float64) Option {
	return func(o *option) {
		o.cellHeight = ht
	}
}

// FontFamilyOption
// ****************************************************************************************************************************************
func FontFamilyOption(fam string) Option {
	return func(o *option) {
		o.fontFamily = fam
	}
}

// FontPathOption
// ****************************************************************************************************************************************
func FontPathOption(path string) Option {
	return func(o *option) {
		o.ttfPath = path
	}
}

// FontSizeOption
// ****************************************************************************************************************************************
func FontSizeOption(size float64) Option {
	return func(o *option) {
		o.fontSize = size
	}
}

// PortraitOption
// ****************************************************************************************************************************************
func PortraitOption() Option {
	return func(o *option) {
		o.orientation = "Portrait"
	}
}

// LandscapeOption
// ****************************************************************************************************************************************
func LandscapeOption() Option {
	return func(o *option) {
		o.orientation = "Landscape"
	}
}

// PageMarginsOption
// ****************************************************************************************************************************************
func PageMarginsOption(left, top, right float64) Option {
	return func(o *option) {
		o.pageLeft, o.pageTop, o.pageRight = left, top, right
	}
}

// PageSizeOption
// ****************************************************************************************************************************************
func PageSizeOption(ps pagesize) Option {
	return func(o *option) {
		o.pageSize = ps
	}
}

// PageSizeA4Option
// ****************************************************************************************************************************************
func PageSizeA4Option() Option {
	return func(o *option) {
		o.pageSize = pagesize{name: "A4", size: fpdf.SizeType{Wd: 210, Ht: 297}}
	}
}

// define the point after cell
//	`ln`	PositionTail, NewLine or Buttom
// ****************************************************************************************************************************************
func PositionOption(ln position) Option {
	return func(o *option) {
		o.position = ln
	}
}

// use page of template file as page template
//	`page`	page idnex start with 0
// ****************************************************************************************************************************************
func TemplateOption(page int) Option {
	return func(o *option) {
		o.template = page
	}
}

// TemplateFileOption
// ****************************************************************************************************************************************
func TemplateFileOption(path string) Option {
	return func(o *option) {
		o.templateFile = path
	}
}

// UnitMM
// ****************************************************************************************************************************************
func UnitMM() Option {
	return func(o *option) {
		o.unit = "mm"
	}
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// align **********************************************************************************************************************************
type align = string

// border *********************************************************************************************************************************
type border = string

// position *******************************************************************************************************************************
type position = int

// pagesize *******************************************************************************************************************************
type pagesize struct {
	name string
	size fpdf.SizeType
}

// option *********************************************************************************************************************************
type option struct {
	cellAlign    align
	cellBorder   border
	cellMargin   float64
	cellWidth    float64
	cellHeight   float64
	borderColor  []int
	textColor    []int
	fontFamily   string
	fontSize     float64
	orientation  string
	pageLeft     float64
	pageTop      float64
	pageRight    float64
	pageSize     pagesize
	position     position
	template     int
	templateFile string
	ttfPath      string
	unit         string
}

// Option
// ****************************************************************************************************************************************
type Option func(*option)

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// applyOption ****************************************************************************************************************************
func applyOption(opts ...Option) (opt *option) {
	opt = new(option)
	for _, o := range opts {
		o(opt)
	}

	return
}
