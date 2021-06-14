// Copyright (c) 2020 @sakura-rip
// Version 1.1 beta
// LastUpdate 2020/08/28

package lineapigo

import (
	ser "github.com/sakura-rip/lineapigo/talkservice"
)

/*
PROFILE FUNCTION
*/

// GetProfile return the object of Profile
func (cl *LineClient) GetProfile() (*ser.Profile, error) {
	prof, err := cl.talk.GetProfile(cl.ctx, ser.SyncReason_UNKNOWN)
	cl.Data.Profile = prof
	return cl.Data.Profile, err
}

// GenerateUserTicket generate user ticket
func (cl *LineClient) GenerateUserTicket() (string, error) {
	res, err := cl.talk.GenerateUserTicket(cl.ctx, 9223372036854775807, 2147483647)
	return res.ID, err
}

// UpdateProfileName change Profile name
func (cl *LineClient) UpdateProfileName(name string) error {
	req := ser.NewUpdateProfileAttributesRequest()
	content := ser.NewProfileContent()
	content.Value = name
	content.Meta = nil
	req.ProfileAttributes = map[ser.Attr]*ser.ProfileContent{ser.Attr_DISPLAY_NAME: content}
	cl.reqSeq++
	err := cl.talk.UpdateProfileAttributes(cl.ctx, cl.reqSeq, req)
	return err
}

// UpdateProfileBio change status message
func (cl *LineClient) UpdateProfileBio(bio string) error {
	req := ser.NewUpdateProfileAttributesRequest()
	content := ser.NewProfileContent()
	content.Value = bio
	content.Meta = nil
	req.ProfileAttributes = map[ser.Attr]*ser.ProfileContent{ser.Attr_STATUS_MESSAGE: content}
	cl.reqSeq++
	err := cl.talk.UpdateProfileAttributes(cl.ctx, cl.reqSeq, req)
	return err
}

/*
Message FUNCTION
*/

// SendMessage send text message
func (cl *LineClient) SendMessage(toMid, text string) (*ser.Message, error) {
	msg := ser.NewMessage()
	msg.Text = text
	msg.To = toMid
	msg.ContentType = ser.ContentType_NONE
	if IsStrInMap(toMid, cl.reqSeqMessage) {
		cl.reqSeqMessage[toMid]++
	} else {
		cl.reqSeqMessage[toMid] = -1
	}
	return cl.talk.SendMessage(cl.ctx, cl.reqSeqMessage[toMid], msg)
}

// SendContact send contact to toMid
func (cl *LineClient) SendContact(toMid, contactMid string) (*ser.Message, error) {
	msg := ser.NewMessage()
	msg.To = toMid
	msg.ContentType = ser.ContentType_CONTACT
	msg.ContentMetadata = map[string]string{"mid": contactMid}
	if IsStrInMap(toMid, cl.reqSeqMessage) {
		cl.reqSeqMessage[toMid]++
	} else {
		cl.reqSeqMessage[toMid] = 1
	}
	return cl.talk.SendMessage(cl.ctx, cl.reqSeqMessage[toMid], msg)
}

// UnsendMessage unsend message
func (cl *LineClient) UnsendMessage(messageID string) error {
	cl.reqSeq++
	err := cl.talk.UnsendMessage(cl.ctx, cl.reqSeq, messageID)
	return err
}

// SendChatChecked read message
func (cl *LineClient) SendChatChecked(groupID, messageID string) error {
	cl.reqSeq++
	err := cl.talk.SendChatChecked(cl.ctx, cl.reqSeq, groupID, messageID, 0)
	return err
}

/*
CHAT FUNCTION
*/

// GetChats get chats
func (cl *LineClient) GetChats(chatsMids []string) ([]*ser.Chat, error) {
	req := ser.NewGetChatsRequest()
	req.ChatMid = chatsMids
	req.WithInvitees = true
	req.WithMembers = true
	res, err := cl.talk.GetChats(cl.ctx, req)
	return res.Chats, err
}

