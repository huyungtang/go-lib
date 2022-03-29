package redis

import (
	"context"
	"encoding/json"

	base "github.com/go-redis/redis/v8"
	"github.com/huyungtang/go-lib/cache"
	"github.com/huyungtang/go-lib/times"
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
	db *base.Client
}

// Init
// ****************************************************************************************************************************************
func (o *Database) Init(dsn string) (err error) {
	var opt *base.Options
	if opt, err = base.ParseURL(dsn); err != nil {
		return
	}

	o.db = base.NewClient(opt)

	return o.db.Ping(context.TODO()).Err()
}

// Set
// ****************************************************************************************************************************************
func (o *Database) Set(key string, val interface{}, expire int64, isOverride bool) (err error) {

	return o.db.Process(context.TODO(), o.setterCmder(key, val, expire, isOverride))
}

// Exists
// ****************************************************************************************************************************************
func (o *Database) Exists(key string) bool {
	res := o.db.Exists(context.TODO(), key)

	return res.Err() == nil && res.Val() == 1
}

// Get
// ****************************************************************************************************************************************
func (o *Database) Get(key string, val interface{}, defa func(interface{}) (int64, error), expire int64) (err error) {
	if defa == nil {
		return o.getCore(key, val, expire)
	} else {
		defer once.Reset()

		once.Do(func() {
			if err = o.getCore(key, val, expire); err == base.Nil {
				var exp int64
				if exp, err = defa(val); err != nil {
					return
				}
				err = o.Set(key, val, exp, true)
			}
		})
	}

	return
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
func (o *Database) setterCmder(key string, val interface{}, expire int64, isOverride bool) (cmd *base.Cmd) {
	args := make([]interface{}, 3, 6)
	args[0] = "SET"
	args[1] = key
	args[2], _ = json.Marshal(val)

	if expire > times.UnixSecond() {
		args = append(args, "EXAT", expire)
	} else if expire > cache.KeepTTL {
		args = append(args, "EX", expire)
	} else if expire == cache.KeepTTL {
		args = append(args, "KEEPTTL")
	}

	if !isOverride {
		args = append(args, "NX")
	}

	return base.NewCmd(context.TODO(), args...)
}

// expireCmder ****************************************************************************************************************************
func (o *Database) expireCmder(key string, secs int64) (cmd *base.Cmd) {
	ctx := context.TODO()
	if secs > times.UnixSecond() {
		cmd = base.NewCmd(ctx, "EXPIREAT", key, secs)
	} else if secs > 0 {
		cmd = base.NewCmd(ctx, "EXPIRE", key, secs)
	} else if secs == cache.Static {
		cmd = base.NewCmd(ctx, "PERSIST", key)
	} else if secs != cache.KeepTTL {
		cmd = base.NewCmd(ctx, "EXPIRE", key, secs)
	}

	return
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
