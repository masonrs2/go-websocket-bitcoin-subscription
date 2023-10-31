package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type Server struct {
	conns map[*websocket.Conn]bool
}

type Response struct {
	MarketData struct {
		CurrentPrice struct {
			Usd float64 `json:"usd"`
		} `json:"current_price"`
	} `json:"market_data"`
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("New connection from client: ", ws.RemoteAddr())
	
	s.conns[ws] = true

	s.readAndBroadcast(ws)
}


func (s *Server) readAndBroadcast(ws *websocket.Conn) {
	buffer := make([]byte, 1024)
	for {
		n, err := ws.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("Error reading from client: ", err.Error())
			continue
		}
		
		msg := buffer[:n]
		fmt.Println(string(msg))
		
		s.broadcast(msg)
		
	}
}


func (s *Server) broadcast(msg []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(msg); err != nil {
				fmt.Println("Error broadcasting message: ", err.Error())
			}
		}(ws)
	}
}

func (s *Server) handleWSSubscription(ws *websocket.Conn) {
	fmt.Println("New subscription from client: ", ws.RemoteAddr())

	for {
		payload := fmt.Sprintf("Bitcoin current price: %f\n", getPriceOfBitcoin())
		ws.Write([]byte(payload))
		time.Sleep(2 * time.Second)
	}

	s.readAndBroadcast(ws)
}

func getPriceOfBitcoin() float64 {
	res, err := http.Get("https://api.coingecko.com/api/v3/coins/bitcoin")
	if err != nil {
		fmt.Println("Error getting price of bitcoin: ", err.Error())
		return 0.0
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println("Error fetching price:", err)
		return 0.0 
	}
	var response Response
	json.Unmarshal(body, &response)
	fmt.Println(response.MarketData.CurrentPrice.Usd)

	return response.MarketData.CurrentPrice.Usd
}

func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.Handle("/subscription", websocket.Handler(server.handleWSSubscription))
	http.ListenAndServe(":3000", nil)
}