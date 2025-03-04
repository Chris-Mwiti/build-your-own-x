import os
from fastapi import APIRouter, FastAPI, WebSocket, Request, HTTPException, WebSocketDisconnect,Depends
from ..sockets.connection import ConnectionManager
from ..utils.utils import get_token
from ..redis.producer import Producer
from ..redis.config import Redis
import uuid

# initialize a new manager to manage the connections
manager = ConnectionManager()

# Initialize a new chat route factory that will contain it specific routes and handlers
chat = APIRouter();

#creation of a new redis connection instance
redis = Redis()

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
async def websocket_endpoint(websocket: WebSocket = WebSocket, token: str = Depends(get_token)):

    #basically this func is used to append all conn instances to a list of conns
    await manager.connect(websocket=websocket)

    # creation of a new redis client instance
    # so for each instance connection redis connection is established
    redis_client = await redis.create_connection()

    #we create a new producer that will create messages that will be queued before entering the consumer
    producer = Producer(redis_client=redis_client) 
    try:
        while True:
            data = await websocket.receive_text()
            print(data)
            stream_data = {}
            stream_data[token] = data

            await producer.add_to_stream(stream_data=stream_data, stream_channel="producer channel")
            await manager.send_personal_message(f"Response: Simulating response from the GPT service", websocket=websocket)
    except WebSocketDisconnect:
        manager.disconnect(websocket=websocket)
