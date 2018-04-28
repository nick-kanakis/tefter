package cmd

import (
	"errors"
	"github.com/nicolasmanic/tefter/model"
	"github.com/nicolasmanic/tefter/repository"
	"reflect"
	"testing"
)

func TestUpdate(t *testing.T) {
	cases := []struct {
		noteDB        mockNoteDBUpdate
		notebookDB    mockNotebookDBUpdate
		id            int64
		notebookTitle string
		noteTitle     string
		tags          []string
		editor        Editor
		expectedErr   error
	}{
		{
			noteDB: mockNoteDBUpdate{
				note: model.NewNote("title", "memo", 1, []string{}),
				err:  nil,
			},
			notebookDB: mockNotebookDBUpdate{
				notebook: &model.Notebook{1, "NotebookTitle", nil},
				err:      nil,
			},
			notebookTitle: "title",
			noteTitle:     "NotebookTitle",
			tags:          []string{"tag1", "tag2"},
			editor:        &fakeEditor{"memo"},
			expectedErr:   nil,
		}, {
			noteDB: mockNoteDBUpdate{
				note: nil,
				err:  errors.New("Unexpected error"),
			},
			notebookDB:    mockNotebookDBUpdate{},
			notebookTitle: "title",
			noteTitle:     "NotebookTitle",
			tags:          []string{},
			editor:        &fakeEditor{"memo"},
			expectedErr:   errors.New("Error while retrieving Note from DB, error msg: Unexpected error"),
		}, {
			noteDB: mockNoteDBUpdate{
				note: model.NewNote("title", "memo", 1, []string{}),
				err:  nil,
			},
			notebookDB: mockNotebookDBUpdate{
				notebook: nil,
				err:      errors.New("Unexpected error"),
			},
			notebookTitle: "title",
			noteTitle:     "NotebookTitle",
			tags:          []string{},
			editor:        &fakeEditor{"memo"},
			expectedErr:   errors.New("Error while constructing updated note, error msg: Unexpected error"),
		},
	}

	for _, c := range cases {
		oldNotebookDB := NotebookDB
		oldNoteDB := NoteDB
		NoteDB = c.noteDB
		NotebookDB = c.notebookDB
		//Restore interface
		defer func() {
			NotebookDB = oldNotebookDB
			NoteDB = oldNoteDB
		}()

		err := update(c.id, c.noteTitle, c.tags, c.notebookTitle, c.editor)
		if !reflect.DeepEqual(c.expectedErr, err) {
			t.Errorf("Expected err to be %q but it was %q", c.expectedErr, err)
		}
	}
}

func TestUpdateJSONNote(t *testing.T) {
	oldNoteDB := NoteDB
	NoteDB = mockNoteDBUpdate{
		note: nil,
		err:  errors.New("Unexpected error"),
	}
	//Restore interface
	defer func() {
		NoteDB = oldNoteDB
	}()
	jNote := &jsonNote{
		ID:            1,
		Title:         "Test Title",
		Tags:          []string{"tag1", "tag2"},
		NotebookTitle: "notebook",
	}
	err := updateJSONNote(jNote)
	expectedErr := errors.New("Error while retrieving Note from DB, error msg: Unexpected error")
	if !reflect.DeepEqual(expectedErr, err) {
		t.Errorf("Expected err to be %q but it was %q", expectedErr, err)
	}
}

func TestConstructUpdatedNote(t *testing.T) {
	oldNotebookDB := NotebookDB
	NotebookDB = mockNotebookDBUpdate{}
	//Restore interface
	defer func() {
		NotebookDB = oldNotebookDB
	}()
	note := model.NewNote("testTitle4", "testMemo", repository.DEFAULT_NOTEBOOK_ID, []string{"tag1", "tag2"})
	constructUpdatedNote(note, "", "", []string{"tag3", "-tag1"}, "NewMemo")

	if len(note.Tags) != 2 {
		t.Error("Failed adding/removing tags")
	}
	if !(note.Tags["tag3"]) || !(note.Tags["tag2"]) {
		t.Error("Failed adding/removing tags")
	}
}

type mockNoteDBUpdate struct {
	repository.NoteRepository
	note *model.Note
	err  error
}

func (mDB mockNoteDBUpdate) GetNote(noteID int64) (*model.Note, error) {
	return mDB.note, mDB.err
}

func (mDB mockNoteDBUpdate) UpdateNote(note *model.Note) error {
	return mDB.err
}

type mockNotebookDBUpdate struct {
	repository.NotebookRepository
	notebook *model.Notebook
	err      error
}

func (mDB mockNotebookDBUpdate) GetNotebookByTitle(notebooksTitle string) (*model.Notebook, error) {
	return mDB.notebook, mDB.err
}
