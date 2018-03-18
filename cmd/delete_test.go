package cmd

import (
	"github.com/nicolasmanic/tefter/repository"
	"testing"
)

func TestDelete(t *testing.T) {
	oldNoteDB := NoteDB
	NoteDB = mockNoteDBDelete{}
	//Restore interface
	defer func() {
		NoteDB = oldNoteDB
	}()

	delete([]int{1, 2, 3}, []string{})
	//empty id slice should not cause any problem
	delete([]int{}, []string{})
}

type mockNoteDBDelete struct {
	repository.NoteRepository
}

func (mDB mockNoteDBDelete) DeleteNotes(noteIDs []int64) error {
	return nil
}
