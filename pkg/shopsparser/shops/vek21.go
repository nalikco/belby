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

type Vek21 struct {
	title string
}

func NewVek21Shop() *Vek21 {
	return &Vek21{
		title: "21 век",
	}
}

func (s *Vek21) GetTitle() string {
	return s.title
}

func (s *Vek21) Find(query string) (entities.Product, error) {
	var products []entities.Product

	response, err := s.makeRequest(query)
	if err != nil {
		return entities.Product{}, err
	}

	products, err = s.parse(response)
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

func (s *Vek21) sort(products []entities.Product) ([]entities.Product, error) {
	sortedProducts := products

	sort.Slice(sortedProducts[:], func(i, j int) bool {
		return sortedProducts[i].Price < sortedProducts[j].Price
	})

	return sortedProducts, nil
}

func (s *Vek21) parse(response string) ([]entities.Product, error) {
	var products []entities.Product

	tokenizer := html.NewTokenizer(strings.NewReader(response))

	divFound := false
	divsCountBeforeClose := 0
	titleFound := false

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
				token := tokenizer.Token()

				if titleFound {
					product.Title = token.Data
					titleFound = false
				}
			}
		case tn == html.StartTagToken:
			token := tokenizer.Token()

			if divFound && token.Data == "div" {
				divsCountBeforeClose++
			}

			if divFound && token.Data == "a" && product.Link == "" {
				linkFound := true
				link := ""
				for _, attribute := range token.Attr {
					if attribute.Key == "data-ga_action" && strings.Contains(attribute.Val, "GoToItem") {
						linkFound = true
					}
					if attribute.Key == "href" {
						link = attribute.Val
					}
				}

				if linkFound {
					product.Link = link
				}
			}

			if divFound && token.Data == "span" {
				priceFound := false
				priceString := ""

				for _, attribute := range token.Attr {
					if attribute.Key == "class" && strings.Contains(attribute.Val, "result__name") {
						titleFound = true
					}
					if attribute.Key == "class" && strings.Contains(attribute.Val, "g-item-data j-item-data") {
						priceFound = true
					}
					if attribute.Key == "data-price" {
						priceString = attribute.Val
					}
				}

				if priceFound {
					price, err := strconv.ParseFloat(priceString, 64)
					if err != nil {
						continue
					}
					product.Price = price

					if product.Title != "" && product.Link != "" {
						products = append(products, product)
					}

					product = entities.Product{
						ShopTitle: s.GetTitle(),
					}
				}
			}

			for _, attribute := range token.Attr {
				if attribute.Key == "class" && strings.Contains(attribute.Val, "result__item") {
					divFound = true
					divsCountBeforeClose = 1
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

func (s *Vek21) makeRequest(query string) (string, error) {
	body, err := request.SendRequest(request.Request{
		Method: http.MethodGet,
		URL:    "https://www.21vek.by/search/",
		Headers: map[string]string{
			"User-Agent": DefaultUserAgent,
		},
		Query: map[string]string{
			"sa":           "",
			"term":         query,
			"order[price]": "desc",
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
