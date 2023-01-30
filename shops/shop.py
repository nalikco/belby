import typing
from entities.product import Product


class Shop:
    title: str = ""

    def find(self, query: str) -> typing.List[Product]:
        """Find products by query"""
        pass
