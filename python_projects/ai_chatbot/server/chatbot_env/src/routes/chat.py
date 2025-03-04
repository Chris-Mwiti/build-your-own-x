import os
from fastapi import APIRouter, FastAPI, WebSocket, Request, HTTPException
import uuid

# Initialize a new chat route factory that will contain it specific routes and handlers
chat = APIRouter();

# used to create tokens that will be used to actually authorize requests to the chat route
@chat.post("/token")
async def token_generator(name: str, request:Request):
    if name == "":
        raise HTTPException(status_code=400, detail={
            "loc": "name",
            "msg": "Enter a valid name"
        })
    
    # create a new token for which the user is going to use
    token = str(uuid.uuid4())

    data = {
        "name": name,
        "token": token
    }
    return data

@chat.post("/refresh_token")
async def refresh_generator(request:Request):
    return None

# a chat web socket initialiization
@chat.websocket("/chat")
async def websocket_endpoint(websocket: WebSocket = WebSocket):
    return None