package gmail

import (
	base "google.golang.org/api/gmail/v1"
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

// result *********************************************************************************************************************************
type result struct {
	*base.Message
	err error
}

// Err
// ****************************************************************************************************************************************
func (o *result) Err() (err error) {
	return o.err
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************