//Package eggshell lightweight, document based database
package eggshell

import (
	"sync"
)

//Driver provides basic catapabilities required to handle a database,
type Driver struct {
	mutex   sync.Mutex
	mutexes map[string]*sync.Mutex
	dir     string // the directory where scribble will create the database
}
