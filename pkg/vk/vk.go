package vk

import (
	"belby/pkg/request"
	"fmt"
	"io"
	"net/http"
)

type Vk struct {
	url         string
	v           string
	groupId     string
	accessToken string
	polling     *Polling
}

func NewVk(groupId, accessToken string) *Vk {
	vk := &Vk{
		url:         "https://api.vk.com/method/",
		v:           "5.131",
		groupId:     groupId,
		accessToken: accessToken,
	}
	vk.polling = NewPolling(vk)

	return vk
}

func (v *Vk) method(method string, query map[string]string) ([]byte, error) {
	query["v"] = v.v
	query["access_token"] = v.accessToken

	body, err := request.SendRequest(request.Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("%s%s", v.url, method),
		Query:  query,
		Callback: func(body io.ReadCloser) ([]byte, error) {
			return io.ReadAll(body)
		},
	})

	return body, err
}

func (v *Vk) Polling() error {
	return v.polling.Run()
}
