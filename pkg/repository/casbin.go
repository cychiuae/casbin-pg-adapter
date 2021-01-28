package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/cychiuae/casbin-pg-adapter/pkg/model"
)

// CasbinRuleRepository is the bridge for adapter and db
type CasbinRuleRepository struct {
	dbSchema  string
	tableName string
	db        *sql.DB
}

// NewCasbinRuleRepository returns a new CasbinRuleRepository
func NewCasbinRuleRepository(dbSchema string, tableName string, db *sql.DB) *CasbinRuleRepository {
	return &CasbinRuleRepository{
		dbSchema:  dbSchema,
		tableName: tableName,
		db:        db,
	}
}

// LoadAllCasbinRules loads all casbin rules from db
func (repository *CasbinRuleRepository) LoadAllCasbinRules() ([]model.CasbinRule, error) {
	rows, err := repository.db.Query(fmt.Sprintf(`
		SELECT p_type, v0, v1, v2, v3, v4, v5 FROM "%s"."%s"
	`, repository.dbSchema, repository.tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return loadPolicyFromRows(rows)
}

// LoadFilteredRules loads casbin rules filtered
func (repository *CasbinRuleRepository) LoadFilteredRules(filter *model.Filter) ([]model.CasbinRule, error) {
	pFilter, gFilter := filteredWhereValues(filter)
	rows, err := repository.db.Query(fmt.Sprintf(`
		SELECT p_type, v0, v1, v2, v3, v4, v5 FROM "%s"."%s"
		 WHERE 
            ( p_type LIKE 'g%%' AND v0 LIKE $1 AND v1 LIKE $2 AND v2 LIKE $3 AND v3 LIKE $4 AND v4 LIKE $5 AND v5 LIKE $6 )
        OR 
            ( p_type LIKE 'p%%' AND v0 LIKE $7 AND v1 LIKE $8 AND v2 LIKE $9 AND v3 LIKE $10 AND v4 LIKE $11 AND v5 LIKE $12 )
	`, repository.dbSchema, repository.tableName),
		gFilter[0], gFilter[1], gFilter[2], gFilter[3], gFilter[4], gFilter[5],
		pFilter[0], pFilter[1], pFilter[2], pFilter[3], pFilter[4], pFilter[5],
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return loadPolicyFromRows(rows)
}

func loadPolicyFromRows(rows *sql.Rows) ([]model.CasbinRule, error) {
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
			INSERT INTO "%s"."%s" (p_type, v0, v1, v2, v3, v4, v5)
			VALUES
				($1, $2, $3, $4, $5, $6, $7)
		`, repository.dbSchema, repository.tableName),
		casbinRule.PType,
		casbinRule.V0,
		casbinRule.V1,
		casbinRule.V2,
		casbinRule.V3,
		casbinRule.V4,
		casbinRule.V5,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
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
	var queryBuilder strings.Builder
	args := make([]interface{}, 0)
	queryBuilder.WriteString(fmt.Sprintf(`
			DELETE FROM "%s"."%s"
			WHERE
				p_type = $1
	`, repository.dbSchema, repository.tableName))
	args = append(args, casbinRule.PType)

	if casbinRule.V0 != "" {
		args = append(args, casbinRule.V0)
		queryBuilder.WriteString(fmt.Sprintf("AND v0 = $%d ", len(args)))
	}
	if casbinRule.V1 != "" {
		args = append(args, casbinRule.V1)
		queryBuilder.WriteString(fmt.Sprintf("AND v1 = $%d ", len(args)))
	}
	if casbinRule.V2 != "" {
		args = append(args, casbinRule.V2)
		queryBuilder.WriteString(fmt.Sprintf("AND v2 = $%d ", len(args)))
	}
	if casbinRule.V3 != "" {
		args = append(args, casbinRule.V3)
		queryBuilder.WriteString(fmt.Sprintf("AND v3 = $%d ", len(args)))
	}
	if casbinRule.V4 != "" {
		args = append(args, casbinRule.V4)
		queryBuilder.WriteString(fmt.Sprintf("AND v4 = $%d ", len(args)))
	}
	if casbinRule.V5 != "" {
		args = append(args, casbinRule.V5)
		queryBuilder.WriteString(fmt.Sprintf("AND v5 = $%d ", len(args)))
	}

	_, err = tx.Exec(
		queryBuilder.String(),
		args...,
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}
	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
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
		TRUNCATE TABLE "%s"."%s"
	`, repository.dbSchema, repository.tableName))
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	values := make([]string, 0)
	for _, casbinRule := range casbinRules {
		value := fmt.Sprintf("(%s)", strings.Join([]string{
			fmt.Sprintf("'%s'", casbinRule.PType),
			fmt.Sprintf("'%s'", casbinRule.V0),
			fmt.Sprintf("'%s'", casbinRule.V1),
			fmt.Sprintf("'%s'", casbinRule.V2),
			fmt.Sprintf("'%s'", casbinRule.V3),
			fmt.Sprintf("'%s'", casbinRule.V4),
			fmt.Sprintf("'%s'", casbinRule.V5),
		}, ","))
		values = append(values, value)
	}

	_, err = tx.Exec(
		fmt.Sprintf(
			`
				INSERT INTO "%s".%s (p_type, v0, v1, v2, v3, v4, v5)
				VALUES %s
			`,
			repository.dbSchema,
			repository.tableName,
			strings.Join(values, ",")),
	)
	if err != nil {
		_ = tx.Rollback()
		return err
	}

	if err = tx.Commit(); err != nil {
		_ = tx.Rollback()
		return err
	}
	return nil
}

func filteredWhereValues(filter *model.Filter) ([]string, []string) {
	p, g := []string{"%", "%", "%", "%", "%", "%"}, []string{"%", "%", "%", "%", "%", "%"}
	for i, token := range filter.P {
		if token != "" {
			p[i] = token
		}
	}
	for i, token := range filter.G {
		if token != "" {
			g[i] = token
		}
	}
	return p, g
}
