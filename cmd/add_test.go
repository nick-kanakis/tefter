package cmd

import (
	"github.com/nicolasmanic/tefter/model"
	"testing"
)

func TestAddNotebookToNoteDefaultNotebook(t *testing.T) {
	note := model.NewNote("title", "memo", 123, []string{})
	err := addNotebookToNote(note, "")

	if err != nil {
		t.Errorf("Error while trying to add note to default notebook, with msg: %v", err)
	}

	if note.NotebookID != DEFAULT_NOTEBOOK_ID {
		t.Error("Error while trying to add note to default notebook")
	}
}

func TestAddNotebookToNoteExistingNotebook(t *testing.T) {
	oldNotebookDB := NotebookDB
	NotebookDB = mockNotebookRepoNotebookExist{}
	//Restore interface
	defer func() {
		NotebookDB = oldNotebookDB
	}()

	note := model.NewNote("title", "memo", 123, []string{})
	err := addNotebookToNote(note, "Existing notebook")
	if err != nil {
		t.Errorf("Error while trying to add note to existing notebook, with msg: %v", err)
	}

	if note.NotebookID != 2 {
		t.Error("Error while trying to add note to existing notebook")
	}

}

func TestAddNotebookToNoteNewNotebook(t *testing.T) {
	oldNotebookDB := NotebookDB
	NotebookDB = mockNotebookRepoNewNotebook{}
	//Restore interface
	defer func() {
		NotebookDB = oldNotebookDB
	}()

	note := model.NewNote("title", "memo", 123, []string{})
	err := addNotebookToNote(note, "New notebook")
	if err != nil {
		t.Errorf("Error while trying to add note to a new notebook, with msg: %v", err)
	}

	if note.NotebookID != 2 {
		t.Error("Error while trying to add note to a new notebook")
	}
}

type mockNotebookRepoNotebookExist struct{}

func (mockNotebookRepoNotebookExist) SaveNotebook(notebook *model.Notebook) (int64, error) {
	panic("not implemented")
}

func (mockNotebookRepoNotebookExist) GetNotebooks(notebooksIDs []int64) ([]*model.Notebook, error) {
	panic("not implemented")
}

func (mockNotebookRepoNotebookExist) GetNotebook(notebooksID int64) (*model.Notebook, error) {
	panic("not implemented")
}

func (mockNotebookRepoNotebookExist) GetNotebookByTitle(notebooksTitle string) (*model.Notebook, error) {
	notebook := model.NewNotebook("existing Notebook")
	notebook.ID = 2
	return notebook, nil
}

func (mockNotebookRepoNotebookExist) UpdateNotebook(notebook *model.Notebook) error {
	panic("not implemented")
}

func (mockNotebookRepoNotebookExist) DeleteNotebooks(notebooksIDs []int64) error {
	panic("not implemented")
}

func (mockNotebookRepoNotebookExist) DeleteNotebook(notebooksID int64) error {
	panic("not implemented")
}

func (mockNotebookRepoNotebookExist) GetAllNotebooksTitle() (map[int64]string, error) {
	panic("not implemented")
}

func (mockNotebookRepoNotebookExist) CloseDB() error {
	panic("not implemented")
}

type mockNotebookRepoNewNotebook struct{}

func (mockNotebookRepoNewNotebook) SaveNotebook(notebook *model.Notebook) (int64, error) {
	return 2, nil
}

func (mockNotebookRepoNewNotebook) GetNotebooks(notebooksIDs []int64) ([]*model.Notebook, error) {
	panic("not implemented")
}

func (mockNotebookRepoNewNotebook) GetNotebook(notebooksID int64) (*model.Notebook, error) {
	panic("not implemented")
}

func (mockNotebookRepoNewNotebook) GetNotebookByTitle(notebooksTitle string) (*model.Notebook, error) {
	return nil, nil
}

func (mockNotebookRepoNewNotebook) UpdateNotebook(notebook *model.Notebook) error {
	panic("not implemented")
}

func (mockNotebookRepoNewNotebook) DeleteNotebooks(notebooksIDs []int64) error {
	panic("not implemented")
}

func (mockNotebookRepoNewNotebook) DeleteNotebook(notebooksID int64) error {
	panic("not implemented")
}

func (mockNotebookRepoNewNotebook) GetAllNotebooksTitle() (map[int64]string, error) {
	panic("not implemented")
}

func (mockNotebookRepoNewNotebook) CloseDB() error {
	panic("not implemented")
}
