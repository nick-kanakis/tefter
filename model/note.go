package model

import "time"

//Note is the note we want to keep.
type Note struct {
	ID          int64
	Title       string
	Memo        string
	Created     time.Time
	LastUpdated time.Time
	Tags        map[string]bool
	Notebook    *Notebook
}

//NewNote returns a new note pointer.
func NewNote(title, memo string, notebook *Notebook, tags []string) *Note {
	note := &Note{
		Title:       title,
		Memo:        memo,
		Created:     time.Now(),
		LastUpdated: time.Now(),
		Tags:        make(map[string]bool),
		Notebook:    notebook,
	}
	note.AddTags(tags)
	return note
}

//UpdateTitle updates the title of the note (also updates the LastUpdate value)
func (note *Note) UpdateTitle(newTitle string) {
	note.Title = newTitle
	noteUpdated(note)
}

//UpdateMemo updates the memo (also updates the LastUpdate value)
func (note *Note) UpdateMemo(newMemo string) {
	note.Memo = newMemo
	noteUpdated(note)
}

//AddTags add one or more tags to note (also updates the LastUpdate value)
func (note *Note) AddTags(tags []string) {
	for _, tag := range tags {
		note.Tags[tag] = true
	}
	noteUpdated(note)
}

//RemoveTags removes tags (also updates the LastUpdate value)
func (note *Note) RemoveTags(tags []string) {
	for _, tag := range tags {
		delete(note.Tags, tag)
	}
	noteUpdated(note)
}

//UpdateNotebook updates the notebook this note belongs to
//(also updates the LastUpdate value)
func (note *Note) UpdateNotebook(notebook *Notebook) {
	note.Notebook = notebook
	noteUpdated(note)
}

//Update the LastUpdated value with the current time
func noteUpdated(note *Note) {
	note.LastUpdated = time.Now()
}
