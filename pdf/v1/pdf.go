package pdf

import (
	"io"

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

// Init
// ****************************************************************************************************************************************
func Init(opts ...Option) (p PDF, err error) {
	out := &context{
		opts: append([]Option{
			AlignOption(AlignML),
			FontFamilyOption("Arial"),
			FontSizeOption(11),
			CellMaringOption(2),
			CellWidthOption(1),
			PageMarginsOption(10, 10, 10),
			PortraitOption(),
			PageSizeA4Option(),
			PageBreakOption(true, 20),
			TemplateOption(-1),
			PositionOption(PositionNewLine),
			UnitMM(),
			BorderColor("#1f1f1f"),
			FontColorOption("#1f1f1f"),
		}, opts...),
	}

	cfg := applyOption(out.opts...)
	out.Fpdf = fpdf.New(cfg.orientation, cfg.unit, cfg.pageSize.name, "/")
	out.Fpdf.SetCellMargin(cfg.cellMargin)
	out.Fpdf.SetMargins(cfg.pageLeft, cfg.pageTop, cfg.pageRight)
	out.Fpdf.SetDrawColor(cfg.borderColor[0], cfg.borderColor[1], cfg.borderColor[2])
	out.Fpdf.SetTextColor(cfg.textColor[0], cfg.textColor[1], cfg.textColor[2])
	out.Fpdf.SetAutoPageBreak(cfg.pageBreak, cfg.pageBottom)

	if cfg.ttfPath != "" {
		var fs []string
		if fs, err = file.ListFiles(cfg.ttfPath, `\.ttf$`); err != nil {
			return
		}

		for _, f := range fs {
			fn := file.GetFilename(f)
			out.Fpdf.AddUTF8Font(fn[0:strings.LastIndex(fn, ".")], "", f)
		}
	}
	out.Fpdf.SetFont(cfg.fontFamily, "", cfg.fontSize)

	if cfg.templateFile != "" {
		out.Importer = gofpdi.NewImporter()
		out.Importer.SetSourceFile(cfg.templateFile)
		for i := 1; i <= out.Importer.GetNumPages(); i++ {
			out.Importer.ImportPage(i, "/MediaBox")
		}

		out.Fpdf.ImportTemplates(out.Importer.PutFormXobjectsUnordered())
		out.Fpdf.ImportObjects(out.Importer.GetImportedObjectsUnordered())
		out.Fpdf.ImportObjPos(out.Importer.GetImportedObjHashPos())
	}

	return out, err
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// PDF
// ****************************************************************************************************************************************
type PDF interface {
	AddPage(...Option) PDF
	AddBarcode128(string, ...Option) PDF
	AddCell(string, ...Option) PDF
	GetXY() (float64, float64)
	SetXY(float64, float64) PDF

	Output(io.Writer) error
	Close() error
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
