import os

from dotenv import load_dotenv

# FOR TEST
# from shops import run_search
from vk import Vk

load_dotenv()

# FOR TEST
# products = run_search("iphone", lambda shop, shops: print(shop))
# for product in products:
#     print(product)

vk = Vk(os.getenv("VK_TOKEN"))

vk.polling(os.getenv("VK_GROUP_ID"))
