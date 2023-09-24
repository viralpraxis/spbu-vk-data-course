package main

import (
	"sync"
)

type TransactionManager struct {
  mu *sync.Mutex
	wal *WAL
	state []byte
}

func NewTransactionManager() *TransactionManager {
	return &TransactionManager{
		mu: &sync.Mutex{},
		wal: NewWAL(),
	}
}

func (tm *TransactionManager) Execute(tx Transaction) TransactionResult {
	var txResult TransactionResult

	tm.mu.Lock()

	switch tx.Command.Name() {
  	case "GET": {
	  	txResult = TransactionResult{Data: tm.state}
			break
	  }
    case "PUT": {
			putCommand, ok := tx.Command.(*PutCommand)
			if !ok {
				panic("Invalid API usage")
			}

			tm.state = putCommand.Data
			tm.wal.AddRecord(tx)
			break
		}
	}

	txResult = TransactionResult{Data: tm.state}

	tm.mu.Unlock()

	return txResult
}
