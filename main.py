import os
from dotenv import load_dotenv
from vk import Vk

load_dotenv()

vk = Vk(os.getenv("VK_TOKEN"))
vk.polling(os.getenv("VK_GROUP_ID"))
