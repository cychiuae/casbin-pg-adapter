package casbinpgadapter

import (
	"database/sql"
	"fmt"
	"log"

	cModel "github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"

	// no-lint
	_ "github.com/lib/pq"

	"github.com/casbin/casbin-pg-adapter/pkg/model"
	"github.com/casbin/casbin-pg-adapter/pkg/repository"
)

// Adapter is a postgresql adaptor for casbin
type Adapter struct {
	db                   *sql.DB
	dbSchema             string
	tableName            string
	casbinRuleRepository *repository.CasbinRuleRepository
}

// NewAdapter returns a new casbin postgresql adapter
func NewAdapter(db *sql.DB, tableName string) (*Adapter, error) {
	return NewAdapterWithDBSchema(db, "public", tableName)
}

func NewAdapterWithDBSchema(db *sql.DB, dbSchema string, tableName string) (*Adapter, error) {
	casbinRuleRepository := repository.NewCasbinRuleRepository(dbSchema, tableName, db)
	adapter := &Adapter{
		db,
		dbSchema,
		tableName,
		casbinRuleRepository,
	}

	if err := adapter.setup(); err != nil {
		return nil, err
	}

	return adapter, nil
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
		CREATE TABLE IF NOT EXISTS "%s"."%s" (
			p_type varchar(256) not null default '',
			v0 		varchar(256) not null default '',
			v1 		varchar(256) not null default '',
			v2 		varchar(256) not null default '',
			v3 		varchar(256) not null default '',
			v4 		varchar(256) not null default '',
			v5 		varchar(256) not null default ''
		)
	`, adapter.dbSchema, adapter.tableName))
	if err != nil {
		_ = tx.Rollback()
		log.Printf("Cannot create table %v", err)
		return err
	}
	columns := [7]string{
		"p_type",
		"v0",
		"v1",
		"v2",
		"v3",
		"v4",
		"v5",
	}
	for _, column := range columns {
		_, err = tx.Exec(fmt.Sprintf(`
			CREATE INDEX IF NOT EXISTS idx_%[2]s_%[3]s ON "%[1]s"."%[2]s" (%[3]s)
		`, adapter.dbSchema, adapter.tableName, column))
		if err != nil {
			log.Printf("Cannot create index for column: %v. Error: %v", column, err)
			_ = tx.Rollback()
			return err
		}
	}
	err = tx.Commit()
	if err != nil {
		log.Printf("Cannot commit transaction %v", err)
		_ = tx.Rollback()
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
