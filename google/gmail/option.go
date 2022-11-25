package gmail

import (
	"net/mail"
	"os"

	"github.com/huyungtang/go-lib/google"
	"google.golang.org/api/gmail/v1"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

const (
	headerSubject string = "Subject"
	headerContent string = "Content-Type"
	headerHtml    string = "text/html; charset=utf-8"
	headerSendTos string = "To"
	headerSendCcs string = "Cc"
	headerSendBcc string = "Bcc"
)

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// GmailModifyScopeOption
// ****************************************************************************************************************************************
func GmailModifyScopeOption(jsonKey string) google.Options {
	return google.ConfigOption(jsonKey, gmail.GmailModifyScope)
}

// AttachOption
// ****************************************************************************************************************************************
func AttachOption(filename string) google.Options {
	return func(o *google.Option) {
		o.AttachOnce.Do(func() {
			o.Attach = make(map[string][]byte)
		})

		if file, err := os.Open(filename); err == nil {
			defer file.Close()

			info, _ := file.Stat()
			buff := make([]byte, info.Size())
			file.Read(buff)
			o.Attach[info.Name()] = buff
		}
	}
}

// SubjectOption
// ****************************************************************************************************************************************
func SubjectOption(subject string) google.Options {
	return func(o *google.Option) {
		o.Header.Set(headerSubject, bEncode(subject))
	}
}

// SendToOption
// ****************************************************************************************************************************************
func SendToOption(addr, name string) google.Options {
	return func(o *google.Option) {
		o.Header.Add(headerSendTos, (&mail.Address{Name: name, Address: addr}).String())
	}
}

// SendCcOption
// ****************************************************************************************************************************************
func SendCcOption(addr, name string) google.Options {
	return func(o *google.Option) {
		o.Header.Add(headerSendCcs, (&mail.Address{Name: name, Address: addr}).String())
	}
}

// SendBccOption
// ****************************************************************************************************************************************
func SendBccOption(addr, name string) google.Options {
	return func(o *google.Option) {
		o.Header.Add(headerSendBcc, (&mail.Address{Name: name, Address: addr}).String())
	}
}

// BodyOption
// ****************************************************************************************************************************************
func BodyOption(body string) google.Options {
	return func(o *google.Option) {
		o.Body = []byte(body)
	}
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
