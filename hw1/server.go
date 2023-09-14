package main

import (
	"net/http"
	"os"
)

const registryFile = "body.txt"
const MAX_BODY_SIZE = 64 * 1024

func getHandler(writer http.ResponseWriter, _ *http.Request) {
	contents, _ := os.ReadFile(registryFile)
	writer.WriteHeader(200)
	writer.Write(contents)
}

func replaceHandler(writer http.ResponseWriter, request *http.Request) {
	var buffer [MAX_BODY_SIZE]byte
	newLen, _ := request.Body.Read(buffer[:])
	os.WriteFile(registryFile, buffer[:newLen], 0777)
	writer.WriteHeader(200)
}

func main() {
	http.HandleFunc("/get", getHandler)
	http.HandleFunc("/replace", replaceHandler)

	if http.ListenAndServe(":8088", nil) != nil {
		println("Unable to serve the address")
		return
	}

	println("Listening for requests")
}
