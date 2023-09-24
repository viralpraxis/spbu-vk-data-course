package main

type Transaction struct {
  Command ICommand
}

// either success or failure, now only success
type TransactionResult struct {
  Data []byte
}

type ICommand interface {
  Name() string
}

type GetCommand struct {}

type PutCommand struct {
  Data []byte
}

func (t GetCommand) Name() string {
  return "GET"
}

func (t PutCommand) Name() string {
  return "PUT"
}
