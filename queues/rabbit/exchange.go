package rabbit

import (
	"github.com/huyungtang/go-lib/queues"
	"github.com/huyungtang/go-lib/strings"
	"github.com/streadway/amqp"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// IExchange
// ****************************************************************************************************************************************
type IExchange interface {
	// Publish
	// 	Options default value
	// 		routing   = ""
	// 		mandatory = *false / true
	// 		immediate = *false / true
	// 		delivery  = *Persistent / Transient
	// 		body      = nil
	Publish(opts ...Option) error

	// Consume
	// 	Options default value
	// 		queue = nil
	Consume(opts ...Option) error

	Close() error
}

// exchange *******************************************************************************************************************************
type exchange struct {
	name    string
	channel *amqp.Channel
	handler []queues.ContextHandler
}

// Publish
// ****************************************************************************************************************************************
func (o *exchange) Publish(opts ...Option) (err error) {
	var body []byte
	routing, mandatory, immediate, delivery := "", false, false, uint8(2)
	for i := 0; i < len(opts); i++ {
		switch opt := opts[i].(type) {
		case *routingOption:
			routing = opt.key
		case *mandatoryOption:
			mandatory = opt.value
		case *immediateOption:
			immediate = opt.value
		case *deliveryOption:
			delivery = opt.mode
		case *bodyOption:
			body = opt.body
		}
	}

	return o.channel.Publish(o.name, routing, mandatory, immediate, amqp.Publishing{
		DeliveryMode: delivery,
		ContentType:  "application/json",
		Body:         body,
	})
}

// Consume
// ****************************************************************************************************************************************
func (o *exchange) Consume(opts ...Option) (err error) {
	errChan := make(chan error)

	for i := 0; i < len(opts); i++ {
		switch opt := opts[i].(type) {
		case *queueOption:
			go func(q *queueOption) {
				if _, err = o.channel.QueueDeclare(q.name, q.durable, q.autoDelete, q.exclusive, q.noWait, q.args); err != nil {
					errChan <- err
					return
				}

				handlers := make(map[string]*routingOption)
				for r := 0; r < len(q.routing); r++ {
					if err = o.channel.QueueBind(q.name, q.routing[r].key, o.name, q.routing[r].noWait, q.routing[r].args); err != nil {

						errChan <- err
						return
					}
					handlers[q.routing[r].key] = q.routing[r]
				}

				var msgs <-chan amqp.Delivery
				cn := strings.Format("%s consumer", q.name)
				if msgs, err = o.channel.Consume(q.name, cn, q.autoAck, q.exclusive, q.noLocal, q.noWait, q.args); err != nil {
					errChan <- err
					return
				}

				for msg := range msgs {
					if h, isOK := handlers[msg.RoutingKey]; isOK {
						ctx := &queues.Context{
							Exchange: msg.Exchange,
							Routing:  msg.RoutingKey,
							Body:     msg.Body,
						}
						ctx.SetHandler(append(o.handler, h.handler)).
							Next()
					}
				}
			}(opt)
		}
	}

	return <-errChan
}

// Close
// ****************************************************************************************************************************************
func (o *exchange) Close() (err error) {
	if o.channel == nil {
		return
	}

	return o.channel.Close()
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
