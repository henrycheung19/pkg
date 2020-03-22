package qeutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDeleteClauseSQLStm(t *testing.T) {
	dc := DeleteClause{
		From: "table",
		Where: []Wh{
			Wh{"in", map[string]interface{}{"in_comp": []string{"hello", "world", "!"}}},
			Wh{">", map[string]interface{}{"gt_comp": 1}},
			Wh{"<", map[string]interface{}{"lt_comp": 2}},
		},
	}
	stm, val, _ := dc.SQLStm()
	assert.Equal(t, "DELETE FROM table WHERE in_comp IN (?,?,?) AND gt_comp > ? AND lt_comp < ?", stm)
	assert.Equal(t, []interface{}{"hello", "world", "!", 1, 2}, val)
}
