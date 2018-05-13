package cmd

import (
	"errors"
	"github.com/nicolasmanic/tefter/model"
	"github.com/nicolasmanic/tefter/repository"
	"os"
	"reflect"
	"testing"
)

func TestExport(t *testing.T) {
	cases := []struct {
		mockNotebookDB mockNotebookDBExport
		mockNoteDB     mockNoteDBExport
		ids            []int
		notebookTitles []string
		tags           []string
		getAll         bool
		expectedErr    error
	}{
		{
			mockNotebookDB: mockNotebookDBExport{},
			mockNoteDB: mockNoteDBExport{
				err: errors.New("Unexpected error"),
			},
			ids:            []int{},
			notebookTitles: []string{},
			tags:           []string{},
			getAll:         true,
			expectedErr:    errors.New("Error while retrieving all notes, error msg: Unexpected error"),
		}, {
			mockNotebookDB: mockNotebookDBExport{},
			mockNoteDB:     mockNoteDBExport{},
			ids:            []int{},
			notebookTitles: []string{},
			tags:           []string{},
			getAll:         true,
			expectedErr:    nil,
		}, {
			mockNotebookDB: mockNotebookDBExport{},
			mockNoteDB:     mockNoteDBExport{},
			ids:            []int{},
			notebookTitles: []string{},
			tags:           []string{},
			getAll:         true,
			expectedErr:    nil,
		},
	}
	for _, c := range cases {
		oldNotebookDB := NotebookDB
		oldNoteDB := NoteDB
		NoteDB = c.mockNoteDB
		NotebookDB = c.mockNotebookDB

		defer func() {
			NotebookDB = oldNotebookDB
			NoteDB = oldNoteDB
			os.Remove("notes.json")
		}()

		err := export(c.ids, c.notebookTitles, c.tags, c.getAll)
		if !reflect.DeepEqual(c.expectedErr, err) {
			t.Errorf("Expected err to be %q but it was %q", c.expectedErr, err)
		}
	}
}

type mockNotebookDBExport struct {
	repository.NotebookRepository
	notebook       *model.Notebook
	notebookTitles map[int64]string
	err            error
}

type mockNoteDBExport struct {
	repository.NoteRepository
	notes []*model.Note
	id    int64
	err   error
}

func (mDB mockNotebookDBExport) GetNotebookByTitle(notebooksTitle string) (*model.Notebook, error) {
	return mDB.notebook, mDB.err
}

func (mDB mockNotebookDBExport) GetAllNotebooksTitle() (map[int64]string, error) {
	return mDB.notebookTitles, mDB.err
}

func (mDB mockNoteDBExport) GetNotesByTag(tags []string) ([]*model.Note, error) {
	return mDB.notes, mDB.err
}

func (mDB mockNoteDBExport) SaveNote(note *model.Note) (int64, error) {
	return 1, nil
}

func (mDB mockNoteDBExport) GetNotes(noteIDs []int64) ([]*model.Note, error) {
	return mDB.notes, mDB.err
}
