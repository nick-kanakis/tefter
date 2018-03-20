package cmd

import (
	"testing"
	"github.com/nicolasmanic/tefter/model"
	"github.com/nicolasmanic/tefter/repository"
)

func TestSearch(t *testing.T) {
	oldNoteDB := NoteDB
	NoteDB = mockNoteDBSearch{}
	//Restore interface
	defer func() {
		NoteDB = oldNoteDB
	}()
	mockOutput := func(notes []*model.Note) {
		if len(notes) != 1 || notes[0].ID != 2 {
			t.Error("SearchNotesByKeyword was not called correctly")
		}
	}
	search("keyword", mockOutput)
}

type mockNoteDBSearch struct {
	repository.NoteRepository
}

func (mDB mockNoteDBSearch) GetNotes(noteIDs []int64) ([]*model.Note, error) {
	note := model.NewNote("testTitle1", "testMemo1", repository.DEFAULT_NOTEBOOK_ID, []string{})
	note.ID = 1
	return []*model.Note{note}, nil
}

func (mDB mockNoteDBSearch) SearchNotesByKeyword(keyword string) ([]*model.Note, error) {
	note := model.NewNote("testTitle2", "testMemo2", repository.DEFAULT_NOTEBOOK_ID, []string{})
	note.ID = 2
	return []*model.Note{note}, nil
}
