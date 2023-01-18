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

// Init
// ****************************************************************************************************************************************
func Init(tkn string, opts ...Options) (c Client, err error) {
	s := &client{
		client:  base.New(tkn),
		msgOpts: make([]base.MsgOption, 0),
	}
	for _, opt := range opts {
		opt(s)
	}

	return s, nil
}

// type defineds **************************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************

// client *********************************************************************************************************************************
type client struct {
	client  *base.Client
	msgOpts []base.MsgOption
}

// Client
// ****************************************************************************************************************************************
type Client interface {
	PostMessage(string, ...MsgOption) error
}

// PostMessage
// ****************************************************************************************************************************************
func (o *client) PostMessage(channel string, opts ...MsgOption) (err error) {
	args := make([]base.MsgOption, len(o.msgOpts), len(o.msgOpts)+1)
	copy(args, o.msgOpts)

	msg := new(msgOption)
	for _, opt := range opts {
		opt(msg)
	}
	args = append(args, base.MsgOptionBlocks(msg.blocks...))

	_, _, err = o.client.PostMessage(channel, args...)

	return
}

// private functions **********************************************************************************************************************
// ****************************************************************************************************************************************
// ****************************************************************************************************************************************
