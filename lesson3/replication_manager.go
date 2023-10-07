package main

import (
	"context"
	"log"

	"nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

type ReplicationManager struct {
	connections []*websocket.Conn
}

func NewReplicationManager() *ReplicationManager {
	return &ReplicationManager{
		connections: make([]*websocket.Conn, 0, 10),
	}
}

func (rm *ReplicationManager) AddNewConnection(connection *websocket.Conn) {
	rm.connections = append(rm.connections, connection)
}

func (rm *ReplicationManager) Notify(journalEntry *JournalEntry) {
	log.Printf("Replicating to %d nodes..", len(rm.connections))

	for _, conn := range rm.connections {
		err := wsjson.Write(context.Background(), conn, journalEntry)
		if err != nil {
			panic(err)
		}

		log.Println("Replicated successfully")
	}
}
