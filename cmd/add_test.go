package cmd

import (
	"github.com/nicolasmanic/tefter/model"
	"github.com/nicolasmanic/tefter/repository"
	"testing"
	"time"
)

func TestAddNotebookToNoteDefaultNotebook(t *testing.T) {
	note := model.NewNote("title", "memo", 123, []string{})
	err := addNotebookToNote(note, "")

	if err != nil {
		t.Errorf("Error while trying to add note to default notebook, with msg: %v", err)
	}

	if note.NotebookID != repository.DEFAULT_NOTEBOOK_ID {
		t.Error("Error while trying to add note to default notebook")
	}
}

func TestAddExistingNotebook(t *testing.T) {
	oldNotebookDB := NotebookDB
	oldNoteDB := NoteDB
	NoteDB = mockNoteDBAdd{}
	NotebookDB = mockNotebookDBAdd{}
	//Restore interface
	defer func() {
		NotebookDB = oldNotebookDB
		NoteDB = oldNoteDB
	}()

	mockEditor := func(text string) string {
		return text
	}
	add("noteTitle", []string{}, "Existing Notebook", mockEditor)
}

func TestAddNewNotebook(t *testing.T) {
	oldNotebookDB := NotebookDB
	oldNoteDB := NoteDB
	NoteDB = mockNoteDBAdd{}
	NotebookDB = mockNotebookDBAdd{}
	//Restore interface
	defer func() {
		NotebookDB = oldNotebookDB
		NoteDB = oldNoteDB
	}()

	mockEditor := func(text string) string {
		return text
	}
	add("noteTitle", []string{}, "New Notebook", mockEditor)
}

func TestAddJSONNote(t *testing.T) {
	oldNotebookDB := NotebookDB
	oldNoteDB := NoteDB
	NoteDB = mockNoteDBAdd{}
	NotebookDB = mockNotebookDBAdd{}
	//Restore interface
	defer func() {
		NotebookDB = oldNotebookDB
		NoteDB = oldNoteDB
	}()

	jNote := &jsonNote{
		ID:            1,
		Title:         "Test Title",
		Memo:          "Test memo",
		Created:       time.Now(),
		LastUpdated:   time.Now(),
		Tags:          []string{"tag1", "tag2"},
		NotebookTitle: "notebook",
	}
	if err := addJSONNote(jNote); err != nil {
		t.Errorf("Failed saving JSONNote, error msg: %v", err)
	}
}

type mockNoteDBAdd struct {
	repository.NoteRepository
}

func (mDB mockNoteDBAdd) SaveNote(note *model.Note) (int64, error) {
	return 0, nil
}

type mockNotebookDBAdd struct {
	repository.NotebookRepository
}

func (mDB mockNotebookDBAdd) GetNotebookByTitle(notebookTitle string) (*model.Notebook, error) {
	if notebookTitle == "Existing Notebook" {
		notebook := model.NewNotebook(notebookTitle)
		return notebook, nil
	}
	return nil, nil
}

func (mDB mockNotebookDBAdd) SaveNotebook(notebook *model.Notebook) (int64, error) {
	return 0, nil
}
