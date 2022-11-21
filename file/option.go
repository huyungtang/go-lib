package file

import (
	"os"
	"regexp"

	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/traditionalchinese"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

var (
	Big5 Options = encodingOption(traditionalchinese.Big5)

	Append   Options = flagOption(os.O_CREATE | os.O_WRONLY | os.O_APPEND)
	Override Options = flagOption(os.O_CREATE | os.O_WRONLY | os.O_TRUNC)
)

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// FilterOption
// ****************************************************************************************************************************************
func FilterOption(pattern string) Options {
	return func(o *Option) {
		o.Filter = regexp.MustCompile(pattern)
	}
}

// OperationOption
// ****************************************************************************************************************************************
func OperationOption(op FileOp) Options {
	return func(o *Option) {
		o.Op = o.Op | op
	}
}

// PathOption
// ****************************************************************************************************************************************
func PathOption(pathes ...string) Options {
	return func(o *Option) {
		o.Path = append(o.Path, pathes...)
	}
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Option
// ****************************************************************************************************************************************
type Option struct {
	encoding.Encoding
	Flag   int
	Filter *regexp.Regexp
	Op     FileOp
	Path   []string
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

// encodingOption *************************************************************************************************************************
func encodingOption(encode encoding.Encoding) Options {
	return func(o *Option) {
		o.Encoding = encode
	}
}

// flagOption *****************************************************************************************************************************
func flagOption(flag int) Options {
	return func(o *Option) {
		o.Flag = flag
	}
}
