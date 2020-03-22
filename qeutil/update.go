package qeutil

import (
	"fmt"

	sq "github.com/Masterminds/squirrel"
)

// UpdateClause .
type UpdateClause struct {
	Update string
	Set    map[string]interface{}
	Where  []Wh
}

// SQLStm return a MySQL query statment from the InsertClause.
func (uc *UpdateClause) SQLStm() (string, []interface{}, error) {
	builder := sq.Update(uc.Update)
	for k, v := range uc.Set {
		builder = builder.Set(k, v)
	}
	for i := range uc.Where {
		builder = builder.Where(uc.Where[i].ToWhBuilder())
	}
	return builder.ToSql()
}

// ToUnlinks return an array of wh's ToStr() function result which can be used to unlink keys in redis.
func (uc *UpdateClause) ToUnlinks() []string {
	whsStr := make([]string, len(uc.Where))
	for i := range uc.Where {
		whsStr[i] = fmt.Sprintf("%v:*[:&]%v", uc.Update, uc.Where[i].ToStr())
	}
	return whsStr
}
