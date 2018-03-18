package cmd

import (
	"github.com/nicolasmanic/tefter/model"
	"github.com/nicolasmanic/tefter/repository"
	"testing"
)

func TestDeleteNotebookNoArguments(t *testing.T) {
	oldNotebookDB := NotebookDB
	NotebookDB = mockNotebookDBDelete{}
	defer func() {
		NotebookDB = oldNotebookDB
		if r := recover(); r == nil {
			t.Errorf("Empty argument should cause the delete cmd to panic")
		}
	}()
	deleteNotebook([]string{})
}

func TestDeleteNotebook(t *testing.T) {
	oldNotebookDB := NotebookDB
	NotebookDB = mockNotebookDBDelete{}
	//Restore interface
	defer func() {
		NotebookDB = oldNotebookDB
	}()
	deleteNotebook([]string{"12"})
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
