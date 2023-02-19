package handler

import (
	"belby/pkg/shopsparser"
	"belby/pkg/vk"
	"encoding/json"
	"fmt"
	"math/rand"
)

func HandleVk(message *vk.Message, vk *vk.Vk) error {
	response, err := vk.Method("messages.send", map[string]string{
		"peer_id":   fmt.Sprintf("%d", message.PeerID),
		"random_id": fmt.Sprintf("%d", rand.Uint32()),
		"message":   "Поиск товара в магазинах...",
	})
	if err != nil {
		fmt.Println(err)
	}

	var responseBody map[string]int64
	if err := json.Unmarshal(response, &responseBody); err != nil {
		fmt.Println(err)
		_, _ = vk.Method("messages.send", map[string]string{
			"peer_id":   fmt.Sprintf("%d", message.PeerID),
			"random_id": fmt.Sprintf("%d", rand.Uint32()),
			"message":   "Не удалось выполнить поиск. Попробуйте ещё раз.",
		})
	}

	messageId, ok := responseBody["response"]
	if !ok {
		_, _ = vk.Method("messages.send", map[string]string{
			"peer_id":   fmt.Sprintf("%d", message.PeerID),
			"random_id": fmt.Sprintf("%d", rand.Uint32()),
			"message":   "Не удалось выполнить поиск. Попробуйте ещё раз.",
		})
	}

	parser := shopsparser.NewShopsParser()

	products, err := parser.Find(message.Text, func(elem, count int) {
		_, _ = vk.Method("messages.edit", map[string]string{
			"peer_id":    fmt.Sprintf("%d", message.PeerID),
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

	_, _ = vk.Method("messages.edit", map[string]string{
		"peer_id":    fmt.Sprintf("%d", message.PeerID),
		"random_id":  fmt.Sprintf("%d", rand.Uint32()),
		"message_id": fmt.Sprintf("%d", messageId),
		"message":    resultMessage,
	})

	return nil
}
