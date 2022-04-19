package redis

import (
	"context"
	"encoding/json"

	redis_ "github.com/go-redis/redis/v8"
	"github.com/huyungtang/go-lib/cache"
	"github.com/huyungtang/go-lib/time"
	"github.com/matryer/resync"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

var once resync.Once

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Database
// ****************************************************************************************************************************************
type Database struct {
	db *redis_.Client
}

// Init
// ****************************************************************************************************************************************
func (o *Database) Init(dsn string) (err error) {
	var opt *redis_.Options
	if opt, err = redis_.ParseURL(dsn); err != nil {
		return
	}

	o.db = redis_.NewClient(opt)

	return o.db.Ping(context.TODO()).Err()
}

// Set
// 	Options default value
// 		expire   = cache.Static
// 		override = true
// ****************************************************************************************************************************************
func (o *Database) Set(key string, val interface{}, opts ...Option) (err error) {
	exp := cache.Static
	ove := true
	for i := 0; i < len(opts); i++ {
		switch opt := opts[i].(type) {
		case *expireOption:
			exp = opt.exp
		case *overrideOption:
			ove = opt.isOverride
		}
	}

	return o.db.Process(context.TODO(), o.setterCmder(key, val, exp, ove))
}

// Exists
// ****************************************************************************************************************************************
func (o *Database) Exists(key string) bool {
	res := o.db.Exists(context.TODO(), key)

	return res.Err() == nil && res.Val() == 1
}

// Get
// 	Options default value
// 		defaFunc    = nil
// 		expire      = cache.Static
// 		renewExpire = cache.KeepTTL
// ****************************************************************************************************************************************
func (o *Database) Get(key string, val interface{}, opts ...Option) (err error) {
	var defa defaultFunc
	sexp := cache.Static
	rexp := cache.KeepTTL
	for i := 0; i < len(opts); i++ {
		switch opt := opts[i].(type) {
		case *defaultOption:
			defa = opt.fn
		case *expireOption:
			if opt.isRenew {
				rexp = opt.exp
			} else {
				sexp = opt.exp
			}
		}
	}

	if defa == nil {
		return o.getCore(key, val, rexp)
	} else {
		defer once.Reset()

		once.Do(func() {
			if err = o.getCore(key, val, rexp); err == redis_.Nil {
				if err = defa(val); err != nil {
					return
				}
				err = o.Set(key, val, ExpireOption(sexp), OverrideOption(true))
			}
		})
	}

	return
}

// Pub
// ****************************************************************************************************************************************
func (o *Database) Pub(ch string, val interface{}) (err error) {
	var bs []byte
	if bs, err = json.Marshal(val); err != nil {
		return
	}

	cmd := o.db.Publish(context.TODO(), ch, bs)

	return cmd.Err()
}

// Sub
// ****************************************************************************************************************************************
func (o *Database) Sub(ch ...string) (ctx <-chan Context, err error) {
	mctx := make(chan Context)

	go func(mc chan<- Context) {
		var msg *redis_.Message
		sub := o.db.Subscribe(context.TODO(), ch...)

		for {
			if msg, err = sub.ReceiveMessage(context.TODO()); err != nil {
				break
			}

			mc <- &messageContext{m: msg}
		}
	}(mctx)

	return mctx, err
}

// Close
// ****************************************************************************************************************************************
func (o *Database) Close() (err error) {
	if o.db != nil {
		err = o.db.Close()
	}

	return
}

// getCore ********************************************************************************************************************************
func (o *Database) getCore(key string, val interface{}, expire int64) (err error) {
	var res string
	ctx := context.TODO()
	if res, err = o.db.Get(ctx, key).Result(); err != nil {
		return
	}

	if err = json.Unmarshal([]byte(res), val); err != nil {
		return
	}

	if c := o.expireCmder(key, expire); c != nil {
		err = o.db.Process(ctx, c)
	}

	return
}

// setterCmder ****************************************************************************************************************************
func (o *Database) setterCmder(key string, val interface{}, expire int64, isOverride bool) (cmd *redis_.Cmd) {
	args := make([]interface{}, 3, 6)
	args[0] = "SET"
	args[1] = key
	args[2], _ = json.Marshal(val)

	if expire > time.UnixSecond() {
		args = append(args, "EXAT", expire)
	} else if expire > cache.KeepTTL {
		args = append(args, "EX", expire)
	} else if expire == cache.KeepTTL {
		args = append(args, "KEEPTTL")
	}

	if !isOverride {
		args = append(args, "NX")
	}

	return redis_.NewCmd(context.TODO(), args...)
}

// expireCmder ****************************************************************************************************************************
func (o *Database) expireCmder(key string, secs int64) (cmd *redis_.Cmd) {
	ctx := context.TODO()
	if secs > time.UnixSecond() {
		cmd = redis_.NewCmd(ctx, "EXPIREAT", key, secs)
	} else if secs > 0 {
		cmd = redis_.NewCmd(ctx, "EXPIRE", key, secs)
	} else if secs == cache.Static {
		cmd = redis_.NewCmd(ctx, "PERSIST", key)
	} else if secs != cache.KeepTTL {
		cmd = redis_.NewCmd(ctx, "EXPIRE", key, secs)
	}

	return
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
