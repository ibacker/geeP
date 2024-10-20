package orm

import (
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"orm/log"
	"orm/session"
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

func OpenDB(t *testing.T) *Engine {
	t.Helper()
	engine, err := NewEngine("sqlite3", "gee.db")
	if err != nil {
		t.Fatal("field to connect err: ", err)
	}
	return engine
}

func TestNewEngine(t *testing.T) {
	engine := OpenDB(t)
	defer engine.Close()
}

type User struct {
	Name string `geeorm:"PRIMARY KEY"`
	Age  int
}

func TestEngine_Transaction_Rollback(t *testing.T) {
	t.Run("rollback", func(t *testing.T) {
		engine, _ := NewEngine("sqlite3", "gee.db")
		defer engine.Close()
		s := engine.NewSession()

		s.Model(&User{}).DropTable()
		_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {

			s.Model(&User{}).CreateTable()
			s.Insert(&User{"Tom", 11})
			return nil, errors.New("Throw ERROR")
		})
		if err == nil || s.HasTable() {
			t.Fatal("failed to rollback")
		}
	})
}
func TestEngine_Transaction_Commit(t *testing.T) {
	t.Run("rollback", func(t *testing.T) {
		engine, _ := NewEngine("sqlite3", "gee.db")
		defer engine.Close()
		s := engine.NewSession()

		s.Model(&User{}).DropTable()
		_, err := engine.Transaction(func(s *session.Session) (result interface{}, err error) {

			s.Model(&User{}).CreateTable()
			s.Insert(&User{"Tom", 11})
			return
		})
		u := &User{}
		_ = s.First(u)
		log.Info(u)
		if err != nil || u.Name != "Tom" {
			t.Fatal("failed to commit")
		}
	})
}
