package eggshell

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

type Cat struct {
	Name string
	Age  int
}

func TestMain(m *testing.M) {
	os.RemoveAll("./testdb")
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
	for i := 0; i < 10; i++ {
		cat := Cat{"topak", i}
		err := db.InsertDocument("cats", cat)
		if err != nil {
			t.Error("expected no error while inserting document, got: ", err)
		}
	}
	documents, err := db.ReadAll("cats")
	if err != nil {
		t.Error("expected no error while reading all, got: ", err)
	}
	for i, document := range documents {
		parsedDoc := Cat{}
		if err := json.Unmarshal([]byte(document), &parsedDoc); err != nil {
			t.Error("expected no error while unmarshaling, got: ", err)
		}
		cat := Cat{"topak", i}
		if parsedDoc != cat {
			t.Error("expected parsed cat to be the same as default cat")
		}
	}
}

func TestInsertAndReadFiltered(t *testing.T) {
	db, _ := createDriver()
	documents, err := db.ReadFiltered("cats", "Age", "5")
	if err != nil {
		t.Error("expected no error while reading all, got: ", err)
	}
	for _, document := range documents {
		expectedCat := Cat{"topak", 5}
		parsedDoc := Cat{}
		if err := json.Unmarshal([]byte(document), &parsedDoc); err != nil {
			t.Error("expected no error while unmarshaling, got: ", err)
		}
		if parsedDoc != expectedCat {
			t.Error("expected parsed cat to be the same as default cat")
		}
		if len(documents) != 1 {
			t.Error("expected 1 cat aged 5, found: ", len(documents))
		}
	}
}

func TestCollectionDelete(t *testing.T) {
	db, _ := createDriver()
	cat := Cat{"topak", 5}
	err := db.InsertDocument("tempcats", cat)
	if err != nil {
		t.Error("expected no error while inserting document, got: ", err)
	}
	db.DeleteCollection("tempcats")
	_, err = os.Stat(db.GetCollectionPath("tempcats"))
	if err == nil {
		t.Error("tempcats folder should have been deleted but is not")
	}
}

func TestGetAllCollections(t *testing.T) {
	db, _ := createDriver()
	cat := Cat{"topak", 5}
	err := db.InsertDocument("optional", cat)
	if err != nil {
		t.Error("expected no error while inserting document, got: ", err)
	}
	//should include cats and optional
	collections := db.GetAllCollections()
	collectionsJoined := strings.Join(collections, " ")
	if !strings.Contains(collectionsJoined, "optional") {
		t.Error("expected optional collection, but it wasnt found")
	}
}

func TestDeleteFiltered(t *testing.T) {
	db, _ := createDriver()
	cat := Cat{"topak", 167}
	err := db.InsertDocument("deletefilter", cat)
	if err != nil {
		t.Error("expected no error while inserting document, got: ", err)
	}
	cat = Cat{"topak", 168}
	err = db.InsertDocument("deletefilter", cat)
	if err != nil {
		t.Error("expected no error while inserting document, got: ", err)
	}
	err = db.DeleteFiltered("deletefilter", "Age", "167")
	if err != nil {
		t.Error("expected no error while deleting by filter, got: ", err)
	}
}

func TestInsertAll(t *testing.T) {
	db, _ := createDriver()
	var cats []interface{}
	for i := 0; i < 10; i++ {
		cat := Cat{"topak", i}
		cats = append(cats, cat)
	}
	err := db.InsertAllDocuments("allcats", cats)
	if err != nil {
		t.Error("expected no error while inserting document, got: ", err)
	}
	documents, err := db.ReadAll("allcats")
	if err != nil {
		t.Error("expected no error while reading all, got: ", err)
	}
	for i, document := range documents {
		parsedDoc := Cat{}
		if err := json.Unmarshal([]byte(document), &parsedDoc); err != nil {
			t.Error("expected no error while unmarshaling, got: ", err)
		}
		cat := Cat{"topak", i}
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
