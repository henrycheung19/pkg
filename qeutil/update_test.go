package qeutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUpdateClauseSQLStm(t *testing.T) {
	uc := UpdateClause{
		Update: "table",
		Set: map[string]interface{}{
			"field1": 1,
			"field2": "2",
		},
		Where: []Wh{
			Wh{"in", map[string]interface{}{"in_comp": []string{"hello", "world", "!"}}},
			Wh{">", map[string]interface{}{"gt_comp": 1}},
			Wh{"<", map[string]interface{}{"lt_comp": 2}},
		},
	}
	stm, val, _ := uc.SQLStm()
	assert.Equal(t, "UPDATE table SET field1 = ?, field2 = ? WHERE in_comp IN (?,?,?) AND gt_comp > ? AND lt_comp < ?", stm)
	assert.Equal(t, []interface{}{1, "2", "hello", "world", "!", 1, 2}, val)
}
