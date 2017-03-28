package overmind

import (
	"encoding/json"
)

const (
	TYPE_NOTICE_NEWMESSAGE = 1
	TYPE_CLUB              = 6
	TYPE_CHAT_NEWMESSAGE   = 7
	TYPE_CHAT_MARKREAD     = 8
	TYPE_LIVE              = 9
	TYPE_BROADCAST         = 253
)

type Job struct {
	Mid        int64           `json:"mid"` // 消息id
	Product    string          `json:"product,omitempty"`
	Collection string          `json:"collection,omitempty"`
	Type       int             `json:"type"`
	Gid        uint64          `json:"gid,omitempty"`
	Uids       []string        `json:"uids,omitempty"`
	Condition  string          `json:"condition"`
	Option     json.RawMessage `json:"option"`
	Timestamp  int64           `json:"time"`
}
