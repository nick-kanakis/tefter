package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/nicolasmanic/tefter/model"
	"time"
)

type sqliteNoteRepository struct {
	*sqlx.DB
}

//NewNoteRepository returns a NoteRepository interface.
func NewNoteRepository(dbPath string) NoteRepository {
	db := connect2DB(dbPath)
	return &sqliteNoteRepository{db}
}

//SaveNote validate note and if everything is as expected persist the object
func (noteRepo *sqliteNoteRepository) SaveNote(note *model.Note) (noteID int64, err error) {
	if note.Title == "" {
		return -1, fmt.Errorf("Note should contain title")
	}
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
		note.NotebookID = 1
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

func (noteRepo *sqliteNoteRepository) GetNotes(noteIDs []int64) (notes []*model.Note, err error) {
	noteIDs = removeDups(noteIDs)
	if len(noteIDs) == 0 {
		return []*model.Note{}, nil
	}
	selectNote := "SELECT id, title, memo, created, lastUpdated, notebook_id FROM note "
	whereNote := "WHERE id IN ("
	args := []interface{}{}

	for _, id := range noteIDs {
		whereNote = whereNote + "?,"
		args = append(args, id)
	}
	whereNote = whereNote[:len(whereNote)-1]
	whereNote = whereNote + ")"
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

func (noteRepo *sqliteNoteRepository) GetNote(noteID int64) (note *model.Note, err error) {
	notes, err := noteRepo.GetNotes([]int64{noteID})
	checkError(err)
	if len(notes) != 1 {
		return nil, fmt.Errorf("Could find note with id: %v", noteID)
	}
	return notes[0], err
}

func (noteRepo *sqliteNoteRepository) UpdateNote(note *model.Note) (err error) {
	if note.Title == "" {
		return fmt.Errorf("Note should contain title")
	}
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
	updateNoteNotebook := `UPDATE notebook_note SET notebook_id = ?
						   WHERE note_id = ?`

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

	tx.MustExec(updateNoteNotebook,
		note.NotebookID,
		note.ID)

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
	whereIdIn := " WHERE id IN ("
	whereNoteIdIn := " WHERE note_id IN ("
	args := []interface{}{}
	for _, id := range noteIDs {
		args = append(args, id)
		whereIdIn += "?,"
		whereNoteIdIn += "?,"
	}

	whereIdIn = whereIdIn[:len(whereIdIn)-1]
	whereIdIn = whereIdIn + ")"
	whereNoteIdIn = whereNoteIdIn[:len(whereNoteIdIn)-1]
	whereNoteIdIn = whereNoteIdIn + ")"

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

	deleteNoteQuery := "DELETE FROM note " + whereIdIn
	deleteTagQuery := "DELETE FROM note_tag " + whereNoteIdIn
	deleteNoteNotebookQuery := "DELETE FROM notebook_note " + whereNoteIdIn

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

func (noteRepo *sqliteNoteRepository) SearchNotesByKeyword(keyword string) (notes []*model.Note, err error) {
	if keyword == "" {
		return nil, fmt.Errorf("Empty search parameter")
	}
	query := `SELECT n.id, n.title, n.memo, n.created, n.lastUpdated, n.notebook_id FROM note n 
			  INNER JOIN note_fts nfs ON n.id = nfs.docid WHERE note_fts MATCH ?`

	err = noteRepo.Select(&notes, query, []interface{}{keyword}...)
	checkError(err)
	return notes, err
}

func (noteRepo *sqliteNoteRepository) SearchNotesByTag(tags []string) (notes []*model.Note, err error) {
	selectNote := `SELECT id, title, memo, created, lastUpdated, notebook_id FROM note n 
				   INNER JOIN note_tag nt ON n.id = nt.note_id `
	whereNote := "WHERE nt.tag IN ("
	args := []interface{}{}

	for _, tag := range tags {
		whereNote = whereNote + "?,"
		args = append(args, tag)
	}

	whereNote = whereNote[:len(whereNote)-1]
	whereNote = whereNote + ") ORDER BY n.created"

	queryNote := selectNote + whereNote
	err = noteRepo.Select(&notes, queryNote, args...)
	checkError(err)

	return notes, err
}

func (noteRepo *sqliteNoteRepository) CloseDB() error {
	return noteRepo.Close()
}
