package cmd

import (
	"github.com/nicolasmanic/tefter/model"
	"github.com/nicolasmanic/tefter/repository"
	"testing"
)

func TestSearch(t *testing.T) {
	oldNotebookDB := NotebookDB
	oldNoteDB := NoteDB
	NoteDB = mockNoteDBSearch{}
	NotebookDB = mockNotebookDBSearch{}
	originalPrintFunc := print
	print = func(notes []*jsonNote) {
		if len(notes) != 1 || notes[0].ID != 2 {
			t.Error("SearchNotesByKeyword was not called correctly")
		}
	}
	//Restore interface
	defer func() {
		NotebookDB = oldNotebookDB
		NoteDB = oldNoteDB
		print = originalPrintFunc
	}()

	searchWrapper(nil, []string{"keyword"})
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

type mockNotebookDBSearch struct {
	repository.NotebookRepository
}

func (mDB mockNotebookDBSearch) GetAllNotebooksTitle() (map[int64]string, error) {
	return map[int64]string{1: "testTitle"}, nil
}
