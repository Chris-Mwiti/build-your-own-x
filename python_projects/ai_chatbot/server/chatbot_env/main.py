from fastapi import FastAPI, Request
import uvicorn
import os
from dotenv import load_dotenv

# load the dotenv file
load_dotenv();

# creates a server locally
api = FastAPI();


#creates a decorator function that will proxy any get operation to the server on this handler
@api.get("/test")
async def root():
    return {"Msg": "the api is online"}



#initialize and execute the server
if __name__ == "__main__":
    if os.environ.get("API_ENV") == "development":
       uvicorn.run("main:api", host="127.0.0.1", port=3500,workers=4,reload=True) 
    else:
        raise Exception("running on a non-development server")