//Package eggshell is lightweight, document based database
package eggshell

import (
	"bufio"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
)

var breakLineBytes []byte = []byte("\n")

//Driver provides basic catapabilities required to handle a database
type Driver struct {
	mutex   sync.Mutex
	mutexes map[string]*sync.Mutex
	Path    string //path which the driver uses to store database files
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
		Path:    path,
		mutexes: make(map[string]*sync.Mutex)}
	return &driver, returnError
}

//InsertDocument inserts document into a collection, given a collection name
//document needs to be a struct which can be marshaled by the json package
func (db *Driver) InsertDocument(collection string, document interface{}) error {
	collectionPath := appendFilePath(db.Path, collection+".json")

	mutex := db.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

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
	if err := writeToFile(f, contentToWrite); err != nil {
		return err
	}
	return nil
}

//InsertAllDocuments inserts all documents in the array into a collection, given a collection name
//documents needs to be an array of type interface and not a struct
func (db *Driver) InsertAllDocuments(collection string, documents []interface{}) error {
	collectionPath := appendFilePath(db.Path, collection+".json")

	mutex := db.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	f, err := os.OpenFile(collectionPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer f.Close()
	var contentToWrite []byte
	for _, document := range documents {
		marshaledDocument, err := json.Marshal(document)
		if err != nil {
			return err
		}
		marshaledDocument = append(marshaledDocument, breakLineBytes...)
		contentToWrite = append(contentToWrite, marshaledDocument...)
	}
	if err := writeToFile(f, contentToWrite); err != nil {
		return err
	}
	return nil
}

//ReadAll reads all documents in a collection
func (db *Driver) ReadAll(collection string) (documents []string, err error) {
	collectionPath := appendFilePath(db.Path, collection+".json")
	collectionFile, err := os.OpenFile(collectionPath, os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}
	defer collectionFile.Close()

	rawDocuments := []string{}

	scanner := bufio.NewScanner(collectionFile)
	for scanner.Scan() {
		rawDocuments = append(rawDocuments, scanner.Text())
	}

	return rawDocuments, nil

}

//ReadFiltered reads documents that match the given filter in a collection
//note that the filter is case sensitive
func (db *Driver) ReadFiltered(collection string, filterKeys, filterValues []string) (documents []string, err error) {
	collectionPath := appendFilePath(db.Path, collection+".json")
	collectionFile, err := os.OpenFile(collectionPath, os.O_RDONLY, 0777)
	if err != nil {
		return nil, err
	}
	defer collectionFile.Close()

	rawDocuments := []string{}

	scanner := bufio.NewScanner(collectionFile)
	for scanner.Scan() {
		line := scanner.Text()
		matches := 0
		for index, filterKey := range filterKeys {
			regex, _ := regexp.Compile("(" + filterKey + ")\"?:\"?(" + filterValues[index] + ")")
			if regex.MatchString(line) {
				matches++
			}
		}
		if matches == len(filterKeys) {
			rawDocuments = append(rawDocuments, scanner.Text())
		}
	}

	return rawDocuments, nil

}

//DeleteFiltered deletes documents that match the given filter in a collection
//note that the filter is case sensitive
func (db *Driver) DeleteFiltered(collection string, filterKeys, filterValues []string) error {
	collectionPath := appendFilePath(db.Path, collection+".json")
	collectionFile, err := os.OpenFile(collectionPath, os.O_RDWR, 0777)
	if err != nil {
		return err
	}

	mutex := db.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	rawDocuments := []string{}

	scanner := bufio.NewScanner(collectionFile)
	for scanner.Scan() {
		line := scanner.Text()
		matches := 0
		for index, filterKey := range filterKeys {
			regex, _ := regexp.Compile("(" + filterKey + ")\"?:\"?(" + filterValues[index] + ")")
			if regex.MatchString(line) {
				matches++
			}
		}
		if matches != len(filterKeys) {
			rawDocuments = append(rawDocuments, scanner.Text())
		}
	}

	collectionFile.Close()

	err = os.Remove(collectionPath)

	collectionFile, err = os.OpenFile(collectionPath,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0777)
	if err != nil {
		return err
	}
	defer collectionFile.Close()

	if err != nil {
		return err
	}

	for _, document := range rawDocuments {
		insertDoc := []byte(document)
		insertDoc = append(insertDoc, breakLineBytes...)
		if err := writeToFile(collectionFile, insertDoc); err != nil {
			return err
		}
	}

	return nil

}

//GetAllCollections gets all collections stored in the database
func (db *Driver) GetAllCollections() []string {

	collectionList := []string{}

	files, err := ioutil.ReadDir(db.Path)
	if err != nil {
		return nil
	}
	for _, file := range files {
		if strings.HasSuffix(file.Name(), ".json") {
			collectionName := strings.Replace(file.Name(), ".json", "", 1)
			collectionList = append(collectionList, collectionName)
		}
	}

	return collectionList

}

//DeleteCollection removes a collection from the database
func (db *Driver) DeleteCollection(collection string) error {

	mutex := db.getOrCreateMutex(collection)
	mutex.Lock()
	defer mutex.Unlock()

	collectionPath := appendFilePath(db.Path, collection+".json")
	return os.Remove(collectionPath)

}

//GetCollectionPath get the path of a collection
func (db *Driver) GetCollectionPath(collection string) string {

	return appendFilePath(db.Path, collection+".json")

}

func writeToFile(file *os.File, content []byte) error {
	if _, err := file.Write(content); err != nil {
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

//thank you sdomino/scribble for this function, it will stay until a better solution is required
func (db *Driver) getOrCreateMutex(collection string) *sync.Mutex {

	db.mutex.Lock()
	defer db.mutex.Unlock()

	m, ok := db.mutexes[collection]

	// if the mutex doesn't exist make it
	if !ok {
		m = &sync.Mutex{}
		db.mutexes[collection] = m
	}

	return m
}
