### Golang Message Server

This message server uses websocket to send messages between clients. Once the client is connected to the message server, it's given an ID(16 chars) that is uses to identify it. The ID is then stored on the clients local storage and uses it until it expires.

Once client A sends a message to client B the message server waits for ACK response from client B to confirm that the message has been delivered. If the ACK response is not received in 30 seconds the message is cached in redis. On the client connecting to message server the messages that were archived are retrived and sent back to the client.

I'll use the messages logged by the message server to explain what happens under the hood:

⇨ http server started on [::]:5000

message-server_1  | 2022/06/10 07:25:47 Size of Connection Pool:  1

message-server_1  | 2022/06/10 07:25:51 Size of Connection Pool:  2

Two clients have been connected, you can tell that from the size of the connection pool. Once a client disconnects its removed from the connection pool.

The client is requesting for a new ID

message-server_1  | Message Received: {ID: Type:1 Body:{Command:newID Data:<nil>}}

Send message request

message-server_1  | Message Received: {ID: Type:1 Body:{Command:sendMessage Data:map[from: text:Test to:Test]}}

On receiving a message, the message server broadcasts the message to all connected clients.

message-server_1  | 2022/06/10 07:25:57 Sending message to all clients in Pool

Ack response received for the message sent

message-server_1  | Message Received: {ID: Type:1 Body:{Command:ack Data:fa62fcdeaa8359aa0f0b503976efcdd1}}

To run the message server

```
cd message-server
docker-compose up --build
```

