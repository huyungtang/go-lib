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

// ParseGlobOption
// ****************************************************************************************************************************************
func ParseGlobOption(pattern string) Options {
	return func(p *tmpl) (err error) {
		p.Template, err = template.ParseGlob(pattern)

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
