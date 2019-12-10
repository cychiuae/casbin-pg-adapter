package casbinpgadapter

import (
	"os"
	"testing"
)

func TestAdapter(t *testing.T) {
	connectionString := os.Getenv("DATABASE_URL")
	adapter, err := NewAdapter(connectionString)
	if adapter == nil || err != nil {
		t.Error("Cannot create adapter")
		return
	}
}
