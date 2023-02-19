package shops

import (
	"belby/pkg/request"
	"encoding/json"
	"io"
	"net/http"
	"sort"
	"strconv"
)

type Element5 struct {
	title string
}

type element5Item struct {
	Name  string `json:"name"`
	Url   string `json:"url"`
	Price string `json:"price"`
}

type element5Results struct {
	Items []element5Item `json:"items"`
}

type element5ResponseBody struct {
	Results element5Results `json:"results"`
}

func NewElement5Shop() *Element5 {
	return &Element5{
		title: "5 элемент",
	}
}

func (s *Element5) GetTitle() string {
	return s.title
}

func (s *Element5) Find(query string) (Product, error) {
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

func (s *Element5) sort(products []Product) ([]Product, error) {
	sortedProducts := products

	sort.Slice(sortedProducts[:], func(i, j int) bool {
		return sortedProducts[i].Price < sortedProducts[j].Price
	})

	return sortedProducts, nil
}

func (s *Element5) parse(response string) ([]Product, error) {
	var products []Product

	var responseBody element5ResponseBody
	err := json.Unmarshal([]byte(response), &responseBody)
	if err != nil {
		return products, err
	}

	for _, item := range responseBody.Results.Items {
		price, err := strconv.ParseFloat(item.Price, 64)
		if err != nil {
			continue
		}

		product := Product{
			ShopTitle: s.GetTitle(),
			Title:     item.Name,
			Price:     price,
			Link:      "https://5element.by" + item.Url,
		}

		products = append(products, product)
	}

	return products, nil
}

func (s *Element5) makeRequest(query string) (string, error) {
	body, err := request.SendRequest(request.Request{
		Method: http.MethodGet,
		URL:    "https://api.multisearch.io/",
		Headers: map[string]string{
			"User-Agent": DefaultUserAgent,
		},
		Query: map[string]string{
			"query":      query,
			"id":         "11432",
			"lang":       "ru",
			"categories": "0",
			"fields":     "true",
			"limit":      "100",
			"filters":    "{}",
			"offset":     "0",
			"offer_type": "product",
			"sort":       "price.asc",
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
