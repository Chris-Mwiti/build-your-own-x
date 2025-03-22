import requests

BASE_URL = "http://127.0.0.1:8000"

def chat_cli():
    user_id = input("Enter your user ID: ")
    print("Chatbot started! Type 'exit' to quit.")

    while True:
        user_message = input("You: ")
        if user_message.lower() == "exit":
            break

        response = requests.post(f"{BASE_URL}/chat/?user_id={user_id}&message={user_message}")
        if response.status_code == 200:
            print("Bot:", response.json()["response"])
        else:
            print("Error:", response.json()["detail"])

if __name__ == "__main__":
    chat_cli()
