import time
import psycopg2
from contextlib import closing
from database import Database, db


class User:
    db: Database = db

    id = 0
    vk_id = 0
    created_at = 0

    @staticmethod
    def get_by_id(user_id: int):
        with closing(psycopg2.connect(User.db.dns)) as conn:
            with conn.cursor() as cursor:
                cursor.execute('SELECT * FROM users WHERE id = %s LIMIT 1', (user_id, ))
                return User.fetch_to_object(cursor.fetchone())

    @staticmethod
    def get_or_create_by_vk_id(vk_id: int):
        with closing(psycopg2.connect(User.db.dns)) as conn:
            with conn.cursor() as cursor:
                cursor.execute('SELECT * FROM users WHERE vk_id = %s LIMIT 1', (vk_id, ))

                result = cursor.fetchone()
                if result:
                    return User.fetch_to_object(result)

                cursor.execute("INSERT INTO users(vk_id, created_at) VALUES (%s, %s) RETURNING id", (vk_id, time.strftime("%Y-%m-%d %H:%M:%S+00")))
                conn.commit()
                return User.get_by_id(cursor.fetchone()[0])

    @staticmethod
    def fetch_to_object(result):
        user = User()
        user.id = result[0]
        user.vk_id = result[1]
        user.created_at = result[2]

        return user
