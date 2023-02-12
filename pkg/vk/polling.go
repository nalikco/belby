package vk

import (
	"belby/pkg/request"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
)

type Polling struct {
	vk      *Vk
	timeOut int
	ts      int
}

func NewPolling(vk *Vk) *Polling {
	return &Polling{
		vk:      vk,
		timeOut: 25,
		ts:      0,
	}
}

func (p *Polling) getServer() (getServerResponse, error) {
	body, err := p.vk.method("groups.getLongPollServer", map[string]string{
		"group_id": p.vk.groupId,
	})
	if err != nil {
		return getServerResponse{}, err
	}

	var responseBody getServerResponse
	err = json.Unmarshal(body, &responseBody)

	return responseBody, err
}

func (p *Polling) requestToServer(server getServerResponse, ts int64) (longPollResponse, error) {
	body, err := request.SendRequest(request.Request{
		Method: http.MethodGet,
		URL: fmt.Sprintf(
			"%s?act=a_check&key=%s&ts=%d&wait=%d",
			server.Response.Server,
			server.Response.Key,
			ts,
			p.timeOut,
		),
		Callback: func(body io.ReadCloser) ([]byte, error) {
			return io.ReadAll(body)
		},
	})
	if err != nil {
		return longPollResponse{}, err
	}

	var responseBody longPollResponse
	err = json.Unmarshal(body, &responseBody)

	return responseBody, err
}

func (p *Polling) Run() error {
	serverResponse, err := p.getServer()
	if err != nil {
		return err
	}

	ts, _ := strconv.ParseInt(serverResponse.Response.Ts, 10, 64)

	for {
		response, _ := p.requestToServer(serverResponse, ts)

		ts = response.Ts

		for _, update := range response.Updates {
			if update.Type == "message_new" {
				fmt.Println(update.Object["message"])
			}
		}
	}
}
