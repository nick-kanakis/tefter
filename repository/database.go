package repository

import (
	_ "github.com/mattn/go-sqlite3"
	"github.com/nicolasmanic/tefter/model"
)

//NoteRepository is an interface for handling DB related tasks for Note
type NoteRepository interface {
	SaveNote(note *model.Note) (int64, error)
	GetNotes(noteIDs []int64) ([]*model.Note, error)
	GetNote(noteID int64) (*model.Note, error)
	GetNotesByTag(tags []string) ([]*model.Note, error)
	UpdateNote(note *model.Note) error
	DeleteNotes(noteIDs []int64) error
	DeleteNote(noteIDs int64) error
	SearchNotesByKeyword(keyword string) ([]*model.Note, error)
	CloseDB() error
}

//NotebookRepository is an interface for handling DB related tasks for Notebook
type NotebookRepository interface {
	SaveNotebook(notebook *model.Notebook) (int64, error)
	GetNotebooks(notebooksIDs []int64) ([]*model.Notebook, error)
	GetNotebook(notebooksID int64) (*model.Notebook, error)
	GetNotebookByTitle(notebookTitle string) (*model.Notebook, error)
	GetAllNotebooksTitle() (map[int64]string, error)
	UpdateNotebook(notebook *model.Notebook) error
	DeleteNotebooks(notebooksIDs []int64) error
	DeleteNotebook(notebooksID int64) error
	CloseDB() error
}

//AccountRepository ia an interface for handling DB related tasks fro Account
type AccountRepository interface {
	CreateAccount(username, password string) error
	GetAccount(username string) (*model.Account, error)
	DeleteAccount(username string) error
	CloseDB() error
}
