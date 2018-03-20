package cmd

import (
	"github.com/nicolasmanic/tefter/repository"
	"testing"
)

func TestDeleteWrapper(t *testing.T) {
	oldNoteDB := NoteDB
	NoteDB = mockNoteDBDelete{}
	//Restore interface
	defer func() {
		NoteDB = oldNoteDB
	}()

	deleteWrapper(nil,[]string{"1", "2"})
	//empty id slice should not cause any problem
	deleteWrapper(nil,[]string{})
}

type mockNoteDBDelete struct {
	repository.NoteRepository
}

func (mDB mockNoteDBDelete) DeleteNotes(noteIDs []int64) error {
	return nil
}