// AcceptChatInvitation join chat
func (cl *LineClient) AcceptChatInvitation(groupMid string) error {
	req := ser.NewAcceptChatInvitationRequest()
	req.ChatMid = groupMid
	req.ReqSeq = cl.reqSeq
	cl.reqSeq++
	_, err := cl.talk.AcceptChatInvitation(cl.ctx, req)
	return err
}

// AcceptChatInvitationByTicket join chat by ticket
func (cl *LineClient) AcceptChatInvitationByTicket(groupMid, ticketID string) error {
	req := ser.NewAcceptChatInvitationByTicketRequest()
	req.ChatMid = groupMid
	req.ReqSeq = cl.reqSeq
	req.TicketId = ticketID
	cl.reqSeq++
	_, err := cl.talk.AcceptChatInvitationByTicket(cl.ctx, req)
	return err
}

// InviteIntoChat invite friend to chat
func (cl *LineClient) InviteIntoChat(chatMid string, targetMids []string) error {
	req := ser.NewInviteIntoChatRequest()
	req.ReqSeq = cl.reqSeq
	cl.reqSeq++
	req.ChatMid = chatMid
	req.TargetUserMids = targetMids
	_, err := cl.talk.InviteIntoChat(cl.ctx, req)
	return err
}

// ReissueChatTicket get chat invitation ticket
func (cl *LineClient) ReissueChatTicket(chatMid string) (string, error) {
	req := ser.NewReissueChatTicketRequest()
	req.GroupMid = chatMid
	req.ReqSeq = cl.reqSeq
	cl.reqSeq++
	res, err := cl.talk.ReissueChatTicket(cl.ctx, req)
	return res.TicketId, err
}

// RejectChatInvitation reject chat
func (cl *LineClient) RejectChatInvitation(chatMid string) error {
	req := ser.NewRejectChatInvitationRequest()
	req.ReqSeq = cl.reqSeq
	cl.reqSeq++
	req.ChatMid = chatMid
	_, err := cl.talk.RejectChatInvitation(cl.ctx, req)
	return err
}

// GetChat return one chatmid
func (cl *LineClient) GetChat(chatID string) (*ser.Chat, error) {
	req := ser.NewGetChatsRequest()
	req.ChatMid = []string{chatID}
	req.WithInvitees = true
	req.WithMembers = true
	res, err := cl.talk.GetChats(cl.ctx, req)
	if len(res.Chats) > 0 {
		return res.Chats[0], err
	}
	return nil, err
}

// UpdateChatName change chat name
func (cl *LineClient) UpdateChatName(chatID, name string) error {
	chat := &ser.Chat{}
	chat.ChatName = name
	req := ser.NewUpdateChatRequest()
	req.Chat = chat
	req.UpdatedAttribute = ser.ChatAttribute_NAME
	_, err := cl.talk.UpdateChat(cl.ctx, req)
	return err

}

// UpdateChatURL change chat url
func (cl *LineClient) UpdateChatURL(chatID string, typeVar bool) error {
	chat := &ser.Chat{}
	chat.Extra.GroupExtra.PreventedJoinByTicket = typeVar
	req := ser.NewUpdateChatRequest()
	req.Chat = chat
	req.UpdatedAttribute = ser.ChatAttribute_PREVENTED_JOIN_BY_TICKET
	_, err := cl.talk.UpdateChat(cl.ctx, req)
	return err
}

// DeleteOtherFromChat kickout some one from chat
func (cl *LineClient) DeleteOtherFromChat(toMid, targetMid string) error {
	req := ser.NewDeleteOtherFromChatRequest()
	req.ChatMid = toMid
	req.ReqSeq = cl.reqSeq
	req.TargetUserMids = []string{targetMid}
	cl.reqSeq++
	_, err := cl.talk.DeleteOtherFromChat(cl.ctx, req)
	return err
}

