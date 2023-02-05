import time
import psycopg2
from contextlib import closing

from database import Database, db


class Message:
    db: Database = db

    id = 0
    user_id = 0
    message = ""
    created_at = 0

    @staticmethod
    def create(user_id: int, message: str):
        with closing(psycopg2.connect(Message.db.dns)) as conn:
            with conn.cursor() as cursor:
                cursor.execute("INSERT INTO messages(user_id, message, created_at) VALUES (%s, %s, %s) RETURNING id",
                               (user_id, message, time.strftime("%Y-%m-%d %H:%M:%S+00")))
                conn.commit()
                return Message.get_by_id(cursor.fetchone()[0])

    @staticmethod
    def get_by_id(message_id: int):
        with closing(psycopg2.connect(Message.db.dns)) as conn:
            with conn.cursor() as cursor:
                cursor.execute('SELECT * FROM messages WHERE id = %s LIMIT 1', (message_id, ))
                return Message.fetch_to_object(cursor.fetchone())

    @staticmethod
    def fetch_to_object(result):
        message = Message()
        message.id = result[0]
        message.user_id = result[1]
        message.message = result[2]
        message.created_at = result[3]

        return message
