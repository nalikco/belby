import typing

from entities.product import Product
from shops.fotomix import Fotomix
from shops.mobilworld import MobilWorld
from shops.onliner import Onliner
from shops.redstore import RedStore
from shops.shopby import ShopBy


def run_search(query: str, callback) -> typing.List[Product]:
    products = []
    shops_count = 5

    callback(1, shops_count)
    red_store = RedStore()
    products += red_store.find(query)[:1]

    callback(2, shops_count)
    onliner_by = Onliner()
    products += onliner_by.find(query)[:1]

    callback(3, shops_count)
    shop_by = ShopBy()
    products += shop_by.find(query)[:1]

    callback(4, shops_count)
    mobil_world = MobilWorld()
    products += mobil_world.find(query)[:1]

    callback(5, shops_count)
    foto_mix = Fotomix()
    products += foto_mix.find(query)[:1]

    products.sort(key=lambda p: p.price)

    return products[:5]
