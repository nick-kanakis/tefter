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

	result := tx.MustExec(`INSERT INTO notebook (title)	VALUES(?)`, notebook.Title)
	notebookID, err = result.LastInsertId()
	checkError(err)
	notebook.ID = notebookID

	err = tx.Commit()
	checkError(err)
	return notebookID, err
}

func (notebookRepo *sqliteNotebookRepository) GetNotebooks(notebooksIDs []int64) (notebooks []*model.Notebook, err error) {
	notebooksIDs = removeDups(notebooksIDs)

	selectNotebook := "SELECT id, title FROM notebook "
	whereIDIn := "WHERE id IN ("
	args := []interface{}{}
	if len(notebooksIDs) == 0 {
		whereIDIn = "WHERE 1"
	} else {
		for _, id := range notebooksIDs {
			whereIDIn = whereIDIn + "?,"
			args = append(args, id)
		}
		whereIDIn = whereIDIn[:len(whereIDIn)-1]
		whereIDIn = whereIDIn + ")"
	}

	querynotebook := selectNotebook + whereIDIn
	err = notebookRepo.Select(&notebooks, querynotebook, args...)
	checkError(err)

	//Use note repository to get the all notes of this notebook.
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

func (notebookRepo *sqliteNotebookRepository) GetNotebookByTitle(notebookTitle string) (notebook *model.Notebook, err error) {

	query := "SELECT id, title FROM notebook WHERE title = ? "
	var notebooks []*model.Notebook
	err = notebookRepo.Select(&notebooks, query, []interface{}{notebookTitle}...)

	if len(notebooks) > 1 {
		return nil, fmt.Errorf(`Found more than one notebooks with the same title,
			 this should not have happened since titles are unique in DB`)
	}

	checkError(err)
	//Use note repository to get the all notes of this notebook.
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

	if len(notebooks) > 0 {
		return notebooks[0], err
	}

	return nil, err
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

	whereIDIn := "WHERE id IN ("

	args := []interface{}{}
	for _, id := range notebooksIDs {
		whereIDIn = whereIDIn + "?,"
		args = append(args, id)
	}
	whereIDIn = whereIDIn[:len(whereIDIn)-1]
	whereIDIn = whereIDIn + ")"
	deleteNotebook := "DELETE FROM notebook " + whereIDIn

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

	//use note repository to delete all notes of deleted notebooks.
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

func (notebookRepo *sqliteNotebookRepository) GetAllNotebooksTitle() (map[int64]string, error) {
	selectNotebook := "SELECT id, title FROM notebook"
	var notebooks = []model.Notebook{}
	err := notebookRepo.Select(&notebooks, selectNotebook, []interface{}{}...)
	checkError(err)

	var notebookNamesMap = make(map[int64]string)
	for _, notebook := range notebooks {
		notebookNamesMap[notebook.ID] = notebook.Title
	}

	return notebookNamesMap, err
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
