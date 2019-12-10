package casbinpgadapter

import (
	"database/sql"
	"fmt"
	"log"
	"runtime"

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
	tableName        string
}

// NewAdapter returns a new casbin postgresql adapter
func NewAdapter(connectionString string, tableName string) (*Adapter, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	adapter := &Adapter{
		connectionString,
		db,
		tableName,
	}
	if err = adapter.setup(); err != nil {
		return nil, err
	}
	runtime.SetFinalizer(adapter, finalizer)
	return adapter, nil
}

func finalizer(a *Adapter) {
	if err := a.db.Close(); err != nil {
		log.Println("Cannot close db connection")
	}
}

func (adapter *Adapter) setup() error {
	if err := adapter.createTableIfNeeded(); err != nil {
		return err
	}
	return nil
}

func (adapter *Adapter) createTableIfNeeded() error {
	tx, err := adapter.db.Begin()
	if err != nil {
		log.Print("Cannot start transaction")
		return nil
	}
	_, err = tx.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS "%v" (
			pType varchar(256) not null default '',
			v0 		varchar(256) not null default '',
			v1 		varchar(256) not null default '',
			v2 		varchar(256) not null default '',
			v3 		varchar(256) not null default '',
			v4 		varchar(256) not null default '',
			v5 		varchar(256) not null default ''
		)
	`, adapter.tableName))
	if err != nil {
		log.Printf("Cannot create table %v", err)
		return err
	}
	columns := [7]string{
		"pType",
		"v0",
		"v1",
		"v2",
		"v3",
		"v4",
		"v5",
	}
	for _, column := range columns {
		tx.Exec(fmt.Sprintf(`
			CREATE INDEX IF NOT EXISTS idx_%[1]v_%[2]v ON %[1]v (%[2]v)
		`, adapter.tableName, column))
	}
	err = tx.Commit()
	if err != nil {
		log.Printf("Cannot commit transaction %v", err)
		return err
	}
	return nil
}
