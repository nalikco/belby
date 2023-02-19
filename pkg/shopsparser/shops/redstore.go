package shops

import (
	"belby/pkg/request"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
)

type RedStore struct {
	title string
}

func NewRedStoreShop() *RedStore {
	return &RedStore{
		title: "RedStore",
	}
}

func (s *RedStore) GetTitle() string {
	return s.title
}

func (s *RedStore) Find(query string) (Product, error) {
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

func (s *RedStore) sort(products []Product) ([]Product, error) {
	sortedProducts := products

	sort.Slice(sortedProducts[:], func(i, j int) bool {
		return sortedProducts[i].Price < sortedProducts[j].Price
	})

	return sortedProducts, nil
}

func (s *RedStore) parse(response string) ([]Product, error) {
	var products []Product

	tokenizer := html.NewTokenizer(strings.NewReader(response))

	divFound := false
	divsCountBeforeClose := 0
	titleFound := false
	priceBlockFound := false
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
			if divFound {
				token := tokenizer.Token()
				if titleFound {
					product.Title = token.Data
					titleFound = false
				}
				if priceFound {
					priceFound = false
					priceString := strings.Replace(token.Data, ".", "", 1)

					price, err := strconv.ParseFloat(strings.TrimSpace(priceString), 64)
					if err != nil {
						continue
					}
					product.Price = price

					products = append(products, product)
					product = Product{
						ShopTitle: s.GetTitle(),
					}
				}
			}
		case tn == html.StartTagToken:
			token := tokenizer.Token()

			if divFound && token.Data == "div" {
				divsCountBeforeClose++
			}

			if divFound && token.Data == "h5" {
				for _, attribute := range token.Attr {
					if attribute.Key == "class" && strings.Contains(attribute.Val, "product-name") {
						titleFound = true
					}
				}
			}

			if divFound && token.Data == "span" {
				for _, attribute := range token.Attr {
					if attribute.Key == "class" && strings.Contains(attribute.Val, "woocommerce-Price-amount") {
						priceBlockFound = true
					}
				}
			}

			if priceBlockFound && token.Data == "bdi" {
				priceBlockFound = false
				priceFound = true
			}

			if titleFound && token.Data == "a" {
				for _, attribute := range token.Attr {
					if attribute.Key == "href" {
						product.Link = attribute.Val
						break
					}
				}
			}

			for _, attribute := range token.Attr {
				if attribute.Key == "class" && strings.Contains(attribute.Val, "type-product") {
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

func (s *RedStore) makeRequest(query string) (string, error) {
	form := url.Values{}
	form.Add("s", query)
	form.Add("post_type", "product")
	form.Add("action", "sr_ajax_search")

	body, err := request.SendRequest(request.Request{
		Method: http.MethodPost,
		URL:    "https://redstore.by/wp-admin/admin-ajax.php",
		Body:   strings.NewReader(form.Encode()),
		Headers: map[string]string{
			"User-Agent":   DefaultUserAgent,
			"Content-Type": "application/x-www-form-urlencoded",
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
