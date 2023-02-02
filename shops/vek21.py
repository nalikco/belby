from urllib.parse import urlencode

import requests
import collections
from bs4 import BeautifulSoup

from entities import Product
from shops.shop import Shop

collections.Callable = collections.abc.Callable


class Vek21(Shop):
    title = "21vek"

    def find(self, query):
        soup = self.send_request(query)

        products = []
        products_html = soup.find_all("li", attrs={"class": "result__item"})
        for product_html in products_html:
            result_link = product_html.find_next("a", attrs={"class": "result__link"})
            if not result_link:
                continue
            result_link = result_link.attrs.get("href")

            result_name = product_html.find_next("span", attrs={"class": "result__name"})
            if not result_name:
                continue
            result_name = result_name.text

            result_price = product_html.find_next("span", attrs={"class": "g-item-data"})
            if not result_price:
                continue

            result_price = float(result_price.text.replace(" ", "").replace(",", "."))

            product = Product(self.title)
            product.title = result_name
            product.price = result_price
            product.link = result_link

            products.append(product)

        return products

    @staticmethod
    def send_request(query):
        url = 'https://www.21vek.by/search/'
        data = {
            'sa': '',
            'term': query,
            'order[price]': 'desc'
        }
        response = requests.get("%s?%s" % (url, urlencode(data)), headers={
            "user-agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 "
                          "Safari/537.36"
        })
        soup = BeautifulSoup(response.text, 'html.parser')

        return soup
