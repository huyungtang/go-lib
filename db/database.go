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

// SqlDB
// ****************************************************************************************************************************************
type SqlDB interface {
	Table(interface{}) Table
	Close() error
	Ping() error
}

// NoSqlDB
// ****************************************************************************************************************************************
type NoSqlDB interface {
	Collection(interface{}) Collection
	Close() error
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
