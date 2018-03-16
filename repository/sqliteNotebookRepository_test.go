package repository

import (
	"github.com/nicolasmanic/tefter/model"
	"os"
	"testing"
)

func TestSaveNotebook(t *testing.T) {
	testRepo := NewNotebookRepository("test.db")
	//tear down test
	defer func() {
		testRepo.CloseDB()
		os.Remove("test.db")
	}()

	mockNotebook := model.NewNotebook("testTitle")
	id, err := testRepo.SaveNotebook(mockNotebook)

	if err != nil {
		t.Errorf("Could not save notebook to DB, error msg: %v", err)
	}

	if id != mockNotebook.ID {
		t.Error("Could not save correctly notebook to DB")
	}
}

func TestGetNotebooks(t *testing.T) {
	testRepo := NewNotebookRepository("test.db")
	testNoteRepo := NewNoteRepository("test.db")
	//tear down test
	defer func() {
		testRepo.CloseDB()
		testNoteRepo.CloseDB()
		os.Remove("test.db")
	}()

	mockNotebook2 := model.NewNotebook("notebook 2")
	mockNotebook3 := model.NewNotebook("notebook 3")
	testRepo.SaveNotebook(mockNotebook2)
	testRepo.SaveNotebook(mockNotebook3)

	mockNote1 := model.NewNote("testTitle", "test Memo", 2, []string{"testTag1", "testTag2"})
	mockNote2 := model.NewNote("testTitle", "test Memo", 3, []string{"testTag3", "testTag4"})

	testNoteRepo.SaveNote(mockNote1)
	testNoteRepo.SaveNote(mockNote2)

	notebooks, err := testRepo.GetNotebooks([]int64{2, 3})

	if err != nil {
		t.Errorf("Could not retrieve notebooks from DB, error msg: %v", err)
	}

	if len(notebooks) != 2 {
		t.Error("Could not retrieve notebooks from DB")
	}

	for _, notebook := range notebooks {
		if len(notebook.Notes) != 1 {
			t.Error("Could not retrieve notebooks from DB")
		}
	}

}

func TestGetNotebookByTitle(t *testing.T) {
	testRepo := NewNotebookRepository("test.db")
	//tear down test
	defer func() {
		testRepo.CloseDB()
		os.Remove("test.db")
	}()

	mockNotebook2 := model.NewNotebook("notebook 2")
	mockNotebook3 := model.NewNotebook("notebook 3")
	testRepo.SaveNotebook(mockNotebook2)
	testRepo.SaveNotebook(mockNotebook3)

	notebook, err := testRepo.GetNotebookByTitle("notebook 2")

	if err != nil {
		t.Errorf("Could not retrieve notebook by title from DB, error msg: %v", err)
	}

	if notebook.Title != "notebook 2" {
		t.Error("Could not retrieve notebook by title from DB")
	}
}

func TestUpdateNotebook(t *testing.T) {
	testRepo := NewNotebookRepository("test.db")
	//tear down test
	defer func() {
		testRepo.CloseDB()
		os.Remove("test.db")
	}()

	mockNotebook := model.NewNotebook("test Title")
	id, _ := testRepo.SaveNotebook(mockNotebook)

	mockNotebook.Title = "Updated Title"

	err := testRepo.UpdateNotebook(mockNotebook)
	if err != nil {
		t.Errorf("Could not update notebook, error msg: %v", err)
	}
	notebook, _ := testRepo.GetNotebook(id)

	if notebook.Title != "Updated Title" {
		t.Error("Could not update notebook")
	}
}

func TestDeleteNotebooks(t *testing.T) {
	testRepo := NewNotebookRepository("test.db")
	testNoteRepo := NewNoteRepository("test.db")
	//tear down test
	defer func() {
		testRepo.CloseDB()
		testNoteRepo.CloseDB()
		os.Remove("test.db")
	}()

	mockNotebook2 := model.NewNotebook("notebook 2")
	mockNotebook3 := model.NewNotebook("notebook 3")
	testRepo.SaveNotebook(mockNotebook2)
	testRepo.SaveNotebook(mockNotebook3)

	mockNote1 := model.NewNote("testTitle", "test Memo", 2, []string{"testTag1", "testTag2"})
	mockNote2 := model.NewNote("testTitle", "test Memo", 3, []string{"testTag3", "testTag4"})

	testNoteRepo.SaveNote(mockNote1)
	testNoteRepo.SaveNote(mockNote2)

	err := testRepo.DeleteNotebooks([]int64{2, 3})
	if err != nil {
		t.Errorf("Could not delete notebooks from DB, error msg: %v", err)
	}

	notes, _ := testNoteRepo.GetNotes([]int64{mockNote1.ID, mockNote2.ID})
	if len(notes) != 0 {
		t.Errorf("Could not delete notes of deleted notebook from DB")
	}

	notebooks, _ := testRepo.GetNotebooks([]int64{2, 3})
	if len(notebooks) != 0 {
		t.Errorf("Could not delete notebook from DB")
	}
}

func TestGetAllNotebooksTitle(t *testing.T) {
	testRepo := NewNotebookRepository("test.db")
	//tear down test
	defer func() {
		testRepo.CloseDB()
		os.Remove("test.db")
	}()

	mockNotebook1 := model.NewNotebook("notebook 1")
	mockNotebook2 := model.NewNotebook("notebook 2")
	id1, _ := testRepo.SaveNotebook(mockNotebook1)
	id2, _ := testRepo.SaveNotebook(mockNotebook2)

	result, err := testRepo.GetAllNotebooksTitle()
	if err != nil {
		t.Errorf("Could not retrieve notebooks title map from DB, error msg: %v", err)
	}
	//2 user defined notebooks and 1 default
	if len(result) != 3 {
		t.Error("Size of notebooks title map is incorrect")
	}

	if result[id1] != "notebook 1" || result[id2] != "notebook 2" {
		t.Error("Incorrect data in Notebook title map")
	}
}
