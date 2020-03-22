package qeutil

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

// DeleteClause .
type DeleteClause struct {
	From  string
	Where []Wh
}

// SQLStm return a MySQL query statment from the DeleteClause.
func (dc *DeleteClause) SQLStm() (string, []interface{}, error) {
	builder := sq.Delete(dc.From)
	for i := range dc.Where {
		builder = builder.Where(dc.Where[i].ToWhBuilder())
	}
	return builder.ToSql()
}

// ToUnlinks return an array of wh's ToStr() function result which can be used to unlink keys in redis.
func (dc *DeleteClause) ToUnlinks() []string {
	whsStr := make([]string, len(dc.Where))
	for i := range dc.Where {
		whsStr[i] = fmt.Sprintf("%v:*[:&]%v", dc.From, dc.Where[i].ToStr())
	}
	return whsStr
}
