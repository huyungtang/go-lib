package goldap

import (
	"errors"
	"net/url"

	base "github.com/go-ldap/ldap/v3"
	"github.com/huyungtang/go-lib/ldap"
	"github.com/huyungtang/go-lib/strings"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Init
// ****************************************************************************************************************************************
func Init(dsn string, opts ...ldap.Options) (client ldap.Client, err error) {
	var u *url.URL
	if u, err = url.Parse(dsn); err != nil {
		return
	}

	var conn *base.Conn
	if conn, err = base.DialURL(strings.Format("%s://%s", u.Scheme, u.Host)); err != nil {
		return
	}

	pswd, _ := u.User.Password()
	opts = append([]ldap.Options{
		ldap.BindRequestOption(u.User.Username(), pswd),
	}, opts...)
	cfg := ldap.ApplyOptions(opts)
	if _, err = conn.SimpleBind(cfg.SimpleBindRequest); err != nil {
		return
	}

	return &database{conn, cfg}, nil
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// database
// ****************************************************************************************************************************************
type database struct {
	*base.Conn
	*ldap.Option
}

// GetUser
// ****************************************************************************************************************************************
func (o *database) GetUser(user string, opts ...ldap.Options) (ety ldap.Entity, err error) {
	var res ldap.Result
	if res, err = o.GetUsers([]string{user}, opts...); err != nil {
		return
	}

	if res.Len() != 1 {
		return
	}

	return res.Index(0), nil
}

// GetUsers
// ****************************************************************************************************************************************
func (o *database) GetUsers(users []string, opts ...ldap.Options) (rtn ldap.Result, err error) {
	fs := make([]string, len(users))
	for i, user := range users {
		fs[i] = strings.Format(("(uid=%s)"), user)
	}

	return o.Search(strings.Format("(&(|%s)(objectClass=person))", strings.Join(fs, "")), opts...)
}

// GetGroups
// ****************************************************************************************************************************************
func (o *database) GetGroups(groups []string, opts ...ldap.Options) (rtn ldap.Result, err error) {
	fs := make([]string, len(groups))
	for i, group := range groups {
		fs[i] = strings.Format(("(cn=%s)"), group)
	}

	return o.Search(strings.Format("(&(|%s)(objectClass=posixGroup))", strings.Join(fs, "")), opts...)
}

// Search
// ****************************************************************************************************************************************
func (o *database) Search(filter string, opts ...ldap.Options) (rtn ldap.Result, err error) {
	opts = append([]ldap.Options{
		ldap.BaseDNOption(o.Option.BaseDN),
		ldap.ScopeWholeSubtreeOption,
	}, opts...)
	cfg := ldap.ApplyOptions(opts)
	req := base.NewSearchRequest(
		cfg.BaseDN,
		cfg.Scope,
		cfg.Alias,
		0,
		0,
		false,
		filter,
		cfg.Attrs, nil)

	var res *base.SearchResult
	if res, err = o.bind().Conn.Search(req); err != nil {
		return
	}

	return &result{res}, nil
}

// Signin
// ****************************************************************************************************************************************
func (o *database) Signin(user, pswd string) (err error) {
	if err = o.Conn.Bind(o.Option.GetUserDN(user), pswd); err != nil {
		return
	}

	var res ldap.Entity
	if res, err = o.GetUser(user, ldap.AttrShadowExpireOption); err != nil {
		return
	}

	if !res.IsValid() {
		return errors.New("user is not valid")
	}

	return
}

// Password
// ****************************************************************************************************************************************
func (o *database) Password(user, oriPswd, newPswd string) (err error) {
	req := base.NewPasswordModifyRequest(o.Option.GetUserDN(user), oriPswd, newPswd)
	_, err = o.Conn.PasswordModify(req)

	return
}

// Close
// ****************************************************************************************************************************************
func (o *database) Close() (err error) {
	if o.Conn != nil {
		o.Conn.Close()
	}

	return
}

// bind ***********************************************************************************************************************************
func (o *database) bind() *database {
	o.Conn.SimpleBind(o.Option.SimpleBindRequest)

	return o
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
