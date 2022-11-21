package config

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

var (
	YAMLOption Options = typeOption("yaml")
	JSONOption Options = typeOption("json")
)

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// EnvironmentOption
// ****************************************************************************************************************************************
func EnvironmentOption(envs ...string) Options {
	return func(o *Option) {
		o.Envs = append(o.Envs, envs...)
	}
}

// NameOption
// ****************************************************************************************************************************************
func NameOption(name string) Options {
	return func(o *Option) {
		o.Name = name
	}
}

// PathOption
// ****************************************************************************************************************************************
func PathOption(pathes ...string) Options {
	return func(o *Option) {
		o.Pathes = append(o.Pathes, pathes...)
	}
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Option
// ****************************************************************************************************************************************
type Option struct {
	Envs     []string
	Name     string
	Pathes   []string
	FileType string
}

// ApplyOptions
// ****************************************************************************************************************************************
func (o *Option) ApplyOptions(opts []Options, defa ...Options) (opt *Option) {
	opts = append(defa, opts...)
	for _, optFn := range opts {
		optFn(o)
	}

	return o
}

// Options
// ****************************************************************************************************************************************
type Options func(*Option)

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// typeOption *****************************************************************************************************************************
func typeOption(tp string) Options {
	return func(o *Option) {
		o.FileType = tp
	}
}
