package shops

import (
	"belby/pkg/request"
	"encoding/json"
	"io"
	"net/http"
	"sort"
	"strconv"
)

type Onliner struct {
	title string
}

type onlinerProductPricesPriceMin struct {
	Amount string `json:"amount"`
}

type onlinerProductPrices struct {
	PriceMin onlinerProductPricesPriceMin `json:"price_min"`
}

type onlinerProduct struct {
	Name    string               `json:"name"`
	HtmlUrl string               `json:"html_url"`
	Prices  onlinerProductPrices `json:"prices"`
}

type onlinerResponseBody struct {
	Products []onlinerProduct `json:"products"`
}

func NewOnlinerShop() *Onliner {
	return &Onliner{
		title: "Onliner",
	}
}

func (s *Onliner) GetTitle() string {
	return s.title
}

func (s *Onliner) Find(query string) (Product, error) {
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

func (s *Onliner) sort(products []Product) ([]Product, error) {
	sortedProducts := products

	sort.Slice(sortedProducts[:], func(i, j int) bool {
		return sortedProducts[i].Price < sortedProducts[j].Price
	})

	return sortedProducts, nil
}

func (s *Onliner) parse(response string) ([]Product, error) {
	var products []Product

	var responseBody onlinerResponseBody
	err := json.Unmarshal([]byte(response), &responseBody)
	if err != nil {
		return products, err
	}

	for _, responseProduct := range responseBody.Products {
		price, err := strconv.ParseFloat(responseProduct.Prices.PriceMin.Amount, 64)
		if err != nil {
			continue
		}

		product := Product{
			ShopTitle: s.GetTitle(),
			Title:     responseProduct.Name,
			Price:     price,
			Link:      responseProduct.HtmlUrl,
		}

		products = append(products, product)
	}

	return products, nil
}

func (s *Onliner) makeRequest(query string) (string, error) {
	body, err := request.SendRequest(request.Request{
		Method: http.MethodGet,
		URL:    "https://catalog.onliner.by/sdapi/catalog.api/search/products",
		Query: map[string]string{
			"query": query,
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
