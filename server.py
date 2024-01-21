import socket
import threading
import signal
import sys
import uuid

class ChatServer:
    def __init__(self, host, port):
        self.server_address = (host, port)
        self.clients = {}
        self.running = True

    def create_server_socket(self):
        server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        server_socket.bind(self.server_address)
        server_socket.listen(5)
        return server_socket

    def handle_client(self, client_socket, client_address):
        try:
            # Receive the username from the client
            username = client_socket.recv(1024).decode().strip()

            # Assign a unique session ID to the user
            session_id = str(uuid.uuid4())
            self.clients[session_id] = {'socket': client_socket, 'username': username}

            # Inform the user about the session ID
            welcome_message = f"Welcome, {username}! Your session ID is {session_id}\n"
            client_socket.sendall(welcome_message.encode())

            # Broadcast "username has joined" message to all clients
            join_message = f"{username} has joined."
            self.broadcast_message(join_message, client_address)

            while True:
                data = client_socket.recv(1024)
                if not data:
                    break

                print(f"Received from {username} ({session_id}): {data.decode()}")

                # Check if the message is a command
                if data.decode().strip().lower() == 'shutdown':
                    self.shutdown_server()
                elif data.decode().strip().lower() == 'list':
                    self.list_clients(client_socket)
                else:
                    # Broadcast the message to all clients
                    self.broadcast_message(data, client_address)
        finally:
            # Remove the user from the clients dictionary and broadcast a "username has left" message
            if session_id in self.clients:
                del self.clients[session_id]
                leave_message = f"{username} has left."
                self.broadcast_message(leave_message, client_address)

            client_socket.close()
            print(f"User {username} ({session_id}) has disconnected.")

    def list_clients(self, client_socket):
        active_clients = [client_info['username'] for client_info in self.clients.values()]
        if active_clients:
            client_socket.sendall(f"Active clients: {', '.join(active_clients)}\n".encode())
        else:
            client_socket.sendall("No active clients.\n".encode())

    def broadcast_message(self, message, sender_address):
        try:
            username = self.clients[sender_address]['username']
            formatted_message = f"{username} ({sender_address}): {message}"
        except KeyError:
            # Handle the case when the client has disconnected
            return

        for client_info in self.clients.values():
            try:
                client_socket = client_info['socket']
                client_socket.sendall(formatted_message.encode())
            except:
                pass  # Handle the case when a client has disconnected

    def command_input_thread(self):
        while self.running:
            command = input("Enter a command: ")
            if command.lower() == 'shutdown':
                self.shutdown_server()
            elif command.lower() == 'list':
                self.list_clients_to_console()

    def list_clients_to_console(self):
        active_clients = [client_info['username'] for client_info in self.clients.values()]
        if active_clients:
            print(f"Active clients: {', '.join(active_clients)}")
        else:
            print("No active clients.")

    def run(self):
        signal.signal(signal.SIGINT, self.handle_shutdown_signal)

        server_socket = self.create_server_socket()
        print(f"Server listening on {self.server_address}")

        # Start a separate thread for command input
        command_thread = threading.Thread(target=self.command_input_thread)
        command_thread.start()

        while self.running:
            connection, client_address = server_socket.accept()
            self.clients[client_address] = {'socket': connection}

            client_handler = threading.Thread(target=self.handle_client, args=(connection, client_address))
            client_handler.start()

    def handle_shutdown_signal(self, signum, frame):
        print("Received shutdown signal. Initiating graceful shutdown.")
        self.shutdown_server()

    def shutdown_server(self):
        self.running = False

        # Inform all clients about the server shutdown
        shutdown_message = "Server is shutting down."
        for client_info in self.clients.values():
            try:
                client_socket = client_info['socket']
                client_socket.sendall(shutdown_message.encode())
                client_socket.close()
            except:
                pass  # Ignore errors during shutdown

        # Cleanup and exit
        sys.exit()

if __name__ == "__main__":
    server = ChatServer('localhost', 12345)
    server.run()
