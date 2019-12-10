package casbinpgadapter

import (
	"database/sql"

	"github.com/cychiuae/casbin-pg-adapter/src/hello"

	// no-lint
	_ "github.com/lib/pq"
)

// HelloWorld returns hello world
func HelloWorld() string {
	return hello.HelloWorld()
}

// Adapter is a postgresql adaptor for casbin
type Adapter struct {
	connectionString string
	db               *sql.DB
}

// NewAdapter returns a new casbin postgresql adapter
func NewAdapter(connectionString string) (*Adapter, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	adapter := &Adapter{
		connectionString,
		db,
	}
	return adapter, nil
}
