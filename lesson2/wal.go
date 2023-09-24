package main

import (
	"container/list"
	"log"
)

type WAL struct {
	logs list.List
}

func NewWAL() *WAL {
	return &WAL{
		logs: *list.New(),
	}
}

func (wal *WAL) AddRecord(tx Transaction) bool {
  wal.logs.PushFront(tx)

  log.Printf("[WAL] inserted new log %s", tx)

  return true
}

// NOTE:
// At this moment the only mutation operation is PUT
// and we only store PUTs in WAL
// so to take the snapshot we just take the latest enrty in list
// Possible improvement: WAL truncation
func (wal *WAL) Snapshot() (bool, State) {
  log.Printf("[SNAPHOT] taking new snapshot")

	var lastEntry *list.Element
	if lastEntry = wal.logs.Front(); lastEntry == nil {
    return false, nil
	}

	lastTX, ok := lastEntry.Value.(Transaction).Command.(*PutCommand)

	if !ok {
		panic("Invalid API usage")
	}

  return true, lastTX.Data
}
