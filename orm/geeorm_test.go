package orm

import (
	_ "github.com/mattn/go-sqlite3"
	"orm/log"
	"testing"
)

func TestNestedGroup(t *testing.T) {
	engine, _ := NewEngine("sqlite3", "gee.db")
	defer engine.Close()

	s := engine.NewSession()
	_, _ = s.Raw("drop table if exists user").Exec()
	_, _ = s.Raw("create table user(name text);").Exec()
	result, _ := s.Raw("insert into user(name) values (?),(?)", "Tom", "same").Exec()
	count, _ := result.RowsAffected()
	log.Info("exec success count: ", count)
}
