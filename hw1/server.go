package main

import (
	"io"
	"net/http"
	"os"
)

const registryFile = "body.txt"

func getHandler(writer http.ResponseWriter, _ *http.Request) {
	contents, err := os.ReadFile(registryFile)
	if err != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	writer.WriteHeader(http.StatusOK)
	writer.Write(contents)
}

func replaceHandler(writer http.ResponseWriter, request *http.Request) {
	buffer, err := io.ReadAll(request.Body)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	os.WriteFile(registryFile, buffer, 0777)
	writer.WriteHeader(http.StatusOK)
}

func main() {
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/replace", replaceHandler)

	if http.ListenAndServe(":8088", nil) != nil {
		println("Error while trying to serve the address")
	}
}
