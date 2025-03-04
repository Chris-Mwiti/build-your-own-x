# A class that will be used to handle websockets connections
from typing import List
from fastapi import WebSocket

class ConnectionManager:
    def __init__(self):
        #create a list that will contain all the active websockets connections
        self.active_connections: List[WebSocket] = []

    async def connect(self, websocket: WebSocket):
        # initialize a new connection
        await websocket.accept()
        # append new connections to the active_connections
        self.active_connections.append(websocket)
    
    def disconnect(self, websocket: WebSocket):
        self.active_connections.remove(websocket)

    async def send_personal_message(self, message: str, websocket: WebSocket):
        await websocket.send_text(message)