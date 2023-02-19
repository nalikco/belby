package vk

import (
	"belby/pkg/request"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"reflect"
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
		timeOut: 5,
		ts:      0,
	}
}

func (p *Polling) getServer() (getServerResponse, error) {
	body, err := p.vk.Method("groups.getLongPollServer", map[string]string{
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

func (p *Polling) Run(callback func(updates []interface{}, vk *Vk)) error {
	serverResponse, err := p.getServer()
	if err != nil {
		return err
	}

	ts, _ := strconv.ParseInt(serverResponse.Response.Ts, 10, 64)

	for {
		response, _ := p.requestToServer(serverResponse, ts)

		ts = response.Ts

		var updates []interface{}

		for _, update := range response.Updates {
			if update.Type == "message_new" {
				message, err := p.ValidateUpdate(update)
				if err != nil {
					fmt.Println(err)
					continue
				}

				updates = append(updates, message)
			}
		}

		callback(updates, p.vk)
	}
}

func (p *Polling) ValidateUpdate(update longPollUpdate) (Message, error) {
	message := Message{}
	messageJson, ok := update.Object["message"]

	mv := reflect.ValueOf(&message).Elem()

	if ok {
		mjv := reflect.ValueOf(messageJson)
		if mjv.Kind() == reflect.Map {
			for _, key := range mjv.MapKeys() {
				for i := 0; i < mv.NumField(); i++ {
					if mv.Type().Field(i).Tag.Get("json") == fmt.Sprintf("%v", key.Interface()) {
						value := mjv.MapIndex(key)

						if reflect.TypeOf(value.Interface()) == reflect.TypeOf(float64(0)) {
							floatValue, _ := strconv.ParseFloat(fmt.Sprintf("%v", value.Interface()), 64)
							reflect.ValueOf(&message).Elem().Field(i).SetInt(int64(floatValue))
						}

						if reflect.TypeOf(value.Interface()) == reflect.TypeOf(true) {
							boolValue, _ := strconv.ParseBool(fmt.Sprintf("%v", value.Interface()))
							reflect.ValueOf(&message).Elem().Field(i).SetBool(boolValue)
						}

						if reflect.TypeOf(value.Interface()) == reflect.TypeOf("") {
							reflect.ValueOf(&message).Elem().Field(i).SetString(fmt.Sprintf("%v", value.Interface()))
						}

						if reflect.TypeOf(value.Interface()) == reflect.TypeOf(map[string]interface{}{}) {
							reflect.ValueOf(&message).Elem().Field(i).Set(value)
						}
					}
				}
			}
		}
	}

	return message, nil
}
