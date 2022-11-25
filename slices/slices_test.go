package slices

import "testing"

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// TestAppend
// ****************************************************************************************************************************************
func TestAppend(t *testing.T) {
	c := New()
	if !c.Append("abc") {
		t.Fail()
	}

	if c.Append("abc") {
		t.Fail()
	}
}

// TestGetString
// ****************************************************************************************************************************************
func TestGetString(t *testing.T) {
	c := New()
	c.Append("a")
	c.Append("b")
	c.Append("b")
	c.Append("b")
	c.Append("c")

	s := c.GetStrings()
	if len(s) != 3 {
		t.Fail()
	}
}

// TestIndexOf
// ****************************************************************************************************************************************
func TestIndexOf(t *testing.T) {
	vals := []string{"a", "b", "c", "d", "e", "f", "g", "h"}
	if idx, exi := IndexOf(len(vals), func(i int) bool { return vals[i] == "d" }); !(exi && idx == 3) {
		t.Fail()
	}

	if idx, exi := IndexOf(len(vals), func(i int) bool { return vals[i] == "z" }); !(!exi && idx == -1) {
		t.Fail()
	}
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
