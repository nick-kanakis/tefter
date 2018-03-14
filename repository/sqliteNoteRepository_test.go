package repository

import (
	"github.com/nicolasmanic/tefter/model"
	"os"
	"testing"
	"time"
)

func TestSaveNoteCompleteData(t *testing.T) {
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

func TestSaveNoteMissingData(t *testing.T) {
	testRepo := NewNoteRepository("test.db")
	//tear down test
	defer func() {
		testRepo.CloseDB()
		os.Remove("test.db")
	}()

	mockNoteShouldBeSaved := model.NewNote("", "test Memo", 0, []string{})
	mockNoteShouldBeSaved.Created = time.Time{}
	mockNoteShouldBeSaved.LastUpdated = time.Time{}
	mockNoteShouldNotBeSaved := model.NewNote("", "", 0, []string{})
	_, err := testRepo.SaveNote(mockNoteShouldBeSaved)

	if err != nil {
		t.Errorf("Could not save note to DB, error msg: %v", err)
	}

	if mockNoteShouldBeSaved.ID != DEFAULT_NOTEPAD_ID {
		t.Error("Could not save note correctly to DB")
	}

	_, err = testRepo.SaveNote(mockNoteShouldNotBeSaved)
	if err.Error() != "Note should contain memo" {
		t.Error("Expected error with message: 'Note should contain memo'")
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
	mockNote3 := model.NewNote("testTitle", "test Memo", 1, []string{})

	id1, _ := testRepo.SaveNote(mockNote1)
	id2, _ := testRepo.SaveNote(mockNote2)
	testRepo.SaveNote(mockNote3)

	notes, err := testRepo.GetNotes([]int64{id1, id2})
	if err != nil {
		t.Errorf("Could not retrieve notes from DB, error msg: %v", err)
	}

	if len(notes) != 2 {
		t.Error("Could not retrieve specific notes from DB")
	}

	if notes[0].ID != id1 || notes[1].ID != id2 {
		t.Error("Could not properly retrieve note from DB")
	}

	if len(notes[0].Tags) != 2 || len(notes[1].Tags) != 2 {
		t.Error("Could not properly retrieve tags of note from DB")
	}

	allNotes, err := testRepo.GetNotes([]int64{})
	if len(allNotes) != 3 {
		t.Error("Could not retrieve all notes from DB")
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

	_, err = testRepo.GetNote(12345)
	if err.Error() != "Could find note with id: 12345"{
		t.Error("Expected error message 'Could find note with id: XXXX'")
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
	mockNote.UpdateMemo("Updated Memo")
	mockNote.UpdateTags([]string{"testTag1"})
	mockNote.Created = time.Time{}
	mockNote.LastUpdated = time.Time{}

	err := testRepo.UpdateNote(mockNote)
	if err != nil {
		t.Errorf("Could not update note, error msg: %v", err)
	}

	updatedNote, _ := testRepo.GetNotes([]int64{mockNote.ID})
	if updatedNote[0].Memo != "Updated Memo" {
		t.Error("Could not update note with changed memo")
	}
	if len(updatedNote[0].Tags) != 1 || !updatedNote[0].Tags["testTag1"] {
		t.Error("Could not update note with changed tags")
	}

	mockNote.UpdateMemo("")
	err = testRepo.UpdateNote(mockNote)
	if err.Error() != "Note should contain memo"{
		t.Error("Expected error with message: 'Note should contain memo'")
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

	_, err = testRepo.SearchNotesByKeyword("")
	if err.Error() != "Empty search parameter"{
		t.Error("Expected error with message: 'Empty search parameter'")
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
