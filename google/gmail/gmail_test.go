package gmail

import (
	"testing"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// TestSendMail
// ****************************************************************************************************************************************
func TestSendMail(t *testing.T) {
	// c, err := viper.Init(
	// 	config.PathOption(file.PathWorking("_testing")),
	// 	config.EnvironmentOption("prod"),
	// )
	// if err != nil {
	// 	t.Error(err)
	// }

	// opt := new(Option)
	// cfg.ApplyOptions([]Options{
	// 	GmailModifyScopeOption(c.GetString("Mail.Credential", "")),
	// 	OAuthTokenOption(c.GetString("Mail.Token", "")),
	// })

	// serv, err := gmail.NewService(context.Background(), opt.GetClientOption())
	// if err != nil {
	// 	t.Error(err)
	// }
	// to1 := &mail.Address{Name: "morris", Address: "huyungtang@Gmail.com"}
	// to2 := &mail.Address{Name: "hu_yt", Address: "hu_yt@hotmail.com"}
	// // tos := []string{to1.String(), to2.String()}

	// // emailTo := "To:" + strings.Join(tos, ",") + "\r\n"
	// // subject := "Subject: " + "Test Email form Gmail API using OAuth" + "\n"
	// // mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	// // msg := []byte(emailTo + subject + mime + "\n" + "emailBody")

	// msg := &gmail.Message{
	// 	// Raw: base64.URLEncoding.EncodeToString(msg),
	// 	Payload: &gmail.MessagePart{
	// 		Headers: []*gmail.MessagePartHeader{
	// 			{Name: "To", Value: base64.URLEncoding.EncodeToString([]byte(to1.String()))},
	// 			{Name: "To", Value: base64.URLEncoding.EncodeToString([]byte(to2.String()))},
	// 			{Name: "Subject", Value: "testing mail user gmail.Message"},
	// 			{Name: "MIME-version", Value: "1.0"},
	// 			{Name: "Content-Type", Value: "text/plain; charset=\"UTF-8\""},
	// 		},
	// 		Body: &gmail.MessagePartBody{
	// 			Data: base64.URLEncoding.EncodeToString([]byte("testing mail user gmail.Message\n")),
	// 		},
	// 	},
	// }

	// _, err = serv.Users.Messages.Send("me", msg).Do()
	// if err != nil {
	// 	t.Error(err)
	// }
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
