package repository

import (
	"os"
	"testing"
)

func TestConnect2DB(t *testing.T) {
	db := connect2DB("testPath")

	defer func() {
		db.Close()
		os.Remove("testPath")
	}()

	if db == nil || db.Ping() != nil {
		t.Error("Could not connect to DB")
	}
}
