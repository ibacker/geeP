package clause

import (
	"reflect"
	"testing"
)

func testSelect(t *testing.T) {
	var cluase Clause
	cluase.Set(LIMIT, 3)
	cluase.Set(SELECT, "User", []string{"*"})
	cluase.Set(WHERE, "Name like ?", "Tom")
	cluase.Set(ORDERBY, "age desc")

	sql, vars := cluase.Build(SELECT, WHERE, ORDERBY, LIMIT)
	t.Log(sql, vars)

	if sql != "SELECT * FROM User WHERE Name like ? ORDER BY age desc LIMIT ?" {
		t.Fatal("failed to build SQL")
	}
	if !reflect.DeepEqual(vars, []interface{}{"Tom", 3}) {
		t.Fatal("failed to build SQLVars")
	}
}

func TestSelect(t *testing.T) {
	t.Run("select", func(t *testing.T) {
		testSelect(t)
	})
}

func testInsert(t *testing.T) {
	var clause Clause
	clause.Set(INSERT, "User", []string{"name", "age"})
	clause.Set(VALUES, []interface{}{"KK", 22}, []interface{}{"BB", 18})
	sql, vars := clause.Build(INSERT, VALUES)
	t.Log(sql, vars)
}
func TestInsert(t *testing.T) {
	t.Run("insert", func(t *testing.T) {
		testInsert(t)
	})
}
