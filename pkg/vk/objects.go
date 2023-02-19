package vk

type Update struct {
	Type   string
	Object interface{}
}

type Message struct {
	Date                  int64                  `json:"date"`
	FromID                int64                  `json:"from_id"`
	ID                    int64                  `json:"id"`
	Out                   int64                  `json:"out"`
	Attachments           map[string]interface{} `json:"attachments"`
	ConversationMessageId int64                  `json:"conversation_message_id"`
	FwdMessages           map[string]interface{} `json:"fwd_messages"`
	Important             bool                   `json:"important"`
	IsHidden              bool                   `json:"is_hidden"`
	PeerID                int64                  `json:"peer_id"`
	RandomID              int64                  `json:"random_id"`
	Text                  string                 `json:"text"`
}
