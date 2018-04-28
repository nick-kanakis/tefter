package cmd

import (
	"errors"
	"github.com/nicolasmanic/tefter/model"
	"github.com/nicolasmanic/tefter/repository"
	"reflect"
	"testing"
)

func TestAdd(t *testing.T) {
	cases := []struct {
		noteDB        mockNoteDBAdd
		notebookDB    mockNotebookDBAdd
		notebookTitle string
		noteTitle     string
		tags          []string
		editor        Editor
		expectedErr   error
	}{
		{
			noteDB: mockNoteDBAdd{
				id:  1,
				err: nil,
			},
			notebookDB: mockNotebookDBAdd{
				notebook: &model.Notebook{1, "NotebookTitle", nil},
				id:       1,
				err:      nil,
			},
			notebookTitle: "NotebookTitle",
			noteTitle:     "noteTitle",
			tags:          []string{},
			editor:        &fakeEditor{"testMemo"},
			expectedErr:   nil,
		}, {
			noteDB: mockNoteDBAdd{
				id:  1,
				err: nil,
			},
			notebookDB:    mockNotebookDBAdd{},
			notebookTitle: "",
			noteTitle:     "noteTitle",
			tags:          []string{},
			editor:        &fakeEditor{"testMemo"},
			expectedErr:   nil,
		}, {
			noteDB: mockNoteDBAdd{},
			notebookDB: mockNotebookDBAdd{
				notebook: nil,
				id:       0,
				err:      errors.New("Unexpected error"),
			},
			notebookTitle: "NotebookTitle",
			noteTitle:     "noteTitle",
			tags:          []string{},
			editor:        &fakeEditor{"testMemo"},
			expectedErr:   errors.New("Error while finding corresponding notebook for note, error msg: Unexpected error"),
		}, {
			noteDB: mockNoteDBAdd{
				id:  0,
				err: errors.New("Unexpected error"),
			},
			notebookDB:    mockNotebookDBAdd{},
			notebookTitle: "NotebookTitle",
			noteTitle:     "noteTitle",
			tags:          []string{},
			editor:        &fakeEditor{"testMemo"},
			expectedErr:   errors.New("Error while saving note, error msg: Unexpected error"),
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

		err := add(c.noteTitle, c.tags, c.notebookTitle, c.editor)
		if !reflect.DeepEqual(c.expectedErr, err) {
			t.Errorf("Expected err to be %q but it was %q", c.expectedErr, err)
		}
	}
}

func TestAddJSONNote(t *testing.T) {
	cases := []struct {
		noteDB      mockNoteDBAdd
		notebookDB  mockNotebookDBAdd
		jNote       *jsonNote
		expectedErr error
	}{
		{
			noteDB: mockNoteDBAdd{
				id:  1,
				err: nil,
			},
			notebookDB: mockNotebookDBAdd{
				notebook: &model.Notebook{1, "NotebookTitle", nil},
				id:       1,
				err:      nil,
			},
			jNote: &jsonNote{
				ID:            1,
				Title:         "Test Title",
				Tags:          []string{"tag1", "tag2"},
				NotebookTitle: "notebook",
			},
			expectedErr: nil,
		}, {
			noteDB: mockNoteDBAdd{},
			notebookDB: mockNotebookDBAdd{
				notebook: nil,
				id:       0,
				err:      errors.New("Unexpected error"),
			},
			jNote: &jsonNote{
				ID:            1,
				Title:         "Test Title",
				Tags:          []string{"tag1", "tag2"},
				NotebookTitle: "notebook",
			},
			expectedErr: errors.New("Error while finding corresponding notebook for note, error msg: Unexpected error"),
		}, {
			noteDB: mockNoteDBAdd{
				id:  0,
				err: errors.New("Unexpected error"),
			},
			notebookDB: mockNotebookDBAdd{
				notebook: &model.Notebook{1, "NotebookTitle", nil},
				id:       1,
				err:      nil,
			},
			jNote: &jsonNote{
				ID:            1,
				Title:         "Test Title",
				Tags:          []string{"tag1", "tag2"},
				NotebookTitle: "notebook",
			},
			expectedErr: errors.New("Error while saving note, error msg: Unexpected error"),
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

		err := addJSONNote(c.jNote)
		if !reflect.DeepEqual(c.expectedErr, err) {
			t.Errorf("Expected err to be %q but it was %q", c.expectedErr, err)
		}
	}
}

type mockNoteDBAdd struct {
	repository.NoteRepository
	id  int64
	err error
}

func (mDB mockNoteDBAdd) SaveNote(note *model.Note) (int64, error) {
	return mDB.id, mDB.err
}

type mockNotebookDBAdd struct {
	repository.NotebookRepository
	notebook *model.Notebook
	id       int64
	err      error
}

func (mDB mockNotebookDBAdd) GetNotebookByTitle(notebookTitle string) (*model.Notebook, error) {
	return mDB.notebook, mDB.err
}

func (mDB mockNotebookDBAdd) SaveNotebook(notebook *model.Notebook) (int64, error) {
	return mDB.id, mDB.err
}

type fakeEditor struct {
	returnedText string
}

func (fe fakeEditor) edit(text string) string {
	return fe.returnedText
}
