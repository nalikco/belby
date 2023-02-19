package handler

import (
	"belby/internal/entities"
	"belby/internal/storages"
	"belby/pkg/shopsparser"
	"belby/pkg/vk"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"math/rand"
)

func HandleVk(storages *storages.Storages, vkMessage vk.Message, vk *vk.Vk) error {
	user, err := storages.Users.GetByVkId(vkMessage.FromID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			user.VkID = vkMessage.FromID

			err := storages.Users.Create(&user)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}

	response, err := vk.Method("messages.send", map[string]string{
		"peer_id":   fmt.Sprintf("%d", vkMessage.PeerID),
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
			"peer_id":   fmt.Sprintf("%d", vkMessage.PeerID),
			"random_id": fmt.Sprintf("%d", rand.Uint32()),
			"message":   "Не удалось выполнить поиск. Попробуйте ещё раз.",
		})
	}

	messageId, ok := responseBody["response"]
	if !ok {
		_, _ = vk.Method("messages.send", map[string]string{
			"peer_id":   fmt.Sprintf("%d", vkMessage.PeerID),
			"random_id": fmt.Sprintf("%d", rand.Uint32()),
			"message":   "Не удалось выполнить поиск. Попробуйте ещё раз.",
		})
	}

	message := entities.Message{
		VkID:    messageId,
		UserID:  user.ID,
		Message: vkMessage.Text,
	}

	if err := storages.Messages.Create(&message); err != nil {
		return err
	}

	parser := shopsparser.NewShopsParser()

	products, err := parser.Find(vkMessage.Text, func(elem, count int) {
		_, _ = vk.Method("messages.edit", map[string]string{
			"peer_id":    fmt.Sprintf("%d", vkMessage.PeerID),
			"random_id":  fmt.Sprintf("%d", rand.Uint32()),
			"message_id": fmt.Sprintf("%d", messageId),
			"message":    fmt.Sprintf("Поиск товара в магазинах (%d из %d)...", elem, count),
		})
	})

	resultMessage := "Для более точного поиска вводите полное наименование продукта.\nВот, что удалось найти по Вашему запросу:\n\n"
	for i, product := range products {
		_ = storages.Products.Create(&entities.Product{
			MessageID: message.ID,
			Shop:      product.ShopTitle,
			Title:     product.Title,
			Price:     product.Price,
			Link:      product.Link,
		})

		productTitle := product.Title
		if len(productTitle) > 40 {
			productTitle = productTitle[0:40] + "..."
		}

		resultMessage += fmt.Sprintf("%d. %s (%.02f руб.): %s\n\n", i+1, productTitle, product.Price, product.Link)
	}

	_, _ = vk.Method("messages.edit", map[string]string{
		"peer_id":    fmt.Sprintf("%d", vkMessage.PeerID),
		"random_id":  fmt.Sprintf("%d", rand.Uint32()),
		"message_id": fmt.Sprintf("%d", messageId),
		"message":    resultMessage,
	})

	return nil
}
