package adapter

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"

	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
)

const (
	tableExistsErrorCode = "ERROR #42P07"
)

// CasbinRule represent casbin_rule table in database
type CasbinRule struct {
	tableName struct{} `pg:"casbin_rules"`

	PType string `pg:"p_type"`
	V0    string `pg:"v0"`
	V1    string `pg:"v1"`
	V2    string `pg:"v2"`
	V3    string `pg:"v3"`
	V4    string `pg:"v4"`
	V5    string `pg:"v5"`
}

// Filter defines the filtering rules for a FilteredAdapter's policy. Empty values
// are ignored, but all others must match the filter.
type Filter struct {
	PType []string
	V0    []string
	V1    []string
	V2    []string
	V3    []string
	V4    []string
	V5    []string
}

type adapter struct {
	db         *pg.DB
	isFiltered bool
}

var _ persist.FilteredAdapter = (*adapter)(nil)

// NewAdapter is the constructor for Adapter.
func NewAdapter(db *pg.DB) persist.FilteredAdapter {
	return &adapter{
		db: db,
	}
}

// LoadPolicy loads policy from database.
func (a *adapter) LoadPolicy(model model.Model) error {
	a.isFiltered = false
	var lines []*CasbinRule

	if _, err := a.db.Query(&lines, `SELECT * FROM casbin_rules`); err != nil {
		return err
	}

	for _, line := range lines {
		loadPolicyLine(line, model)
	}

	return nil
}

// SavePolicy saves policy to database.
func (a *adapter) SavePolicy(model model.Model) error {
	err := a.dropTable()
	if err != nil {
		return err
	}
	err = a.createTable()
	if err != nil {
		return err
	}

	var lines []*CasbinRule

	for ptype, ast := range model["p"] {
		for _, rule := range ast.Policy {
			line := a.savePolicyLine(ptype, rule)
			lines = append(lines, line)
		}
	}

	for ptype, ast := range model["g"] {
		for _, rule := range ast.Policy {
			line := a.savePolicyLine(ptype, rule)
			lines = append(lines, line)
		}
	}

	err = a.db.Insert(&lines)
	return err
}

// AddPolicy adds a policy rule to the storage.
func (a *adapter) AddPolicy(sec string, ptype string, rule []string) error {
	line := a.savePolicyLine(ptype, rule)
	err := a.db.Insert(line)
	return err
}

// RemovePolicy removes a policy rule from the storage.
func (a *adapter) RemovePolicy(sec string, ptype string, rule []string) error {
	line := a.savePolicyLine(ptype, rule)
	err := a.rawDelete(line)
	return err
}

// RemoveFilteredPolicy removes policy rules that match the filter from the storage.
func (a *adapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	line := &CasbinRule{PType: ptype}

	idx := fieldIndex + len(fieldValues)
	if fieldIndex <= 0 && idx > 0 {
		line.V0 = fieldValues[0-fieldIndex]
	}
	if fieldIndex <= 1 && idx > 1 {
		line.V1 = fieldValues[1-fieldIndex]
	}
	if fieldIndex <= 2 && idx > 2 {
		line.V2 = fieldValues[2-fieldIndex]
	}
	if fieldIndex <= 3 && idx > 3 {
		line.V3 = fieldValues[3-fieldIndex]
	}
	if fieldIndex <= 4 && idx > 4 {
		line.V4 = fieldValues[4-fieldIndex]
	}
	if fieldIndex <= 5 && idx > 5 {
		line.V5 = fieldValues[5-fieldIndex]
	}

	err := a.rawDelete(line)
	return err
}

func (a *adapter) savePolicyLine(ptype string, rule []string) *CasbinRule {
	line := &CasbinRule{PType: ptype}

	l := len(rule)
	if l > 0 {
		line.V0 = rule[0]
	}
	if l > 1 {
		line.V1 = rule[1]
	}
	if l > 2 {
		line.V2 = rule[2]
	}
	if l > 3 {
		line.V3 = rule[3]
	}
	if l > 4 {
		line.V4 = rule[4]
	}
	if l > 5 {
		line.V5 = rule[5]
	}

	return line
}

