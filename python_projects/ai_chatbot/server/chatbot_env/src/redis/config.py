import os
from dotenv import load_dotenv
import redis

load_dotenv()

class Redis():
    def __init__(self):
        print("initialization of server redis instance")
        self.connection_url = os.environ['REDIS_SERVER_URL']  
    
    async def create_connection(self):
        self.connection = redis.from_url(
            self.connection_url, db=0
        ) 
        print(self.connection)
        return self.connection