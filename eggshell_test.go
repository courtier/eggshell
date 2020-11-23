package eggshell

import (
	"encoding/json"
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
	_, err := createDriver()
	if err != nil {
		t.Error("expected no error while inserting document, got: ", err)
	}
}

func TestInsertAndRead(t *testing.T) {
	db, _ := createDriver()
	cat := Cat{"topak", 5}
	for i := 0; i < 10; i++ {
		err := db.InsertDocument("cats", cat)
		if err != nil {
			t.Error("expected no error while inserting document, got: ", err)
		}
	}
	documents, err := db.ReadAll("cats")
	if err != nil {
		t.Error("expected no error while reading all, got: ", err)
	}
	for _, document := range documents {
		parsedDoc := Cat{}
		if err := json.Unmarshal([]byte(document), &parsedDoc); err != nil {
			t.Error("expected no error while unmarshaling, got: ", err)
		}
		if parsedDoc != cat {
			t.Error("expected parsed cat to be the same as default cat")
		}
	}
}

func createDriver() (db *Driver, err error) {
	path := "testdb"
	fDb, fErr := CreateDriver(path)
	if fErr != nil {
		return nil, fErr
	}
	return fDb, nil
}