func loadPolicyLine(line *CasbinRule, model model.Model) {
	const prefixLine = ", "
	var sb strings.Builder

	sb.WriteString(line.PType)
	if len(line.V0) > 0 {
		sb.WriteString(prefixLine)
		sb.WriteString(line.V0)
	}
	if len(line.V1) > 0 {
		sb.WriteString(prefixLine)
		sb.WriteString(line.V1)
	}
	if len(line.V2) > 0 {
		sb.WriteString(prefixLine)
		sb.WriteString(line.V2)
	}
	if len(line.V3) > 0 {
		sb.WriteString(prefixLine)
		sb.WriteString(line.V3)
	}
	if len(line.V4) > 0 {
		sb.WriteString(prefixLine)
		sb.WriteString(line.V4)
	}
	if len(line.V5) > 0 {
		sb.WriteString(prefixLine)
		sb.WriteString(line.V5)
	}

	persist.LoadPolicyLine(sb.String(), model)
}

func (a *adapter) createTable() error {
	err := a.db.CreateTable((*CasbinRule)(nil), &orm.CreateTableOptions{
		Temp: false,
	})
	if err != nil {
		errorCode := err.Error()[0:12]
		if errorCode != tableExistsErrorCode {
			return err
		}
	}
	return nil
}

func (a *adapter) dropTable() error {
	err := a.db.DropTable((*CasbinRule)(nil), &orm.DropTableOptions{})
	if err != nil {
		return err
	}

	return nil
}

func (a *adapter) rawDelete(line *CasbinRule) (err error) {
	queryArgs := []interface{}{line.PType}
	query := fmt.Sprintf("DELETE FROM %s WHERE p_type = ?", "casbin_rules")
	if line.V0 != "" {
		query += " AND v0 = ?"
		queryArgs = append(queryArgs, line.V0)
	}
	if line.V1 != "" {
		query += " AND v1 = ?"
		queryArgs = append(queryArgs, line.V1)
	}
	if line.V2 != "" {
		query += " AND v2 = ?"
		queryArgs = append(queryArgs, line.V2)
	}
	if line.V3 != "" {
		query += " AND v3 = ?"
		queryArgs = append(queryArgs, line.V3)
	}
	if line.V4 != "" {
		query += " AND v4 = ?"
		queryArgs = append(queryArgs, line.V4)
	}
	if line.V5 != "" {
		query += " AND v5 = ?"
		queryArgs = append(queryArgs, line.V5)
	}
	_, err = a.db.Exec(query, queryArgs...)
	if err != nil {
		return
	}
	return
}

// IsFiltered returns true if the loaded policy has been filtered.
func (a *adapter) IsFiltered() bool {
	return a.isFiltered
}

// LoadFilteredPolicy loads only policy rules that match the filter.
func (a *adapter) LoadFilteredPolicy(model model.Model, filter interface{}) error {
	if filter == nil {
		return a.LoadPolicy(model)
	}
	var lines []CasbinRule

	filterValue, ok := filter.(*Filter)
	if !ok {
		return errors.New("invalid filter type")
	}

	err := a.db.Model(&lines).WhereGroup(a.filterQuery(filterValue)).Select()

	if err != nil {
		return err
	}

	for _, line := range lines {
		loadPolicyLine(&line, model)
	}
	a.isFiltered = true

	return nil
}

// filterQuery builds the query to match the rule filter to use within a scope.
func (a *adapter) filterQuery(filter *Filter) func(q *orm.Query) (*orm.Query, error) {
	return func(q *orm.Query) (*orm.Query, error) {
		if len(filter.PType) > 0 {
			q = q.Where("p_type in (?)", pg.In(filter.PType))
		}
		if len(filter.V0) > 0 {
			q = q.Where("v0 in (?)", pg.In(filter.V0))
		}
		if len(filter.V1) > 0 {
			q = q.Where("v1 in (?)", pg.In(filter.V1))
		}
		if len(filter.V2) > 0 {
			q = q.Where("v2 in (?)", pg.In(filter.V2))
		}
		if len(filter.V3) > 0 {
			q = q.Where("v3 in (?)", pg.In(filter.V3))
		}
		if len(filter.V4) > 0 {
			q = q.Where("v4 in (?)", pg.In(filter.V4))
		}
		if len(filter.V5) > 0 {
			q = q.Where("v5 in (?)", pg.In(filter.V5))
		}
		return q, nil
	}
}
