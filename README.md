# UBM
User Behaviour Model implementation. Simple library for tracking user activity.

[![MIT licensed](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/dzen-it/ubm/master/LICENSE) [![GoDoc](https://godoc.org/github.com/dzen-it/ubm?status.svg)](https://godoc.org/github.com/dzen-it/ubm)

## Installation

`go get -u github.com/dzen-it/ubm`

## Quick Start

Create a server

```go
db, err := ubm.NewMongoDB("db.example.com:27017", "dbname")
if err != nil {
	panic(err)
}
srv := ubm.NewServerCoala(5683, db)

go func() {
	srv.Serve()
}()
```

Create client.

```go
client := ubm.NewClientCoala("127.0.0.1:5683")
```

Add action for user.

```go
if err := client.AddAction("user_id", "action-1"); err != nil {
	panic(err)
}
```

Get different info about actions for user.

```go
action, err := client.GetAction("user_id", "action-1")
if err!=nil{
    panic(err)
}

fmt.Printf("Last call: %v, total calls: %v\n", action.LastCall, action.Count)
// Display: "Last call: 2018-04-13 13:13:09.876, total calls: 42"

lastAction, err := client.GetLastAction("user_id")
if err!=nil{
    panic(err)
}

fmt.Printf("Last action: %v, last calls: %v\n", lastAction.Name, action.LastCall)
// Display: "Last action: action-1, last call: 2018-04-13 13:13:09.876"
```

Manage action states using triggers.

```go
// Sets a enabling action-1 and a disabling action-2
client.SetTrigger("user_id", "action-1", "action-2")

fmt.Println("action-1 is blocked:", client.TriggerStatus("user_id", "action-1"))
// Display: "action-1 is blocked: true"

// Add a disabling action "action-2"
if err := client.AddAction("user_id", "action-2"); err != nil {
	panic(err)
}

fmt.Println("action-1 is blocked:", client.TriggerStatus("user_id", "action-1"))
// Display: "action-1 is blocked: false"

```