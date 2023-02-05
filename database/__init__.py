import os
import time
import psycopg2
from contextlib import closing
from dotenv import load_dotenv

from entities.message import fetch_message_to_object
from entities.user import fetch_user_to_object


class Database:
    dns = ""

    def __init__(self, host: str, username: str, password: str, database: str):
        self.dns = "postgres://%s:%s@%s:5432/%s" % (username, password, host, database)

    def get_or_create_user_by_vk_id(self, vk_id: int):
        with closing(psycopg2.connect(self.dns)) as conn:
            with conn.cursor() as cursor:
                cursor.execute('SELECT * FROM users WHERE vk_id = %s LIMIT 1', (vk_id, ))

                result = cursor.fetchone()
                if result:
                    return fetch_user_to_object(result)

                cursor.execute("INSERT INTO users(vk_id, created_at) VALUES (%s, %s) RETURNING id", (vk_id, time.strftime("%Y-%m-%d %H:%M:%S+00")))
                conn.commit()
                return self.get_user_by_id(cursor.fetchone()[0])

    def get_user_by_id(self, user_id: int):
        with closing(psycopg2.connect(self.dns)) as conn:
            with conn.cursor() as cursor:
                cursor.execute('SELECT * FROM users WHERE id = %s LIMIT 1', (user_id, ))
                return fetch_user_to_object(cursor.fetchone())

    def create_message(self, user_id: int, message: str):
        with closing(psycopg2.connect(self.dns)) as conn:
            with conn.cursor() as cursor:
                cursor.execute("INSERT INTO messages(user_id, message, created_at) VALUES (%s, %s, %s) RETURNING id",
                               (user_id, message, time.strftime("%Y-%m-%d %H:%M:%S+00")))
                conn.commit()
                return self.get_message_by_id(cursor.fetchone()[0])

    def get_message_by_id(self, message_id: int):
        with closing(psycopg2.connect(self.dns)) as conn:
            with conn.cursor() as cursor:
                cursor.execute('SELECT * FROM messages WHERE id = %s LIMIT 1', (message_id, ))
                return fetch_message_to_object(cursor.fetchone())


load_dotenv()
db = Database(
    os.getenv("DB_HOST"),
    os.getenv("DB_USER"),
    os.getenv("DB_PASS"),
    os.getenv("DB_NAME"),
)
