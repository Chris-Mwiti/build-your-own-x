import os
from dotenv import load_dotenv
import redis

load_dotenv()

class Redis():
    def __init__(self):
        #initialiaze connection
        self.REDIS_URL = os.environ['REDIS_URL']
        self.REDIS_PASSWORD = os.environ['REDIS_PASSWORD']
        self.REDIS_USER = os.environ['REDIS_USER']
        self.connection_url = os.environ['REDIS_SERVER_URL']

    async def create_connection(self):
        self.connection = redis.from_url(self.connection_url,db=0)
        return self.connection