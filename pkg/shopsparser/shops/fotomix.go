package shops

import (
	"belby/pkg/request"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"sort"
	"strconv"
	"strings"
)

type Fotomix struct {
	title string
}

func NewFotomixShop() *Fotomix {
	return &Fotomix{
		title: "Fotomix",
	}
}

func (s *Fotomix) GetTitle() string {
	return s.title
}

func (s *Fotomix) Find(query string) (Product, error) {
	var products []Product

	response, err := s.makeRequest(query)
	if err != nil {
		return Product{}, err
	}

	products, err = s.parse(response)
	if err != nil {
		return Product{}, nil
	}

	sortedProducts, err := s.sort(products)
	if err != nil {
		return Product{}, nil
	}

	if len(sortedProducts) > 0 {
		return sortedProducts[0], nil
	} else {
		return Product{}, nil
	}
}

func (s *Fotomix) sort(products []Product) ([]Product, error) {
	sortedProducts := products

	sort.Slice(sortedProducts[:], func(i, j int) bool {
		return sortedProducts[i].Price < sortedProducts[j].Price
	})

	return sortedProducts, nil
}

func (s *Fotomix) parse(response string) ([]Product, error) {
	var products []Product

	tokenizer := html.NewTokenizer(strings.NewReader(response))

	divFound := false
	divsCountBeforeClose := 0
	priceFound := false

	product := Product{
		ShopTitle: s.GetTitle(),
	}
	for {
		tn := tokenizer.Next()

		switch {
		case tn == html.ErrorToken:
			return products, nil
		case tn == html.TextToken:
			if divFound && priceFound {
				priceFound = false
				token := tokenizer.Token()

				priceString := strings.ReplaceAll(strings.ReplaceAll(token.Data, " р.", ""), " ", "")
				price, err := strconv.ParseFloat(priceString, 64)
				if err != nil {
					continue
				}
				product.Price = price

				products = append(products, product)
				product = Product{
					ShopTitle: s.GetTitle(),
				}
			}
		case tn == html.StartTagToken:
			token := tokenizer.Token()

			if divFound && token.Data == "div" {
				divsCountBeforeClose++
			}

			if divFound && token.Data == "a" && product.Link == "" {
				for _, attribute := range token.Attr {
					if attribute.Key == "href" {
						product.Link = attribute.Val
						break
					}
				}
			}

			if divFound && token.Data == "span" {
				for _, attribute := range token.Attr {
					if attribute.Key == "class" && attribute.Val == "d-inline-block" {
						priceFound = true
					}
				}
			}

			for _, attribute := range token.Attr {
				if attribute.Key == "class" && strings.Contains(attribute.Val, "product-layout product-items product-grid") {
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

func (s *Fotomix) makeRequest(query string) (string, error) {
	body, err := request.SendRequest(request.Request{
		Method: http.MethodGet,
		URL:    "https://fotomix.by/search/",
		Headers: map[string]string{
			"User-Agent": DefaultUserAgent,
		},
		Query: map[string]string{
			"search":       query,
			"sub_category": "true",
			"sort":         "p.price",
			"order":        "ASC",
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
