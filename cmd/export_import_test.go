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

	note1 := model.NewNote("testTitle", "testMemo", DEFAULT_NOTEBOOK_ID, []string{})
	note1.ID = 1
	note2 := model.NewNote("testTitle2", "testMemo2", DEFAULT_NOTEBOOK_ID, []string{})
	note2.ID = 2
	notes := []*model.Note{note1, note2}
	export2JSON(notes)
	importNotes([]string{"notes.json"})
}

func TestImportNoArguments(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("Calling import notes with no arguments should cause panic")
		}
	}()
	importNotes([]string{})
}

type mockNotebookDBExportImport struct {
	repository.NotebookRepository
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

type mockNoteDBExportImport struct {
	repository.NoteRepository
}

func (mDB mockNoteDBExportImport) SaveNote(note *model.Note) (int64, error) {
	return 1, nil
}
