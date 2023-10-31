README
Bitcoin Price WebSocket Server

This Go program runs a WebSocket server that broadcasts the current price of Bitcoin to all connected clients. It fetches the latest Bitcoin price from the CoinGecko API every 2 seconds and sends updates to clients via WebSocket connections.
Usage

1. Run the server with go run main.go.
2. Connect to the server from a WebSocket client. For example, in a web browser's JavaScript console, you can use the following code to connect and log incoming messages:

let socket = new WebSocket("ws://localhost:3000/subscription");

socket.onmessage = (event) => { 
  console.log("received from the server: ", event.data) 
};

This will print the current price of Bitcoin every 2 seconds.# go-websocket-subsciption-chat
