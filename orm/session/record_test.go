package session

import (
	"orm/log"
	"testing"
)

//type User struct {
//	Name string `geeorm:"PRIMARY KEY"`
//	Age  int
//}
//
//func (User) TableName() string {
//	return "UserNm"
//}

var (
	user1 = &User{"Tom", 14}
	user2 = &User{"Jack", 13}
	user3 = &User{"Jerry", 15}
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

func TestRecord(t *testing.T) {
	s := testRecordInit(t)
	affected, err := s.Insert(user3)
	if err != nil || affected != 1 {
		t.Fatalf("failed to insert record")
	}
	var users []User
	s.Find(&users)
	log.Info(users)
}
