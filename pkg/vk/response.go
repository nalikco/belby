package vk

type getServerResponse struct {
	Response getServerResponseBody `json:"response"`
}

type getServerResponseBody struct {
	Key    string `json:"key"`
	Server string `json:"server"`
	Ts     string `json:"ts"`
}

type longPollResponse struct {
	Ts      int64             `json:"ts"`
	Updates []longPollUpdates `json:"updates"`
}

type longPollUpdates struct {
	GroupId int64                  `json:"group_id"`
	Type    string                 `json:"type"`
	EventID string                 `json:"event_id"`
	V       string                 `json:"v"`
	Object  map[string]interface{} `json:"object"`
}
