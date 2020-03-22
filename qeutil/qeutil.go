package qeutil

import (
	"bytes"
	"errors"
	"fmt"
	"strings"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
)

var (
	// ErrNotExist is returned when the required resources not exist in DB.
	ErrNotExist = errors.New("required resources not exist in DB")

	// ErrNotChanged is returned when the sql is accpeted but no row is affected..
	ErrNotChanged = errors.New("no row affected")
)

// // WhereSQL returns the SQL statement of where clause.
// func WhereSQL(w map[string][2]string) []byte {
// 	buf := bytes.Buffer{}
// 	for k, v := range w {
// 		if buf.Len() > 0 {
// 			buf.WriteString(" AND ")
// 		}
// 		buf.WriteString(k)
// 		buf.WriteString(" ")
// 		buf.WriteString(v[0])
// 		buf.WriteString(" ")
// 		buf.WriteString(v[1])
// 	}

// 	return buf.Bytes()
// }

// // CacheHelper .
// type CacheHelper struct {
// 	Redis  *rediscli.Client
// 	Object string
// 	Keys   map[string][2]string
// }

// // ClearCache remove all keys in redis that matching the given conditions.
// func (ch *CacheHelper) ClearCache() error {
// 	buf := bytes.Buffer{}
// 	redisCli := ch.Redis.Client()

// 	for k, v := range ch.Keys {
// 		// Prepare key
// 		buf.Reset()
// 		if ch.Object == "" {
// 			panic("CacheHelper object cannot be empty.")
// 		}
// 		buf.WriteString(ch.Object)
// 		buf.WriteString(":*[:&]")
// 		buf.WriteString(k)
// 		buf.WriteString(v[0])
// 		buf.WriteString(v[1])
// 		buf.WriteString("*")

// 		// Scan and unlink key
// 		iter := redisCli.Scan(0, string(bytes.ToLower(buf.Bytes())), 0).Iterator()
// 		for iter.Next() {
// 			err := redisCli.Unlink(iter.Val()).Err()
// 			if err != nil {
// 				return err
// 			}
// 		}
// 		if err := iter.Err(); err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// ExistInDB check if the required resources exists in DB.
func ExistInDB(db *sqlx.DB, target string, wheres []Wh) (bool, error) {
	builder := sq.Select("*").From(target)
	for i := range wheres {
		builder = builder.Where(wheres[i].ToWhBuilder())
	}
	stm, val, err := builder.ToSql()
	if err != nil {
		return false, err
	}

	buf := bytes.Buffer{}
	buf.WriteString("SELECT EXISTS (")
	buf.WriteString(stm)
	buf.WriteString(")")
	var exist bool
	if err := db.Get(&exist, buf.String(), val...); err != nil {
		return exist, err
	}
	return exist, nil
}

// // Code for research in future
// func (qc *QueryClause) QueryWithCache(obj interface{}, db *sqlx.DB, redis *rediscli.Client, cacheExpiry time.Duration) (interface{}, error) {
// 	result := reflect.MakeSlice(reflect.SliceOf(reflect.TypeOf(obj)), 0, 0).Interface()
// 	fmt.Println("arrtype:", reflect.TypeOf(result))
// 	fmt.Println("objtype:", reflect.TypeOf(obj))

// 	// Try to get record from redis
// 	if err := redis.Get(qc.CacheKey(), &result); err != nil {
// 		// Missed for whatever reason, try to get from DB
// 		if rows, err := db.Queryx(qc.SQLStm()); err == nil {
// 			// Returned query result, scan struct
// 			inf := []interface{}{}
// 			defer rows.Close()
// 			for rows.Next() {
// 				item := obj
// 				fmt.Println("itemtype:", reflect.TypeOf(item))
// 				// item.Elem().Set(reflect.ValueOf(obj))
// 				if err := rows.StructScan(item); err != nil {
// 					fmt.Println(err)
// 					return nil, ErrNotFoundInDB
// 				}
// 				fmt.Println("item:", item, "type:", reflect.TypeOf(item))
// 				fmt.Println("Indirect:", reflect.ValueOf(item).Elem())
// 				inf = append(inf, item)
// 				result = reflect.Append(reflect.ValueOf(result), reflect.ValueOf(reflect.ValueOf(item).Pointer()))
// 			}
// 			fmt.Println("inf:", inf)
// 		} else {
// 			return nil, ErrDBWentWrong
// 		}

// 		if err == cache.ErrCacheMiss {
// 			// Try to write result to cache
// 			// fmt.Println("result:", reflect.TypeOf(reflect.Indirect(result)))
// 			// if err := redis.Set(qc.CacheKey(), result, cacheExpiry); err != nil {
// 			// 	return nil, err
// 			// }
// 		} else {
// 			// Log unknown Redis errors
// 			return nil, err
// 		}
// 	}
// 	return result, nil
// }

// Wh contains the operator and values of a MySQL where clause.
type Wh struct {
	Operator string
	Values   map[string]interface{}
}

const (
	// Eq representing the equal operator in MySQL
	Eq string = "="
	// Gt representing the greater than operator in MySQL
	Gt string = ">"
	// Lt representing the less than operator in MySQL
	Lt string = "<"
	// GtEq representing the greater than or Equal operator in MySQL
	GtEq string = ">="
	// LtEq representing the less than or Equal operator in MySQL
	LtEq string = "<="
	// NotEq representing the not equal operator in MySQL
	NotEq string = "<>"
	// Like representing the not equal operator in MySQL
	Like string = "like"
	// In representing the in operator in MySQL
	In string = "in"
)

// ToStr returns a string format of where clause.
func (wh *Wh) ToStr() string {
	var key string
	var val interface{}
	for k, v := range wh.Values {
		key = k
		val = v
	}
	switch wh.Operator {
	case "In":
		return strings.ReplaceAll(fmt.Sprintf("%v=%v", key, val), " ", ",")
	default:
		return fmt.Sprintf("%v%v%v", key, wh.Operator, val)
	}
}

// ToWhBuilder transforms the where clause to a squirrel where pred.
func (wh *Wh) ToWhBuilder() interface{} {
	switch wh.Operator {
	case In:
		fallthrough
	case Eq:
		return wh.Values
	case Gt:
		return sq.Gt(wh.Values)
	case Lt:
		return sq.Lt(wh.Values)
	case GtEq:
		return sq.GtOrEq(wh.Values)
	case LtEq:
		return sq.LtOrEq(wh.Values)
	case NotEq:
		return sq.NotEq(wh.Values)
	case Like:
		return sq.Like(wh.Values)
	default:
		return nil
	}
}
