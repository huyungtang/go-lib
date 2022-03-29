package config

import (
	"encoding/json"
	"errors"
	"reflect"
	"strconv"

	"github.com/huyungtang/go-lib/encrypt"
	"github.com/huyungtang/go-lib/reflects"
	"github.com/huyungtang/go-lib/strings"
	"github.com/spf13/viper"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

var errOptionNotMatched = errors.New("given option not matched")
var errNotStruct = errors.New("not a struct")

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Init
// 	Options default value
// 		name   = "config"
// 		type   = "yaml"
// 		path   = "" ex. "/etc/app/config"
// 		suffix = "" ex. "CFG_ENV"
// ****************************************************************************************************************************************
func Init(opts ...Option) (cfg Config, err error) {
	c := &config{
		db:           viper.New(),
		configName:   "config",
		configType:   "yaml",
		suffixOption: make([]*suffixOption, 0),
	}
	c.db.AutomaticEnv()
	c.db.AddConfigPath(".")

	for i := 0; i < len(opts); i++ {
		switch opt := opts[i].(type) {
		case *nameOption:
			c.configName = opt.name
			c.db.SetConfigName(opt.name)
		case *typeOption:
			c.configType = opt.tp
			c.db.SetConfigType(opt.tp)
		case *pathOption:
			c.db.AddConfigPath(opt.path)
		case *suffixOption:
			c.suffixOption = append(c.suffixOption, opt)
		}
	}

	if err = c.db.ReadInConfig(); err != nil {
		return
	}

	for i := 0; i < len(c.suffixOption); i++ {
		if err = c.MergeConfig(c.suffixOption[i]); err != nil {
			return
		}
	}

	return c, nil
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Config
// ****************************************************************************************************************************************
type Config interface {
	MergeConfig(Option) error
	GetBool(string, bool) bool
	GetInt(string, int) int
	GetIntSlice(string, []int) []int
	GetInt64(string, int64) int64
	GetInt64Slice(string, []int64) []int64
	GetUInt64(string, uint64) uint64
	GetString(string, string) string
	GetStringSlice(string, []string) []string
	GetStruct(interface{}, ...Option) error
}

// config *********************************************************************************************************************************
type config struct {
	db           *viper.Viper
	configName   string
	configType   string
	suffixOption []*suffixOption
}

// MergeConfig
// ****************************************************************************************************************************************
func (o *config) MergeConfig(opt Option) (err error) {
	suf, isOK := opt.(*suffixOption)
	if !isOK {
		return errOptionNotMatched
	}

	if s := o.db.GetString(suf.varName); s != "" {
		o.db.SetConfigName(strings.Format("%s-%s", o.configName, s))

		return o.db.MergeInConfig()
	}

	return
}

// GetBool
// ****************************************************************************************************************************************
func (o *config) GetBool(key string, defa bool) bool {
	return o.getCore(key, defa).(bool)
}

// GetInt
// ****************************************************************************************************************************************
func (o *config) GetInt(key string, defa int) int {
	return o.getCore(key, defa).(int)
}

// GetIntSlice
// ****************************************************************************************************************************************
func (o *config) GetIntSlice(key string, defa []int) []int {
	return o.getCore(key, defa).([]int)
}

// GetInt64
// ****************************************************************************************************************************************
func (o *config) GetInt64(key string, defa int64) int64 {
	return o.getCore(key, defa).(int64)
}

// GetInt64Slice
// ****************************************************************************************************************************************
func (o *config) GetInt64Slice(key string, defa []int64) []int64 {
	return o.getCore(key, defa).([]int64)
}

// GetUInt64
// ****************************************************************************************************************************************
func (o *config) GetUInt64(key string, defa uint64) uint64 {
	return o.getCore(key, defa).(uint64)
}

// GetString
// ****************************************************************************************************************************************
func (o *config) GetString(key, defa string) (s string) {
	s, _ = encrypt.Decrypt(o.getCore(key, defa).(string))

	return
}

// GetStringSlice
// ****************************************************************************************************************************************
func (o *config) GetStringSlice(key string, defa []string) []string {
	return o.getCore(key, defa).([]string)
}

// GetStruct
// ****************************************************************************************************************************************
func (o *config) GetStruct(dto interface{}, opts ...Option) (err error) {
	node := ""

	for i := 0; i < len(opts); i++ {
		switch opt := opts[i].(type) {
		case *pathOption:
			node = opt.path
		}
	}

	tp := reflects.TypeOf(dto)
	if tp.Kind() != reflect.Struct {
		return errNotStruct
	}

	val := reflects.ValueOf(dto)
	for i := 0; i < tp.NumField(); i++ {
		vf := val.Field(i)
		if !vf.CanSet() {
			continue
		}

		ff := tp.Field(i)
		ft := reflects.GetTags(ff, "config")
		if _, isIgnore := ft["ignore"]; isIgnore {
			continue
		}

		if _, isOk := ft["key"]; !isOk {
			ft["key"] = ff.Name
		}
		ft["key"] = strings.Join(".", true, node, ft["key"])

		switch vf.Type().Kind() {
		case reflect.Bool:
			vf.SetBool(o.GetBool(ft["key"], ft["defa"] == "true"))
		case reflect.String:
			vf.SetString(o.GetString(ft["key"], ft["defa"]))
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			defa, _ := strconv.ParseInt(ft["defa"], 10, 64)
			vf.SetInt(int64(o.GetInt(ft["key"], int(defa))))
		case reflect.Struct:
			sv := reflect.New(vf.Type()).Interface()
			o.GetStruct(sv, PathOption(ft["key"]))
			vf.Set(reflects.ValueOf(sv))
		case reflect.Slice:
			switch vf.Type().Elem().Kind() {
			case reflect.String:
				var defa []string
				if d := ft["defa"]; d != "" {
					json.Unmarshal([]byte(d), &defa)
				}
				vf.Set(reflect.ValueOf(o.GetStringSlice(ft["key"], defa)))
			case reflect.Int:
				var defa []int
				if d := ft["defa"]; d != "" {
					json.Unmarshal([]byte(d), &defa)
				}
				vf.Set(reflect.ValueOf(o.GetIntSlice(ft["key"], defa)))
			}
		}
	}

	return
}

// getCore *******************************************************************************************************************************
func (o *config) getCore(key string, defa interface{}) interface{} {
	if !o.db.IsSet(key) {
		return defa
	}

	switch reflects.KindOf(defa) {
	case reflect.Bool:
		return o.db.GetBool(key)
	case reflect.Int:
		return o.db.GetInt(key)
	case reflect.Int64:
		return o.db.GetInt64(key)
	case reflect.Uint64:
		return o.db.GetUint64(key)
	case reflect.String:
		return o.db.GetString(key)
	case reflect.Slice:
		switch reflects.TypeOf(defa).Elem().Kind() {
		case reflect.Int:
			val := o.db.GetIntSlice(key)
			rtn := make([]int, len(val))
			for i := 0; i < len(val); i++ {
				rtn[i] = int(val[i])
			}

			return rtn
		case reflect.Int64:
			val := o.db.GetIntSlice(key)
			rtn := make([]int64, len(val))
			for i := 0; i < len(val); i++ {
				rtn[i] = int64(val[i])
			}

			return rtn
		case reflect.String:
			return o.db.GetStringSlice(key)
		}
	}

	return nil
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
