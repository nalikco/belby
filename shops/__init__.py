import typing

from entities.product import Product
from shops.electrosila import ElectroSila
from shops.element5 import Element5
from shops.fotomix import Fotomix
from shops.mobilworld import MobilWorld
from shops.onliner import Onliner
from shops.redstore import RedStore
from shops.shopby import ShopBy
from shops.vek21 import Vek21


def run_search(query: str, callback) -> typing.List[Product]:
    products = []
    shops_count = 8

    try:
        callback(1, shops_count)
        red_store = RedStore()
        products += red_store.find(query)[:1]
    except:
        pass

    try:
        callback(2, shops_count)
        onliner_by = Onliner()
        products += onliner_by.find(query)[:1]
    except ValueError:
        pass

    try:
        callback(3, shops_count)
        shop_by = ShopBy()
        products += shop_by.find(query)[:1]
    except ValueError:
        pass

    try:
        callback(4, shops_count)
        mobil_world = MobilWorld()
        products += mobil_world.find(query)[:1]
    except ValueError:
        pass

    try:
        callback(5, shops_count)
        foto_mix = Fotomix()
        products += foto_mix.find(query)[:1]
    except ValueError:
        pass

    try:
        callback(6, shops_count)
        vek_21 = Vek21()
        products += vek_21.find(query)[:1]
    except ValueError:
        pass

    try:
        callback(7, shops_count)
        element_5 = Element5()
        products += element_5.find(query)[:1]
    except ValueError:
        pass

    try:
        callback(8, shops_count)
        electro_sila = ElectroSila()
        products += electro_sila.find(query)[:1]
    except ValueError:
        pass

    products.sort(key=lambda p: p.price)

    return products[:5]
