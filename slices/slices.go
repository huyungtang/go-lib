package slices

import (
	"github.com/huyungtang/go-lib/reflect"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// New
// ****************************************************************************************************************************************
func New() Slice {
	return &sliceContainer{
		store: make(map[any]bool),
	}
}

// IndexOf
// ****************************************************************************************************************************************
func IndexOf(lns int, compare func(int) bool) (idx int, exi bool) {
	for i := 0; i < lns; i++ {
		if compare(i) {
			return i, true
		}
	}

	return -1, false
}

// Reverse
// ****************************************************************************************************************************************
func Reverse(s any) {
	if !reflect.IsSlice(s) {
		return
	}

	vs := reflect.ValueOf(s)
	swap := reflect.Swapper(vs.Interface())
	for i, j := 0, vs.Len()-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Slice
// ****************************************************************************************************************************************
type Slice interface {
	Append(any) bool
	GetString() []string
	GetUInt64() []uint64
}

// sliceContainer *************************************************************************************************************************
type sliceContainer struct {
	store map[any]bool
}

// Append
// ****************************************************************************************************************************************
func (o *sliceContainer) Append(v any) (isApp bool) {
	if _, isOK := o.store[v]; !isOK {
		o.store[v], isApp = true, true
	}

	return
}

// GetString
// ****************************************************************************************************************************************
func (o *sliceContainer) GetString() (strs []string) {
	strs = make([]string, 0, len(o.store))
	for k := range o.store {
		if v, isOK := k.(string); isOK {
			strs = append(strs, v)
		}
	}

	return
}

// GetUInt64
// ****************************************************************************************************************************************
func (o *sliceContainer) GetUInt64() (vals []uint64) {
	vals = make([]uint64, 0, len(o.store))
	for k := range o.store {
		if v, isOk := k.(uint64); isOk {
			vals = append(vals, v)
		}
	}

	return
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
