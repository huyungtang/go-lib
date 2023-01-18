package gorm

import "errors"

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

var (
	errInvalidDriver = errors.New("cannot identify the database driver")
	errInvalidData   = errors.New("value is not a valid data type")
)

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************