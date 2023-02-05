import random
from urllib.parse import urlencode

import requests

from entities import Message, User
from shops import run_search


class Vk:
    url: str = "https://api.vk.com/method/"
    v: str = "5.131"
    access_token: str = ""
    ts: int = 0
    timeout: int = 25

    def __init__(self, access_token: str):
        self.access_token = access_token

    def method(self, method: str, data=None):
        if data is None:
            data = {}

        params = {
            "access_token": self.access_token,
            "v": self.v,
        }
        url = "%s%s?%s" % (self.url, method, urlencode(params | data))

        return requests.get(url)

    def polling(self, group_id: str):
        response = self.method("groups.getLongPollServer", {"group_id": group_id})

        try:
            response_json = response.json()["response"]
            self.ts = response_json["ts"]
            server = response_json["server"]
            key = response_json["key"]

            while True:
                response = requests.get("%s?act=a_check&key=%s&ts=%s&wait=%d" % (server, key, self.ts, self.timeout))
                response_json = response.json()
                self.ts = response_json["ts"]

                if len(response_json["updates"]) == 0:
                    continue

                user = User.get_or_create_by_vk_id(int(response_json["updates"][0]["object"]["message"]["from_id"]))
                message = Message.create(user.id, response_json["updates"][0]["object"]["message"]["text"])
                print(message.created_at)

                send_message_response = self.method("messages.send", {
                    "peer_id": response_json["updates"][0]["object"]["message"]["from_id"],
                    "random_id": random.randint(10000000000, 99999999999),
                    "message": "Получение данных из магазинов"
                })
                if not send_message_response.json()["response"]:
                    continue

                message_id = send_message_response.json()["response"]
                products = run_search(response_json["updates"][0]["object"]["message"]["text"],
                                      lambda shop, shops: self.method("messages.edit", {
                                          "peer_id": response_json["updates"][0]["object"]["message"]["from_id"],
                                          "random_id": random.randint(10000000000, 99999999999),
                                          "message": "Получение данных из магазинов (%d/%d)" % (shop, shops),
                                          "message_id": message_id
                                      }))

                message = "\n\nЛучшие результаты:\n\n\n"
                i = 1
                for product in products:
                    message += "%d: %s..., %0.2f BYN: %s\n\n" % (
                        i,
                        product.title[:30],
                        product.price,
                        product.link
                    )

                    i += 1

                self.method("messages.edit", {
                    "peer_id": response_json["updates"][0]["object"]["message"]["from_id"],
                    "random_id": random.randint(10000000000, 99999999999),
                    "message": message,
                    "message_id": message_id
                })
        except KeyError:
            print("connection error")
