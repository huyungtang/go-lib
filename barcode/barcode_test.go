package barcode

import (
	"os"
	"testing"

	"github.com/huyungtang/go-lib/file"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// TestCode128
// ****************************************************************************************************************************************
func TestCode128(t *testing.T) {
	reader, err := Code128("0123456789", 280, 60)
	if err != nil {
		t.Error(reader)
	}

	bs := make([]byte, reader.Len())
	if _, err = reader.Read(bs); err != nil {
		t.Error(err)
	}

	var f *os.File
	path := file.PathWorking("_testing", "code128.png")
	if f, err = os.Create(path); err != nil {
		t.Error(err)
	}
	defer f.Close()

	f.Write(bs)

	if file.IsExist(path) != file.IsFile {
		t.Fail()
	}
}

// TestQRCODE
// ****************************************************************************************************************************************
func TestQRCODE(t *testing.T) {
	reader, err := QRCode("0123456789", 160)
	if err != nil {
		t.Error(reader)
	}

	bs := make([]byte, reader.Len())
	if _, err = reader.Read(bs); err != nil {
		t.Error(err)
	}

	var f *os.File
	path := file.PathWorking("_testing", "qrcode.png")
	if f, err = os.Create(path); err != nil {
		t.Error(err)
	}
	defer f.Close()

	f.Write(bs)

	if file.IsExist(path) != file.IsFile {
		t.Fail()
	}
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
