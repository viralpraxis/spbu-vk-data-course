package main

import (
	"log"
	"sync"
)

type Journal struct { // aka WAL
  entries []*JournalEntry
  mu *sync.RWMutex
}

type JournalEntry struct {
  Source string    `json:"Source"`
  Id uint64        `json:"Id"`
  Payload string   `json:"Payload"`
}

func NewJournal() *Journal {
  return &Journal{
    entries: make([]*JournalEntry, 0),
    mu: &sync.RWMutex{},
  }
}

func (wal *Journal) AddRecord(source string, id uint64, payload []byte) *JournalEntry {
  var journalEntry = &JournalEntry{
    Source: source,
    Id: id,
    Payload: string(payload),
  }

  wal.mu.Lock()
  wal.entries = append(wal.entries, journalEntry)
  wal.mu.Unlock()

  log.Printf("[Journal] inserted new log %s", payload)

  return journalEntry
}
