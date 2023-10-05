package db

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
type Identity interface {
	SetId(interface{})
}

// Paged
// ****************************************************************************************************************************************
type Paged interface {
	GetPageIndex() int
	GetPagedSize() int
	GetDataDTO() interface{}
	SetCount(int64)
}

// Created
// ****************************************************************************************************************************************
type Created interface {
	Create()
}

// Updated
// ****************************************************************************************************************************************
type Updated interface {
	Update()
}

// Deleted
// ****************************************************************************************************************************************
type Deleted interface {
	Delete()
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
