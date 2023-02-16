package template

import (
	"html/template"
	"io"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// FuncsOption
// ****************************************************************************************************************************************
func FuncsOption(funcs template.FuncMap) Options {
	return func(p *tmpl) (err error) {
		p.Template = p.init().Template.Funcs(funcs)

		return
	}
}

// ParseGlobOption
// ****************************************************************************************************************************************
func ParseGlobOption(pattern string) Options {
	return func(p *tmpl) (err error) {
		p.Template, err = p.init().Template.ParseGlob(pattern)

		return
	}
}

// WriterOption
// ****************************************************************************************************************************************
func WriterOption(w io.Writer) Options {
	return func(p *tmpl) (err error) {
		p.buf = w

		return
	}
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Options
// ****************************************************************************************************************************************
type Options func(*tmpl) error

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
