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

	LockForUpdate() Table

	Select(any, ...any) Table
	Omit(...string) Table

	Join(string, ...any) Table
	Preload(string, ...any) Table

	Available() Table
	Where(any, ...any) Table
	Having(any, ...any) Table

	Order(any) Table
	Group(string) Table

	Offset(int) Table
	Limit(int) Table

	Create(any) error
	Get(any) error
	Update(any) error
	Delete(any) error
	Count() (int64, error)

	Exec(string, ...any) error

	Begin() Transaction

	ResetClauses()

	CreateColumns() []string
	UpdateColumns() []string
	RowsAffected() int64

	SubQuery() any
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
	GetById(any, any) error
	GetPagedList(any) error
}

// TableName
// ****************************************************************************************************************************************
type TableName interface {
	TableName() string
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
