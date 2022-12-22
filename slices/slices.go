package slices

import (
	refl "reflect"

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
		store: make(map[interface{}]bool),
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
func Reverse(s interface{}) {
	if !reflect.IsSlice(s) {
		return
	}

	vs := reflect.ValueOf(s)
	swap := refl.Swapper(vs.Interface())
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
	Append(interface{}) bool
	GetStrings() []string
}

// sliceContainer *************************************************************************************************************************
type sliceContainer struct {
	store map[interface{}]bool
}

// Append
// ****************************************************************************************************************************************
func (o *sliceContainer) Append(v interface{}) (isApp bool) {
	if _, isOK := o.store[v]; !isOK {
		o.store[v], isApp = true, true
	}

	return
}

// GetStrings
// ****************************************************************************************************************************************
func (o *sliceContainer) GetStrings() (strs []string) {
	strs = make([]string, 0, len(o.store))
	for k := range o.store {
		if v, isOK := k.(string); isOK {
			strs = append(strs, v)
		}
	}

	return
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
