package gorm

import (
	"database/sql/driver"
	"errors"

	"github.com/huyungtang/go-lib/encrypt"
	"github.com/huyungtang/go-lib/strings"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Identity
// ****************************************************************************************************************************************
type Identity struct {
	Id uint64 `gorm:"column:id;primaryKey"`
}

// String
// ****************************************************************************************************************************************
type String string

// Scan
// ****************************************************************************************************************************************
func (o *String) Scan(val interface{}) (err error) {
	var s string
	if s, err = decryptString(val); err != nil {
		return
	}

	*o = String(s)

	return
}

// Value
// ****************************************************************************************************************************************
func (o String) Value() (val driver.Value, err error) {
	return encryptString(o.String())
}

// Set
// ****************************************************************************************************************************************
func (o *String) Set(s string) {
	*o = String(s)
}

// String
// ****************************************************************************************************************************************
func (o String) String() string {
	return string(o)
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// decryptString **************************************************************************************************************************
func decryptString(val interface{}) (s string, err error) {
	bs, isOK := val.([]byte)
	if !isOK {
		return "", errors.New("value is not a valid data type")
	}

	if len(bs) == 0 {
		return
	}

	s = strings.Format("$a%d$%s", encrypt.DefaCost, string(bs))
	if s, err = encrypt.Decrypt(s); err != nil {
		return
	}

	return
}

// encryptString **************************************************************************************************************************
func encryptString(s string) (val driver.Value, err error) {
	val = ""
	if len(s) > 0 {
		if s, err = encrypt.Encrypt(s, encrypt.DefaCost); err != nil {
			return
		}

		val = s[5:]
	}

	return
}
