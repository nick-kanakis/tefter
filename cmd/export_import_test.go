package cmd

import (
	"github.com/nicolasmanic/tefter/model"
	"github.com/nicolasmanic/tefter/repository"
	"os"
	"testing"
)

func TestExportImport(t *testing.T) {
	oldNotebookDB := NotebookDB
	oldNoteDB := NoteDB
	NoteDB = mockNoteDBExportImport{}
	NotebookDB = mockNotebookDBExportImport{}

	defer func() {
		NotebookDB = oldNotebookDB
		NoteDB = oldNoteDB
		os.Remove("notes.json")
	}()
	jsonNotes, err:=retrieveJSONNotes([]int{1}, []string{"test"}, []string{"test"}, false)
	if err!= nil{
		t.Errorf("retrieveJSONNotes failed, error msg: %v", err)
	}
	export2File(jsonNotes)
	importNotes("notes.json")
}

func TestImportNoArguments(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Calling import notes with no arguments should cause panic")
		}
	}()
	importNotesWrapper(nil, []string{})
}

type mockNotebookDBExportImport struct {
	repository.NotebookRepository
}

type mockNoteDBExportImport struct {
	repository.NoteRepository
}

func (mDB mockNotebookDBExportImport) GetNotebookByTitle(notebooksTitle string) (*model.Notebook, error) {
	notebook := model.NewNotebook(notebooksTitle)
	notebook.ID = 1
	note := model.NewNote("testTitle", "testMemo", notebook.ID, []string{})
	note.ID = 1
	notebook.AddNote(note)
	return notebook, nil
}

func (mDB mockNotebookDBExportImport) GetAllNotebooksTitle() (map[int64]string, error) {
	return map[int64]string{1: "testTitle"}, nil
}

func (mDB mockNoteDBExportImport) GetNotesByTag(tags []string) ([]*model.Note, error) {
	note1 := model.NewNote("testTitle", "testMemo", repository.DEFAULT_NOTEBOOK_ID, []string{})
	note1.ID = 1
	note2 := model.NewNote("testTitle2", "testMemo2", repository.DEFAULT_NOTEBOOK_ID, []string{})
	note2.ID = 2
	return []*model.Note{note1, note2}, nil
}

func (mDB mockNoteDBExportImport) SaveNote(note *model.Note) (int64, error) {
	return 1, nil
}

func (mDB mockNoteDBExportImport) GetNotes(noteIDs []int64) ([]*model.Note, error) {
	note := model.NewNote("testTitle4", "testMemo", repository.DEFAULT_NOTEBOOK_ID, []string{})
	note.ID = 4
	return []*model.Note{note}, nil
}
