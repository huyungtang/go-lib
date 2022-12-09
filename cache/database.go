package cache

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

const (
	Expired int64 = -2
	KeepTTL int64 = 0
	Static  int64 = -1
)

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Database
// ****************************************************************************************************************************************
type Database interface {
	// Exists(key)
	Exists(string) bool

	// Get(key, value, DefaultOption, ExpireOption, *KeepTTLOption)
	Get(string, interface{}, ...Options) error

	// Set(key, value, ExpireOption, *StaticOption)
	Set(string, interface{}, ...Options) error

	// Push(key, value, ExpireOption, *LPushOpiton, *StaticOption)
	Push(string, interface{}, ...Options) error

	// Del(keys)
	Del(...string) error

	Close() error
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
