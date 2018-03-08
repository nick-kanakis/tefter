package repository

import (
	"fmt"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/nicolasmanic/tefter/model"
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
	note.ID = noteID
	checkError(err)

	for tag := range note.Tags {
		sanitizedTag := strings.TrimSpace(tag)
		sanitizedTag = strings.ToLower(sanitizedTag)
		tx.MustExec(`INSERT INTO note_tag (
			note_id, tag)
			VALUES(?,?)`,
			noteID,
			sanitizedTag)
	}

	err = tx.Commit()
	checkError(err)

	return noteID, err
}

func (noteRepo *sqliteNoteRepository) GetNotes(noteIDs []int64) ( notes []*model.Note, err error) {
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

	queryTag := "SELECT tag FROM note_tag WHERE note_id = ?"
	for _, note := range notes {
		tags := []string{}
		err = noteRepo.Select(&tags, queryTag, note.ID)
		checkError(err)
		note.Tags = make (map[string]bool)
		for _, tag := range tags {
			note.Tags[tag] = true
		}
	}

	return notes, err
}

func (noteRepo *sqliteNoteRepository) DeleteNotes(noteIDs []int64) error {
	return nil
}

func (noteRepo *sqliteNoteRepository) SearchNotesByKeyword(keyword string) ([]*model.Note, error) {
	return nil, nil
}

func (noteRepo *sqliteNoteRepository) SearchNoteByTag(tags []string) ([]*model.Note, error) {
	return nil, nil
}

func (noteRepo *sqliteNoteRepository) UpdateNote(note *model.Note) (*model.Note, error) {
	return nil, nil
}

func (noteRepo *sqliteNoteRepository) CloseDB() error{
	return noteRepo.Close()
}
