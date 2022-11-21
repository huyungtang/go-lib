package http

import (
	base "net/http"
	"net/url"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

var (
	urlencodedOption Options = HeaderOption("Content-Type", "application/x-www-form-urlencoded")
)

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// HandlerOption
// ****************************************************************************************************************************************
func HandlerOption(handler ContextHandler) Options {
	return func(o *Option) {
		o.Handlers = append(o.Handlers, handler)
	}
}

// HeaderOption
// ****************************************************************************************************************************************
func HeaderOption(key, value string) Options {
	return func(o *Option) {
		o.Headers = append(o.Headers, []string{key, value})
	}
}

// HostOption
// ****************************************************************************************************************************************
func HostOption(host *url.URL) Options {
	return func(o *Option) {
		o.Host = host
	}
}

// ParameterOption
// ****************************************************************************************************************************************
func ParameterOption(key, val string) Options {
	return func(o *Option) {
		o.Params = append(o.Params, []string{key, val})
	}
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Option
// ****************************************************************************************************************************************
type Option struct {
	Host    *url.URL
	Headers [][]string
	Params  [][]string
	Ckies   []*base.Cookie

	Handler  int
	Handlers []ContextHandler
}

// ApplyOptions
// ****************************************************************************************************************************************
func (o *Option) ApplyOptions(opts []Options, defa ...Options) (opt *Option) {
	opts = append(defa, opts...)
	for _, optFn := range opts {
		optFn(o)
	}

	return o
}

// Options
// ****************************************************************************************************************************************
type Options func(*Option)

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
