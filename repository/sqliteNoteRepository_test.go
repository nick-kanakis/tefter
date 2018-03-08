package repository

import (
	"testing"
	"github.com/nicolasmanic/tefter/model"
	"os"
)
//TODO clean up DB or mock it
func TestSaveNote(t *testing.T){
	testRepo := NewNoteRepository("test.db")	
	//tear down test
	defer func(){
		testRepo.CloseDB()
		os.Remove("test.db")
	}()
	
	mockNote := model.NewNote("testTitle", "test Memo", 1, []string{"testTag1", "testTag2"})
	id, err:= testRepo.SaveNote(mockNote)

	if err != nil{
		t.Errorf("Could not save note to DB, error msg: %v", err)
	}

	if id != mockNote.ID {
		t.Error("Could not save correctly note to DB")
	}
}

func TestGetNotes(t *testing.T){
	testRepo := NewNoteRepository("test.db")	
	//tear down test
	defer func(){
		testRepo.CloseDB()
		os.Remove("test.db")
	}()

	mockNote1 := model.NewNote("testTitle", "test Memo", 1, []string{"testTag1", "testTag2"})
	mockNote2 := model.NewNote("testTitle", "test Memo", 1, []string{"testTag3", "testTag4"})

	id1, _:= testRepo.SaveNote(mockNote1)
	id2, _:= testRepo.SaveNote(mockNote2)

	notes, err:=testRepo.GetNotes([]int64{id1, id2})
	if err != nil{
		t.Errorf("Could not retrieve note from DB, error msg: %v", err)
	}

	if len(notes) != 2{
		t.Error("Could not retrieve note from DB")
	}
	
	if notes[0].ID != id1 ||notes[1].ID != id2{
		t.Error("Could not properly retrieve note from DB")
	}

	if notes[0].ID != id1 ||notes[1].ID != id2{
		t.Error("Could not properly retrieve tags of note from DB")
	}

	if  len(notes[0].Tags) != 2 || len(notes[1].Tags) != 2{
		t.Error("Could not properly retrieve tags of note from DB")
	}
}