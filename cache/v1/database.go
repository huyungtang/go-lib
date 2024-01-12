package cache

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Database
// ****************************************************************************************************************************************
type Database interface {
	/*
		Delete(cachedKeys)
	*/
	Delete(...string) error

	/*
		Exists(cachedKey)
	*/
	Exists(string) (bool, error)

	/*
		Get(cachedKey, value, Options)

		Options: *KeepTTLOption(), DefaultOption(), ExpireOption()
	*/
	Get(string, interface{}, ...Options) error

	/*
		GetSlice(cachedKey, slice_value)
	*/
	GetSlice(string, interface{}, ...Options) error

	/*
		Increase(cachedKey, value)

		value: int64 or float64
	*/
	Increase(string, interface{}, ...Options) error

	/*
		Push(cachedKey, value, Options)

		Options: *DirectionRightOption(), ExpireOption()
	*/
	Push(string, interface{}, ...Options) error

	/*
		Set(cachedKey, value, Options)

		Options: *StaticOption(), SkipOverrideOption(), UpdateOnlyOption(), ExpireOption()
	*/
	Set(string, interface{}, ...Options) error

	Close() error
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
