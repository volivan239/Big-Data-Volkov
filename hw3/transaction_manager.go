package main

import (
	jsonpatch "github.com/evanphx/json-patch"
	"log"
)

var snapshot = "{}"
var Journal []Transaction
var vclock = make(map[string]uint64)

var transactionManagerEventChannel = make(chan Transaction)

type Transaction struct {
	Source  string
	Id      uint64
	Payload string
}

func transactionManager() {
	for transaction := range transactionManagerEventChannel {
		if transaction.Id <= vclock[transaction.Source] {
			continue
		}

		Journal = append(Journal, transaction)

		patch, err := jsonpatch.DecodePatch([]byte(transaction.Payload))
		if err != nil {
			log.Printf("Failed to decode patch: %v", err)
			continue
		}

		new_state, err := patch.Apply([]byte(snapshot))
		if err != nil {
			log.Printf("Failed to apply patch: %v", err)
			continue
		}
		snapshot = string(new_state)

		vclock[transaction.Source] = transaction.Id

		for _, connection := range downstreamConnections {
			go sendByWsConnection(connection, transaction)
		}
	}
}
