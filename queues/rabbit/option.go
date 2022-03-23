package rabbit

import (
	"encoding/json"

	"github.com/huyungtang/go-lib/queues"
	"github.com/streadway/amqp"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

var (
	ExchangeDirect  Option = &modeOption{stringOption{value: "direct"}}
	ExchangeTopic   Option = &modeOption{stringOption{value: "topic"}}
	ExchangeHeaders Option = &modeOption{stringOption{value: "headers"}}
	ExchangeFanout  Option = &modeOption{stringOption{value: "fanout"}}

	DeliveryTransient  Option = &deliveryOption{mode: 1}
	DeliveryPersistent Option = &deliveryOption{mode: 2}
)

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// AutoAckOption
// ****************************************************************************************************************************************
func AutoAckOption(v bool) Option {
	return &autoAckOption{boolOption{value: v}}
}

// AutoDeleteOption
// ****************************************************************************************************************************************
func AutoDeleteOption(v bool) Option {
	return &autoDeleteOption{boolOption{value: v}}
}

// BodyOption
// ****************************************************************************************************************************************
func BodyOption(body interface{}) Option {
	dto := new(bodyOption)
	dto.body, _ = json.Marshal(body)

	return dto
}

// DurableOption
// ****************************************************************************************************************************************
func DurableOption(v bool) Option {
	return &durableOption{boolOption{value: v}}
}

// ExclusiveOption ************************************************************************************************************************
func ExclusiveOption(v bool) Option {
	return &exclusiveOption{boolOption{value: v}}
}

// HandlerOption
// ****************************************************************************************************************************************
func HandlerOption(handler queues.ContextHandler) Option {
	return &handlerOption{handler: handler}
}

// ImmediateOption
// ****************************************************************************************************************************************
func ImmediateOption(v bool) Option {
	return &immediateOption{boolOption{value: v}}
}

// InternalOption
// ****************************************************************************************************************************************
func InternalOption(v bool) Option {
	return &internalOption{boolOption{value: v}}
}

// MandatoryOption
// ****************************************************************************************************************************************
func MandatoryOption(v bool) Option {
	return &mandatoryOption{boolOption{value: v}}
}

// NoLocalOption
// ****************************************************************************************************************************************
func NoLocalOption(v bool) Option {
	return &noLocalOption{boolOption{value: v}}
}

// NoWaitOption
// ****************************************************************************************************************************************
func NoWaitOption(v bool) Option {
	return &noWaitOption{boolOption{value: v}}
}

// QueueOption
// 	Options default value
// 		durable    = *false / true
// 		autoDelete = *false / true
// 		exclusive  = *false / true
// 		noWait     = *false / true
// 		autoAck    = false / *true
// 		noLocal    = *false / true
// 		table      = nil
// 		routing    = nil
// ****************************************************************************************************************************************
func QueueOption(name string, opts ...Option) Option {
	o := &queueOption{
		name:       name,
		durable:    false,
		autoDelete: false,
		exclusive:  false,
		noWait:     false,
		autoAck:    true,
		routing:    make([]*routingOption, 0),
	}
	for i := 0; i < len(opts); i++ {
		switch opt := opts[i].(type) {
		case *durableOption:
			o.durable = opt.value
		case *autoDeleteOption:
			o.autoDelete = opt.value
		case *exclusiveOption:
			o.exclusive = opt.value
		case *noWaitOption:
			o.noWait = opt.value
		case *autoAckOption:
			o.autoAck = opt.value
		case *noLocalOption:
			o.noLocal = opt.value
		case *tableOption:
			o.args = opt.value
		case *routingOption:
			o.routing = append(o.routing, opt)
		}
	}

	return o
}

// RoutingOption
// 	Options default value
// 		noWait     = *false / true
// 		table      = nil
// ****************************************************************************************************************************************
func RoutingOption(key string, handler queues.ContextHandler, opts ...Option) Option {
	o := &routingOption{
		key:     key,
		noWait:  false,
		handler: handler,
	}
	for i := 0; i < len(opts); i++ {
		switch opt := opts[i].(type) {
		case *noWaitOption:
			o.noWait = opt.value
		case *tableOption:
			o.args = opt.value
		}
	}

	return o
}

// TableOption
// ****************************************************************************************************************************************
func TableOption(args map[string]interface{}) Option {
	return &tableOption{value: args}
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Option
// ****************************************************************************************************************************************
type Option interface {
	Option()
}

// option *********************************************************************************************************************************
type option struct {
}

// option
// ****************************************************************************************************************************************
func (o *option) Option() {
}

// boolOption *****************************************************************************************************************************
type boolOption struct {
	option
	value bool
}

// autoAckOption **************************************************************************************************************************
type autoAckOption struct {
	boolOption
}

// autoDeleteOption ***********************************************************************************************************************
type autoDeleteOption struct {
	boolOption
}

// durableOption **************************************************************************************************************************
type durableOption struct {
	boolOption
}

// exclusiveOption ************************************************************************************************************************
type exclusiveOption struct {
	boolOption
}

// immediateOption ************************************************************************************************************************
type immediateOption struct {
	boolOption
}

// internalOption *************************************************************************************************************************
type internalOption struct {
	boolOption
}

// mandatoryOption ************************************************************************************************************************
type mandatoryOption struct {
	boolOption
}

// noLocalOption **************************************************************************************************************************
type noLocalOption struct {
	boolOption
}

// noWaitOption ***************************************************************************************************************************
type noWaitOption struct {
	boolOption
}

// stringOption ***************************************************************************************************************************
type stringOption struct {
	option
	value string
}

// modeOption *****************************************************************************************************************************
type modeOption struct {
	stringOption
}

// bodyOption *****************************************************************************************************************************
type bodyOption struct {
	option
	body []byte
}

// deliveryOption *************************************************************************************************************************
type deliveryOption struct {
	option
	mode uint8
}

// handlerOption **************************************************************************************************************************
type handlerOption struct {
	option
	handler queues.ContextHandler
}

// queueOption ****************************************************************************************************************************
type queueOption struct {
	option
	name       string
	durable    bool
	autoDelete bool
	exclusive  bool
	noWait     bool
	args       amqp.Table
	autoAck    bool
	noLocal    bool
	routing    []*routingOption
}

// routingOption **************************************************************************************************************************
type routingOption struct {
	option
	key     string
	noWait  bool
	args    amqp.Table
	handler queues.ContextHandler
}

// tableOption ****************************************************************************************************************************
type tableOption struct {
	option
	value amqp.Table
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
