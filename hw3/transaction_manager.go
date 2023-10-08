package main

import (
	jsonpatch "github.com/evanphx/json-patch"
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

		patch, _ := jsonpatch.DecodePatch([]byte(transaction.Payload))
		new_state, _ := patch.Apply([]byte(snapshot))
		snapshot = string(new_state)

		vclock[transaction.Source] = transaction.Id

		for _, connection := range wsConnections {
			go sendByWsConnection(connection, transaction)
		}
	}
}
