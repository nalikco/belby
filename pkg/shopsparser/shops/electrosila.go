package shops

import (
	"belby/internal/entities"
	"belby/pkg/request"
	"fmt"
	"golang.org/x/net/html"
	"golang.org/x/text/encoding/charmap"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type ElectroSila struct {
	title string
}

func NewElectroSilaShop() *ElectroSila {
	return &ElectroSila{
		title: "ЭлектроСила",
	}
}

func (s *ElectroSila) GetTitle() string {
	return s.title
}

func (s *ElectroSila) Find(query string) (entities.Product, error) {
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

func (s *ElectroSila) sort(products []entities.Product) ([]entities.Product, error) {
	sortedProducts := products

	sort.Slice(sortedProducts[:], func(i, j int) bool {
		return sortedProducts[i].Price < sortedProducts[j].Price
	})

	return sortedProducts, nil
}

func (s *ElectroSila) parse(response string) ([]entities.Product, error) {
	var products []entities.Product

	tokenizer := html.NewTokenizer(strings.NewReader(response))

	divFound := false
	divsCountBeforeClose := 0
	titleFound := false
	priceFound := false
	priceDecimalFound := 0

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
				if priceDecimalFound != 0 {
					priceDecimalFound++
				}
				if priceDecimalFound == 3 {
					decimals, err := strconv.ParseFloat(token.Data, 64)
					if err != nil {
						continue
					}
					product.Price += decimals

					priceDecimalFound = 0

					products = append(products, product)
					product = entities.Product{
						ShopTitle: s.GetTitle(),
					}
				}
				if priceFound {
					priceFound = false

					price, err := strconv.ParseFloat(token.Data, 64)
					if err != nil {
						continue
					}
					product.Price = price
					priceDecimalFound = 1
				}
			}
		case tn == html.StartTagToken:
			token := tokenizer.Token()

			if divFound && token.Data == "div" {
				for _, attribute := range token.Attr {
					if attribute.Key == "class" && attribute.Val == "price" {
						priceFound = true
					}
				}
				divsCountBeforeClose++
			}

			if divFound && token.Data == "strong" {
				titleFound = true
			}

			if divFound && token.Data == "a" && product.Link == "" {
				for _, attribute := range token.Attr {
					if attribute.Key == "href" {
						product.Link = attribute.Val
						break
					}
				}
			}

			for _, attribute := range token.Attr {
				if attribute.Key == "class" && attribute.Val == "tov_prew_search" {
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

func (s *ElectroSila) makeRequest(query string) (string, error) {
	body, err := request.SendRequest(request.Request{
		Method: http.MethodGet,
		URL:    fmt.Sprintf("https://sila.by/search/%s/sort/6", query),
		Headers: map[string]string{
			"User-Agent": DefaultUserAgent,
		},
		Callback: func(body io.ReadCloser) ([]byte, error) {
			reader := charmap.Windows1251.NewDecoder().Reader(io.Reader(body))
			return io.ReadAll(reader)
		},
	})
	if err != nil {
		return "", err
	}

	return string(body), nil
}
