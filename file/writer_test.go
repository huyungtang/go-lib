package file

import "testing"

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// TestWriter
// ****************************************************************************************************************************************
func TestWriter(t *testing.T) {
	writer, err := InitWriter(PathWorking("_testing/big5w.txt"), Big5)
	if err != nil {
		t.Error(err)
	}
	defer writer.Close()

	writer.Writeln("測試 big 編碼錯字")

	writer.Write("瑜")
	writer.Write("瑜")
	writer.Writeln("瑜")

	writer.Write("婷")
	writer.Write("婷")
	writer.Writeln("婷")
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
