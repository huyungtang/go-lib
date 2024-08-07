package template

import (
	base "html/template"
	"io"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// OutputOption
// ****************************************************************************************************************************************
func OutputOption(w io.Writer) Options {
	return func(t *template) (err error) {
		t.out = w

		return
	}
}

// FuncsOption
// ****************************************************************************************************************************************
func FuncsOption(funcs base.FuncMap) Options {
	return func(t *template) (err error) {
		t.tmplOnce.Do(func() {
			t.Template, err = base.New("empty").Parse("must use ExecuteTemplate")
		})

		if err != nil {
			return
		}

		t.Template = t.Template.Funcs(funcs)

		return
	}
}

// ParseFiles
// ****************************************************************************************************************************************
func ParseFiles(filename ...string) Options {
	return func(t *template) (err error) {
		done := false
		t.tmplOnce.Do(func() {
			defer func() {
				done = true
			}()

			t.Template, err = base.ParseFiles(filename...)
		})

		if done || err != nil {
			return
		}

		t.Template, err = t.Template.ParseFiles(filename...)

		return
	}
}

// ParseGlobOption
// ****************************************************************************************************************************************
func ParseGlobOption(pattern string) Options {
	return func(t *template) (err error) {
		done := false
		t.tmplOnce.Do(func() {
			defer func() {
				done = true
			}()

			t.Template, err = base.ParseGlob(pattern)
		})

		if done || err != nil {
			return
		}

		t.Template, err = t.Template.ParseGlob(pattern)

		return
	}
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Options
// ****************************************************************************************************************************************
type Options func(*template) error

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
