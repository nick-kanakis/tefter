package cmd

import (
	"github.com/nicolasmanic/tefter/model"
	"github.com/nicolasmanic/tefter/repository"
	"testing"
)

func TestDeleteNotebooksNoArguments(t *testing.T) {
	oldNotebookDB := NotebookDB
	NotebookDB = mockNotebookDBDelete{}
	defer func() {
		NotebookDB = oldNotebookDB
		if r := recover(); r == nil {
			t.Errorf("Empty argument should cause the delete cmd to panic")
		}
	}()
	deleteNotebooksWrapper(nil, []string{})
}

func TestDeleteNotebooks(t *testing.T) {
	oldNotebookDB := NotebookDB
	NotebookDB = mockNotebookDBDelete{}
	//Restore interface
	defer func() {
		NotebookDB = oldNotebookDB
	}()
	deleteNotebooks([]string{"12"})
}

type mockNotebookDBDelete struct {
	repository.NotebookRepository
}

func (mDB mockNotebookDBDelete) DeleteNotebook(notebooksID int64) error {
	return nil
}

func (mDB mockNotebookDBDelete) GetNotebookByTitle(notebookTitle string) (*model.Notebook, error) {
	notebook := model.NewNotebook(notebookTitle)
	notebook.ID = 1
	return notebook, nil
}
