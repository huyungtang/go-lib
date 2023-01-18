package slack

import (
	base "github.com/slack-go/slack"
)

// constants & variables ******************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// public functions ***********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// NewAttachment
// ****************************************************************************************************************************************
func NewAttachment(opts ...AttachOption) MsgOption {
	attach := base.Attachment{}
	for _, opt := range opts {
		opt(&attach)
	}

	return func(mo *msgOption) {
		mo.attaches = append(mo.attaches, attach)
	}
}

// AttachTitleOption
// ****************************************************************************************************************************************
func AttachTitleOption(text string) AttachOption {
	return func(a *base.Attachment) {
		a.Title = text
	}
}

// AttachTextOption
// ****************************************************************************************************************************************
func AttachTextOption(text string) AttachOption {
	return func(a *base.Attachment) {
		a.Text = text
	}
}

// AttachColorOption
// ****************************************************************************************************************************************
func AttachColorOption(color string) AttachOption {
	return func(a *base.Attachment) {
		a.Color = color
	}
}

// AttachFieldOption
// ****************************************************************************************************************************************
func AttachFieldOption(title, value string, short bool) AttachOption {
	return func(a *base.Attachment) {
		a.Fields = append(a.Fields, base.AttachmentField{
			Title: title,
			Value: value,
			Short: short,
		})
	}
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// AttachOption
// ****************************************************************************************************************************************
type AttachOption func(*base.Attachment)

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
