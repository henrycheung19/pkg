package qeutil

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSelectClauseSQLStm(t *testing.T) {
	offset := 50
	limit := 10
	sc := SelectClause{
		Select: []string{"id"},
		From:   "table",
		Where: []Wh{
			Wh{"in", map[string]interface{}{"in_comp": []string{"hello", "world", "!"}}},
			Wh{">", map[string]interface{}{"gt_comp": 1}},
			Wh{"<", map[string]interface{}{"lt_comp": 2}},
		},
		GroupBy: []string{"grp"},
		Having:  "1<>0",
		OrderBy: []string{"ord"},
		Limit:   &limit,
		Offset:  &offset,
	}
	stm, val, _ := sc.SQLStm()
	assert.Equal(t, "SELECT id FROM table WHERE in_comp IN (?,?,?) AND gt_comp > ? AND lt_comp < ? GROUP BY grp HAVING 1<>0 ORDER BY ord LIMIT 10 OFFSET 50", stm)
	assert.Equal(t, []interface{}{"hello", "world", "!", 1, 2}, val)
}

func TestSelectClauseCacheKey(t *testing.T) {
	offset := 50
	limit := 10
	sc := SelectClause{
		Select: []string{"id"},
		From:   "table",
		Where: []Wh{
			Wh{"in", map[string]interface{}{"in_comp": []string{"hello", "world", "!"}}},
			Wh{">", map[string]interface{}{"gt_comp": 1}},
			Wh{"<", map[string]interface{}{"lt_comp": 2}},
		},
		GroupBy: []string{"grp"},
		Having:  "1<>0",
		OrderBy: []string{"ord"},
		Limit:   &limit,
		Offset:  &offset,
	}
	assert.Equal(t, "table:where:in_comp=[hello,world,!]&gt_comp>1&lt_comp<2:grp:grp:hav:1<>0:ord:ord:lim:10:off:50", sc.CacheKey())
}
