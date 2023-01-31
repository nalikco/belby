import os

from dotenv import load_dotenv
from vk import Vk

load_dotenv()

# FOR TEST
# run_search("macbook", lambda shop, shops: print(shop))

vk = Vk(os.getenv("VK_TOKEN"))

vk.polling(os.getenv("VK_GROUP_ID"))
