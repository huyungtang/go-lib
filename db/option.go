package db

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// ApplyOptions
// ****************************************************************************************************************************************
func ApplyOptions(opts []Options) (opt *Option) {
	opt = new(Option)
	for _, optFn := range opts {
		optFn(opt)
	}

	return
}

// SkipDefaultTransactionOption
// ****************************************************************************************************************************************
func SkipDefaultTransactionOption(s bool) Options {
	return func(d *Option) {
		d.SkipDefaultTransaction = s
	}
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Option DO NOT DIRECT TO USE THIS
// ****************************************************************************************************************************************
type Option struct {
	SkipDefaultTransaction bool
}

// Options
// ****************************************************************************************************************************************
type Options func(*Option)

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
