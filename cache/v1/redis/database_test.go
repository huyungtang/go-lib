package redis

import (
	"sync"
	"testing"
	"time"

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

// TestIncrease
// ****************************************************************************************************************************************
func TestIncrease(t *testing.T) {
	db, err := Init(dsn)
	if err != nil {
		t.Error()
	}

	defer func() {
		if err = db.Close(); err != nil {
			t.Error()
		}
	}()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(200 * time.Millisecond)
			defer wg.Done()
			var val int64
			if err = db.Increase("testing_increase_int", &val, cache.IncreaseByOption(5), cache.ExpireOption(60)); err != nil {
				t.Error(err)
			}
			t.Log("val  -> ", val)

			var val1 float64
			if err = db.Increase("testing_increase_float64", &val1, cache.IncreaseByFloatOption(1.3), cache.ExpireOption(60)); err != nil {
				t.Error(err)
			}
			t.Log("val1 -> ", val1)
		}()
	}

	wg.Wait()

}

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
