package repository

import (
	"os"
	"testing"
	"github.com/nicolasmanic/tefter/model"
)

func TestSaveNote(t *testing.T) {
	testRepo := NewNoteRepository("test.db")
	//tear down test
	defer func() {
		testRepo.CloseDB()
		os.Remove("test.db")
	}()

	mockNote := model.NewNote("testTitle", "test Memo", 1, []string{"testTag1", "testTag2"})
	id, err := testRepo.SaveNote(mockNote)

	if err != nil {
		t.Errorf("Could not save note to DB, error msg: %v", err)
	}

	if id != mockNote.ID {
		t.Error("Could not save correctly note to DB")
	}
}

func TestGetNotes(t *testing.T) {
	testRepo := NewNoteRepository("test.db")
	//tear down test
	defer func() {
		testRepo.CloseDB()
		os.Remove("test.db")
	}()

	mockNote1 := model.NewNote("testTitle", "test Memo", 1, []string{"testTag1", "testTag2"})
	mockNote2 := model.NewNote("testTitle", "test Memo", 1, []string{"testTag3", "testTag4"})

	id1, _ := testRepo.SaveNote(mockNote1)
	id2, _ := testRepo.SaveNote(mockNote2)

	notes, err := testRepo.GetNotes([]int64{id1, id2})
	if err != nil {
		t.Errorf("Could not retrieve note from DB, error msg: %v", err)
	}

	if len(notes) != 2 {
		t.Error("Could not retrieve note from DB")
	}

	if notes[0].ID != id1 || notes[1].ID != id2 {
		t.Error("Could not properly retrieve note from DB")
	}

	if len(notes[0].Tags) != 2 || len(notes[1].Tags) != 2 {
		t.Error("Could not properly retrieve tags of note from DB")
	}

	if notes[0].Created.IsZero() || notes[0].Created.IsZero() {
		t.Error("Could not properly retrieve tags of note from DB")
	}
}

func TestGetNote(t *testing.T) {
	testRepo := NewNoteRepository("test.db")
	//tear down test
	defer func() {
		testRepo.CloseDB()
		os.Remove("test.db")
	}()

	mockNote := model.NewNote("testTitle", "test Memo", 1, []string{"testTag1", "testTag2"})
	id, _ := testRepo.SaveNote(mockNote)

	note, err := testRepo.GetNote(id)
	if err != nil {
		t.Errorf("Could not retrieve note from DB, error msg: %v", err)
	}

	if note.ID != id {
		t.Error("Could not properly retrieve note from DB")
	}
}

func TestUpdateNote(t *testing.T) {
	testRepo := NewNoteRepository("test.db")
	//tear down test
	defer func() {
		testRepo.CloseDB()
		os.Remove("test.db")
	}()

	mockNote := model.NewNote("testTitle", "test Memo", 1, []string{"testTag1", "testTag2"})
	testRepo.SaveNote(mockNote)
	mockNote.Memo = "Updated Memo"
	delete(mockNote.Tags, "testTag1")
	mockNote.Tags["testTag3"] = true

	err := testRepo.UpdateNote(mockNote)
	if err != nil {
		t.Errorf("Could not update note, error msg: %v", err)
	}

	updatedNote, _ := testRepo.GetNotes([]int64{mockNote.ID})
	if updatedNote[0].Memo != "Updated Memo" {
		t.Error("Could not update note")
	}
	if len(updatedNote[0].Tags) != 2 || !updatedNote[0].Tags["testTag3"] {
		t.Error("Could not update note")
	}
}

func TestDeleteNotes(t *testing.T) {
	testRepo := NewNoteRepository("test.db")
	//tear down test
	defer func() {
		testRepo.CloseDB()
		os.Remove("test.db")
	}()

	mockNote1 := model.NewNote("testTitle", "test Memo", 1, []string{"testTag1", "testTag2"})
	mockNote2 := model.NewNote("testTitle", "test Memo", 1, []string{"testTag3", "testTag4"})
	mockNote3 := model.NewNote("testTitle", "test Memo", 1, []string{"testTag5"})

	id1, _ := testRepo.SaveNote(mockNote1)
	id2, _ := testRepo.SaveNote(mockNote2)
	id3, _ := testRepo.SaveNote(mockNote3)

	err := testRepo.DeleteNotes([]int64{id1, id2})
	if err != nil {
		t.Errorf("Could not delete notes, error msg: %v", err)
	}

	err = testRepo.DeleteNote(id3)
	if err != nil {
		t.Errorf("Could not delete notes, error msg: %v", err)
	}

	notes, _ := testRepo.GetNotes([]int64{id1, id2, id3})
	if len(notes) != 0 {
		t.Errorf("Could not delete notes")
	}
}

func TestSearchNotesByKeyword(t *testing.T) {
	testRepo := NewNoteRepository("test.db")
	//tear down test
	defer func() {
		testRepo.CloseDB()
		os.Remove("test.db")
	}()

	mockNote1 := model.NewNote("title1", "test Memo 1", 1, []string{"testTag1", "testTag2"})
	mockNote2 := model.NewNote("title2", "test Memo 2", 1, []string{"testTag3", "testTag4"})
	mockNote3 := model.NewNote("title3", "test Memo 3", 1, []string{"testTag5"})

	testRepo.SaveNote(mockNote1)
	testRepo.SaveNote(mockNote2)
	testRepo.SaveNote(mockNote3)

	allNotes, err := testRepo.SearchNotesByKeyword("Memo")
	if err != nil {
		t.Errorf("Could not search notes by keyword, err msg: %v", err)
	}
	if len(allNotes) != 3 {
		t.Error("Could not search notes by keyword")
	}

	subSetOfNotes, err := testRepo.SearchNotesByKeyword("title1")
	if err != nil {
		t.Errorf("Could not search notes by keyword, err msg: %v", err)
	}
	if len(subSetOfNotes) != 1 {
		t.Error("Could not search notes by keyword")
	}
}

func TestSearchNoteByTag(t *testing.T) {
	testRepo := NewNoteRepository("test.db")
	//tear down test
	defer func() {
		testRepo.CloseDB()
		os.Remove("test.db")
	}()

	mockNote1 := model.NewNote("testTitle", "test Memo", 1, []string{"testTag1", "testTag2"})
	mockNote2 := model.NewNote("testTitle", "test Memo", 1, []string{"testTag3", "testTag4"})
	mockNote3 := model.NewNote("testTitle", "test Memo", 1, []string{"testTag5"})

	id1, _ := testRepo.SaveNote(mockNote1)
	testRepo.SaveNote(mockNote2)
	id3, _ := testRepo.SaveNote(mockNote3)

	notes, err := testRepo.SearchNotesByTag([]string{"testTag1", "testTag5"})

	if err != nil {
		t.Errorf("Could not search notes by tag, err msg: %v", err)
	}
	if len(notes) != 2 {
		t.Errorf("Could not search notes by tag")
	}
	if !((notes[0].ID != id1 || notes[0].ID != id3) && (notes[1].ID != id1 || notes[1].ID != id3)) {
		t.Errorf("Could not search notes by tag")
	}

}
