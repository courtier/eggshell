//Package eggshell lightweight, document based database
package eggshell

import (
	"path/filepath"
	"sync"
)

//Driver provides basic catapabilities required to handle a database,
type Driver struct {
	mutex   sync.Mutex
	mutexes map[string]*sync.Mutex
	dir     string // the directory where scribble will create the database
}

//CreateDriver creates a new driver that uses the path specified as the database's main directory
func CreateDriver() (*Driver, error) {
	return
}

func cleanFilePath(filePath string) string {
	return filepath.FromSlash(filePath)
}
