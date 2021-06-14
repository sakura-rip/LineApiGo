// Copyright (c) 2020 @sakura-rip
// Version 1.1 beta
// LastUpdate 2020/08/28

package lineapigo

import (
	"encoding/json"

	ser "github.com/sakura-rip/lineapigo/talkservice"
)

// IsStrInMap check string in map
func IsStrInMap(str string, dic map[string]int32) bool {
	_, isIn := dic[str]
	return isIn
}

// JSONMention struct of json for parse Mention
type JSONMention struct {
	MENTION struct {
		MENTIONEES []struct {
			S string `json:"S"`
			E string `json:"E"`
			M string `json:"M"`
		} `json:"MENTIONEES"`
	} `json:"MENTION"`
}

// ParseMention get list of mid
func ParseMention(msg ser.Message) []string {
	var mentions []string
	var strut JSONMention
	// change map to byte
	bytes, _ := json.Marshal(msg.ContentMetadata)
	// Parse byte(map) to json struct
	err := json.Unmarshal(bytes, &strut)
	// if parse failed, return No element list
	if err != nil {
		return mentions
	}
	for _, mid := range strut.MENTION.MENTIONEES {
		mentions = append(mentions, mid.M)
	}
	return mentions
}

// MyChat self made struct
type MyChat struct {
	chatName              string
	chatMid               string
	createdTime           int64
	preventedJoinByTicket bool
	invitationTicket      string
	memberMids            []string
	inviteeMids           []string
	creatorMid            string
}

// GetChatUtil get chat util from chat return MyChat
func (cl *LineClient) GetChatUtil(chat *ser.Chat) *MyChat {
	struc := &MyChat{}
	struc.chatMid = chat.ChatMid
	struc.chatName = chat.ChatName
	struc.preventedJoinByTicket = chat.Extra.GroupExtra.PreventedJoinByTicket
	struc.invitationTicket = chat.Extra.GroupExtra.InvitationTicket
	var memberMids []string
	for mid := range chat.Extra.GroupExtra.MemberMids {
		memberMids = append(memberMids, mid)
	}
	struc.memberMids = memberMids
	var inviteeMids []string
	for mid := range chat.Extra.GroupExtra.InviteeMids {
		inviteeMids = append(inviteeMids, mid)
	}
	struc.inviteeMids = inviteeMids
	struc.creatorMid = chat.Extra.GroupExtra.Creator
	return struc
}
