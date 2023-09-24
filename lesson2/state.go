package main

// At this monent, database only operates within single `data cell`
// so its state is represented with single BLOB.
type State = []byte
