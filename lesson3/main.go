package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	_ "embed"

	websocket "nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
)

//go:embed static/test.html
var testHTML string

var Source string
var port string

var tm *TransactionManager

var snapshot string
var clock uint64
var vClocks map[string]uint64
var peers []string

func main ()  {
  initialize()
  startReplicasStatePoller()
  bindAndServe()
}

func initialize() {
  tm = NewTransactionManager()
  vClocks = make(map[string]uint64)

  port = os.Getenv("PORT")
  if len(port) == 0 {
    port = "8000"
  }

  Source = os.Getenv("NAME")
  if len(Source) == 0 {
    Source = port
  }

  peersParam := os.Getenv("PEERS")
  if len(peersParam) == 0 {
    peers = []string{}
  } else {
    peers = strings.Split(peersParam, ",")
  }
}

func startReplicasStatePoller() {
  for _, peer := range peers {
    go func (ppeer string) {
      ctx := context.Background()
      log.Printf("Going to connect to peer %s", ppeer)

      var conn *websocket.Conn
      var err error
      var nodeAddr = fmt.Sprintf("ws://%s/ws", ppeer)

      for {
        if conn, _, err = websocket.Dial(ctx, nodeAddr, nil); err != nil {
          log.Printf("Failed to connect to %s, retrying..\n", nodeAddr)
          time.Sleep(3 * time.Second)
          continue
        }
        log.Printf("Connected to peer %s", ppeer)
        break
      }

      var journalEntry JournalEntry
      for {
        err = wsjson.Read(ctx, conn, &journalEntry)
        if err != nil {
          panic(err)
        }

        if journalEntry.Source == Source {
          log.Print("Skipping, same Source")
          continue
        }

        if journalEntry.Id > vClocks[ppeer] + 1 {
          log.Print("Skipping, id and vlock mismatch")
          continue
        }

        vClocks[ppeer] += 1

        tm.Execute([]byte(journalEntry.Payload), journalEntry.Source, journalEntry.Id)
      }
    }(peer)
  }
}

func bindAndServe() {
  http.HandleFunc("/test", testHandler) // HTTP GET /test
  http.HandleFunc("/replace", replaceHandler) // HTTP POST /replace
  http.HandleFunc("/get", getHandler) // HTTP GET /get
  http.HandleFunc("/vclock", vclockHandler) // HTTP GET /vclock
  http.HandleFunc("/ws", wsHandler)

  log.Printf("Listening on %s/tcp, node name: %s, peers: %s", port, Source, peers)

  if err := http.ListenAndServe(":" + port, nil); err != nil {
    log.Fatal(err)
  }
}

func testHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.Error(w, "Not Found", http.StatusNotFound)
  }

  w.Header().Set("Content-Type", "text/html")
  io.WriteString(w, testHTML)
}

func replaceHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    http.Error(w, "Not Found", http.StatusNotFound)
  }

  defer r.Body.Close()
  data, _ := io.ReadAll(r.Body)

  clock += 1
  tm.Execute(data, Source, clock)

  w.Header().Set("Content-Type", "text/plain")
  w.Write(tm.state)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.Error(w, "Not Found", http.StatusNotFound)
  }

  w.Header().Set("Content-Type", "text/plain")
  w.Write(tm.state)
}

func vclockHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.Error(w, "Not Found", http.StatusNotFound)
  }

  w.Header().Set("Content-Type", "text/plain")
  jsonStr, err := json.Marshal(vClocks)
  if err != nil {
    panic(err)
  }
  w.Write(jsonStr)
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
  connection, err := websocket.Accept(w, r, &websocket.AcceptOptions{
    InsecureSkipVerify: true,
    OriginPatterns:     []string{"*"},
  })
  if err != nil {
    panic(err)
  }

  for _, journalEntry := range tm.Journal.entries {
    wsjson.Write(context.Background(), connection, *journalEntry)
  }

  tm.replicationManager.AddNewConnection(connection)
}
