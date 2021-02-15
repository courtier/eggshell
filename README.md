## Eggshell [![GoDoc](https://godoc.org/github.com/courtier/eggshell?status.svg)](https://godoc.org/github.com/courtier/eggshell) [![Go Report Card](https://goreportcard.com/badge/github.com/courtier/eggshell)](https://goreportcard.com/report/github.com/courtier/eggshell) [![Eggshell](https://circleci.com/gh/courtier/eggshell.svg?style=svg)](https://circleci.com/gh/courtier/eggshell)


### Lightweight (and probably brittle) JSON based database

Inspired by scribble, aims to be better than it, at least for my usages

### To-Do (in order of priority):
- [X] Insert document
- [X] Read all documents in collection
- [X] Delete collection
- [X] Read all documents matching a parameter
- [X] Delete all documents matching a parameter
- [X] Insert all documents in an array so as to improve performance by not opening the file for each document
- [ ] Edit document
- [ ] Encryption (will slow down the processes)
- [ ] Indexing using Snowflake IDs (sorted)
- [ ] Logging?

### Features:
- [X] Supports windows (i think)

### Changelog:
- Added filtering by multiple keys and values
- Fixed couple thing