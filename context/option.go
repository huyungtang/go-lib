package context

import (
	base "context"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// WithValueOption
// ****************************************************************************************************************************************
func WithValueOption(key, val interface{}) Option {
	return func(ctx *context) {
		ctx.context = base.WithValue(ctx.context, key, val)
	}
}

// HandlerOption
// ****************************************************************************************************************************************
func HandlerOption(fn handler) Option {
	return func(ctx *context) {
		ctx.handlers = append(ctx.handlers, fn)
	}
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Option
// ****************************************************************************************************************************************
type Option func(*context)

// handler ********************************************************************************************************************************
type handler func(Context) error

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
