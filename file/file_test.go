package file

import (
	"testing"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// TestPathSavename
// ****************************************************************************************************************************************
func TestPathSavename(t *testing.T) {
	if len(PathSavename("/root/", "", 1)) != 41 ||
		len(PathSavename("/root/", "abcd1234", 1)) != 41 ||
		len(PathSavename("/root/", "", 5)) != len(PathSavename("/root/", "", 6)) {
		t.Fail()
	}
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
