# Mox – MongoDB ODM for Go

> Shape MongoDB like clay — structured, expressive, effortless.

Mox is a lightweight and idiomatic ODM (Object Document Mapper) for MongoDB in Go, built for developers who want clean data modeling and fast development.

## Features

- Struct-based modeling with `model.Base`
- Automatic timestamps (CreatedAt, UpdatedAt)
- Lifecycle hooks (`BeforeSave`, `AfterFind`)
- CRUD operations with transactions
- Session and transaction management
- BSON utilities for ID handling

## Installation

```bash
go get github.com/paywithclay/mox@latest
```

## Quickstart

```go
package main

import (
	"context"
	"log"
	
	"github.com/paywithclay/mox"
	"github.com/paywithclay/mox/model"
)

type User struct {
	model.Base
	Name string `bson:"name"`
}

func (u *User) CollectionName() string { return "users" }

func main() {
	// Connect to MongoDB
	client, err := mox.Connect("mongodb://localhost:27017", "mydb")
	if err != nil {
		log.Fatal(err)
	}
	defer client.Client.Disconnect(context.Background())

	// Create a user
	user := &User{Name: "John Doe"}
	err = client.Collection(user).Save(context.Background(), user)
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("Created user with ID: %v", user.ID)
}
```

## Documentation

### Models
Embed `model.Base` in your structs:
```go
type Product struct {
	model.Base
	Name  string  `bson:"name"`
	Price float64 `bson:"price"`
}
```

### Hooks
Implement hook interfaces:
```go
func (p *Product) BeforeSave() error {
	p.UpdatedAt = time.Now()
	return nil
}
```

## License

MIT