import os
from dotenv import load_dotenv


class Database:
    dns = ""

    def __init__(self, host: str, username: str, password: str, database: str):
        self.dns = "postgres://%s:%s@%s:5432/%s" % (username, password, host, database)


load_dotenv()
db = Database(
    os.getenv("DB_HOST"),
    os.getenv("DB_USER"),
    os.getenv("DB_PASS"),
    os.getenv("DB_NAME"),
)
