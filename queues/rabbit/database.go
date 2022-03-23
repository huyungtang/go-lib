package rabbit

import (
	"github.com/huyungtang/go-lib/queues"
	"github.com/matryer/resync"
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

// Database
// ****************************************************************************************************************************************
type Database struct {
	once    resync.Once
	db      *amqp.Connection
	exchgs  map[string]IExchange
	handler []queues.ContextHandler
}

// Init
// 	Options default value
// 		handler = nil
// ****************************************************************************************************************************************
func (o *Database) Init(dsn string, opts ...Option) (err error) {
	if o.db, err = amqp.Dial(dsn); err != nil {
		return
	}

	o.exchgs = make(map[string]IExchange)
	o.handler = make([]queues.ContextHandler, 0)
	for i := 0; i < len(opts); i++ {
		switch opt := opts[i].(type) {
		case *handlerOption:
			o.handler = append(o.handler, opt.handler)
		}
	}

	return
}

// ExchangeDeclare
// ****************************************************************************************************************************************
func (o *Database) ExchangeDeclare(name string, opts ...Option) (e IExchange, err error) {
	defer o.once.Reset()

	o.once.Do(func() {
		if _, isOK := o.exchgs[name]; !isOK {
			exchg := &exchange{name: name, handler: o.handler}
			if exchg.channel, err = o.db.Channel(); err != nil {
				return
			}

			var args amqp.Table
			mode, durable, autoDelete, internal, noWait := "direct", true, false, false, false
			for i := 0; i < len(opts); i++ {
				switch opt := opts[i].(type) {
				case *modeOption:
					mode = opt.value
				case *durableOption:
					durable = opt.value
				case *autoDeleteOption:
					autoDelete = opt.value
				case *internalOption:
					internal = opt.value
				case *noWaitOption:
					noWait = opt.value
				case *tableOption:
					args = opt.value
				case *handlerOption:
					exchg.handler = append(o.handler, opt.handler)
				}
			}

			if err = exchg.channel.ExchangeDeclare(name, mode, durable, autoDelete, internal, noWait, args); err != nil {
				return
			}

			o.exchgs[name] = exchg
		}
	})

	return o.exchgs[name], err
}

// Close
// ****************************************************************************************************************************************
func (o *Database) Close() (err error) {
	if o.db == nil {
		return
	}

	return o.db.Close()
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
