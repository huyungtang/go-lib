package numbers

import (
	"testing"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// TestMaxInt
// ****************************************************************************************************************************************
func TestMaxInt(t *testing.T) {
	is := []int{3, 2, 6, 4, 5, 1, 7, 8, 9, 0}
	if MaxInt(is[0], is[1:]...) != 9 {
		t.Fail()
	}
}

// TestMinInt
// ****************************************************************************************************************************************
func TestMinInt(t *testing.T) {
	is := []int{3, 2, 6, 4, 5, 1, 7, 8, 9, 0}
	if MinInt(is[0], is[1:]...) != 0 {
		t.Fail()
	}
}

// TestRand
// ****************************************************************************************************************************************
func TestRand(t *testing.T) {
	mi, mx := 20, 100
	rn := Rand(mi, mx)
	for i := 0; i < 10; i++ {
		if r := rn.NextInt(); r < mi || r > mx {
			t.Fail()
			break
		}
	}
}

// TestGetPathParent
// ****************************************************************************************************************************************
func TestGetPathParent(t *testing.T) {
	if r, s := GetPathParent(1061300, 100); r != 1060000 || s != 100 {
		t.Fail()
	}

	if r, s := GetPathParent(1060000, 100); r != 1000000 || s != 10000 {
		t.Fail()
	}

	if r, s := GetPathParent(1061300, 1000); r != 1061000 || s != 1 {
		t.Fail()
	}
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
