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
	Ts      int64  `json:"ts"`
	Updates string `json:"updates"`
}
