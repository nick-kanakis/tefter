package cmd

import (
	"errors"
	"github.com/nicolasmanic/tefter/model"
	"github.com/nicolasmanic/tefter/repository"
	"reflect"
	"testing"
)

func TestDeleteNotebooks(t *testing.T) {
	cases := []struct {
		mDB         mockNotebookDBDelete
		expectedErr error
		titles      []string
	}{
		{
			mDB: mockNotebookDBDelete{
				notebook: model.NewNotebook("title"),
				err:      nil,
			},
			expectedErr: nil,
			titles:      []string{"title1", "title2"},
		}, {
			mDB:         mockNotebookDBDelete{},
			expectedErr: errors.New("No argument passed, at least one notebook title should be provided"),
			titles:      []string{},
		}, {
			mDB: mockNotebookDBDelete{
				notebook: nil,
				err:      errors.New("Unexpected error"),
			},
			expectedErr: errors.New("Could not retrieve notebook for title: title1 error msg: Unexpected error"),
			titles:      []string{"title1"},
		},
	}

	for _, c := range cases {
		oldNotebookDB := NotebookDB
		NotebookDB = c.mDB
		defer func() {
			NotebookDB = oldNotebookDB
		}()

		err := deleteNotebooks(c.titles)
		if !reflect.DeepEqual(err, c.expectedErr) {
			t.Errorf("Expected err to be %q but it was %q", c.expectedErr, err)
		}
	}
}

type mockNotebookDBDelete struct {
	repository.NotebookRepository
	notebook *model.Notebook
	err      error
}

func (mDB mockNotebookDBDelete) DeleteNotebook(notebooksID int64) error {
	return mDB.err
}

func (mDB mockNotebookDBDelete) GetNotebookByTitle(notebookTitle string) (*model.Notebook, error) {
	/*notebook := model.NewNotebook(notebookTitle)
	notebook.ID = 1*/
	return mDB.notebook, mDB.err
}
