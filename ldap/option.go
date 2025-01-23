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
func AttributesOption(attrs ...string) Option {
	return func(o *Context) {
		o.Attrs = append(o.Attrs, attrs...)
	}
}

// BaseDNOption
// ****************************************************************************************************************************************
func BaseDNOption(dn string) Option {
	return func(o *Context) {
		o.BaseDN = dn
	}
}

// BindRequestOption
// ****************************************************************************************************************************************
func BindRequestOption(user, pswd string) Option {
	return func(o *Context) {
		o.SimpleBindRequest = &base.SimpleBindRequest{
			Username:           user,
			Password:           pswd,
			AllowEmptyPassword: pswd == "",
		}
	}
}

// GroupFilterOption
// ****************************************************************************************************************************************
func GroupFilterOption(filter string) Option {
	return func(o *Context) {
		o.GroupFilter = filter
	}
}

// UserFilterOption
// ****************************************************************************************************************************************
func UserFilterOption(filter string) Option {
	return func(o *Context) {
		o.UserFilter = filter
	}
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Context
// ****************************************************************************************************************************************
type Context struct {
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
func (o *Context) ApplyOptions(opts []Option, defa ...Option) (opt *Context) {
	opts = append(defa, opts...)
	for _, optFn := range opts {
		optFn(o)
	}

	return o
}

// GetUserDN
// ****************************************************************************************************************************************
func (o *Context) GetUserDN(user string) string {
	return strings.Format("uid=%s,cn=users,%s", user, o.BaseDN)
}

// Option
// ****************************************************************************************************************************************
type Option func(*Context)

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// scopeOption ****************************************************************************************************************************
func scopeOption(scope int) Option {
	return func(o *Context) {
		o.Scope = scope
	}
}

// aliasesOption **************************************************************************************************************************
func aliasesOption(alias int) Option {
	return func(o *Context) {
		o.Alias = alias
	}
}