// DeleteSelfFromChat leave from chat
func (cl *LineClient) DeleteSelfFromChat(toMid string) error {
	req := ser.NewDeleteSelfFromChatRequest()
	req.ChatMid = toMid
	req.ReqSeq = cl.reqSeq
	cl.reqSeq++
	req.LastSeenMessageDeliveredTime = 0
	req.LastSeenMessageId = ""
	req.LastMessageDeliveredTime = 0
	req.LastMessageId = ""
	_, err := cl.talk.DeleteSelfFromChat(cl.ctx, req)
	return err
}

// GetAllChatMids get all chat mids
func (cl *LineClient) GetAllChatMids() (*ser.GetAllChatMidsResponse, error) {
	req := ser.NewGetAllChatMidsRequest()
	req.WithInvitedChats = true
	req.WithMemberChats = true
	return cl.talk.GetAllChatMids(cl.ctx, req, ser.SyncReason_UNKNOWN)
}

// CancelChatInvitation cancel user
func (cl *LineClient) CancelChatInvitation(groupID, targetMid string) error {
	req := ser.NewCancelChatInvitationRequest()
	req.ChatMid = groupID
	req.ReqSeq = cl.reqSeq
	cl.reqSeq++
	req.TargetUserMids = []string{targetMid}
	_, err := cl.talk.CancelChatInvitation(cl.ctx, req)
	return err
}

// FindChatByTicket find chat by ticket
func (cl *LineClient) FindChatByTicket(ticketID string) (*ser.Chat, error) {
	req := ser.NewFindChatByTicketRequest()
	req.TicketId = ticketID
	res, err := cl.talk.FindChatByTicket(cl.ctx, req)
	return res.Chat, err
}

/*
Settings
*/

// GetSettings get settings
func (cl *LineClient) GetSettings() (*ser.Settings, error) {
	return cl.talk.GetSettings(cl.ctx, ser.SyncReason_UNKNOWN)
}

// UpdateSettings update settings
func (cl *LineClient) UpdateSettings(attr []ser.SettingsAttributes, settings *ser.Settings) error {
	_, err := cl.talk.UpdateSettingsAttributes2(cl.ctx, cl.reqSeq, attr, settings)
	return err
}

// DisableE2ee disable e2ee
func (cl *LineClient) DisableE2ee() error {
	set := ser.NewSettings()
	set.E2eeEnable = false
	err := cl.UpdateSettings([]ser.SettingsAttributes{ser.SettingsAttributes_E2EE_ENABLE}, set)
	return err
}

/*
Contact
*/

// FindAndAddContactsByMid add friend
func (cl *LineClient) FindAndAddContactsByMid(targetMid string) (map[string]*ser.Contact, error) {
	return cl.talk.FindAndAddContactsByMid(cl.ctx, cl.reqSeq, targetMid, ser.MIDType_MID, `{"screen":"homeTab","spec":"native"}`)
}

// GetContacts get contact with list
func (cl *LineClient) GetContacts(targetMid []string) ([]*ser.Contact, error) {
	return cl.talk.GetContacts(cl.ctx, targetMid)
}

// GetContact get contact with list
func (cl *LineClient) GetContact(targetMid string) (*ser.Contact, error) {
	return cl.talk.GetContact(cl.ctx, targetMid)
}

// CoteHan update display name over ridden
func (cl *LineClient) CoteHan(mid, cote string) error {
	return cl.talk.UpdateContactSetting(cl.ctx, cl.reqSeq, mid, ser.ContactFlag_CONTACT_SETTING_DISPLAY_NAME_OVERRIDE, cote)
}

// GetAllContactIds get all list of mid
func (cl *LineClient) GetAllContactIds() ([]string, error) {
	return cl.talk.GetAllContactIds(cl.ctx, ser.SyncReason_UNKNOWN)
}

/*
OTHER
*/

// Noop nothing to do
func (cl *LineClient) Noop() {
	cl.talk.Noop(cl.ctx)
}
