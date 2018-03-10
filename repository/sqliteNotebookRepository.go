package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"github.com/nicolasmanic/tefter/model"
)

type sqliteNotebookRepository struct {
	dbPath string
	*sqlx.DB
}

//NewNotebookRepository returns a NotebookRepository interface
func NewNotebookRepository(dbPath string) NotebookRepository {
	db := connect2DB(dbPath)
	return &sqliteNotebookRepository{dbPath, db}
}

func (notebookRepo *sqliteNotebookRepository) SaveNotebook(notebook *model.Notebook) (notebookID int64, err error) {
	if notebook.Title == "" {
		return -1, fmt.Errorf("Notebook should contain title")
	}

	tx, err := notebookRepo.Beginx()
	if err != nil {
		return -1, err
	}

	defer func() {
		if r := recover(); r != nil {
			panicErr, _ := r.(error)
			tx.Rollback()
			notebookID = -1
			err = panicErr
		}
	}()

	result := tx.MustExec(`INSERT INTO notebook (title)	VALUES(?)`,
		notebook.Title)
	notebookID, err = result.LastInsertId()
	checkError(err)
	notebook.ID = notebookID

	err = tx.Commit()
	checkError(err)
	return notebookID, err
}

func (notebookRepo *sqliteNotebookRepository) GetNotebooks(notebooksIDs []int64) (notebooks []*model.Notebook, err error) {
	notebooksIDs = removeDups(notebooksIDs)
	if len(notebooksIDs) == 0 {
		return []*model.Notebook{}, nil
	}

	selectNotebook := "SELECT id, title FROM notebook "
	whereIdIn := "WHERE id IN ("
	args := []interface{}{}

	for _, id := range notebooksIDs {
		whereIdIn = whereIdIn + "?,"
		args = append(args, id)
	}
	whereIdIn = whereIdIn[:len(whereIdIn)-1]
	whereIdIn = whereIdIn + ")"
	querynotebook := selectNotebook + whereIdIn

	err = notebookRepo.Select(&notebooks, querynotebook, args...)
	checkError(err)

	noteRepo := NewNoteRepository(notebookRepo.dbPath)
	defer noteRepo.CloseDB()

	for _, notebook := range notebooks {
		noteIDs := notebookRepo.getNoteIDs(notebook.ID)
		notes, err := noteRepo.GetNotes(noteIDs)
		checkError(err)
		notebook.Notes = make(map[int64]*model.Note)
		for _, note := range notes {
			notebook.Notes[note.ID] = note
		}
	}

	return notebooks, err
}

func (notebookRepo *sqliteNotebookRepository) GetNotebook(notebookID int64) (*model.Notebook, error) {
	notebooks, err := notebookRepo.GetNotebooks([]int64{notebookID})
	checkError(err)

	if len(notebooks) != 1 {
		return nil, fmt.Errorf("Could find notebook with id: %v", notebookID)
	}
	return notebooks[0], err
}

func (notebookRepo *sqliteNotebookRepository) UpdateNotebook(notebook *model.Notebook) (err error) {
	if notebook.Title == "" {
		return fmt.Errorf("Notebook should contain title")
	}

	tx, err := notebookRepo.Beginx()
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

	updateNotebookQuery := `UPDATE notebook SET	title = ? WHERE id = ?`
	tx.MustExec(updateNotebookQuery, notebook.Title, notebook.ID)

	err = tx.Commit()
	checkError(err)
	return err
}

func (notebookRepo *sqliteNotebookRepository) DeleteNotebooks(notebooksIDs []int64) error {
	notebooksIDs = removeDups(notebooksIDs)
	if len(notebooksIDs) == 0 {
		return nil
	}

	whereIdIn := "WHERE id IN ("
	args := []interface{}{}

	for _, id := range notebooksIDs {
		whereIdIn = whereIdIn + "?,"
		args = append(args, id)
	}
	whereIdIn = whereIdIn[:len(whereIdIn)-1]
	whereIdIn = whereIdIn + ")"
	deleteNotebook := "DELETE FROM notebook " + whereIdIn

	tx, err := notebookRepo.Beginx()
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

	tx.MustExec(deleteNotebook, args...)
	err = tx.Commit()
	checkError(err)

	noteRepo := NewNoteRepository(notebookRepo.dbPath)
	defer noteRepo.CloseDB()
	for _, notebookID := range notebooksIDs {
		noteIDs := notebookRepo.getNoteIDs(notebookID)
		err := noteRepo.DeleteNotes(noteIDs)
		checkError(err)
	}

	return err
}

func (notebookRepo *sqliteNotebookRepository) DeleteNotebook(notebookID int64) error {
	return notebookRepo.DeleteNotebooks([]int64{notebookID})
}

func (notebookRepo *sqliteNotebookRepository) CloseDB() error {
	return notebookRepo.Close()
}

func (notebookRepo *sqliteNotebookRepository) getNoteIDs(notebookID int64) []int64 {
	query := "SELECT note_id FROM notebook_note WHERE notebook_id = ?"
	noteIDs := []int64{}
	err := notebookRepo.Select(&noteIDs, query, []interface{}{notebookID}...)
	checkError(err)
	return noteIDs
}
