package model

import (
	"testing"
	"time"
)

func TestUpdateTitle(t *testing.T) {
	note := createMockedNote()
	note.UpdateTitle("NewTitle")

	if note.Title != "NewTitle" {
		t.Error("Failed updating title")
	}
}

func TestUpdateMemo(t *testing.T) {
	note := createMockedNote()
	note.UpdateMemo("NewMemo")

	if note.Memo != "NewMemo" {
		t.Error("Failed updating memo")
	}
}

func TestAddTags(t *testing.T) {
	note := createMockedNote()
	note.AddTags([]string{"tag3"})

	if len(note.Tags) != 3 {
		t.Error("Failed adding tag")
	}
}

func TestRemoveTags(t *testing.T) {
	note := createMockedNote()
	note.RemoveTags([]string{"tag2"})

	if len(note.Tags) != 1 {
		t.Error("Failed adding tag")
	}
}

func TestUpdateTags(t *testing.T) {
	note := createMockedNote()
	note.UpdateTags([]string{"tag3"})

	if !note.Tags["tag3"] {
		t.Error("Failed updating tags")
	}
}
func TestUpdateNotebook(t *testing.T) {
	note := createMockedNote()
	note.UpdateNotebook(123)

	if note.NotebookID != 123 {
		t.Error("Failed updating notebook")
	}
}

func TestNoteUpdated(t *testing.T) {
	note := createMockedNote()
	originalTime := note.LastUpdated
	//allow some time to pass
	time.Sleep(100)
	noteUpdated(note)

	if !originalTime.Before(note.LastUpdated) {
		t.Error("Failed updating `LastUpdated` value")
	}
}

func createMockedNote() *Note {
	return NewNote("testTitle", "testMemo", 111, []string{"tag1", "tag2"})
}
