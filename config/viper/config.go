package viper

import (
	"encoding/json"
	"errors"
	base "reflect"
	"strconv"
	"time"

	"github.com/huyungtang/go-lib/config"
	"github.com/huyungtang/go-lib/encrypt"
	"github.com/huyungtang/go-lib/reflect"
	"github.com/huyungtang/go-lib/strings"
	"github.com/spf13/viper"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Init
// ****************************************************************************************************************************************
func Init(opts ...config.Options) (cfgs config.Config, err error) {
	v := viper.New()
	v.AutomaticEnv()
	v.AddConfigPath(".")

	cfg := new(config.Option).
		ApplyOptions(opts,
			config.NameOption("config"),
		)

	v.SetConfigName(cfg.Name)
	v.SetConfigType(cfg.FileType)
	for _, path := range cfg.Pathes {
		v.AddConfigPath(path)
	}

	if err = v.ReadInConfig(); err != nil {
		return
	}

	for _, env := range cfg.Envs {
		v.SetConfigName(strings.Format("%s-%s", cfg.Name, env))
		if err = v.MergeInConfig(); err != nil {
			return
		}
	}

	return &database{
		Viper:  v,
		Option: cfg,
	}, nil
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// database *******************************************************************************************************************************
type database struct {
	*viper.Viper
	*config.Option
}

// GetBool
// ****************************************************************************************************************************************
func (o *database) GetBool(key string, defa bool) bool {
	return o.getCore(key, defa).(bool)
}

// GetDuration
// ****************************************************************************************************************************************
func (o *database) GetDuration(key string, defa string) int64 {
	if o.IsSet(key) {
		return o.Viper.GetDuration(key).Nanoseconds()
	}

	if dur, err := time.ParseDuration(defa); err == nil {
		return dur.Nanoseconds()
	}

	return 0
}

// GetInt
// ****************************************************************************************************************************************
func (o *database) GetInt(key string, defa int) int {
	return o.getCore(key, defa).(int)
}

// GetIntSlice
// ****************************************************************************************************************************************
func (o *database) GetIntSlice(key string, defa []int) []int {
	return o.getCore(key, defa).([]int)
}

// GetString
// ****************************************************************************************************************************************
func (o *database) GetString(key, defa string) string {
	return o.getCore(key, defa).(string)
}

// GetStringSlice
// ****************************************************************************************************************************************
func (o *database) GetStringSlice(key string, defa []string) []string {
	return o.getCore(key, defa).([]string)
}

// GetStruct
// ****************************************************************************************************************************************
func (o *database) GetStruct(dto interface{}, opts ...config.Options) (err error) {
	tp := reflect.TypeOf(dto)
	if tp.Kind() != base.Struct {
		return errors.New("target is not a struct")
	}

	cfg := new(config.Option).
		ApplyOptions(opts,
			config.PathOption(""),
		)

	val := reflect.ValueOf(dto)
	for i := 0; i < tp.NumField(); i++ {
		if !val.Field(i).CanSet() {
			continue
		}

		tags := reflect.GetTags(tp.Field(i), "config")
		if _, isIgnore := tags["ignore"]; isIgnore {
			continue
		}

		if _, isOk := tags["key"]; !isOk {
			tags["key"] = tp.Field(i).Name
		}
		keys := append(cfg.Pathes, tags["key"])
		tags["key"] = strings.Join(strings.OmitEmpty(keys), ".")

		switch val.Field(i).Type().Kind() {
		case base.Bool:
			val.Field(i).SetBool(o.GetBool(tags["key"], tags["defa"] == "true"))
		case base.Int, base.Int8, base.Int16, base.Int32, base.Int64:
			defa, _ := strconv.ParseInt(tags["defa"], 10, 64)
			val.Field(i).SetInt(int64(o.GetInt(tags["key"], int(defa))))
		case base.String:
			val.Field(i).SetString(o.GetString(tags["key"], tags["defa"]))
		case base.Struct:
			n := base.New(val.Field(i).Type()).Interface()
			if e := o.GetStruct(n, config.PathOption(tags["key"])); e == nil {
				val.Field(i).Set(reflect.ValueOf(n))
			}
		case base.Slice:
			switch val.Field(i).Type().Elem().Kind() {
			case base.String:
				var defa []string
				if e := o.getTagDefa(&defa, tags); e == nil {
					val.Field(i).Set(base.ValueOf(o.GetStringSlice(tags["key"], defa)))
				}
			case base.Int:
				var defa []int
				if e := o.getTagDefa(&defa, tags); e == nil {
					val.Field(i).Set(base.ValueOf(o.GetIntSlice(tags["key"], defa)))
				}
			}
		}
	}

	return
}

// MergeInConfig
// ****************************************************************************************************************************************
func (o *database) MergeInConfig(opts ...config.Options) (err error) {
	cfg := new(config.Option).
		ApplyOptions(
			opts,
			config.NameOption(o.Name),
		)
	for _, env := range cfg.Envs {
		o.Viper.SetConfigName(strings.Format("%s-%s", cfg.Name, env))
	}

	return o.Viper.MergeInConfig()
}

// getCore ********************************************************************************************************************************
func (o *database) getCore(key string, defa interface{}) interface{} {
	if o.IsSet(key) {
		switch reflect.KindOf(defa) {
		case base.Bool:
			return o.Viper.GetBool(key)
		case base.Int:
			return o.Viper.GetInt(key)
		case base.String:
			s, _ := encrypt.Decrypt(o.Viper.GetString(key))

			return s
		case base.Slice:
			switch reflect.TypeOf(defa).Elem().Kind() {
			case base.Int:
				return o.Viper.GetIntSlice(key)
			case base.String:
				return o.Viper.GetStringSlice(key)
			}
		}
	}

	return defa
}

// getTagDefa *****************************************************************************************************************************
func (o *database) getTagDefa(dto interface{}, tags map[string]string) (err error) {
	if d := tags["defa"]; d != "" {
		err = json.Unmarshal([]byte(d), dto)
	}

	return
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
