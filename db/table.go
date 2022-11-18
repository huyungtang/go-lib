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

// Table
// ****************************************************************************************************************************************
type Table interface {
	Query

	Select(interface{}, ...interface{}) Table
	Omit(...string) Table

	Join(string, ...interface{}) Table

	Where(interface{}, ...interface{}) Table
	Having(interface{}, ...interface{}) Table

	Order(interface{}) Table
	Group(string) Table

	Offset(int) Table
	Limit(int) Table

	Create(interface{}) error
	Get(interface{}) error
	Update(interface{}) error
	Delete() error
	Count() (int64, error)

	Begin() Transaction

	ResetClauses()
}

// Transaction
// ****************************************************************************************************************************************
type Transaction interface {
	Table

	Rollback()
	Commit() error
}

// Query
// ****************************************************************************************************************************************
type Query interface {
	GetById(interface{}, interface{}) error
	GetPagedList(interface{}) error
}

// TableName
// ****************************************************************************************************************************************
type TableName interface {
	TableName() string
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
