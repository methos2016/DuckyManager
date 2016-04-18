package main

import (
	"encoding/json"
	"os"
)

// Save is used when writing the database or the config file to the disk
func Save(path string, toSave interface{}) (err error) {
	b, err := json.Marshal(toSave)
	if err != nil {
		return
	}

	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0777)
	if err != nil {
		return
	}
	defer func() { err = f.Close() }()

	if _, err = f.Write(b); err != nil {
		return
	}

	return
}
