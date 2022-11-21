package redis

import (
	"context"
	"encoding/json"
	"time"

	base "github.com/go-redis/redis/v9"
	"github.com/huyungtang/go-lib/cache"
	"github.com/huyungtang/go-lib/times"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Init
// ****************************************************************************************************************************************
func Init(dsn string, opts ...cache.Option) (db cache.Database, err error) {
	var opt *base.Options
	if opt, err = base.ParseURL(dsn); err != nil {
		return
	}

	client := base.NewClient(opt)
	if err = client.Ping(context.Background()).Err(); err != nil {
		return
	}

	return &database{client}, nil
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// database *******************************************************************************************************************************
type database struct {
	*base.Client
}

// Exists
// ****************************************************************************************************************************************
func (o *database) Exists(key string) bool {
	cmd := o.Client.Exists(context.Background(), key)

	return cmd.Err() == nil && cmd.Val() == 1
}

// Get
// ****************************************************************************************************************************************
func (o *database) Get(key string, val interface{}, opts ...cache.Options) (err error) {
	cfg := new(cache.Option).
		ApplyOptions(opts,
			cache.KeepTTLOption,
		)

	ctx := context.Background()
	if o.Exists(key) {
		var res *base.StringCmd
		if _, err = o.Pipelined(ctx, func(p base.Pipeliner) (er error) {
			res = p.Get(ctx, key)
			if er = res.Err(); er != nil {
				return
			}

			if ex := getExpireCmder(cfg, ctx, key); ex != nil {
				er = p.Process(ctx, ex)
			}

			return
		}); err == nil {
			return json.Unmarshal([]byte(res.Val()), val)
		}
	} else if cfg.DefaFn != nil {
		var exp cache.Options
		if exp, err = cfg.DefaFn(val); err != nil {
			return
		}

		err = o.Set(key, val, exp)
	} else {
		err = base.Nil
	}

	return
}

// Set
// ****************************************************************************************************************************************
func (o *database) Set(key string, val interface{}, opts ...cache.Options) (err error) {
	cfg := new(cache.Option).ApplyOptions(opts)
	args := base.SetArgs{
		Mode: cfg.Override,
	}
	args.ExpireAt, args.TTL, args.KeepTTL = getExpire(cfg)

	val, _ = json.Marshal(val)
	if err = o.SetArgs(context.Background(), key, val, args).Err(); err == base.Nil {
		return nil
	}

	return
}

// Close
// ****************************************************************************************************************************************
func (o *database) Close() (err error) {
	if o.Client == nil {
		return
	}

	return o.Client.Close()
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// getExpire ******************************************************************************************************************************
func getExpire(opt *cache.Option) (exp time.Time, ttl time.Duration, keep bool) {
	now := times.Now()
	if opt.Expire > now.UnixSecond() {
		exp = times.Unix(opt.Expire, 0).Time()
	} else if opt.Expire > cache.KeepTTL {
		ttl = times.Duration(0, 0, opt.Expire)
	} else if opt.Expire == cache.KeepTTL {
		keep = true
	}

	return
}

// getExpireCmder *************************************************************************************************************************
func getExpireCmder(opt *cache.Option, ctx context.Context, key string) (cmd base.Cmder) {
	now := times.Now()
	if opt.Expire > now.UnixSecond() {
		cmd = base.NewCmd(ctx, "EXPIREAT", key, opt.Expire)
	} else if opt.Expire > 0 {
		cmd = base.NewCmd(ctx, "EXPIRE", key, opt.Expire)
	} else if opt.Expire == cache.Static {
		cmd = base.NewCmd(ctx, "PERSIST", key)
	} else if opt.Expire != cache.KeepTTL {
		cmd = base.NewCmd(ctx, "EXPIRE", key, opt.Expire)
	}

	return
}
