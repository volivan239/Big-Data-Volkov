package main

import (
	"context"
	"log"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

var wsConnections []*websocket.Conn

func registerWsConnection(connection *websocket.Conn) {
	wsConnections = append(wsConnections, connection)
	go listenWsConnection(connection)
}

func listenWsConnection(connection *websocket.Conn) {
	for {
		var transaction Transaction
		err := wsjson.Read(context.Background(), connection, &transaction)
		if err != nil {
			log.Print(err)
		} else {
			log.Printf("Received transaction #%d from %s (payload = %v)\n", transaction.Id, transaction.Source, transaction.Payload)
			transactionManagerEventChannel <- transaction
		}
	}
}

func sendByWsConnection(connection *websocket.Conn, transaction Transaction) {
	err := wsjson.Write(context.Background(), connection, transaction)
	if err != nil {
		log.Print(err)
	}
}
