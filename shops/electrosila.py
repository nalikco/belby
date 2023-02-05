import sys
from urllib.parse import urlencode

import requests
import collections
from bs4 import BeautifulSoup

from entities import Product
from shops.shop import Shop

collections.Callable = collections.abc.Callable


class ElectroSila(Shop):
    title = "ElectroSila"

    def find(self, query):
        soup = self.send_request(query)

        products = []
        products_html = soup.find_all("div", attrs={"class": "tov_prew_search"})
        for product_html in products_html:
            result_link = product_html.find_next("a")
            if not result_link:
                continue
            result_link = result_link.attrs.get("href")

            result_name = product_html.find_next("strong")
            if not result_name:
                continue
            result_name = result_name.text

            result_price_block = product_html.find_next("div", attrs={"class": "price"})
            if not result_price_block:
                continue

            result_prices = result_price_block.find_next("div")
            if not result_prices:
                continue

            result_price = result_prices.find_all("b")[:2]
            if not result_price:
                continue

            price_string = ""
            i = 1
            for tag in result_price:
                price_string += tag.text.replace(".", "")
                if i == 1:
                    price_string += "."
                i += 1

            product = Product(self.title)
            product.title = result_name
            product.price = float(price_string)
            product.link = result_link

            products.append(product)

        return products

    @staticmethod
    def send_request(query):
        url = 'https://sila.by/search/%s/sort/6' % query
        response = requests.get(url, headers={
            "user-agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 "
                          "Safari/537.36"
        })
        soup = BeautifulSoup(response.text, 'html.parser')

        return soup
