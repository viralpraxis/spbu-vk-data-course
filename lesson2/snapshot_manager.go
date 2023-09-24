package main

import "sync"

type SnapshotManager struct {
  mu *sync.Mutex
}

func NewSnapshotManager() *SnapshotManager {
  return &SnapshotManager{
    mu: &sync.Mutex{},
  }
}

func (sm *SnapshotManager) TakeSnapshot(txManager *TransactionManager) State {
  sm.mu.Lock()
  var snapshot State

  ok, snapshot := txManager.wal.Snapshot()
  if !ok {
    sm.mu.Unlock()
    return nil
  }

  sm.mu.Unlock()

  return snapshot
}
