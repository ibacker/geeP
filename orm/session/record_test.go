package session

import (
	"orm/log"
	"testing"
)

type Person struct {
	Name   string `geeorm:"PRIMARY KEY"`
	Gender int
}

func (Person) TableName() string {
	return "PersonTabel"
}

var (
	user1   = &User{"Tom", 14}
	user2   = &User{"Jack", 13}
	person1 = &Person{"Jerry", 1}
)

func testRecordInit(t *testing.T) *Session {
	t.Helper()
	s := NewSession().Model(&User{})
	err1 := s.DropTable()
	err2 := s.CreateTable()
	_, err3 := s.Insert(user1)
	if err1 != nil || err2 != nil || err3 != nil {
		t.Fatalf("failed to init table")
	}
	return s
}

func TestRecordInsert2(t *testing.T) {
	t.Helper()
	s := NewSession().Model(&Person{})
	s.DropTable()
	s.CreateTable()
	s.Insert(person1)
	s.Insert(user2)
}

func TestRecord(t *testing.T) {
	s := testRecordInit(t)
	affected, err := s.Insert(user2)
	if err != nil || affected != 1 {
		t.Fatalf("failed to insert record")
	}
	var users []User
	s.Find(&users)
	log.Info(users)
}

func TestUpdate(t *testing.T) {
	t.Helper()
	s := NewSession().Model(&Person{})
	s.DropTable()
	s.CreateTable()
	s.Insert(person1)
	s.Update("Gender", 10)
	var persons []Person
	s.Find(&persons)
	log.Info(persons)
	s.Delete()
	count, _ := s.Count()
	log.Info("count", count)
}
