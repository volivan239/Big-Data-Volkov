package main

import (
	"bytes"
	"io"
	"net/http"
)

func getHandler(writer http.ResponseWriter, _ *http.Request) {
	responseListener := make(chan []byte)
	eventChannel <- Transaction{ResponseListener: responseListener}
	response := <-responseListener
	writer.WriteHeader(http.StatusOK)
	writer.Write(response)
}

func replaceHandler(writer http.ResponseWriter, request *http.Request) {
	buffer, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	responseListener := make(chan []byte)
	eventChannel <- Transaction{Request: buffer, ResponseListener: responseListener}
	response := <-responseListener

	if !bytes.Equal(buffer, response) {
		writer.WriteHeader(http.StatusInternalServerError)
	} else {
		writer.WriteHeader(http.StatusOK)
	}
}

func main() {
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/replace", replaceHandler)

	go transactionManager()
	go snapshotManager()

	if http.ListenAndServe(":8088", nil) != nil {
		println("Error while trying to serve the address")
	}
}
