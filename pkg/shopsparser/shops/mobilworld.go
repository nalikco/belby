package shops

import (
	"belby/internal/entities"
	"belby/pkg/request"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type MobilWorld struct {
	title string
}

func NewMobilWorldShop() *MobilWorld {
	return &MobilWorld{
		title: "MobilWorld",
	}
}

func (s *MobilWorld) GetTitle() string {
	return s.title
}

func (s *MobilWorld) Find(query string) (entities.Product, error) {
	var products []entities.Product

	response, err := s.makeRequest(query)
	if err != nil {
		return entities.Product{}, err
	}

	products, err = s.parse(query, response)
	if err != nil {
		return entities.Product{}, nil
	}

	sortedProducts, err := s.sort(products)
	if err != nil {
		return entities.Product{}, nil
	}

	if len(sortedProducts) > 0 {
		return sortedProducts[0], nil
	} else {
		return entities.Product{}, nil
	}
}

func (s *MobilWorld) sort(products []entities.Product) ([]entities.Product, error) {
	sortedProducts := products

	sort.Slice(sortedProducts[:], func(i, j int) bool {
		return sortedProducts[i].Price < sortedProducts[j].Price
	})

	return sortedProducts, nil
}

func (s *MobilWorld) parse(query, response string) ([]entities.Product, error) {
	var products []entities.Product

	tokenizer := html.NewTokenizer(strings.NewReader(response))

	divFound := false
	divsCountBeforeClose := 0
	priceFound := false

	product := entities.Product{
		ShopTitle: s.GetTitle(),
	}
	for {
		tn := tokenizer.Next()

		switch {
		case tn == html.ErrorToken:
			return products, nil
		case tn == html.TextToken:
			if divFound {
				if priceFound {
					priceFound = false

					token := tokenizer.Token()

					price, err := strconv.ParseFloat(strings.ReplaceAll(token.Data, " ", ""), 64)
					if err != nil {
						continue
					}
					product.Price = price

					if strings.Contains(product.Title, query) {
						products = append(products, product)
					}
					product = entities.Product{
						ShopTitle: s.GetTitle(),
					}
				}
			}
		case tn == html.StartTagToken:
			token := tokenizer.Token()

			if divFound && token.Data == "div" {
				divsCountBeforeClose++
			}

			if divFound && token.Data == "a" && product.Link == "" {
				linkFound := false
				for _, attribute := range token.Attr {
					if attribute.Key == "class" && attribute.Val == "dark_link" {
						linkFound = true
						break
					}
				}
				if linkFound {
					for _, attribute := range token.Attr {
						if attribute.Key == "href" {
							product.Link = "https://mobilworld.by" + attribute.Val
							break
						}
					}
				}
			}

			if divFound && token.Data == "span" && product.Price == 0 && product.Title != "" {
				for _, attribute := range token.Attr {
					if attribute.Key == "class" && attribute.Val == "price_value" {
						priceFound = true
					}
				}
			}

			for _, attribute := range token.Attr {
				if attribute.Key == "class" && attribute.Val == "catalog_item_wrapp item" {
					divFound = true
					divsCountBeforeClose = 1
				}
			}
		case tn == html.SelfClosingTagToken:
			if divFound {
				token := tokenizer.Token()
				if token.Data == "img" && product.Title == "" {
					for _, attribute := range token.Attr {
						if attribute.Key == "title" {
							product.Title = attribute.Val
							break
						}
					}
				}
			}
		case tn == html.EndTagToken:
			if divFound {
				token := tokenizer.Token()
				if token.Data == "div" {
					divsCountBeforeClose--

					if divsCountBeforeClose == 0 {
						divFound = false
					}
				}
			}
		}
	}
}

func (s *MobilWorld) makeRequest(query string) (string, error) {
	body, err := request.SendRequest(request.Request{
		Method: http.MethodGet,
		URL:    "https://mobilworld.by/catalog/",
		Headers: map[string]string{
			"User-Agent": DefaultUserAgent,
		},
		Query: map[string]string{
			"q":     query,
			"sort":  "PRICE",
			"order": "asc",
			"s":     "Найти",
		},
		Callback: func(body io.ReadCloser) ([]byte, error) {
			return io.ReadAll(body)
		},
	})
	if err != nil {
		return "", err
	}

	return string(body), nil
}
