package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/nicolasmanic/tefter/model"
	"time"
)

//DEFAULT_NOTEBOOK_ID is always set to 1. It defines the default notebook that is present at the DB. This notebook
//should never be deleted.
const DEFAULT_NOTEBOOK_ID = 1

type sqliteNoteRepository struct {
	*sqlx.DB
}

//NewNoteRepository returns a NoteRepository interface.
func NewNoteRepository(dbPath string) NoteRepository {
	db := connect2DB(dbPath)
	return &sqliteNoteRepository{db}
}

//SaveNote persist a note to DB. For a note to be valid the memo field must not be empty.
//All other fields can be auto-completed.
//Default values of note fields are:
// title: ""
//Created: current time
//LastUpdated: current time
//NotepadId: 1 (Default notepad)
func (noteRepo *sqliteNoteRepository) SaveNote(note *model.Note) (noteID int64, err error) {
	if note.Memo == "" {
		return -1, fmt.Errorf("Note should contain memo")
	}
	if note.Created.IsZero() {
		note.Created = time.Now().UTC()
	}
	if note.LastUpdated.IsZero() {
		note.LastUpdated = time.Now().UTC()
	}
	//If notebook id is 0 set it to the default notebook.
	if note.NotebookID == 0 {
		note.NotebookID = DEFAULT_NOTEBOOK_ID
	}

	tx, err := noteRepo.Beginx()
	if err != nil {
		return -1, err
	}

	defer func() {
		if r := recover(); r != nil {
			panicErr, _ := r.(error)
			tx.Rollback()
			noteID = -1
			err = panicErr
		}
	}()

	result := tx.MustExec(`INSERT INTO note (
		title, memo, created, lastUpdated, notebook_id) 
		VALUES(?, ?, ?, ?, ?)`,
		note.Title,
		note.Memo,
		note.Created,
		note.LastUpdated,
		note.NotebookID)

	noteID, err = result.LastInsertId()
	checkError(err)
	note.ID = noteID

	tx.MustExec(`INSERT INTO notebook_note (note_id, notebook_id)
	VALUES (?, ?)`, note.ID, note.NotebookID)

	tagInsertStmt, err := tx.Preparex(`INSERT INTO note_tag (note_id, tag) VALUES(?,?)`)
	checkError(err)

	for tag := range note.Tags {
		tagInsertStmt.MustExec(noteID, tag)
	}

	err = tx.Commit()
	checkError(err)

	return noteID, err
}

//GetNotes return a slice of notes based on the given slice of ids,
//if ids slice is empty all notes are returned
func (noteRepo *sqliteNoteRepository) GetNotes(noteIDs []int64) (notes []*model.Note, err error) {
	noteIDs = removeDups(noteIDs)
	selectNote := "SELECT id, title, memo, created, lastUpdated, notebook_id FROM note "
	var whereNote string
	args := []interface{}{}
	if len(noteIDs) != 0 {
		whereNote = "WHERE id IN ("
		for _, id := range noteIDs {
			whereNote = whereNote + "?,"
			args = append(args, id)
		}
		whereNote = whereNote[:len(whereNote)-1]
		whereNote = whereNote + ") ORDER BY created desc"
	} else {
		whereNote = "WHERE 1 ORDER BY created desc"
	}

	querynote := selectNote + whereNote
	err = noteRepo.Select(&notes, querynote, args...)
	checkError(err)

	selectTagStmt, err := noteRepo.Preparex("SELECT tag FROM note_tag WHERE note_id = ?")
	checkError(err)

	for _, note := range notes {
		tags := []string{}
		err = selectTagStmt.Select(&tags, note.ID)
		checkError(err)
		note.Tags = make(map[string]bool)
		for _, tag := range tags {
			note.Tags[tag] = true
		}
	}

	return notes, err
}

//GetNote returns a single note based on an id, returns error if note with id doesn't exist
func (noteRepo *sqliteNoteRepository) GetNote(noteID int64) (note *model.Note, err error) {
	notes, err := noteRepo.GetNotes([]int64{noteID})
	checkError(err)
	if len(notes) != 1 {
		return nil, fmt.Errorf("Could find note with id: %v", noteID)
	}
	return notes[0], err
}

