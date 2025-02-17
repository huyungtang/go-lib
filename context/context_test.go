package context

import (
	"testing"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// TestContext
// ****************************************************************************************************************************************
func TestContext(t *testing.T) {
	key := new(struct{})
	ctx := Handler(
		HandlerOption(func(ctx Context) error {
			ctx.WithValue(key, 10)

			if e := ctx.Next(); e != nil {
				return e
			}

			ctx.WithValue(key, 19)

			return ctx.Next()
		}),
		HandlerOption(func(ctx Context) error {
			ctx.WithValue(key, 20)

			if e := ctx.Next(); e != nil {
				return e
			}

			ctx.WithValue(key, 29)

			return ctx.Next()
		}),
		HandlerOption(func(ctx Context) error {
			ctx.WithValue(key, 30)

			return ctx.Next()
		}),
	)

	if e := ctx.Next(); e != nil {
		t.Error(e)
	}

	if v := ctx.Value(key); v != 19 {
		t.Fail()
	}
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
