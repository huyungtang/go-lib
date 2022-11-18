package goldap

import (
	"strconv"

	base "github.com/go-ldap/ldap/v3"
	"github.com/huyungtang/go-lib/ldap"
	"github.com/huyungtang/go-lib/times"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// entity *********************************************************************************************************************************
type entity struct {
	*base.Entry
}

// Attr
// ****************************************************************************************************************************************
func (o *entity) Attr(name string) string {
	return o.GetAttributeValue(name)
}

// AttrInt
// ****************************************************************************************************************************************
func (o *entity) AttrInt(name string) (i int64, e error) {
	if s := o.GetAttributeValue(name); s != "" {
		return strconv.ParseInt(s, 10, 64)
	}

	return
}

// Attrs
// ****************************************************************************************************************************************
func (o *entity) Attrs(name string) []string {
	return o.GetAttributeValues(name)
}

// DN
// ****************************************************************************************************************************************
func (o *entity) DN() string {
	return o.Entry.DN
}

// IsValid
// ****************************************************************************************************************************************
func (o *entity) IsValid() bool {
	if exp, err := o.AttrInt(ldap.ShadowExpire); err == nil && (exp == -1 || exp >= times.Today().UnixDay()) {
		return true
	}

	return false
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
