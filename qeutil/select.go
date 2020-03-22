package qeutil

import (
	"bytes"
	"fmt"
	"strconv"

	sq "github.com/Masterminds/squirrel"
)

// SelectClause .
type SelectClause struct {
	Select  []string
	From    string
	Where   []Wh
	GroupBy []string
	Having  string
	OrderBy []string
	Limit   *int
	Offset  *int
}

// SQLStm return a MySQL query statment from the SelectClause.
func (sc *SelectClause) SQLStm() (string, []interface{}, error) {
	if len(sc.Select) == 0 {
		sc.Select = []string{"*"}
	}
	builder := sq.Select(sc.Select...).From(sc.From)

	for i := range sc.Where {
		builder = builder.Where(sc.Where[i].ToWhBuilder())
	}

	builder = builder.GroupBy(sc.GroupBy...)
	if sc.Having != "" {
		builder = builder.Having(sc.Having)
	}
	builder = builder.OrderBy(sc.OrderBy...)
	if sc.Limit != nil {
		builder = builder.Limit(uint64(*sc.Limit))
		if sc.Offset != nil {
			builder = builder.Offset(uint64(*sc.Offset))
		}
	}
	return builder.ToSql()
}

// CacheKey return a cache key from the SelectClause.
func (sc *SelectClause) CacheKey() string {
	buf := bytes.Buffer{}

	// Main key
	if sc.From == "" {
		panic("Target table not given.")
	} else {
		buf.WriteString(sc.From)
	}

	if len(sc.Where) > 0 {
		buf.WriteString(":where:")
		whereBuf := bytes.Buffer{}
		for i := range sc.Where {
			if whereBuf.Len() > 0 {
				whereBuf.WriteString("&")
			}
			whereBuf.WriteString(sc.Where[i].ToStr())
		}
		buf.Write(whereBuf.Bytes())
	}

	if len(sc.GroupBy) > 0 {
		buf.WriteString(":grp:")
		for i := range sc.GroupBy {
			if i > 0 {
				buf.WriteString("&")
			}
			buf.WriteString(sc.GroupBy[i])
		}
	}

	if sc.Having != "" {
		buf.WriteString(":hav:")
		buf.WriteString(sc.Having)
	}

	if len(sc.OrderBy) > 0 {
		buf.WriteString(":ord:")
		for i := range sc.OrderBy {
			if i > 0 {
				buf.WriteString("&")
			}
			buf.WriteString(sc.OrderBy[i])
		}
	}

	if sc.Limit != nil {
		buf.WriteString(":lim:")
		buf.WriteString(strconv.Itoa(*sc.Limit))
		if sc.Offset != nil {
			buf.WriteString(":off:")
			buf.WriteString(strconv.Itoa(*sc.Offset))
		}
	}

	return string(bytes.ToLower(buf.Bytes()))
}

// ToUnlinks return an array of wh's ToStr() function result which can be used to unlink keys in redis.
func (sc *SelectClause) ToUnlinks() []string {
	whsStr := make([]string, len(sc.Where))
	for i := range sc.Where {
		whsStr[i] = fmt.Sprintf("%v:*[:&]%v", sc.From, sc.Where[i].ToStr())
	}
	return whsStr
}
