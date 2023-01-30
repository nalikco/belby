import requests
import collections
from bs4 import BeautifulSoup

from entities import Product
from shops.shop import Shop

collections.Callable = collections.abc.Callable


class RedStore(Shop):
    title = "RedStore"

    def find(self, query):
        soup = self.send_request(query)

        products = []
        products_html = soup.find_all("div", attrs={"class": "type-product"})
        for product_html in products_html:
            product = Product(self.title)

            product_html_name = product_html.find_next("h5", attrs={"class": "product-name"})
            if not product_html_name:
                raise Exception("product name not found")

            product_html_price = product_html.find_next("span", attrs={"class": "woocommerce-Price-amount amount"})
            if not product_html_price:
                raise Exception("product price not found")

            product_html_link = product_html.find_next("a", attrs={"class": "thumb-hover scale"})
            if not product_html_link:
                raise Exception("product link not found")

            product.title = product_html_name.text
            product.price = float(product_html_price.text.split("\xa0")[0].replace(".", ""))
            product.link = product_html_link.attrs.get("href")

            products.append(product)

        return products

    @staticmethod
    def send_request(query):
        url = 'https://redstore.by/wp-admin/admin-ajax.php'
        data = {
            's': query,
            'post_type': 'product',
            'action': 'sr_ajax_search'
        }
        response = requests.post(url, data=data)
        soup = BeautifulSoup(response.text, 'html.parser')

        return soup
