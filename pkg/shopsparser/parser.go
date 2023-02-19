package shopsparser

import (
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

func (p *ShopsParser) sort(products []shops.Product) []shops.Product {
	var sortedProducts []shops.Product

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

func (p *ShopsParser) Find(query string, callback func(elem, count int)) ([]shops.Product, error) {
	var products []shops.Product

	for i, shop := range p.ShopsList {
		callback(i+1, len(p.ShopsList))
		productsFromShop, err := shop.Find(query)
		if err != nil {
			return products, nil
		}

		products = append(products, productsFromShop)
	}

	sortedProducts := p.sort(products)

	if len(sortedProducts) > 5 {
		return sortedProducts[:5], nil
	} else {
		return sortedProducts, nil
	}
}
