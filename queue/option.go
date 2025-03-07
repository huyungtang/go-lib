package queue

import (
	"encoding/json"
	"sync"

	amqp "github.com/rabbitmq/amqp091-go"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

var (
	EmptyBody Options = messageOption("plain/text", []byte{})

	ExchangeDirect  Options = KindOption("direct")
	ExchangeTopic   Options = KindOption("topic")
	ExchangeHeaders Options = KindOption("headers")
	ExchangeFanout  Options = KindOption("fanout")
)

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// ArgumentsOption
// ****************************************************************************************************************************************
func ArgumentsOption(args map[string]any) Options {
	return func(qo *Option) {
		qo.Arguments = args
	}
}

// AutoAckOption
// ****************************************************************************************************************************************
func AutoAckOption(a bool) Options {
	return func(qo *Option) {
		qo.AutoAck = a
	}
}

// AutoDeleteOption
// ****************************************************************************************************************************************
func AutoDeleteOption(del bool) Options {
	return func(qo *Option) {
		qo.AutoDelete = del
	}
}

// DurableOption
// ****************************************************************************************************************************************
func DurableOption(dur bool) Options {
	return func(qo *Option) {
		qo.Durable = dur
	}
}

// ExclusiveOption
// ****************************************************************************************************************************************
func ExclusiveOption(e bool) Options {
	return func(qo *Option) {
		qo.Exclusive = e
	}
}

// HandlerOption
// ****************************************************************************************************************************************
func HandlerOption(handlers ...ConsumerHandler) Options {
	return func(qo *Option) {
		qo.Handlers = append(qo.Handlers, handlers...)
	}
}

// ImmediateOption
// ****************************************************************************************************************************************
func ImmediateOption(i bool) Options {
	return func(qo *Option) {
		qo.Immediate = i
	}
}

// InternalOption
// ****************************************************************************************************************************************
func InternalOption(i bool) Options {
	return func(qo *Option) {
		qo.Internal = i
	}
}

// JsonMessageOption
// ****************************************************************************************************************************************
func JsonMessageOption(dto any) Options {
	bs, _ := json.Marshal(dto)

	return messageOption("application/json", bs)
}

// KindOption
// ****************************************************************************************************************************************
func KindOption(k string) Options {
	return func(qo *Option) {
		qo.Kind = k
	}
}

// MandatoryOption
// ****************************************************************************************************************************************
func MandatoryOption(m bool) Options {
	return func(qo *Option) {
		qo.Mandatory = m
	}
}

// NameOption
// ****************************************************************************************************************************************
func NameOption(name string) Options {
	return func(qo *Option) {
		qo.Name = name
	}
}

// NoLocalOption
// ****************************************************************************************************************************************
func NoLocalOption(n bool) Options {
	return func(qo *Option) {
		qo.NoLocal = n
	}
}

// NoWait
// ****************************************************************************************************************************************
func NoWait(w bool) Options {
	return func(qo *Option) {
		qo.NoWait = w
	}
}

// RoutingOption
// ****************************************************************************************************************************************
func RoutingOption(rs ...string) Options {
	return func(qo *Option) {
		qo.Routing = append(qo.Routing, rs...)
	}

}

// StringMessageOption
// ****************************************************************************************************************************************
func StringMessageOption(str string) Options {
	return messageOption("plain/text", []byte(str))
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Option
// ****************************************************************************************************************************************
type Option struct {
	Arguments  amqp.Table
	AutoAck    bool
	AutoDelete bool
	Durable    bool
	Exclusive  bool
	Handlers   []ConsumerHandler
	Immediate  bool
	Internal   bool
	Kind       string
	Mandatory  bool
	Messages   []*Content
	Name       string
	NoLocal    bool
	NoWait     bool
	Routing    []string

	msgOnce sync.Once
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

// ConsumerHandler
// ****************************************************************************************************************************************
type ConsumerHandler func(Message)

// Content
// ****************************************************************************************************************************************
type Content struct {
	ContentType string
	Body        []byte
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// messageOption **************************************************************************************************************************
func messageOption(ct string, body []byte) Options {
	return func(qo *Option) {
		qo.msgOnce.Do(func() {
			qo.Messages = make([]*Content, 0, 1)
		})

		qo.Messages = append(qo.Messages, &Content{ct, body})
	}

}
