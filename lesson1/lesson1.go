package main

import (
	"io"
	"log"
	"net/http"
	"os"
	"sync"
)

// go 1.21.1

var mu *sync.RWMutex
const storageFilepath = "/tmp/yk-hw-1"

func main ()  {
  http.HandleFunc("/replace", replaceHandler) // HTTP POST /replace
  http.HandleFunc("/get", getHandler) // HTTP GET /get

  port := os.Getenv("PORT")
  if len(port) == 0 {
    port = "13377"
  }

  mu = &sync.RWMutex{}

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

  mu.Lock()
  if err := os.WriteFile(storageFilepath, data, 0644); err != nil {
    log.Fatal(err)
  }
  mu.Unlock()

  w.Header().Set("Content-Type", "text/plain")
  io.WriteString(w, "OK")
}

func getHandler(w http.ResponseWriter, r *http.Request) {
  if r.Method != http.MethodGet {
    http.Error(w, "Not Found", http.StatusNotFound)
  }

  mu.RLock()
  data, err := os.ReadFile(storageFilepath)
  mu.RUnlock()
  if err != nil && !os.IsNotExist(err) {
    log.Fatal(err)
  } else if err != nil {
    return // return initial state (empty octet stream)
  }

  w.Header().Set("Content-Type", "application/octet-stream")
  w.Write(data)
}
