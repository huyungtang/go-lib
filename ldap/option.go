package ldap

import (
	base "github.com/go-ldap/ldap/v3"
	"github.com/huyungtang/go-lib/strings"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

const (
	ShadowExpire = "shadowExpire"
)

var (
	AttrShadowExpireOption = AttributesOption(ShadowExpire)

	ScopeBaseObjectOption   = scopeOption(0)
	ScopeSingleLevelOption  = scopeOption(1)
	ScopeWholeSubtreeOption = scopeOption(2)

	DerefAliasesNeverOption       = aliasesOption(0)
	DerefAliasesInSearchingOption = aliasesOption(1)
	DerefAliasesFindingBaseObj    = aliasesOption(2)
	DerefAliasesAlways            = aliasesOption(3)
)

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// AttributesOption
// ****************************************************************************************************************************************
func AttributesOption(attrs ...string) Options {
	return func(o *Option) {
		o.Attrs = append(o.Attrs, attrs...)
	}
}

// BaseDNOption
// ****************************************************************************************************************************************
func BaseDNOption(dn string) Options {
	return func(o *Option) {
		o.BaseDN = dn
	}
}

// BindRequestOption
// ****************************************************************************************************************************************
func BindRequestOption(user, pswd string) Options {
	return func(o *Option) {
		o.SimpleBindRequest = &base.SimpleBindRequest{
			Username:           user,
			Password:           pswd,
			AllowEmptyPassword: pswd == "",
		}
	}
}

// GroupFilterOption
// ****************************************************************************************************************************************
func GroupFilterOption(filter string) Options {
	return func(o *Option) {
		o.GroupFilter = filter
	}
}

// UserFilterOption
// ****************************************************************************************************************************************
func UserFilterOption(filter string) Options {
	return func(o *Option) {
		o.UserFilter = filter
	}
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Option
// ****************************************************************************************************************************************
type Option struct {
	*base.SimpleBindRequest
	Alias       int
	Attrs       []string
	BaseDN      string
	Scope       int
	UserFilter  string
	GroupFilter string
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

// GetUserDN
// ****************************************************************************************************************************************
func (o *Option) GetUserDN(user string) string {
	return strings.Format("uid=%s,cn=users,%s", user, o.BaseDN)
}

// Options
// ****************************************************************************************************************************************
type Options func(*Option)

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// scopeOption ****************************************************************************************************************************
func scopeOption(scope int) Options {
	return func(o *Option) {
		o.Scope = scope
	}
}

// aliasesOption **************************************************************************************************************************
func aliasesOption(alias int) Options {
	return func(o *Option) {
		o.Alias = alias
	}
}
