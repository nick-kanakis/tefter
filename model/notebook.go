package model

//Notebook holds a list of "similar" notes
type Notebook struct {
	ID    int64
	Title string
	Notes map[int64]*Note
}

//NewNotebook returns a Notebook pointer
func NewNotebook(title string) *Notebook {
	return &Notebook{
		ID:    0,
		Notes: make(map[int64]*Note),
		Title: title,
	}
}

//AddNote adds note to Notebook list, also update the note with current Notebook
func (nBook *Notebook) AddNote(note *Note) bool {
	if doesNoteExist(nBook, note.ID) {
		return false
	}
	nBook.Notes[note.ID] = note
	note.UpdateNotebook(nBook.ID)
	return true
}

//RemoveNote remove note from Notebook list.
//WARNING: when removing a note we should add it to a default notebook.
func (nBook *Notebook) RemoveNote(noteID int64) bool {
	if !doesNoteExist(nBook, noteID) {
		return false
	}
	delete(nBook.Notes, noteID)
	return true
}

func doesNoteExist(notebook *Notebook, noteID int64) bool {
	_, exists := notebook.Notes[noteID]
	return exists
}
