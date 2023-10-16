package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"os"
)

var transactionCounter uint64

func getHandler(writer http.ResponseWriter, _ *http.Request) {
	writer.WriteHeader(http.StatusOK)
	writer.Write([]byte(snapshot))
}

func replaceHandler(writer http.ResponseWriter, request *http.Request) {
	buffer, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	transactionCounter += 1
	transactionManagerEventChannel <- Transaction{
		Source:  Source,
		Id:      transactionCounter,
		Payload: string(buffer),
	}

	writer.WriteHeader(http.StatusOK)
}

func testHandler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	index, _ := os.ReadFile("index.html")
	writer.Write(index)
}

func vclockHandler(writer http.ResponseWriter, request *http.Request) {
	writer.WriteHeader(http.StatusOK)
	vclockJson, _ := json.Marshal(vclock)
	writer.Header().Add("content-type", "application/json")
	writer.Write(vclockJson)
}

func wsHandler(writer http.ResponseWriter, request *http.Request) {
	connection, err := websocket.Accept(writer, request, &websocket.AcceptOptions{
		InsecureSkipVerify: true,
		OriginPatterns:     []string{"*"},
	})

	if err != nil {
		log.Print(err)
	}

	for _, transaction := range Journal {
		wsjson.Write(request.Context(), connection, transaction)
	}

	log.Printf("New ws connection with %s", request.RemoteAddr)
	downstreamConnections = append(downstreamConnections, connection)
}

var Source string
var peers []string

func main() {
	addr := os.Args[1]
	Source = os.Args[2]
	peers = os.Args[3:]

	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/replace", replaceHandler)
	http.HandleFunc("/test", testHandler)
	http.HandleFunc("/vclock", vclockHandler)
	http.HandleFunc("/ws", wsHandler)

	go transactionManager()

	for _, peer := range peers {
		go listenUpstreamWsConnection(peer)
	}

	if http.ListenAndServe(addr, nil) != nil {
		log.Panic("Error while trying to serve the address")
	}
}
