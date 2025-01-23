package gmail

import (
	"bytes"
	"context"
	"encoding/base64"
	"mime"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"regexp"

	"github.com/huyungtang/go-lib/google"
	"github.com/huyungtang/go-lib/strings"
	base "google.golang.org/api/gmail/v1"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

var (
	bEncodeReg = regexp.MustCompile(`=\?utf-8\?b\?[^?]+\?=`)
)

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// Init
// ****************************************************************************************************************************************
func Init(opts ...google.Option) (serv Service, err error) {
	var s *base.Service
	cfg := new(google.Context).ApplyOptions(opts)
	if s, err = base.NewService(context.Background(), cfg.GetClientOption()); err != nil {
		return
	}

	return &service{s}, nil
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// service ********************************************************************************************************************************
type service struct {
	*base.Service
}

// Service
// ****************************************************************************************************************************************
type Service interface {
	Send(...google.Option) google.Message
}

// Send
// ****************************************************************************************************************************************
func (o *service) Send(opts ...google.Option) google.Message {
	var buff bytes.Buffer
	writer := multipart.NewWriter(&buff)
	defer writer.Close()

	cfg := (&google.Context{
		Header: textproto.MIMEHeader{
			headerContent: {strings.Format("multipart/mixed;\n  boundary=\"%s\"", writer.Boundary())},
		},
	}).ApplyOptions(opts)

	writer.CreatePart(cfg.Header)

	bodyContent, _ := writer.CreatePart(textproto.MIMEHeader{headerContent: {headerHtml}})
	bodyContent.Write(cfg.Body)

	for k, v := range cfg.Attach {
		attach, _ := writer.CreatePart(textproto.MIMEHeader{
			"Content-Type":              {strings.Format("%s; name=%s", http.DetectContentType(v), k)},
			"Content-Transfer-Encoding": {"base64"},
			"Content-Disposition":       {strings.Format("attachment; filename=%s", k)},
		})
		attach.Write(v)
	}

	res := new(result)
	res.Message, res.err = o.Users.Messages.Send("me", &base.Message{Raw: base64.URLEncoding.EncodeToString(buff.Bytes())}).Do()

	return res
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// bEncode ********************************************************************************************************************************
func bEncode(str string) string {
	return strings.Join(bEncodeReg.FindAllString(mime.BEncoding.Encode("utf-8", str), -1), "\n ")
}
