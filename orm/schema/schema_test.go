package schema

import (
	"orm/dialect"
	"testing"
)

type User struct {
	Name string `orm:"primary key"`
	Age  int
}

// User实现TableName 接口
func (u *User) TableName() string {
	return "userTable"
}

var TestSchemaDial, _ = dialect.GetDialect("sqlite3")

func TestSchema(t *testing.T) {
	schema := Parse(&User{}, TestSchemaDial)
	if schema.Name != "userTable" || len(schema.Fields) != 2 {
		t.Fatal("failed to parse schema")
	}
}
