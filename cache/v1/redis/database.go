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
	cmderLPush   string = "LPUSH"
	cmderRPush   string = "RPUSH"
	cmderLLen    string = "LLEN"
	cmderLPop    string = "LPOP"
	cmderRPop    string = "RPOP"
	cmderIncr    string = "INCRBY"
	cmderIncrF   string = "INCRBYFLOAT"
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

	instance := &database{Client: client}
	instance.Option = cache.ApplyOptions([]cache.Options{}, opts...)
	if instance.IsDebug {
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
	*cache.Option
}

// Delete
// ****************************************************************************************************************************************
func (o *database) Delete(key ...string) (err error) {
	cmd := o.Client.Del(context.Background(), key...)

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
		var cmds []base.Cmder
		p.Process(context.Background(), base.NewStringCmd(context.Background(), cmderGet, key))
		if cmds, e = p.Exec(context.Background()); e == nil {
			if e = parseValue(val, cmds[0].(*base.StringCmd)); e != nil {
				return
			}
		} else if e == base.Nil && cfg.DefaFn != nil {
			var os []cache.Options
			if os, e = cfg.DefaFn(val); e != nil {
				return
			}

			p.Process(context.Background(), setCore(key, val, cache.ApplyOptions(os)))

			cfg = cache.ApplyOptions(os)
		}

		if expCmd := expireCore(key, cfg); expCmd != nil {
			p.Process(context.Background(), expCmd)
		}

		return
	})

	if err == base.Nil {
		return nil
	}

	return
}

// GetSlice
// ****************************************************************************************************************************************
func (o *database) GetSlice(key string, val interface{}, opts ...cache.Options) (err error) {
	cfg := cache.ApplyOptions([]cache.Options{cache.KeepTTLOption(), cache.DirectionLeftOption()}, opts...)
	var cmds []base.Cmder
	if cmds, err = o.Client.Pipelined(context.Background(), func(p base.Pipeliner) (e error) {
		p.Process(context.Background(), base.NewIntCmd(context.Background(), cmderLLen, key))
		if cmds, err = p.Exec(context.Background()); err == nil && cmds[0].(*base.IntCmd).Val() > 0 {
			args := make([]interface{}, 3)
			switch cfg.Direction {
			case cache.DirectionLeft:
				args[0] = cmderLPop
			case cache.DirectionRight:
				args[0] = cmderRPop
			}
			copy(args[1:3], []interface{}{key, cmds[0].(*base.IntCmd).Val()})

			p.Process(context.Background(), base.NewSliceCmd(context.Background(), args...))

			if cmdExp := expireCore(key, cfg); cmdExp != nil {
				p.Process(context.Background(), cmdExp)
			}
		}

		return
	}); err != nil {
		return
	}

	if len(cmds) > 0 {
		result := cmds[0].(*base.SliceCmd)
		if l := len(result.Val()); l > 0 {
			switch vs := val.(type) {
			case *[]string:
				*vs = make([]string, l)
				for i, v := range result.Val() {
					(*vs)[i] = v.(string)
				}
			case *[]int64:
				*vs = make([]int64, l)
				for i, v := range result.Val() {
					(*vs)[i] = v.(int64)
				}
			}
		}
	}

	return
}

// Increase
// ****************************************************************************************************************************************
func (o *database) Increase(key string, val interface{}, opts ...cache.Options) (err error) {
	cfg := cache.ApplyOptions([]cache.Options{
		cache.IncreaseByOption(1),
		cache.IncreaseByFloatOption(1),
		cache.StaticOption(),
	}, opts...)

	var isInt64 bool
	if _, isMatch := val.(*int64); isMatch {
		isInt64 = true
	} else if _, isMatch := val.(*float64); isMatch {
		isInt64 = false
	} else {
		return
	}

	var cmds []base.Cmder
	if cmds, err = o.Pipelined(context.Background(), func(p base.Pipeliner) (e error) {
		p.Process(context.Background(), base.NewIntCmd(context.Background(), cmderExists, key))
		if cmds, e = p.Exec(context.Background()); e != nil {
			return
		} else if cmds[0].(*base.IntCmd).Val() == 0 && cfg.DefaFn != nil {
			var os []cache.Options
			if os, e = cfg.DefaFn(val); e != nil {
				return
			}

			cf := cache.ApplyOptions(os)
			p.Process(context.Background(), setCore(key, val, cf))
			if cmd := expireCore(key, cf); cmd != nil {
				p.Process(context.Background(), cmd)
			}

			if _, e = p.Exec(context.Background()); e != nil {
				return
			}
		}

		if isInt64 {
			p.Process(context.Background(), base.NewCmd(context.Background(), cmderIncr, key, cfg.IncrInt))
		} else {
			p.Process(context.Background(), base.NewCmd(context.Background(), cmderIncrF, key, cfg.IncrFloat))
		}

		if cmd := expireCore(key, cfg); cmd != nil {
			e = p.Process(context.Background(), cmd)
		}

		return
	}); err != nil {
		return
	}

	if len(cmds) > 0 {
		if isInt64 {
			v, _ := val.(*int64)
			*v, _ = cmds[0].(*base.Cmd).Int64()
		} else {
			v, _ := val.(*float64)
			*v, _ = cmds[0].(*base.Cmd).Float64()
		}
	}

	return
}

// Push
// ****************************************************************************************************************************************
func (o *database) Push(key string, val interface{}, opts ...cache.Options) (err error) {
	cfg := cache.ApplyOptions([]cache.Options{cache.DirectionRightOption()}, opts...)
	_, err = o.Client.Pipelined(context.Background(), func(p base.Pipeliner) (e error) {
		args := make([]interface{}, 3)
		switch cfg.Direction {
		case cache.DirectionLeft:
			args[0] = cmderLPush
		case cache.DirectionRight:
			args[0] = cmderRPush
		}
		copy(args[1:3], []interface{}{key, valueOf(val)})

		p.Process(context.Background(), base.NewStatusCmd(context.Background(), args...))

		if cmd := expireCore(key, cfg); cmd != nil {
			p.Process(context.Background(), cmd)
		}

		return
	})

	return
}

// Set
// ****************************************************************************************************************************************
func (o *database) Set(key string, val interface{}, opts ...cache.Options) (err error) {
	cfg := cache.ApplyOptions([]cache.Options{cache.StaticOption()}, opts...)
	if _, err = o.Client.Pipelined(context.Background(), func(p base.Pipeliner) (e error) {
		p.Process(context.Background(), setCore(key, val, cfg))

		if cmd := expireCore(key, cfg); cmd != nil {
			p.Process(context.Background(), cmd)
		}

		return
	}); err == base.Nil {
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
	} else if reflect.IsPointer(v) {
		return reflect.ValueOf(v).Interface()
	}

	return v
}

// parseValue *****************************************************************************************************************************
func parseValue(val interface{}, cmd *base.StringCmd) (err error) {
	switch t := val.(type) {
	case *string:
		*t = cmd.Val()
	case *int:
		*t, err = cmd.Int()
	case *int64:
		*t, err = cmd.Int64()
	case *uint64:
		*t, err = cmd.Uint64()
	default:
		err = json.Unmarshal([]byte(cmd.Val()), val)
	}

	return
}
