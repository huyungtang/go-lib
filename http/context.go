package http

import (
	"io"
	base "net/http"
	"net/url"
	"strings"

	"github.com/huyungtang/go-lib/reflect"
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

// context ********************************************************************************************************************************
type context struct {
	*Option
	*base.Response
	body []byte
	err  error
}

// Get
// ****************************************************************************************************************************************
func (o *context) Get(path string, opts ...Options) {
	ctx := reflect.Clone(o).(*context)
	ctx.requestCore(base.MethodGet, path, opts)
}

// Post
// ****************************************************************************************************************************************
func (o *context) Post(path string, opts ...Options) {
	ctx := reflect.Clone(o).(*context)
	opts = append(opts, urlencodedOption)
	ctx.requestCore(base.MethodGet, path, opts)
}

// requestCore ****************************************************************************************************************************
func (o *context) requestCore(method, path string, opts []Options) {
	handler := HandlerOption(func(sc Session) {
		if ctx, isOK := sc.(*context); isOK {
			vals := &url.Values{}
			for _, pars := range ctx.Params {
				vals.Add(pars[0], pars[1])
			}

			var body io.Reader
			switch method {
			case base.MethodGet:
				ctx.Host.RawQuery = vals.Encode()
			case base.MethodPost:
				body = strings.NewReader(vals.Encode())
			}

			var req *base.Request
			ctx.Host.Path = path
			if req, ctx.err = base.NewRequest(method, ctx.Host.String(), body); ctx.err != nil {
				return
			}

			for _, hdrs := range ctx.Headers {
				req.Header.Add(hdrs[0], hdrs[1])
			}

			for _, ckes := range ctx.Ckies {
				req.AddCookie(ckes)
			}

			var resp *base.Response
			c := new(base.Client)
			if resp, ctx.err = c.Do(req); ctx.err != nil {
				return
			}
			defer resp.Body.Close()

			if resp.Body != nil {
				ctx.body, ctx.err = io.ReadAll(resp.Body)
			}
			ctx.Response = resp
			o.Ckies = resp.Cookies()
		}

		sc.Next()
	})

	opts = append([]Options{handler}, opts...)
	o.Option = ApplyOptions(opts, o.Option)
	o.Handler = -1
	o.Next()
}

// ContextHandler
// ****************************************************************************************************************************************
type ContextHandler func(Session)

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
