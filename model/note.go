package model

import "time"

//Note is the note we want to keep.
type Note struct {
	ID          int64 `db:"id"`
	Title       string `db:"title"`
	Memo        string `db:"memo"`
	Created     time.Time `db:"created"`
	LastUpdated time.Time `db:"lastUpdated"`
	Tags        map[string]bool
	NotebookID  int64 `db:"notebook_id"`
}

//NewNote returns a new note pointer.
func NewNote(title, memo string, notebookID int64, tags []string) *Note {
	note := &Note{
		Title:       title,
		Memo:        memo,
		Created:     time.Now(),
		LastUpdated: time.Now(),
		Tags:        make(map[string]bool),
		NotebookID:  notebookID,
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
func (note *Note) UpdateNotebook(notebookID int64) {
	note.NotebookID = notebookID
	noteUpdated(note)
}

//Update the LastUpdated value with the current time
func noteUpdated(note *Note) {
	note.LastUpdated = time.Now()
}
