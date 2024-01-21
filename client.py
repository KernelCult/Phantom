import socket
import threading

class ChatClient:
    def __init__(self, server_address):
        self.server_address = server_address

    def run(self):
        client_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        client_socket.connect(self.server_address)

        # Get the username from the user
        username = input("Enter your username: ")
        client_socket.sendall(username.encode())

        # Receive and print the welcome message from the server
        welcome_message = client_socket.recv(1024).decode()
        print(welcome_message)

        # Start a separate thread to continuously receive messages from the server
        receive_thread = threading.Thread(target=self.receive_messages, args=(client_socket,))
        receive_thread.start()

        # Main thread handles sending messages to the server
        while True:
            message = input("Enter your message or command: ")
            client_socket.sendall(message.encode())

    def receive_messages(self, client_socket):
        while True:
            data = client_socket.recv(1024)
            if not data:
                break
            print(data.decode())

if __name__ == "__main__":
    client = ChatClient(('localhost', 12345))
    client.run()
