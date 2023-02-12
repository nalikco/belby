package shopsparser

import (
	"belby/internal/entities"
	"belby/pkg/shopsparser/shops"
	"sort"
)

type ShopsParser struct {
	ShopsList []shops.Shop
}

func NewShopsParser() ShopsParser {
	p := ShopsParser{}
	p.initializeShops()

	return p
}

func (p *ShopsParser) initializeShops() {
	var list []shops.Shop

	list = append(list, shops.NewOnlinerShop())
	list = append(list, shops.NewElectroSilaShop())
	list = append(list, shops.NewElement5Shop())
	list = append(list, shops.NewFotomixShop())
	list = append(list, shops.NewMobilWorldShop())
	list = append(list, shops.NewRedStoreShop())
	list = append(list, shops.NewShopByShop())
	list = append(list, shops.NewVek21Shop())

	p.ShopsList = list
}

func (p *ShopsParser) sort(products []entities.Product) []entities.Product {
	var sortedProducts []entities.Product

	for _, product := range products {
		if product.Title != "" && product.Link != "" && product.Price != 0 {
			sortedProducts = append(sortedProducts, product)
		}
	}

	sort.Slice(sortedProducts[:], func(i, j int) bool {
		return sortedProducts[i].Price < sortedProducts[j].Price
	})

	return sortedProducts
}

func (p *ShopsParser) Find(query string) ([]entities.Product, error) {
	var products []entities.Product

	for _, shop := range p.ShopsList {
		productsFromShop, err := shop.Find(query)
		if err != nil {
			return products, nil
		}

		products = append(products, productsFromShop)
	}

	return p.sort(products), nil
}
