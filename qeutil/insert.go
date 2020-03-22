package qeutil

import (
	sq "github.com/Masterminds/squirrel"
)

// InsertClause .
type InsertClause struct {
	Into   string
	Values map[string]interface{}
}

// SQLStm return a MySQL query statment from the InsertClause.
func (ic *InsertClause) SQLStm() (string, []interface{}, error) {
	var values []interface{}
	builder := sq.Insert(ic.Into)

	for k, v := range ic.Values {
		builder = builder.Columns(k)
		values = append(values, v)
	}
	builder = builder.Values(values...)
	return builder.ToSql()
}
