package gmail

import (
	"net/mail"

	"github.com/huyungtang/go-lib/google"
	"google.golang.org/api/gmail/v1"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

const (
	mailTo  mailReceiver = "to"
	mailCc  mailReceiver = "cc"
	mailBcc mailReceiver = "bcc"
)

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// GmailModifyScopeOption
// ****************************************************************************************************************************************
func GmailModifyScopeOption(jsonKey string) google.Options {
	return google.ConfigOption(jsonKey, gmail.GmailModifyScope)
}

// MailSendToOption
// ****************************************************************************************************************************************
func MailSendToOption(addr, name string) google.Options {
	return mailReceiverOption(mailTo, addr, name)
}

// MailSendCcOption
// ****************************************************************************************************************************************
func MailSendCcOption(addr, name string) google.Options {
	return mailReceiverOption(mailCc, addr, name)
}

// MailSendBccOption
// ****************************************************************************************************************************************
func MailSendBccOption(addr, name string) google.Options {
	return mailReceiverOption(mailBcc, addr, name)
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// mailReceiver ***************************************************************************************************************************
type mailReceiver = string

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// mailReceiverOption *********************************************************************************************************************
func mailReceiverOption(key mailReceiver, addr, name string) google.Options {
	return func(o *google.Option) {
		m := &mail.Address{Name: name, Address: addr}
		switch key {
		case mailTo:
			o.MailTo = append(o.MailTo, m)
		case mailCc:
			o.MailCc = append(o.MailCc, m)
		case mailBcc:
			o.MailBcc = append(o.MailBcc, m)
		}
	}
}
