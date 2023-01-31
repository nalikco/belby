from urllib.parse import urlencode

import requests
import collections
from bs4 import BeautifulSoup

from entities import Product
from shops.shop import Shop

collections.Callable = collections.abc.Callable


class Fotomix(Shop):
    title = "Fotomix"

    def find(self, query):
        soup = self.send_request(query)

        products = []
        products_html = soup.find_all("div", attrs={"class": "product-layout"})

        for product_html in products_html:
            try:
                link = product_html.find_next("a")
                price_text = product_html.find_next("div", attrs={"class": "price"}).find_next("span").text

                product = Product(self.title)
                product.title = link.find_next("img").attrs.get("title")
                product.price = float(price_text.replace(" р.", "").replace(" ", ""))
                product.link = link.attrs.get("href")

                products.append(product)

            except ValueError:
                pass

        return products

    @staticmethod
    def send_request(query):
        url = 'https://fotomix.by/search/'
        data = {
            'search': query,
            'sub_category': 'true',
            'sort': "p.price",
            "order": "ASC"
        }
        response = requests.get("%s?%s" % (url, urlencode(data)), headers={
            "user-agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 "
                          "Safari/537.36"
        })
        soup = BeautifulSoup(response.text, 'html.parser')

        return soup
