//Package eggshell lightweight, document based database
package eggshell

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"sync"
)

var breakLineBytes []byte = []byte("\n")

//Driver provides basic catapabilities required to handle a database,
type Driver struct {
	mutex   sync.Mutex
	mutexes map[string]*sync.Mutex
	path    string
}

//CreateDriver creates a new driver that uses the path specified as the database's main directory
func CreateDriver(path string) (*Driver, error) {
	path = cleanFilePath(path)
	info, err := os.Stat(path)
	var returnError error = nil
	if err != nil {
		returnError = os.Mkdir(path, 0777)
	} else if !info.IsDir() {
		returnError = errors.New("path exists, but is not directory")
	}
	driver := Driver{
		path:    path,
		mutexes: make(map[string]*sync.Mutex)}
	return &driver, returnError
}

//InsertDocument inserts document into a collection
func InsertDocument(db *Driver, collection string, document interface{}) error {
	collectionPath := appendFilePath(db.path, collection+".json")
	f, err := os.OpenFile(collectionPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	contentToWrite, err := json.Marshal(document)
	if err != nil {
		return err
	}
	contentToWrite = append(contentToWrite, breakLineBytes...)
	if _, err := f.Write(contentToWrite); err != nil {
		return err
	}
	return nil
}

func appendFilePath(filepath, appendation string) string {
	return cleanFilePath(filepath + "/" + appendation)
}

func cleanFilePath(filePath string) string {
	return filepath.FromSlash(filePath)
}
