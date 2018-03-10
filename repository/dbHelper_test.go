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

func TestRemoveDups(t *testing.T) {
	input := []int64{1, 2, 3, 4, 4, 4}
	result := removeDups(input)

	if len(result) != 4 {
		t.Error("Could not remove duplicates")
	}
}
