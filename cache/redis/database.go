package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/huyungtang/go-lib/cache"
	"github.com/huyungtang/go-lib/times"
	base "github.com/redis/go-redis/v9"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Init
// ****************************************************************************************************************************************
func Init(dsn string, opts ...cache.Context) (db cache.Database, err error) {
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
func (o *database) Get(key string, val any, opts ...cache.Option) (err error) {
	cfg := new(cache.Context).
		ApplyOptions(
			opts,
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
		var exp cache.Option
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
func (o *database) Set(key string, val any, opts ...cache.Option) (err error) {
	cfg := new(cache.Context).
		ApplyOptions(
			opts,
			cache.StaticOption,
		)

	val, _ = json.Marshal(val)
	if err = o.SetArgs(context.Background(), key, val, getArgs(cfg)).Err(); err == base.Nil {
		return nil
	}

	return
}

// Push
// ****************************************************************************************************************************************
func (o *database) Push(key string, val any, opts ...cache.Option) (err error) {
	cfg := new(cache.Context).
		ApplyOptions(
			opts,
			cache.RPushOption,
			cache.StaticOption,
		)

	ctx := context.Background()
	if _, err = o.Pipelined(ctx, func(p base.Pipeliner) (e error) {
		cmd := base.NewCmd(ctx, cfg.Cmder, key, val)
		if e = p.Process(ctx, cmd); e != nil {
			return
		}
		if e = p.Process(ctx, getExpireCmder(cfg, ctx, key)); e != nil {
			return
		}

		return
	}); err != nil {
		return
	}

	return
}

// Pop
// ****************************************************************************************************************************************
func (o *database) Pop(key string, val any, opts ...cache.Option) (err error) {
	cfg := new(cache.Context).
		ApplyOptions(
			opts,
			cache.LPopOption,
			cache.PopCountOption(1),
		)

	var cmd *base.StringSliceCmd
	if cfg.Cmder == "RPOP" {
		cmd = o.RPopCount(context.Background(), key, int(cfg.Count))
	} else {
		cmd = o.LPopCount(context.Background(), key, int(cfg.Count))
	}

	if err = cmd.Err(); err != nil {
		return
	}

	return cmd.ScanSlice(val)
}

// Del
// ****************************************************************************************************************************************
func (o *database) Del(keys ...string) (err error) {
	cmd := o.Client.Del(context.Background(), keys...)

	return cmd.Err()
}

// Ping
// ****************************************************************************************************************************************
func (o *database) Ping() (err error) {
	return o.Client.Ping(context.Background()).Err()
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

// getArgs ********************************************************************************************************************************
func getArgs(cfg *cache.Context) (args base.SetArgs) {
	args = base.SetArgs{
		Mode: cfg.Override,
	}
	args.ExpireAt, args.TTL, args.KeepTTL = getExpire(cfg)

	return
}

// getExpire ******************************************************************************************************************************
func getExpire(opt *cache.Context) (exp time.Time, ttl time.Duration, keep bool) {
	if opt.Expire > times.UnixSecond(times.Now()) {
		exp = times.FromUnix(opt.Expire, 0)
	} else if opt.Expire > cache.KeepTTL {
		ttl = times.Duration(0, 0, opt.Expire)
	} else if opt.Expire == cache.KeepTTL {
		keep = true
	}

	return
}

// getExpireCmder *************************************************************************************************************************
func getExpireCmder(opt *cache.Context, ctx context.Context, key string) (cmd base.Cmder) {
	if opt.Expire > times.UnixSecond(times.Now()) {
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
