package redis

import (
	"testing"

	"github.com/huyungtang/go-lib/cache"
	"github.com/huyungtang/go-lib/config"
	"github.com/huyungtang/go-lib/config/viper"
	"github.com/huyungtang/go-lib/file"
	"github.com/huyungtang/go-lib/strings"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// TestGet
// ****************************************************************************************************************************************
func TestGet(t *testing.T) {
	c, err := viper.Init(
		config.PathOption(file.PathWorking("../_testing")),
		config.EnvironmentOption("prod"),
	)
	if err != nil {
		t.Error(err)
	}

	var db cache.Database
	if db, err = Init(c.GetString("Redis", "")); err != nil {
		t.Error(err)
	}
	defer db.Close()

	var str string
	key1 := "testing:key1"
	if err = db.Get(key1, &str,
		cache.DefaultOption(func(i interface{}) (cache.Options, error) {
			if s, isOK := i.(*string); isOK {
				*s = "default value"
			}
			return cache.ExpireOption(30), nil
		}),
		cache.KeepTTLOption,
	); err != nil {
		t.Error(err)
	}

	if str != "default value" {
		t.Fail()
	}
}

// TestPop
// ****************************************************************************************************************************************
func TestPop(t *testing.T) {
	c, err := viper.Init(
		config.PathOption(file.PathWorking("../_testing")),
		config.EnvironmentOption("prod"),
	)
	if err != nil {
		t.Error(err)
	}

	var db cache.Database
	if db, err = Init(c.GetString("Redis", "")); err != nil {
		t.Error(err)
	}
	defer db.Close()

	key := "testing.stringslice"
	val := []rune{65, 66, 67, 68, 69, 70}
	for _, v := range val {
		if err = db.Push(key, string(v)); err != nil {
			t.Error(err)
		}
	}

	strs := make([]string, 0)
	if err = db.Pop(key, &strs, cache.PopCountOption(uint64(len(val)))); err != nil {
		t.Error(err)
	}

	if string(val) != strings.Join(strs, "") {
		t.Fail()
	}

}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
