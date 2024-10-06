package dialect

import "reflect"

// dialect 解耦不同数据库 和 应用间 类型定义差异

var dialectsMap = map[string]Dialect{}

type Dialect interface {
	// 将 go 语言的类型转换为数据库的数据类型
	DataTypeOf(typ reflect.Value) string
	//判断某个表是否存在
	TableExistSQL(tableName string) (string, []interface{})
}

// RegisterDialect 注册
func RegisterDialect(name string, dialect Dialect) {
	dialectsMap[name] = dialect
}

// GetDialect 获取
func GetDialect(name string) (dialect Dialect, ok bool) {
	dialect, ok = dialectsMap[name]
	return
}
