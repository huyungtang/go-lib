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

// Context
// ****************************************************************************************************************************************
func Init(opts ...Option) Context {
	ctx := &context{
		context: base.Background(),
		handler: -1,
	}
	for _, optFn := range opts {
		optFn(ctx)
	}

	return ctx
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// context ********************************************************************************************************************************
type context struct {
	context base.Context

	handler  int
	handlers []handler

	messages []Message
}

// Context
// ****************************************************************************************************************************************
type Context interface {
	Value(interface{}) interface{}
	WithValue(interface{}, interface{})
	Next() error
	SetInfo(string)
	SetWarn(string)
	SetError(string)
	Messages() []Message
}

// Value
// ****************************************************************************************************************************************
func (o *context) Value(key interface{}) interface{} {
	return o.context.Value(key)
}

// WithValue
// ****************************************************************************************************************************************
func (o *context) WithValue(key, val interface{}) {
	o.context = base.WithValue(o.context, key, val)
}

// Next
// ****************************************************************************************************************************************
func (o *context) Next() (err error) {
	o.handler++
	if o.handler >= len(o.handlers) || o.handlers[o.handler] == nil {
		return
	}

	return o.handlers[o.handler](o)
}

// Info
// ****************************************************************************************************************************************
func (o *context) SetInfo(message string) {
	o.message(MessageKindInfo, message)
}

// Warn
// ****************************************************************************************************************************************
func (o *context) SetWarn(message string) {
	o.message(MessageKindWarn, message)
}

// Error
// ****************************************************************************************************************************************
func (o *context) SetError(message string) {
	o.message(MessageKindError, message)
}

// Message
// ****************************************************************************************************************************************
func (o *context) Messages() []Message {
	return o.messages
}

// message ********************************************************************************************************************************
func (o *context) message(kind messageKind, content string) {
	o.messages = append(o.messages, &message{
		kind:    kind,
		content: content,
	})
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
