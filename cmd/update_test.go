package cmd

import  (
	"fmt"
	"testing"
	"github.com/nicolasmanic/tefter/model"
	"github.com/nicolasmanic/tefter/repository"
)

func TestUpdate(t *testing.T){
	oldNotebookDB := NotebookDB
	oldNoteDB := NoteDB
	NoteDB = mockNoteDBUpdate{}
	NotebookDB = mockNotebookDBUpdate{}
	//Restore interface
	defer func() {
		NotebookDB = oldNotebookDB
		NoteDB = oldNoteDB
	}()
	mockEditor := func(text string) string{
		return text
	}
	update(1, "NewTitle", []string{}, "notebook", mockEditor)
}

type mockNoteDBUpdate struct{
	repository.NoteRepository
}

func (mDB mockNoteDBUpdate) GetNote(noteID int64) (*model.Note, error){
	note:= model.NewNote("testTitle4", "testMemo", DEFAULT_NOTEBOOK_ID, []string{})
	note.ID = 2
	return note, nil
}

func (mDB mockNoteDBUpdate) UpdateNote(note *model.Note) error{
	if note.Title != "NewTitle" {
		return fmt.Errorf("Failed to update Note, expected title: %v, got %v", "NewTitle", note.Title)
	}
	return nil
}

type mockNotebookDBUpdate struct{
	repository.NotebookRepository
}

func (mDB mockNotebookDBUpdate) GetNotebookByTitle(notebooksTitle string) (*model.Notebook, error){
	notebook := model.NewNotebook(notebooksTitle)
	notebook.ID = 1
	note:= model.NewNote("testTitle", "testMemo", notebook.ID, []string{})
	note.ID = 1
	notebook.AddNote(note)
	return notebook, nil
}

