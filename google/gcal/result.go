package gcal

import base "google.golang.org/api/calendar/v3"

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
	*base.Event
	err error
}

// EventId
// ****************************************************************************************************************************************
func (o *result) EventId() (id string) {
	if o.Event != nil {
		id = o.Event.Id
	}

	return
}

// Err
// ****************************************************************************************************************************************
func (o *result) Err() (err error) {
	return o.err
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
