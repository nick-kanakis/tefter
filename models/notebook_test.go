package models

import "testing"

func TestAddNote(t *testing.T){
	notebook:=createMockedNotebook()
	notebook.AddNote(createMockedNote())

	if len(notebook.Notes) != 1{
		t.Error("Could not add note to notebook")
	}

	note := notebook.Notes[0]
	if note.Notebook != notebook{
		t.Error("Could update reference to notebook in note")
	}
}

func TestRemoveNote(t *testing.T){
	notebook:=createMockedNotebook()
	note1:=createMockedNote()
	note1.ID =1
	notebook.AddNote(note1)
	note2:=createMockedNote()
	note1.ID =2
	notebook.AddNote(note2)

	if len(notebook.Notes) != 2{
		t.Error("Could not add note to notebook")
	}

	notebook.RemoveNote(1)
	if len(notebook.Notes) != 1{
		t.Error("Could not delete note to notebook")
	}
}

func TestDoesNoteExist(t *testing.T){
	notebook:=createMockedNotebook()
	note1:=createMockedNote()
	note1.ID =1
	notebook.AddNote(note1)
	note2:=createMockedNote()
	note1.ID =2
	notebook.AddNote(note2)

	if !doesNoteExist(notebook, 1){
		t.Error("Failed finding note in notebook")
	}
}

func createMockedNotebook() *Notebook{
	return NewNotebook("testNotebook")
}