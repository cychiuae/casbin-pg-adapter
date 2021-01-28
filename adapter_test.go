package casbinpgadapter

import (
	"database/sql"
	"os"
	"testing"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
	"github.com/cychiuae/casbin-pg-adapter/pkg/model"
)

// TestAdapter is a very bad all-in-one integration test to test the adapter
func TestAdapter(t *testing.T) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatalf("Fail to open db %v", err)
		return
	}

	enforcer, err := casbin.NewEnforcer("./example/model.conf", "./example/policy.csv")
	if err != nil {
		t.Fatal("Cannot create enforcer")
		return
	}
	adapter, err := NewAdapter(db, "casbin")
	if err != nil {
		t.Fatalf("Cannot create adapter %v", err)
		return
	}
	if err = adapter.SavePolicy(enforcer.GetModel()); err != nil {
		t.Fatalf("Cannot initial policy %v", err)
		return
	}

	adapter, err = NewAdapter(db, "casbin")
	if err != nil {
		t.Fatalf("Cannot create adapter %v", err)
		return
	}
	enforcer, err = casbin.NewEnforcer("./example/model.conf", adapter)
	if err != nil {
		t.Fatalf("Cannot create enforcer %v", err)
		return
	}
	enforcerPolicy := enforcer.GetPolicy()
	want := [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}}
	if !util.Array2DEquals(enforcerPolicy, want) {
		t.Fatalf("Want %v but got %v", want, enforcerPolicy)
		return
	}

	enforcer.EnableAutoSave(false)
	if _, err = enforcer.AddPolicy("alice", "data1", "write"); err != nil {
		t.Fatalf("Cannot add policy")
		return
	}
	if err = enforcer.LoadPolicy(); err != nil {
		t.Fatalf("Cannot load policy")
		return
	}
	enforcerPolicy = enforcer.GetPolicy()
	if !util.Array2DEquals(enforcerPolicy, want) {
		t.Fatalf("Want %v but got %v", want, enforcerPolicy)
		return
	}

	enforcer.EnableAutoSave(true)

	if _, err = enforcer.AddPolicy("alice", "data1", "write"); err != nil {
		t.Fatalf("Cannot add policy")
		return
	}
	if err = enforcer.LoadPolicy(); err != nil {
		t.Fatalf("Cannot load policy")
		return
	}
	enforcerPolicy = enforcer.GetPolicy()
	want = [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}, {"alice", "data1", "write"}}
	if !util.Array2DEquals(enforcerPolicy, want) {
		t.Fatalf("Want %v but got %v", want, enforcerPolicy)
		return
	}

	if _, err = enforcer.RemovePolicy("alice", "data1", "write"); err != nil {
		t.Fatalf("Cannot remove policy")
		return
	}
	if err = enforcer.LoadPolicy(); err != nil {
		t.Fatalf("Cannot load policy")
		return
	}
	enforcerPolicy = enforcer.GetPolicy()
	want = [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}, {"data2_admin", "data2", "read"}, {"data2_admin", "data2", "write"}}
	if !util.Array2DEquals(enforcerPolicy, want) {
		t.Fatalf("Want %v but got %v", want, enforcerPolicy)
		return
	}

	if _, err = enforcer.RemoveFilteredPolicy(0, "data2_admin"); err != nil {
		t.Fatalf("Cannot remove filtered policy")
		return
	}
	if err = enforcer.LoadPolicy(); err != nil {
		t.Fatalf("Cannot load policy")
		return
	}
	enforcerPolicy = enforcer.GetPolicy()
	want = [][]string{{"alice", "data1", "read"}, {"bob", "data2", "write"}}
	if !util.Array2DEquals(enforcerPolicy, want) {
		t.Fatalf("Want %v but got %v", want, enforcerPolicy)
		return
	}
}

// TestAdapter is a very bad all-in-one integration test to test the adapter
func TestFilteredAdapter(t *testing.T) {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		t.Fatalf("Fail to open db %v", err)
		return
	}

	enforcer, err := casbin.NewEnforcer("./example/model.conf", "./example/policy.csv")
	if err != nil {
		t.Fatal("Cannot create enforcer")
		return
	}
	adapter, err := NewAdapter(db, "casbin")
	if err != nil {
		t.Fatalf("Cannot create adapter %v", err)
		return
	}
	if err = adapter.SavePolicy(enforcer.GetModel()); err != nil {
		t.Fatalf("Cannot initial policy %v", err)
		return
	}

	filetedAdapter, err := NewFilteredAdapter(db, "casbin")
	if err != nil {
		t.Fatalf("Cannot create filtered adapter %v", err)
		return
	}
	enforcer, err = casbin.NewEnforcer("./example/model.conf", filetedAdapter)
	if err != nil {
		t.Fatalf("Cannot create enforcer %v", err)
		return
	}
	err = enforcer.LoadFilteredPolicy(&model.Filter{P: []string{"alice"}})
	if err != nil {
		t.Fatalf("can not load filtered policy: %v", err)
	}

	enforcerPolicy := enforcer.GetPolicy()

	want := [][]string{{"alice", "data1", "read"}}
	if !util.Array2DEquals(enforcerPolicy, want) {
		t.Fatalf("Want %v but got %v", want, enforcerPolicy)
		return
	}

	if _, err = enforcer.AddPolicy("alice", "data1", "write"); err != nil {
		t.Fatalf("Cannot add policy")
		return
	}
	if err = enforcer.LoadFilteredPolicy(&model.Filter{P: []string{"alice"}}); err != nil {
		t.Fatalf("Cannot load policy: %v", err)
		return
	}
	enforcerPolicy = enforcer.GetPolicy()
	want = [][]string{{"alice", "data1", "read"}, {"alice", "data1", "write"}}
	if !util.Array2DEquals(enforcerPolicy, want) {
		t.Fatalf("Want %v but got %v", want, enforcerPolicy)
		return
	}
}
