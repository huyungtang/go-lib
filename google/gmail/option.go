package gmail

import (
	"net/mail"
	"os"

	"github.com/huyungtang/go-lib/google"
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
	headerFrom    string = "From"
	headerReplyTo string = "Reply-To"
)

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// AttachOption
// ****************************************************************************************************************************************
func AttachOption(filename string) google.Option {
	return func(o *google.Context) {
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
func SubjectOption(subject string) google.Option {
	return func(o *google.Context) {
		o.Header.Set(headerSubject, bEncode(subject))
	}
}

// SendToOption
// ****************************************************************************************************************************************
func SendToOption(addr, name string) google.Option {
	return func(o *google.Context) {
		o.Header.Add(headerSendTos, (&mail.Address{Name: name, Address: addr}).String())
	}
}

// SendCcOption
// ****************************************************************************************************************************************
func SendCcOption(addr, name string) google.Option {
	return func(o *google.Context) {
		o.Header.Add(headerSendCcs, (&mail.Address{Name: name, Address: addr}).String())
	}
}

// SendBccOption
// ****************************************************************************************************************************************
func SendBccOption(addr, name string) google.Option {
	return func(o *google.Context) {
		o.Header.Add(headerSendBcc, (&mail.Address{Name: name, Address: addr}).String())
	}
}

// FromOption
// ****************************************************************************************************************************************
func FromOption(addr, name string) google.Option {
	return func(o *google.Context) {
		o.Header.Add(headerFrom, (&mail.Address{Name: name, Address: addr}).String())
	}
}

// ReplyToOption
// ****************************************************************************************************************************************
func ReplyToOption(mail string) google.Option {
	return func(o *google.Context) {
		o.Header.Add(headerReplyTo, mail)
	}
}

// BodyOption
// ****************************************************************************************************************************************
func BodyOption(body string) google.Option {
	return func(o *google.Context) {
		o.Body = []byte(body)
	}
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
