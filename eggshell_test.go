package eggshell

import (
	"os"
	"testing"
)

type Cat struct {
	Name string
	Age  int
}

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}

func TestCreateDriver(t *testing.T) {
	path := "db"
	db, err := CreateDriver(path)
	if err != nil {
		t.Error("expected no error while creating driver, got: ", err)
	}
	cat := Cat{"topak", 5}
	err = InsertDocument(db, "cats", cat)
	if err != nil {
		t.Error("expected no error while inserting document, got: ", err)
	}
}
