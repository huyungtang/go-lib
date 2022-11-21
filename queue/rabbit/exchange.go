package rabbit

import (
	"context"

	"github.com/huyungtang/go-lib/queue"
	"github.com/huyungtang/go-lib/strings"
	amqp "github.com/rabbitmq/amqp091-go"
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

// exchange *******************************************************************************************************************************
type exchange struct {
	*amqp.Channel
	name     string
	handlers []queue.ConsumerHandler
}

// Publish
// ****************************************************************************************************************************************
func (o *exchange) Publish(routing string, opts ...queue.Options) (err error) {
	cfg := new(queue.Option).ApplyOptions(opts)
	for _, msg := range cfg.Messages {
		if err = o.Channel.
			PublishWithContext(
				context.Background(),
				o.name,
				routing,
				cfg.Mandatory,
				cfg.Immediate,
				amqp.Publishing{ContentType: msg.ContentType, Body: msg.Body},
			); err != nil {
			break
		}
	}

	return
}

// QueueBind
// ****************************************************************************************************************************************
func (o *exchange) QueueBind(name string, opts ...queue.Options) (err error) {
	cfg := new(queue.Option).
		ApplyOptions(opts,
			queue.DurableOption(true),
		)

	if _, err = o.Channel.QueueDeclare(name, cfg.Durable, cfg.AutoDelete, cfg.Exclusive, cfg.NoWait, cfg.Arguments); err != nil {
		return
	}

	for _, routing := range cfg.Routing {
		if err = o.Channel.QueueBind(name, routing, o.name, cfg.NoWait, cfg.Arguments); err != nil {
			return
		}
	}

	return
}

// Consume
// ****************************************************************************************************************************************
func (o *exchange) Consume(queueName string, errChan chan<- error, opts ...queue.Options) {
	cfg := new(queue.Option).
		ApplyOptions(opts,
			queue.AutoAckOption(true),
			queue.HandlerOption(o.handlers...),
		)

	msgs, err := o.Channel.
		Consume(
			queueName,
			strings.Format("%s-consumer", queueName),
			cfg.AutoAck,
			cfg.Exclusive,
			cfg.NoLocal,
			cfg.NoWait,
			cfg.Arguments)
	if err != nil {
		errChan <- err

		return
	}

	for msg := range msgs {
		ctx := &message{Delivery: &msg, queue: queueName, handlers: cfg.Handlers, handler: -1}
		ctx.Next()
	}
}

// Close
// ****************************************************************************************************************************************
func (o *exchange) Close() (err error) {
	if o.Channel != nil {
		err = o.Channel.Close()
	}

	return
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
