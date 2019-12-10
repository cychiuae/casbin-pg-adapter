package casbinpgadapter

import (
	"database/sql"
	"fmt"
	"log"
	"runtime"

	cModel "github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"

	// no-lint
	_ "github.com/lib/pq"

	"github.com/cychiuae/casbin-pg-adapter/pkg/model"
	"github.com/cychiuae/casbin-pg-adapter/pkg/repository"
)

// Adapter is a postgresql adaptor for casbin
type Adapter struct {
	connectionString     string
	db                   *sql.DB
	tableName            string
	casbinRuleRepository *repository.CasbinRuleRepository
}

// NewAdapter returns a new casbin postgresql adapter
func NewAdapter(connectionString string, tableName string) (*Adapter, error) {
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	casbinRuleRepository := repository.NewCasbinRuleRepository(tableName, db)
	adapter := &Adapter{
		connectionString,
		db,
		tableName,
		casbinRuleRepository,
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
		return err
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

// LoadPolicy loads all policy rules from the storage.
func (adapter *Adapter) LoadPolicy(cmodel cModel.Model) error {
	casbinRules, err := adapter.casbinRuleRepository.LoadAllCasbinRules()
	if err != nil {
		return err
	}

	for _, casbinRule := range casbinRules {
		persist.LoadPolicyLine(casbinRule.ToPolicyLine(), cmodel)
	}

	return nil
}

// SavePolicy saves all policy rules to the storage.
func (adapter *Adapter) SavePolicy(cmodel cModel.Model) error {
	casbinRules := make([]model.CasbinRule, 0)
	for pType, ast := range cmodel["p"] {
		for _, rule := range ast.Policy {
			casbinRule := model.NewCasbinRuleFromPTypeAndRule(pType, rule)
			casbinRules = append(casbinRules, casbinRule)
		}
	}
	for pType, ast := range cmodel["g"] {
		for _, rule := range ast.Policy {
			casbinRule := model.NewCasbinRuleFromPTypeAndRule(pType, rule)
			casbinRules = append(casbinRules, casbinRule)
		}
	}
	if err := adapter.casbinRuleRepository.ReplaceAllCasbinRules(casbinRules); err != nil {
		return err
	}
	return nil
}

// AddPolicy adds a policy rule to the storage.
// This is part of the Auto-Save feature.
func (adapter *Adapter) AddPolicy(sec string, ptype string, rule []string) error {
	casbinRule := model.NewCasbinRuleFromPTypeAndRule(ptype, rule)
	err := adapter.casbinRuleRepository.InsertCasbinRule(casbinRule)
	return err
}

// RemovePolicy removes a policy rule from the storage.
// This is part of the Auto-Save feature.
func (adapter *Adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	casbinRule := model.NewCasbinRuleFromPTypeAndRule(ptype, rule)
	err := adapter.casbinRuleRepository.DeleteCasbinRule(casbinRule)
	return err
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
// This is part of the Auto-Save feature.
func (adapter *Adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	casbinRule := model.NewCasbinRuleFromPTypeAndFilter(ptype, fieldIndex, fieldValues...)
	err := adapter.casbinRuleRepository.DeleteCasbinRule(casbinRule)
	return err
}
