from urllib.parse import urlencode

import requests
import collections
from bs4 import BeautifulSoup

from entities import Product
from shops.shop import Shop

collections.Callable = collections.abc.Callable


class MobilWorld(Shop):
    title = "MobilWorld"

    def find(self, query):
        soup = self.send_request(query)

        products = []
        products_html = soup.find_all("div", attrs={"class": "item_block"})
        for product_html in products_html:
            item_info = product_html.find_next("div", attrs={"class": "item_info TYPE_1"})
            item_title = item_info.find_next("div", attrs={"class": "item-title"})
            item_link = item_title.find_next("a", attrs={"class": "dark_link"})
            item_title = item_title.find_next("span")
            price_block = product_html.find_next("div", attrs={"class": "price"})
            price = price_block.find_next("span", attrs={"class": "price_value"})

            product = Product(self.title)
            product.title = item_title.text
            product.price = float(price.text.replace(" ", ""))
            product.link = "https://mobilworld.by" + item_link.attrs.get("href")

            products.append(product)

        return products

    @staticmethod
    def send_request(query):
        url = 'https://mobilworld.by/catalog/'
        data = {
            'q': query,
            's': 'Найти'
        }
        response = requests.get("%s?%s" % (url, urlencode(data)), headers={
            "user-agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 "
                          "Safari/537.36"
        })
        soup = BeautifulSoup(response.text, 'html.parser')

        return soup
