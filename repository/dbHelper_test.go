package repository

import (
	"os"
	"testing"
)

func TestConnect2DB(t *testing.T) {
	db := connect2DB("test.db")

	defer func() {
		db.Close()
		os.Remove("test.db")
	}()

	if db == nil || db.Ping() != nil {
		t.Error("Could not connect to DB")
	}
}
