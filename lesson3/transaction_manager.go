package main

import (
	"sync"

	jsonpatch "github.com/evanphx/json-patch"
)

type TransactionManager struct {
  mu *sync.Mutex
  Journal *Journal
  replicationManager *ReplicationManager
  state []byte
  clock *uint64
}

func NewTransactionManager() *TransactionManager {
  return &TransactionManager{
    mu: &sync.Mutex{},
    Journal: NewJournal(),
    state: []byte("{}"),
    replicationManager: NewReplicationManager(),
  }
}

func (tm *TransactionManager) Execute(tx Transaction, source string, id uint64) {
  tm.mu.Lock()
  defer func() { tm.mu.Unlock() }()

  patch, err := jsonpatch.DecodePatch(tx)
  if err != nil {
    panic(err)
  }
  newState, err := patch.Apply(tm.state)
  if err != nil {
    panic(err)
  }
  tm.state = newState

  journalEntry := tm.Journal.AddRecord(source, id, tx)

  tm.replicationManager.Notify(journalEntry)
}
