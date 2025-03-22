from fastapi import FastAPI, HTTPException
import redis
import os
from dotenv import load_dotenv
from openai import OpenAI

load_dotenv(dotenv_path="./.env")
# Initialize Redis
redis_client = redis.Redis(host="localhost", port=6379, decode_responses=True)
print(redis_client)

# Initialize FastAPI
app = FastAPI()

# Get OpenAI API Key
OPENAI_API_KEY = os.getenv("OPENAI_API_KEY")
client = OpenAI(
    api_key=OPENAI_API_KEY,
    base_url="https://api.deepseek.com"
)


@app.post("/chat/")
async def chat(user_id: str, message: str):
    """Receive a user message and generate a response."""
    if not OPENAI_API_KEY:
        raise HTTPException(status_code=500, detail="OpenAI API key is missing")

    # Store user message in Redis
    redis_client.rpush(f"chat:{user_id}", f"You: {message}")

    # Generate AI response
    try:
        response = client.chat.completions.create(
            model="deepseek-chat",
            messages=[{"role": "user", "content": message}],
            max_tokens=1024,
            temperature=0.7,
            stream=False
        )
        bot_reply = response["choices"][0]["message"]["content"]
    except Exception as e:
        print(e)
        raise HTTPException(status_code=500, detail=str(e))

    # Store bot response in Redis
    redis_client.rpush(f"chat:{user_id}", f"Bot: {bot_reply}")

    return {"response": bot_reply}


@app.get("/history/")
async def get_history(user_id: str):
    """Fetch conversation history for a user."""
    history = redis_client.lrange(f"chat:{user_id}", 0, -1)
    return {"history": history}
