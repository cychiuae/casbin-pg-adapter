package casbinpgadapter

import (
	"database/sql"
	"errors"

	casbinModel "github.com/casbin/casbin/v2/model"
	"github.com/cychiuae/casbin-pg-adapter/pkg/model"
)

// FilteredAdapter is the filtered file adapter for Casbin. It can load policy
// from file or save policy to file and supports loading of filtered policies.
type FilteredAdapter struct {
	*Adapter
	filtered bool
}

// NewFilteredAdapter is the constructor for FilteredAdapter.
func NewFilteredAdapter(db *sql.DB, tableName string) (*FilteredAdapter, error) {
	a := FilteredAdapter{filtered: false}
	var err error
	a.Adapter, err = NewAdapter(db, tableName)
	return &a, err
}

// NewFilteredAdapterWithDBSchema return a pointer for FilteredAdapter which has schema dbSchema
func NewFilteredAdapterWithDBSchema(db *sql.DB, dbSchema string, tableName string) (*FilteredAdapter, error) {
	a := FilteredAdapter{filtered: false}
	var err error
	a.Adapter, err = NewAdapterWithDBSchema(db, dbSchema, tableName)
	return &a, err
}

// LoadPolicy loads all policy rules from the storage.
func (a *FilteredAdapter) LoadPolicy(model casbinModel.Model) error {
	a.filtered = false
	return a.Adapter.LoadPolicy(model)
}

// LoadFilteredPolicy loads only policy rules that match the filter.
func (a *FilteredAdapter) LoadFilteredPolicy(mod casbinModel.Model, filter interface{}) error {
	mod.ClearPolicy()
	if filter == nil {
		return a.LoadPolicy(mod)
	}

	filterValue, ok := filter.(*model.Filter)
	if !ok {
		return errors.New("invalid filter type")
	}
	err := a.loadFilteredPolicyFile(mod, filterValue)
	if err == nil {
		a.filtered = true
	}
	return err
}

func (a *FilteredAdapter) loadFilteredPolicyFile(model casbinModel.Model, filter *model.Filter) error {
	casbinRules, err := a.casbinRuleRepository.LoadFilteredRules(filter)
	if err != nil {
		return err
	}

	for _, casbinRule := range casbinRules {
		rule := casbinRule.ToStringSlice()
		sec := rule[0][0:1]
		model.AddPolicy(sec, rule[0], rule[1:])
	}
	return nil
}

// IsFiltered returns true if the loaded policy has been filtered.
func (a *FilteredAdapter) IsFiltered() bool {
	return a.filtered
}

// SavePolicy saves all policy rules to the storage.
func (a *FilteredAdapter) SavePolicy(model casbinModel.Model) error {
	if a.filtered {
		return errors.New("cannot save a filtered policy")
	}
	return a.Adapter.SavePolicy(model)
}
