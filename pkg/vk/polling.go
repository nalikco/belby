package vk

import (
	"belby/pkg/request"
	"belby/pkg/shopsparser"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
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
			object, err := p.ValidateUpdate(update)
			if err != nil {
				fmt.Println(err)
				continue
			}

			response, err := p.vk.method("messages.send", map[string]string{
				"peer_id":   fmt.Sprintf("%d", object.PeerID),
				"random_id": fmt.Sprintf("%d", rand.Uint32()),
				"message":   "Поиск товара в магазинах...",
			})
			if err != nil {
				fmt.Println(err)
			}

			var responseBody map[string]int64
			if err := json.Unmarshal(response, &responseBody); err != nil {
				fmt.Println(err)
				_, _ = p.vk.method("messages.send", map[string]string{
					"peer_id":   fmt.Sprintf("%d", object.PeerID),
					"random_id": fmt.Sprintf("%d", rand.Uint32()),
					"message":   "Не удалось выполнить поиск. Попробуйте ещё раз.",
				})
			}

			messageId, ok := responseBody["response"]
			if !ok {
				_, _ = p.vk.method("messages.send", map[string]string{
					"peer_id":   fmt.Sprintf("%d", object.PeerID),
					"random_id": fmt.Sprintf("%d", rand.Uint32()),
					"message":   "Не удалось выполнить поиск. Попробуйте ещё раз.",
				})
			}

			parser := shopsparser.NewShopsParser()

			products, err := parser.Find(object.Text, func(elem, count int) {
				_, _ = p.vk.method("messages.edit", map[string]string{
					"peer_id":    fmt.Sprintf("%d", object.PeerID),
					"random_id":  fmt.Sprintf("%d", rand.Uint32()),
					"message_id": fmt.Sprintf("%d", messageId),
					"message":    fmt.Sprintf("Поиск товара в магазинах (%d из %d)...", elem, count),
				})
			})

			resultMessage := "Для более точного поиска вводите полное наименование продукта.\nВот, что удалось найти по Вашему запросу:\n\n"
			for i, product := range products {
				productTitle := product.Title
				if len(productTitle) > 40 {
					productTitle = productTitle[0:40] + "..."
				}

				resultMessage += fmt.Sprintf("%d. %s (%.02f руб.): %s\n\n", i+1, productTitle, product.Price, product.Link)
			}

			_, _ = p.vk.method("messages.edit", map[string]string{
				"peer_id":    fmt.Sprintf("%d", object.PeerID),
				"random_id":  fmt.Sprintf("%d", rand.Uint32()),
				"message_id": fmt.Sprintf("%d", messageId),
				"message":    resultMessage,
			})
		}
	}
}

func (p *Polling) ValidateUpdate(update longPollUpdate) (Message, error) {
	if update.Type == "message_new" {

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

			return message, nil
		}
	}

	return Message{}, nil
}
