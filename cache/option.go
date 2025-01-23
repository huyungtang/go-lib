package cache

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

var (
	ExpiredOption Option = ExpireOption(Expired)
	KeepTTLOption Option = ExpireOption(KeepTTL)
	StaticOption  Option = ExpireOption(Static)

	SkipOverrideOption Option = overrideOption("NX")
	UpdateOnlyOption   Option = overrideOption("XX")

	LPushOption Option = cmderOption("LPUSH")
	RPushOption Option = cmderOption("RPUSH")

	LPopOption Option = cmderOption("LPOP")
	RPopOption Option = cmderOption("RPOP")
)

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// DefaultOption
// ****************************************************************************************************************************************
func DefaultOption(fn DefaultFn) Option {
	return func(co *Context) {
		co.DefaFn = fn
	}
}

// ExpireOption
// ****************************************************************************************************************************************
func ExpireOption(sec int64) Option {
	return func(co *Context) {
		co.Expire = sec
	}
}

// PopCountOption
// ****************************************************************************************************************************************
func PopCountOption(count uint64) Option {
	return func(o *Context) {
		if count == 0 {
			o.Count = 1
		} else {
			o.Count = count
		}
	}
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Context
// ****************************************************************************************************************************************
type Context struct {
	DefaFn   DefaultFn
	Expire   int64
	Override string
	Cmder    string
	Count    uint64
}

// ApplyOptions
// ****************************************************************************************************************************************
func (o *Context) ApplyOptions(opts []Option, defa ...Option) (opt *Context) {
	opts = append(defa, opts...)
	for _, optFn := range opts {
		optFn(o)
	}

	return o
}

// Option
// ****************************************************************************************************************************************
type Option func(*Context)

// DefaultFn
// ****************************************************************************************************************************************
type DefaultFn func(interface{}) (Option, error)

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// overrideOption *************************************************************************************************************************
func overrideOption(over string) Option {
	return func(co *Context) {
		co.Override = over
	}
}

// cmderOption ****************************************************************************************************************************
func cmderOption(dir string) Option {
	return func(o *Context) {
		o.Cmder = dir
	}
}
