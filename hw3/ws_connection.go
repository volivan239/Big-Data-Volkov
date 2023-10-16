package main

import (
	"context"
	"fmt"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"time"
)

var downstreamConnections []*websocket.Conn

func listenUpstreamWsConnection(peer string) {
	for {
		connection, _, err := websocket.Dial(context.Background(), fmt.Sprintf("ws://%s/ws", peer), nil)
		if err != nil {
			log.Print(err)
			time.Sleep(10 * time.Second)
			continue
		}
		log.Printf("Established ws connection with %s", peer)
		for {
			var transaction Transaction
			err := wsjson.Read(context.Background(), connection, &transaction)
			if err != nil {
				log.Print(err)
				connection.Close(400, "")
				log.Print("Closed connection with %s", peer)
				time.Sleep(10 * time.Second)
				break
			} else {
				log.Printf("Received transaction #%d from %s (payload = %v)\n", transaction.Id, transaction.Source, transaction.Payload)
				transactionManagerEventChannel <- transaction
			}
		}
	}
}

func sendByWsConnection(connection *websocket.Conn, transaction Transaction) {
	err := wsjson.Write(context.Background(), connection, transaction)
	if err != nil {
		log.Print(err)
	}
}
