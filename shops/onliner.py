import requests
import collections
import json
from urllib.parse import urlencode
from entities import Product
from shops.shop import Shop

collections.Callable = collections.abc.Callable


class Onliner(Shop):
    title = "Onliner"

    def find(self, query):
        products_json = json.loads(self.send_request(query))["products"]
        products = []

        for product_json in products_json:
            if product_json["prices"]:
                product = Product(self.title)
                product.title = product_json["name"]
                product.link = product_json["html_url"]
                product.price = float(product_json["prices"]["price_min"]["amount"])

                products.append(product)

        return products

    @staticmethod
    def send_request(query):
        url = 'https://catalog.onliner.by/sdapi/catalog.api/search/products'
        data = {
            'query': query,
        }
        query_params = urlencode(data)
        response = requests.get("%s?%s" % (url, query_params))

        return response.text
