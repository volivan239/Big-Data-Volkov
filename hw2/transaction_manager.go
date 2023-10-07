package main

import "sync"

var snapshot = []byte("Initial value")
var previousSnapshot []byte
var journal [][]byte
var previousJournal [][]byte

var mutex sync.Mutex
var eventChannel = make(chan Transaction)

type Transaction struct {
	Request          []byte
	ResponseListener chan []byte
}

func transactionManager() {
	for transaction := range eventChannel {
		mutex.Lock()

		if transaction.Request != nil {
			journal = append(journal, transaction.Request)
		}

		value := snapshot
		if len(journal) > 0 {
			value = journal[len(journal)-1]
		}

		if transaction.ResponseListener != nil {
			transaction.ResponseListener <- value
		}

		mutex.Unlock()
	}
}
