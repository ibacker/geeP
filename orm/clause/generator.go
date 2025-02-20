package clause

import (
	"fmt"
	"strings"
)

type generator func(values ...interface{}) (string, []interface{})

var generators map[Type]generator

func init() {
	generators = make(map[Type]generator)
	generators[INSERT] = _insert
	generators[SELECT] = _select
	generators[VALUES] = _values
	generators[WHERE] = _where
	generators[LIMIT] = _limit
	generators[ORDERBY] = _orderBy
	generators[UPDATE] = _update
	generators[DELETE] = _delete
	generators[COUNT] = _count
}

func genBindVars(num int) string {
	var vars []string
	for i := 0; i < num; i++ {
		vars = append(vars, "?")
	}
	return strings.Join(vars, ", ")
}

func _values(values ...interface{}) (string, []interface{}) {
	// INSERT into
	var bindStr string
	var sql strings.Builder
	var vars []interface{}

	sql.WriteString("VALUES ")
	for i, value := range values {
		v := value.([]interface{})
		if bindStr == "" {
			bindStr = genBindVars(len(v))
		}
		sql.WriteString(fmt.Sprintf("(%v)", bindStr))
		if i+1 != len(values) {
			sql.WriteString(", ")
		}
		vars = append(vars, v...)
	}
	return sql.String(), vars
}

func _insert(values ...interface{}) (string, []interface{}) {
	//INSERT INTO $tableNm ($fields)
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("INSERT INTO %s (%v)", tableName, fields), []interface{}{}
}

func _select(values ...interface{}) (string, []interface{}) {
	//SELECT $fields FROM $tableName
	tableName := values[0]
	fields := strings.Join(values[1].([]string), ",")
	return fmt.Sprintf("SELECT %v FROM %s", fields, tableName), []interface{}{}
}

func _limit(limit ...interface{}) (string, []interface{}) {
	//LIMIT $num
	return "LIMIT ?", limit
}

func _where(values ...interface{}) (string, []interface{}) {
	//WHERE $desc
	desc, vars := values[0], values[1:]
	return fmt.Sprintf("WHERE %s", desc), vars
}

func _orderBy(values ...interface{}) (string, []interface{}) {
	return fmt.Sprintf("ORDER BY %s", values[0]), []interface{}{}
}

// _update 两个参数，第一个参数为表名 第二个参数为 map 结构，表示更新的 key value
func _update(values ...interface{}) (string, []interface{}) {
	tableName := values[0]
	m := values[1].(map[string]interface{})
	var keys []string
	var vas []interface{}
	for k, v := range m {
		keys = append(keys, k+" = ?")
		vas = append(vas, v)
	}
	return fmt.Sprintf("UPDATE %s SET %s", tableName, strings.Join(keys, ", ")), vas
}

func _delete(values ...interface{}) (string, []interface{}) {
	// DELETE FROM $tableNm
	return fmt.Sprintf("DELETE FROM %s", values[0]), []interface{}{}
}

func _count(values ...interface{}) (string, []interface{}) {
	// SELECT count(*) from $tableNm
	return _select(values[0], []string{"count(*)"})
}
