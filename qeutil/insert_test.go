package qeutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInsertClauseSQLStm(t *testing.T) {
	ic := InsertClause{
		Into: "table",
		Values: map[string]interface{}{
			"field1": 1,
			"field2": "2",
		},
	}
	stm, val, _ := ic.SQLStm()
	assert.Equal(t, "INSERT INTO table (field1,field2) VALUES (?,?)", stm)
	assert.Equal(t, []interface{}{1, "2"}, val)
}
