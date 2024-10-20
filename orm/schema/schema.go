package schema

import (
	"go/ast"
	"orm/dialect"
	"reflect"
)

// Field 代表数据库表的一列（一个字段）
type Field struct {
	// 列名
	Name string
	// 类型
	Type string
	// 注释
	Tag string
}

// Schema 代表一张数据库表
type Schema struct {
	// 被映射的对象
	Model interface{}
	// 表名
	Name string
	// 字段列表
	Fields     []*Field
	FieldNames []string
	fieldMap   map[string]*Field
}

func (schema *Schema) GetField(name string) *Field {
	return schema.fieldMap[name]
}

// ITableName 接口，表名
type ITableName interface {
	TableName() string
}

// Parse 将 dest 被映射的对象转换为 schema 实例
func Parse(dest interface{}, d dialect.Dialect) *Schema {
	modelType := reflect.Indirect(reflect.ValueOf(dest)).Type()

	var tableName string
	// 表名
	t, ok := dest.(ITableName)
	if !ok {
		// 结构体的名称作为表名
		tableName = modelType.Name()
	} else {
		tableName = t.TableName()
	}

	schema := &Schema{
		Model:    dest,
		Name:     tableName,
		fieldMap: make(map[string]*Field),
	}

	for i := 0; i < modelType.NumField(); i++ {
		p := modelType.Field(i)
		if !p.Anonymous && ast.IsExported(p.Name) {
			field := &Field{
				Name: p.Name,
				Type: d.DataTypeOf(reflect.Indirect(reflect.New(p.Type))),
			}

			if v, ok := p.Tag.Lookup("orm"); ok {
				field.Tag = v
			}

			schema.Fields = append(schema.Fields, field)
			schema.FieldNames = append(schema.FieldNames, field.Name)
			schema.fieldMap[field.Name] = field
		}
	}
	return schema
}

// RecordValues 通过反射  取对象值
func (schema *Schema) RecordValues(dest interface{}) []interface{} {
	destValue := reflect.Indirect(reflect.ValueOf(dest))
	var fieldValues []interface{}
	for _, field := range schema.Fields {
		fieldValues = append(fieldValues, destValue.FieldByName(field.Name).Interface())
	}
	return fieldValues
}
