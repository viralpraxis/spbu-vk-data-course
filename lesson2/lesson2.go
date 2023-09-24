package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// go 1.21.1

var tm *TransactionManager
var sm *SnapshotManager

var snapshot State

func main ()  {
  tm = NewTransactionManager()
  sm = NewSnapshotManager()

  runSnapshotter()

  bindAndServe()
}

func runSnapshotter() {
  ticker := time.NewTicker(60 * time.Second)

  quit := make(chan struct{})
  go func() {
    for {
      select {
        case <- ticker.C:
          snapshot = sm.TakeSnapshot(tm) // thread-safe impl
        case <- quit:
          ticker.Stop()
          log.Print("Snapshots stopped")
          return
        }
      }
   }()
}

func bindAndServe() {
  http.HandleFunc("/replace", replaceHandler) // HTTP POST /replace
  http.HandleFunc("/get", getHandler) // HTTP GET /get

  port := os.Getenv("PORT")
  if len(port) == 0 {
    port = "13377"
  }

  log.Printf("Listening on %s/tcp", port)

  if err := http.ListenAndServe(":" + port, nil); err != nil {
    log.Fatal(err)
  }
}

func replaceHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodPost {
    http.Error(w, "Not Found", http.StatusNotFound)
  }

  defer r.Body.Close()
  data, _ := io.ReadAll(r.Body)

  tx := Transaction{Command: &PutCommand{Data: data}}
  txResult := tm.Execute(tx)

  w.Header().Set("Content-Type", "text/plain")
  w.Write(txResult.Data)
}

func getHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.Error(w, "Not Found", http.StatusNotFound)
  }

  tx := Transaction{Command: &GetCommand{}}
  txResult := tm.Execute(tx)

  w.Header().Set("Content-Type", "application/octet-stream")
  w.Write(txResult.Data)
}
