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

type ShopBy struct {
	title string
}

func NewShopByShop() *ShopBy {
	return &ShopBy{
		title: "Shop.by",
	}
}

func (s *ShopBy) GetTitle() string {
	return s.title
}

func (s *ShopBy) Find(query string) (entities.Product, error) {
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

func (s *ShopBy) sort(products []entities.Product) ([]entities.Product, error) {
	sortedProducts := products

	sort.Slice(sortedProducts[:], func(i, j int) bool {
		return sortedProducts[i].Price < sortedProducts[j].Price
	})

	return sortedProducts, nil
}

func (s *ShopBy) parse(response string) ([]entities.Product, error) {
	var products []entities.Product

	tokenizer := html.NewTokenizer(strings.NewReader(response))

	divFound := false
	divsCountBeforeClose := 0
	titleFound := false
	priceFound := false
	priceSecondFound := false

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

				if priceSecondFound {
					priceSecondFound = false

					price, err := strconv.ParseFloat(
						strings.ReplaceAll(strings.ReplaceAll(token.Data, ",", "."), " ", ""),
						64,
					)
					if err != nil {
						continue
					}

					product.Price = price
					products = append(products, product)
					product = entities.Product{
						ShopTitle: s.GetTitle(),
					}
				}

				if priceFound {
					priceFound = false

					price, err := strconv.ParseFloat(
						strings.ReplaceAll(strings.ReplaceAll(token.Data, ",", "."), " ", ""),
						64,
					)
					if err != nil {
						priceSecondFound = true
						continue
					}

					product.Price = price
					products = append(products, product)
					product = entities.Product{
						ShopTitle: s.GetTitle(),
					}
				}

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

			if divFound && token.Data == "a" {
				linkFound := false
				link := ""
				for _, attribute := range token.Attr {
					if attribute.Key == "itemprop" && attribute.Val == "url" {
						linkFound = true
					}
					if attribute.Key == "href" {
						link = "https://shop.by" + attribute.Val
					}
				}

				if linkFound {
					product.Link = link
				}
			}

			if divFound && token.Data == "span" {
				for _, attribute := range token.Attr {
					if attribute.Key == "itemprop" && attribute.Val == "name" {
						titleFound = true
					}
					if attribute.Key == "class" && attribute.Val == "PriceBlock__PriceValue" {
						priceFound = true
					}
				}
			}

			for _, attribute := range token.Attr {
				if attribute.Key == "class" && attribute.Val == "ModelList__ModelBlockRow" {
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

func (s *ShopBy) makeRequest(query string) (string, error) {
	body, err := request.SendRequest(request.Request{
		Method: http.MethodGet,
		URL:    "https://shop.by/find/",
		Headers: map[string]string{
			"User-Agent": DefaultUserAgent,
		},
		Query: map[string]string{
			"sort":     "price--number",
			"findtext": query,
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