//UpdateNote updates an existing note. For a note to be valid the memo field must not be empty.
func (noteRepo *sqliteNoteRepository) UpdateNote(note *model.Note) (err error) {
	if note.Memo == "" {
		return fmt.Errorf("Note should contain memo")
	}
	if note.Created.IsZero() {
		note.Created = time.Now().UTC()
	}
	if note.LastUpdated.IsZero() {
		note.LastUpdated = time.Now().UTC()
	}

	tx, err := noteRepo.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			panicErr, _ := r.(error)
			tx.Rollback()
			err = panicErr
		}
	}()

	updateNoteQuery := `UPDATE note SET
		title = ?, memo = ?, created = ?, lastUpdated = ?, notebook_id =?  
		WHERE id = ?`
	deleteNoteNotebook := `DELETE FROM notebook_note WHERE note_id = ?`
	insertNoteNotebook := `INSERT INTO notebook_note (note_id, notebook_id) VALUES (?, ?)`

	deleteNoteTagQuery := `DELETE FROM note_tag WHERE note_id = ?`
	insertNoteTagStmt, err := tx.Preparex(`INSERT INTO note_tag (note_id, tag) VALUES(?,?)`)
	checkError(err)

	tx.MustExec(updateNoteQuery,
		note.Title,
		note.Memo,
		note.Created,
		note.LastUpdated,
		note.NotebookID,
		note.ID)

	tx.MustExec(deleteNoteNotebook, note.ID)

	tx.MustExec(insertNoteNotebook,
		note.ID,
		note.NotebookID)

	tx.MustExec(deleteNoteTagQuery, note.ID)

	for tag := range note.Tags {
		insertNoteTagStmt.MustExec(note.ID, tag)
	}

	err = tx.Commit()
	checkError(err)
	return err
}

func (noteRepo *sqliteNoteRepository) DeleteNotes(noteIDs []int64) (err error) {
	noteIDs = removeDups(noteIDs)
	whereIDIn := " WHERE id IN ("
	whereNoteIDIn := " WHERE note_id IN ("
	args := []interface{}{}
	for _, id := range noteIDs {
		args = append(args, id)
		whereIDIn += "?,"
		whereNoteIDIn += "?,"
	}

	whereIDIn = whereIDIn[:len(whereIDIn)-1]
	whereIDIn = whereIDIn + ")"
	whereNoteIDIn = whereNoteIDIn[:len(whereNoteIDIn)-1]
	whereNoteIDIn = whereNoteIDIn + ")"

	tx, err := noteRepo.Beginx()
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			panicErr, _ := r.(error)
			tx.Rollback()
			err = panicErr
		}
	}()

	deleteNoteQuery := "DELETE FROM note " + whereIDIn
	deleteTagQuery := "DELETE FROM note_tag " + whereNoteIDIn
	deleteNoteNotebookQuery := "DELETE FROM notebook_note " + whereNoteIDIn

	tx.MustExec(deleteNoteQuery, args...)
	tx.MustExec(deleteTagQuery, args...)
	tx.MustExec(deleteNoteNotebookQuery, args...)

	err = tx.Commit()
	checkError(err)
	return err
}

func (noteRepo *sqliteNoteRepository) DeleteNote(noteID int64) (err error) {
	return noteRepo.DeleteNotes([]int64{noteID})
}

//SearchNotesByKeyword searches the DB for notes containing the keyword. Keyword cannot be empty
//also keyword must be a complete word, partial words can not be matched
func (noteRepo *sqliteNoteRepository) SearchNotesByKeyword(keyword string) (notes []*model.Note, err error) {
	if keyword == "" {
		return nil, fmt.Errorf("Empty search parameter")
	}
	query := `SELECT n.id, n.title, n.memo, n.created, n.lastUpdated, n.notebook_id FROM note n 
			  INNER JOIN note_fts nfs ON n.id = nfs.docid WHERE note_fts MATCH ? ORDER BY n.created desc`

	err = noteRepo.Select(&notes, query, []interface{}{keyword}...)
	checkError(err)
	selectTagStmt, err := noteRepo.Preparex("SELECT tag FROM note_tag WHERE note_id = ?")
	checkError(err)

	for _, note := range notes {
		tags := []string{}
		err = selectTagStmt.Select(&tags, note.ID)
		checkError(err)
		note.Tags = make(map[string]bool)
		for _, tag := range tags {
			note.Tags[tag] = true
		}
	}
	return notes, err
}

//GetNotesByTag returns all notes tagged with one or more of tags given as inputs
func (noteRepo *sqliteNoteRepository) GetNotesByTag(tags []string) (notes []*model.Note, err error) {
	selectNote := `SELECT id, title, memo, created, lastUpdated, notebook_id FROM note n 
				   INNER JOIN note_tag nt ON n.id = nt.note_id `
	whereNote := "WHERE nt.tag IN ("
	args := []interface{}{}

	for _, tag := range tags {
		whereNote = whereNote + "?,"
		args = append(args, tag)
	}

	whereNote = whereNote[:len(whereNote)-1]
	whereNote = whereNote + ") ORDER BY n.created desc"

	queryNote := selectNote + whereNote
	err = noteRepo.Select(&notes, queryNote, args...)
	checkError(err)

	selectTagStmt, err := noteRepo.Preparex("SELECT tag FROM note_tag WHERE note_id = ?")
	checkError(err)

	for _, note := range notes {
		tags := []string{}
		err = selectTagStmt.Select(&tags, note.ID)
		checkError(err)
		note.Tags = make(map[string]bool)
		for _, tag := range tags {
			note.Tags[tag] = true
		}
	}

	return notes, err
}

func (noteRepo *sqliteNoteRepository) CloseDB() error {
	return noteRepo.Close()
}
