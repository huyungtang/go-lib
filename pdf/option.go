package pdf

import (
	"github.com/go-pdf/fpdf"
	"github.com/huyungtang/go-lib/file"
	"github.com/huyungtang/go-lib/strings"
	"github.com/phpdave11/gofpdi"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// FontPathOption
// ****************************************************************************************************************************************
func FontPathOption(path string) PageOption {
	return func(ctx *context) {
		if fs, err := file.ListFiles(path, `\.ttf$`); err == nil {
			for _, f := range fs {
				fn := file.GetFilename(f)
				ctx.fpdf.AddUTF8Font(fn[0:strings.LastIndex(fn, ".")], "", f)
			}
		}
	}
}

// set page margins and auto page-break
//	`auto`	auto page-break
// ****************************************************************************************************************************************
func PageMarginsOption(left, top, right, bottom float64, auto bool) PageOption {
	return func(ctx *context) {
		ctx.fpdf.SetMargins(left, top, right)
		ctx.fpdf.SetAutoPageBreak(auto, bottom)
	}
}

// PageSizeA4Option
// ****************************************************************************************************************************************
func PageSizeA4Option(landscape bool) PageOption {
	return pageSizeOption(fpdf.SizeType{Wd: 210, Ht: 297}, landscape)
}

// PageSizeA5Option
// ****************************************************************************************************************************************
func PageSizeA5Option(landscape bool) PageOption {
	return pageSizeOption(fpdf.SizeType{Wd: 148, Ht: 210}, landscape)
}

// TemplateOption
// ****************************************************************************************************************************************
func TemplateOption(page int) PageOption {
	return func(ctx *context) {
		if page > -1 && ctx.pageIndex != ctx.fpdf.PageNo() && ctx.importer != nil && page < ctx.importer.GetNumPages() {
			w, h := ctx.fpdf.GetPageSize()
			tn, sX, sY, tX, tY := ctx.importer.UseTemplate(page, 0, 0, w, h)
			ctx.fpdf.UseImportedTemplate(tn, sX, sY, tX, tY)
		}
	}
}

// TemplatesOption
// ****************************************************************************************************************************************
func TemplatesOption(path string) PageOption {
	return func(ctx *context) {
		if ctx.importer == nil {
			ctx.importer = gofpdi.NewImporter()
			ctx.importer.SetSourceFile(path)
			for i := 1; i <= ctx.importer.GetNumPages(); i++ {
				ctx.importer.ImportPage(i, "/MediaBox")
			}

			ctx.fpdf.ImportTemplates(ctx.importer.PutFormXobjectsUnordered())
			ctx.fpdf.ImportObjects(ctx.importer.GetImportedObjectsUnordered())
			ctx.fpdf.ImportObjPos(ctx.importer.GetImportedObjHashPos())
		}
	}
}

// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// BorderColorOption
// ****************************************************************************************************************************************
func BorderColorOption(rgb string) CellOption {
	return func(ctx *context) {
		if r, g, b, err := rgbColor(rgb); err == nil {
			ctx.fpdf.SetDrawColor(r, g, b)
		}
	}
}

// AlignOption
// ****************************************************************************************************************************************
func AlignOption(align align) CellOption {
	return func(ctx *context) {
		ctx.cellAlign = align
	}
}

// BorderOption
// ****************************************************************************************************************************************
func BorderOption(border border) CellOption {
	return func(ctx *context) {
		ctx.cellBorder = border
	}
}

// MarginsOption
// ****************************************************************************************************************************************
func MarginsOption(m float64) CellOption {
	return func(ctx *context) {
		ctx.fpdf.SetCellMargin(m)
	}
}

// HeightOption
// ****************************************************************************************************************************************
func HeightOption(ht float64) CellOption {
	return func(ctx *context) {
		ctx.cellHeight = ht
	}
}

// WidthOption
// ****************************************************************************************************************************************
func WidthOption(wd float64) CellOption {
	return func(ctx *context) {
		ctx.cellWidth = wd

		if ctx.cellWidth <= 1 {
			w, _, _ := ctx.fpdf.PageSize(ctx.fpdf.PageNo())
			l, _, r, _ := ctx.fpdf.GetMargins()

			ctx.cellWidth = (w - l - r) * wd
		}
	}
}

// FontFamilyOption
// ****************************************************************************************************************************************
func FontFamilyOption(style string) CellOption {
	return func(ctx *context) {
		ctx.fontFamily = style
		ctx.fpdf.SetFont(ctx.fontFamily, "", ctx.fontSize)
	}
}

// FontSizeOption
// ****************************************************************************************************************************************
func FontSizeOption(size float64) CellOption {
	return func(ctx *context) {
		ctx.fontSize = size
		ctx.fpdf.SetFont(ctx.fontFamily, "", ctx.fontSize)
	}
}

// LocationOption
// ****************************************************************************************************************************************
func LocationOption(x, y float64) CellOption {
	return func(ctx *context) {
		ctx.fpdf.SetXY(x, y)
	}
}

// PositionTailOption
// ****************************************************************************************************************************************
func PositionTailOption() CellOption {
	return positionOption(0)
}

// PositionNewLineOption
// ****************************************************************************************************************************************
func PositionNewLineOption() CellOption {
	return positionOption(1)
}

// PositionBottomOption
// ****************************************************************************************************************************************
func PositionBottomOption() CellOption {
	return positionOption(2)
}

// TextColorOption
// ****************************************************************************************************************************************
func TextColorOption(rgb string) CellOption {
	return func(ctx *context) {
		if r, g, b, err := rgbColor(rgb); err == nil {
			ctx.fpdf.SetTextColor(r, g, b)
		}
	}
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// CellOption
// ****************************************************************************************************************************************
type CellOption = option

// PageOption
// ****************************************************************************************************************************************
type PageOption = option

// option *********************************************************************************************************************************
type option func(*context)

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// pageSizeOption *************************************************************************************************************************
func pageSizeOption(size fpdf.SizeType, landscape bool) option {
	return func(ctx *context) {
		if landscape {
			size.Wd, size.Ht = size.Ht, size.Wd
		}
		ctx.fpdf.AddPageFormat("P", size)
	}
}

// positionOption *************************************************************************************************************************
func positionOption(pos int) CellOption {
	return func(ctx *context) {
		ctx.position = pos
	}
}

// rgbColor *******************************************************************************************************************************
func rgbColor(rgb string) (r, g, b int, err error) {
	err = strings.Parse(strings.ToUpper(rgb), "#%2X%2X%2X", &r, &g, &b)

	return
}
