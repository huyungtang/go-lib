package redis

import (
	"testing"

	"github.com/huyungtang/go-lib/cache/v1"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

const (
	dsn string = "redis://192.168.0.31:30948/15"
)

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// TestExists
// ****************************************************************************************************************************************
func TestExists(t *testing.T) {
	db, err := Init(dsn, cache.DebugOption())
	if err != nil {
		t.Error()
	}

	defer func() {
		if err = db.Close(); err != nil {
			t.Error()
		}
	}()

	if _, err = db.Exists("testing:isExists"); err != nil {
		t.Error()
	}
}

// TestSet
// ****************************************************************************************************************************************
func TestSet(t *testing.T) {
	db, err := Init(dsn, cache.DebugOption())
	if err != nil {
		t.Error()
	}
	defer func() {
		if err = db.Close(); err != nil {
			t.Error()
		}
	}()

	if err = db.Set("testing:mapData", &struct{ A string }{A: "222"}, cache.ExpireOption(1000)); err != nil {
		t.Error(err)
	}
}

// TestGet
// ****************************************************************************************************************************************
func TestGet(t *testing.T) {
	db, err := Init(dsn, cache.DebugOption())
	if err != nil {
		t.Error()
	}
	defer func() {
		if err = db.Close(); err != nil {
			t.Error()
		}
	}()

	var val int
	if err = db.Get("testing:mapData", val, cache.DefaultOption(func(v interface{}) ([]cache.Options, error) {
		val = 3
		return []cache.Options{cache.ExpireOption(100)}, nil
	})); err != nil {
		t.Error(err)
	}

	t.Log("value -> ", val)
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
