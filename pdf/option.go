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

// FontPath
// ****************************************************************************************************************************************
func FontPath(path string) PageOption {
	return func(ctx *pdfContext) {
		if fs, err := file.ListFiles(path, `\.ttf$`); err == nil {
			for _, f := range fs {
				fn := file.GetFilename(f)
				ctx.fpdf.AddUTF8Font(fn[0:strings.LastIndex(fn, ".")], "", f)
			}
		}
	}
}

// set page margins and auto page-break
//
//	`auto`	auto page-break
//
// ****************************************************************************************************************************************
func PageMargins(left, top, right, bottom float64, auto bool) PageOption {
	return func(ctx *pdfContext) {
		ctx.fpdf.SetMargins(left, top, right)
		ctx.fpdf.SetAutoPageBreak(auto, bottom)
	}
}

// PageSizeA4
// ****************************************************************************************************************************************
func PageSizeA4(landscape bool) PageOption {
	return pageSize(fpdf.SizeType{Wd: 210, Ht: 297}, landscape)
}

// PageSizeA5
// ****************************************************************************************************************************************
func PageSizeA5(landscape bool) PageOption {
	return pageSize(fpdf.SizeType{Wd: 148, Ht: 210}, landscape)
}

// Template
// ****************************************************************************************************************************************
func Template(page int) PageOption {
	return func(ctx *pdfContext) {
		if page > -1 && ctx.pageIndex != ctx.fpdf.PageNo() && ctx.importer != nil && page < ctx.importer.GetNumPages() {
			w, h := ctx.fpdf.GetPageSize()
			tn, sX, sY, tX, tY := ctx.importer.UseTemplate(page, 0, 0, w, h)
			ctx.fpdf.UseImportedTemplate(tn, sX, sY, tX, tY)
		}
	}
}

// Templates
// ****************************************************************************************************************************************
func Templates(path string) PageOption {
	return func(ctx *pdfContext) {
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

// BorderColor
// ****************************************************************************************************************************************
func BorderColor(rgb string) CellOption {
	return func(ctx *pdfContext) {
		if r, g, b, err := rgbColor(rgb); err == nil {
			ctx.fpdf.SetDrawColor(r, g, b)
		}
	}
}

// Align
// ****************************************************************************************************************************************
func Align(align align) CellOption {
	return func(ctx *pdfContext) {
		ctx.cellAlign = align
	}
}

// Border
// ****************************************************************************************************************************************
func Border(border border) CellOption {
	return func(ctx *pdfContext) {
		ctx.cellBorder = border
	}
}

// Margins
// ****************************************************************************************************************************************
func Margins(m float64) CellOption {
	return func(ctx *pdfContext) {
		ctx.fpdf.SetCellMargin(m)
	}
}

// Height
// ****************************************************************************************************************************************
func Height(ht float64) CellOption {
	return func(ctx *pdfContext) {
		ctx.cellHeight = ht
	}
}

// Width
// ****************************************************************************************************************************************
func Width(wd float64) CellOption {
	return func(ctx *pdfContext) {
		ctx.cellWidth = wd

		if ctx.cellWidth <= 1 {
			w, _, _ := ctx.fpdf.PageSize(ctx.fpdf.PageNo())
			l, _, r, _ := ctx.fpdf.GetMargins()

			ctx.cellWidth = (w - l - r) * wd
		}
	}
}

// FontFamily
// ****************************************************************************************************************************************
func FontFamily(style string) CellOption {
	return func(ctx *pdfContext) {
		ctx.fontFamily = style
		ctx.fpdf.SetFont(ctx.fontFamily, "", ctx.fontSize)
	}
}

// FontSize
// ****************************************************************************************************************************************
func FontSize(size float64) CellOption {
	return func(ctx *pdfContext) {
		ctx.fontSize = size
		ctx.fpdf.SetFont(ctx.fontFamily, "", ctx.fontSize)
	}
}

// Location
// ****************************************************************************************************************************************
func Location(x, y float64) CellOption {
	return func(ctx *pdfContext) {
		ctx.fpdf.SetXY(x, y)
	}
}

// Tail
// ****************************************************************************************************************************************
func Tail() CellOption {
	return position(0)
}

// NewLine
// ****************************************************************************************************************************************
func NewLine() CellOption {
	return position(1)
}

// Below
// ****************************************************************************************************************************************
func Below() CellOption {
	return position(2)
}

// TextColor
// ****************************************************************************************************************************************
func TextColor(rgb string) CellOption {
	return func(ctx *pdfContext) {
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
type option func(*pdfContext)

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// pageSize *************************************************************************************************************************
func pageSize(size fpdf.SizeType, landscape bool) option {
	return func(ctx *pdfContext) {
		if landscape {
			size.Wd, size.Ht = size.Ht, size.Wd
		}
		ctx.fpdf.AddPageFormat("P", size)
	}
}

// position *************************************************************************************************************************
func position(pos int) CellOption {
	return func(ctx *pdfContext) {
		ctx.position = pos
	}
}

// rgbColor *******************************************************************************************************************************
func rgbColor(rgb string) (r, g, b int, err error) {
	err = strings.Parse(strings.ToUpper(rgb), "#%2X%2X%2X", &r, &g, &b)

	return
}
