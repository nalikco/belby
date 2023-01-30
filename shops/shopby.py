from urllib.parse import urlencode

import requests
import collections
from bs4 import BeautifulSoup

from entities import Product
from shops.shop import Shop

collections.Callable = collections.abc.Callable


class ShopBy(Shop):
    title = "ShopBy"

    def find(self, query):
        soup = self.send_request(query)

        products = []
        products_html = soup.find_all("div", attrs={"class": "ModelList__ModelBlockRow"})

        for product_html in products_html:
            product_html_name = product_html.find_next("div", attrs={"class": "ModelList__NameBlock"})
            if not product_html_name:
                raise Exception("product name not found")

            product_html_price = product_html.find_next("div", attrs={"class": "ModelList__PriceBlock"})
            if not product_html_price:
                raise Exception("product price not found")

            product_price_parts = product_html_price.text.strip().split("\n")[0].replace(" ", "").split("\xa0")
            product_price = 0

            try:
                if len(product_price_parts) == 1 or len(product_price_parts) == 2:
                    product_price = float(product_price_parts[0].replace(",", "."))
                if len(product_price_parts) == 3:
                    product_price = float(product_price_parts[1].replace(",", "."))
            except ValueError:
                continue

            product_html_link = product_html.find_next("a", attrs={"class": "ModelList__LinkModel"})
            if product_html_link:
                product = Product(self.title)
                product.title = product_html_name.text.strip()
                product.price = product_price
                product.link = "%s%s" % ("https://shop.by", product_html_link.attrs.get("href"))

                products.append(product)

        return products

    @staticmethod
    def send_request(query):
        url = 'https://shop.by/find/'
        data = {
            'findtext': query,
        }
        query_params = urlencode(data)
        response = requests.get("%s?%s" % (url, query_params))
        soup = BeautifulSoup(response.text, 'html.parser')

        return soup
