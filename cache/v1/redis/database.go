package redis

import (
	"context"
	"encoding/json"
	"net"

	base "github.com/go-redis/redis/v9"
	"github.com/huyungtang/go-lib/cache/v1"
	"github.com/huyungtang/go-lib/logger"
	"github.com/huyungtang/go-lib/reflect"
	"github.com/huyungtang/go-lib/times"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

const (
	cmderDelete  string = "DEL"
	cmderExists  string = "EXISTS"
	cmderExpire  string = "EXPIRE"
	cmderGet     string = "GET"
	cmderPersist string = "PERSIST"
	cmderSet     string = "SET"
)

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Init
// ****************************************************************************************************************************************
func Init(dsn string, opts ...cache.Options) (db cache.Database, err error) {
	var opt *base.Options
	if opt, err = base.ParseURL(dsn); err != nil {
		return
	}

	client := base.NewClient(opt)
	if err = client.Ping(context.Background()).Err(); err != nil {
		return
	}

	instance := &database{client}
	cfg := cache.ApplyOptions([]cache.Options{}, opts...)
	if cfg.IsDebug {
		client.AddHook(instance)
	}

	return instance, nil
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// database *******************************************************************************************************************************
type database struct {
	*base.Client
}

// Delete
// ****************************************************************************************************************************************
func (o *database) Delete(key string) (err error) {
	cmd := o.Client.Del(context.Background(), key)

	return cmd.Err()
}

// Exists
// ****************************************************************************************************************************************
func (o *database) Exists(key string) (exs bool, err error) {
	cmd := o.Client.Exists(context.Background(), key)

	return (cmd.Err() == nil && cmd.Val() == 1), cmd.Err()
}

// Get
// ****************************************************************************************************************************************
func (o *database) Get(key string, val interface{}, opts ...cache.Options) (err error) {
	cfg := cache.ApplyOptions([]cache.Options{cache.KeepTTLOption()}, opts...)
	_, err = o.Client.Pipelined(context.Background(), func(p base.Pipeliner) (e error) {
		cmd := base.NewStringCmd(context.Background(), cmderGet, key)
		p.Process(context.Background(), cmd)
		if _, e = p.Exec(context.Background()); e == nil {
			if e = json.Unmarshal([]byte(cmd.Val()), val); e != nil {
				return
			}

			if c := expireCore(key, cfg); c != nil {
				p.Process(context.Background(), c)
			}

			return
		} else if e == base.Nil && cfg.DefaFn != nil {
			var os []cache.Options
			if os, e = cfg.DefaFn(val); e != nil {
				return
			}

			cfg = cache.ApplyOptions(os)
			p.Process(context.Background(), setCore(key, val, cfg))

			if cmd := expireCore(key, cfg); cmd != nil {
				p.Process(context.Background(), cmd)
			}
		}

		return
	})

	if err == base.Nil {
		return nil
	}

	return
}

// Set
// ****************************************************************************************************************************************
func (o *database) Set(key string, val interface{}, opts ...cache.Options) (err error) {
	cfg := cache.ApplyOptions([]cache.Options{cache.StaticOption()}, opts...)
	_, err = o.Client.Pipelined(context.Background(), func(p base.Pipeliner) (e error) {
		p.Process(context.Background(), setCore(key, val, cfg))

		if cmd := expireCore(key, cfg); cmd != nil {
			p.Process(context.Background(), cmd)
		}

		return
	})

	if err == base.Nil {
		return nil
	}

	return
}

// DialHook implements redis.Hook.
func (*database) DialHook(next base.DialHook) base.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

// ProcessHook implements redis.Hook.
func (*database) ProcessHook(next base.ProcessHook) base.ProcessHook {
	return func(ctx context.Context, cmd base.Cmder) error {
		logger.Printf("process: %s\n", cmd.String())
		return next(ctx, cmd)
	}
}

// ProcessPipelineHook implements redis.Hook.
func (*database) ProcessPipelineHook(next base.ProcessPipelineHook) base.ProcessPipelineHook {
	return func(ctx context.Context, cmds []base.Cmder) error {
		logger.Print("process: ")
		for i, cmd := range cmds {
			if i == 0 {
				logger.Println(cmd.String())
			} else {
				logger.Printf("         %s\n", cmd.String())
			}
		}
		return next(ctx, cmds)
	}
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// expireCore *****************************************************************************************************************************
func expireCore(key string, cfg *cache.Option) (cmd *base.StatusCmd) {
	args := make([]interface{}, 0, 3)
	if cfg.Expire == cache.ExpirationStatic {
		args = append(args, cmderPersist, key)
	} else if cfg.Expire == cache.ExpirationExpired {
		args = append(args, cmderDelete, key)
	} else if now := times.Now().Unix(); cfg.Expire > now {
		args = append(args, cmderExpire, key, cfg.Expire-now)
	} else if cfg.Expire > 0 {
		args = append(args, cmderExpire, key, cfg.Expire)
	}

	if len(args) > 0 {
		return base.NewStatusCmd(context.Background(), args...)
	}

	return
}

// setCore ********************************************************************************************************************************
func setCore(key string, val interface{}, cfg *cache.Option) (cmd *base.StatusCmd) {
	args := make([]interface{}, 3, 4)
	copy(args[0:3], []interface{}{cmderSet, key, valueOf(val)})

	switch cfg.Update {
	case cache.SetSkipOverride:
		args = append(args, "NX")
	case cache.SetExistOnly:
		args = append(args, "XX")
	}

	return base.NewStatusCmd(context.Background(), args...)
}

// valueOf ********************************************************************************************************************************
func valueOf(v interface{}) interface{} {
	if reflect.IsObject(v) {
		b, _ := json.Marshal(v)

		return string(b)
	}

	return v
}
