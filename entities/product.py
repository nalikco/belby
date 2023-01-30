class Product:
    title = ""
    price = 0
    shop = ""
    link = ""

    def __init__(self, shop):
        self.shop = shop

    def __str__(self):
        return "%s (%s), %0.2f BYN: %s" % (self.title, self.shop, self.price, self.link)
