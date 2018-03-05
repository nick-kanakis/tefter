package models

//Notebook holds a list of "similar" notes
type Notebook struct{
	Notes map[int64]*Note 
	Title string 
}

//NewNotebook returns a Notebook pointer
func NewNotebook(title string) *Notebook{
	return &Notebook{
		Notes: make(map[int64]*Note),
		Title: title,
	}
}

//AddNote adds note to Notebook list, also update the note with current Notebook
func (nBook *Notebook) AddNote(note *Note) bool{
	if doesNoteExist(nBook, note.ID){
		return false
	}
	nBook.Notes[note.ID] = note
	note.UpdateNotebook(nBook)
	return true
}

//RemoveNote remove note from Notebook list.
//WARNING: when removing a note we should add it to a default notebook.
func (nBook *Notebook) RemoveNote(noteID int64) bool{
	if !doesNoteExist(nBook, noteID){
		return false
	}
	delete(nBook.Notes, noteID)
	return true
}

func doesNoteExist(notebook *Notebook, noteID int64) bool{
	_, exists:= notebook.Notes[noteID]
	return exists
}