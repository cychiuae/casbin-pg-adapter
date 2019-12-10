package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/cychiuae/casbin-pg-adapter/pkg/model"
)

// CasbinRuleRepository is the bridge for adapter and db
type CasbinRuleRepository struct {
	tableName string
	db        *sql.DB
}

// NewCasbinRuleRepository returns a new CasbinRuleRepository
func NewCasbinRuleRepository(tableName string, db *sql.DB) *CasbinRuleRepository {
	return &CasbinRuleRepository{
		tableName: tableName,
		db:        db,
	}
}

// LoadAllCasbinRules loads all casbin rules from db
func (repository *CasbinRuleRepository) LoadAllCasbinRules() ([]model.CasbinRule, error) {
	rows, err := repository.db.Query(fmt.Sprintf(`
		SELECT pType, v0, v1, v2, v3, v4, v5 FROM %s
	`, repository.tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	casbinRules := make([]model.CasbinRule, 0)
	for rows.Next() {
		var pType string
		var v0 string
		var v1 string
		var v2 string
		var v3 string
		var v4 string
		var v5 string
		scanErr := rows.Scan(
			&pType,
			&v0,
			&v1,
			&v2,
			&v3,
			&v4,
			&v5,
		)
		if scanErr != nil {
			return nil, scanErr
		}
		casbinRule := model.CasbinRule{
			PType: pType,
			V0:    v0,
			V1:    v1,
			V2:    v2,
			V3:    v3,
			V4:    v4,
			V5:    v5,
		}
		casbinRules = append(casbinRules, casbinRule)
	}
	return casbinRules, nil
}

// InsertCasbinRule insert a casbin rule into db
func (repository *CasbinRuleRepository) InsertCasbinRule(casbinRule model.CasbinRule) error {
	tx, err := repository.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		fmt.Sprintf(`
			INSERT INTO %s (pType, v0, v1, v2, v3, v4, v5)
			VALUES
				($1, $2, $3, $4, $5, $6, $7)
		`, repository.tableName),
		casbinRule.PType,
		casbinRule.V0,
		casbinRule.V1,
		casbinRule.V2,
		casbinRule.V3,
		casbinRule.V4,
		casbinRule.V5,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// DeleteCasbinRule insert a casbin rule into db
func (repository *CasbinRuleRepository) DeleteCasbinRule(casbinRule model.CasbinRule) error {
	tx, err := repository.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(
		fmt.Sprintf(`
			DELETE FROM %s
			WHERE
				pType = $1 AND
				v0 = $1 AND
				v1 = $2 AND
				v2 = $3 AND
				v3 = $4 AND
				v4 = $5 AND
				v5 = $6 AND
		`, repository.tableName),
		casbinRule.PType,
		casbinRule.V0,
		casbinRule.V1,
		casbinRule.V2,
		casbinRule.V3,
		casbinRule.V4,
		casbinRule.V5,
	)
	if err != nil {
		tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

// ReplaceAllCasbinRules replaces the existing db with casbinRules
func (repository *CasbinRuleRepository) ReplaceAllCasbinRules(casbinRules []model.CasbinRule) error {
	tx, err := repository.db.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(fmt.Sprintf(`
		TRUNCATE TABLE %s
	`, repository.tableName))
	if err != nil {
		tx.Rollback()
		return err
	}

	values := make([]string, 0)
	for _, casbinRule := range casbinRules {
		value := fmt.Sprintf("(%s)", strings.Join([]string{
			casbinRule.PType,
			casbinRule.V0,
			casbinRule.V1,
			casbinRule.V2,
			casbinRule.V3,
			casbinRule.V4,
			casbinRule.V5,
		}, ","))
		values = append(values, value)
	}

	_, err = tx.Exec(
		fmt.Sprintf(
			`
				INSERT INTO %s (pType, v0, v1, v2, v3, v4, v5)
				VALUES %s
			`,
			repository.tableName,
			strings.Join(values, ",")),
	)
	if err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}
