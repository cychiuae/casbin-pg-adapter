# Casbin Postgres Adapter

Casbin Postgres Adapter is the postgres adapter for [Casbin](Casbin)

## Installation
```sh
$ go get github.com/cychiuae/casbin-pg-adapter
```

## Example
```go
package main

import (
  "os"

  "github.com/casbin/casbin/v2"
  "github.com/cychiuae/casbin-pg-adapter"
)

func main() {
  connectionString := "postgresql://postgres:@localhost:5432/postgres?sslmode=disable"
  tableName := "casbin"
  adapter, err := casbinpgadapter.NewAdapter(connectionString, tableName)
  if err != nil {
    panic(err)
  }

  enforcer, err := casbin.NewEnforcer("./examples/model.conf", adapter)
  if err != nil {
    panic(err)
  }

  // Load stored policy from database
  enforcer.LoadPolicy()

  // Do permission checking
  enforcer.Enforce("alice", "data1", "write")

  // Do some mutations
  enforcer.AddPolicy("alice", "data2", "write")
  enforcer.RemovePolicy("alice", "data1", "write")

  // Persist policy to database
  enforcer.SavePolicy()
}
```
