package cmd

import (
	"errors"
	"github.com/nicolasmanic/tefter/model"
	"github.com/nicolasmanic/tefter/repository"
	"reflect"
	"testing"
)

func TestUpdateNotebook(t *testing.T) {
	cases := []struct {
		mDB         mockNotebookDBUpdateNotebook
		expectedErr error
		newTitle    string
	}{
		{
			mDB: mockNotebookDBUpdateNotebook{
				notebook: nil,
				err:      nil,
			},
			newTitle:    "",
			expectedErr: errors.New("Notebook title should not be empty"),
		}, {
			mDB: mockNotebookDBUpdateNotebook{
				notebook: model.NewNotebook("notebookTitle"),
				err:      nil,
			},
			newTitle:    "newNotebookTitle",
			expectedErr: nil,
		}, {
			mDB: mockNotebookDBUpdateNotebook{
				notebook: nil,
				err:      nil,
			},
			newTitle:    "newNotebookTitle",
			expectedErr: errors.New("No notebook with title: oldTitle"),
		}, {
			mDB: mockNotebookDBUpdateNotebook{
				notebook: nil,
				err:      errors.New("Unexpected error"),
			},
			newTitle:    "newNotebookTitle",
			expectedErr: errors.New("Error while retrieving notebook by title, error msg: Unexpected error"),
		},
	}

	for _, c := range cases {
		oldNotebookDB := NotebookDB
		NotebookDB = c.mDB
		//Restore interface
		defer func() {
			NotebookDB = oldNotebookDB
		}()
		err := updateNotebook("oldTitle", c.newTitle)
		if !reflect.DeepEqual(c.expectedErr, err) {
			t.Errorf("Expected err to be %q but it was %q", c.expectedErr, err)
		}
	}
}

type mockNotebookDBUpdateNotebook struct {
	repository.NotebookRepository
	notebook *model.Notebook
	err      error
}

func (mDB mockNotebookDBUpdateNotebook) UpdateNotebook(notebook *model.Notebook) error {
	return mDB.err
}

func (mDB mockNotebookDBUpdateNotebook) GetNotebookByTitle(notebookTitle string) (*model.Notebook, error) {
	/*notebook := model.NewNotebook(notebookTitle)
	notebook.ID = 1
	note := model.NewNote("testTitle", "testMemo", notebook.ID, []string{})
	note.ID = 3
	notebook.AddNote(note)*/
	return mDB.notebook, mDB.err
}
