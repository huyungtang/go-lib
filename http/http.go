package http

import (
	"encoding/json"
	"io"
	"io/ioutil"
	base "net/http"
	"net/url"
	"reflect"
	"strings"

	"github.com/huyungtang/go-lib/file"
	"github.com/huyungtang/go-lib/reflects"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// New
// ****************************************************************************************************************************************
func New(host string, opts ...Option) Context {
	ctx := &context{
		client:  new(base.Client),
		host:    host,
		hdls:    make([]contextHandler, 0),
		headers: make([]*headerOption, 0),
		idx:     -1,
	}

	for i := 0; i < len(opts); i++ {
		switch opt := opts[i].(type) {
		case *handlerOption:
			ctx.hdls = append(ctx.hdls, opt.hdl)
		case *headerOption:
			ctx.headers = append(ctx.headers, opt)
		}
	}

	return ctx
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// contextHandler *************************************************************************************************************************
type contextHandler func(Context) error

// Context
// ****************************************************************************************************************************************
type Context interface {
	Next()
	Get(url string) Response
	Post(url string, opts ...Option) Response
	Err() error
}

// Response
// ****************************************************************************************************************************************
type Response interface {
	StatusCode() int
	Status() string
	Body() string
	Parse(interface{}) error
	Err() error
}

// context ********************************************************************************************************************************
type context struct {
	client  *base.Client
	host    string
	err     error
	idx     int
	hdls    []contextHandler
	headers []*headerOption
	body    []byte
	resp    *base.Response
	cookies []*base.Cookie
}

// Next
// ****************************************************************************************************************************************
func (o *context) Next() {
	o.idx++
	if o.idx >= len(o.hdls) || o.hdls[o.idx] == nil {
		return
	}

	if err := o.hdls[o.idx](o); err != nil {
		o.err = err
		return
	}
}

// Get
// ****************************************************************************************************************************************
func (o *context) Get(uri string) Response {
	ctx := o.clone()
	ctx.hdls = append(ctx.hdls, ctx.requestCore(base.MethodGet, ctx.getUri(uri), nil))
	ctx.Next()

	if ctx.Err() == nil {
		o.cookies = ctx.resp.Cookies()
	}

	return ctx
}

// Post
// ****************************************************************************************************************************************
func (o *context) Post(uri string, opts ...Option) Response {
	ctx := o.clone()
	ctp := &headerOption{key: "Content-Type", value: "application/x-www-form-urlencoded"}
	val := &url.Values{}
	for i := 0; i < len(opts); i++ {
		switch opt := opts[i].(type) {
		case *paramOption:
			val.Add(opt.key, opt.value)
		case *contentOption:
			ctp.value = opt.value
		}
	}
	ctx.headers = append(ctx.headers, ctp)

	ctx.hdls = append(ctx.hdls, ctx.requestCore(base.MethodGet, ctx.getUri(uri), strings.NewReader(val.Encode())))
	ctx.Next()

	if ctx.Err() == nil {
		o.cookies = ctx.resp.Cookies()
	}

	return ctx
}

// Err
// ****************************************************************************************************************************************
func (o *context) Err() (err error) {
	return o.err
}

// StatusCode
// ****************************************************************************************************************************************
func (o *context) StatusCode() int {
	return o.resp.StatusCode
}

// StatusCode
// ****************************************************************************************************************************************
func (o *context) Status() string {
	return o.resp.Status
}

// Body
// ****************************************************************************************************************************************
func (o *context) Body() string {
	if o.resp.StatusCode != base.StatusOK {
		return o.resp.Status
	}

	return string(o.body)
}

// Parse
// ****************************************************************************************************************************************
func (o *context) Parse(dto interface{}) (err error) {
	return json.Unmarshal(o.body, dto)
}

// getUri *********************************************************************************************************************************
func (o *context) getUri(uri string) (rtn string) {
	val, err := url.Parse(o.host)
	if err != nil {
		return
	}
	val.Path = file.PathJoin(val.Path, uri)

	return val.String()
}

// clone **********************************************************************************************************************************
func (o *context) clone() *context {
	tar := reflect.New(reflects.TypeOf(o))
	tar.Elem().Set(reflects.ValueOf(o))

	return tar.Interface().(*context)
}

// requestCore ****************************************************************************************************************************
func (o *context) requestCore(method, uri string, body io.Reader) contextHandler {
	return func(ctx Context) (err error) {
		var req *base.Request
		if req, err = base.NewRequest(method, uri, body); err != nil {
			o.err = err
			return
		}

		for i := 0; i < len(o.headers); i++ {
			req.Header.Add(o.headers[i].key, o.headers[i].value)
		}

		for i := 0; i < len(o.cookies); i++ {
			req.AddCookie(o.cookies[i])
		}

		if o.resp, err = o.client.Do(req); err != nil {
			o.err = err
			return
		}
		defer o.resp.Body.Close()

		if o.body, err = ioutil.ReadAll(o.resp.Body); err != nil {
			o.err = err
			return
		}

		return
	}
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
