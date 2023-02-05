import json
from urllib.parse import urlencode

import requests
import collections

from entities import Product
from shops.shop import Shop

collections.Callable = collections.abc.Callable


class Element5(Shop):
    title = "5element"

    def find(self, query):
        products = []

        try:
            products_json = json.loads(self.send_request(query))["results"]["items"]
            for product_json in products_json:
                product = Product(self.title)
                product.title = product_json["name"]
                product.link = "https://5element.by" + product_json["url"]
                product.price = float(product_json["price"])

                products.append(product)
        except KeyError:
            pass

        return products

    @staticmethod
    def send_request(query):
        url = 'https://api.multisearch.io/'
        data = {
            'query': query,
            'id': '11432',
            'lang': 'ru',
            'categories': '0',
            'fields': 'true',
            'limit': '100',
            'filters': '{}',
            'offset': '0',
            'offer_type': 'product',
            'sort': 'price.asc',
        }
        query_params = urlencode(data)
        response = requests.get("%s?%s" % (url, query_params))

        return response.text
