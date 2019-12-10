package casbinpgadapter

import (
	"os"
	"testing"
)

func TestAdapter(t *testing.T) {
	connectionString := os.Getenv("DATABASE_URL")
	adapter, err := NewAdapter(connectionString, "casbin")
	if adapter == nil || err != nil {
		t.Errorf("Cannot create adapter adatper: %v err: %v", adapter, err)
	}
}
