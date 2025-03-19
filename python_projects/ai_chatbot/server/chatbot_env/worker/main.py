from src.redis.config import Redis
import asyncio

async def main():
    # create a new instance of the redis connection manager
    redis = Redis()

    print("connecting to redis instance....")
    redis = await redis.create_connection()
    print("connection instance is complete")
    print(redis)

    #set the key values pairs of the redis instance
    isSet = redis.set("key", "value")

    if isSet == False:
        raise Exception("Unable to set value")




if __name__ == "__main__":
    asyncio.run(main())