package repository

import (
	//FIXME: should be moved to main
	_"github.com/mattn/go-sqlite3"
	"github.com/nicolasmanic/tefter/model"
)

//FIXME: Is there a cleaner way to implement separation of conserns for Repository layer???

//NoteRepository is an interface for handling DB related tasks for Note
type NoteRepository interface {
	CreateNote(note model.Note) (int64, error)
	GetNotes(noteIDs ...int64) ([]model.Note, error)
	DeleteNotes(noteIDs ...int64) error
	SearchNotesByKeyword(keyword string) ([]model.Note, error)
	SearchNoteByTag(tags ...string) ([]model.Note, error)
	UpdateNote(note model.Note) (*model.Note, error)
}

//NotebookRepository is an interface for handling DB related tasks for Notebook
type NotebookRepository interface {
	CreateNotebook(notebook model.Notebook) (int64, error)
	GetNotebooks(notebooksIDs ...int64) ([]model.Notebook, error)
	DeleteNotebooks(notebooksIDs ...int64) error
	UpdateNotebook(notebook model.Notebook) (*model.Notebook, error)
}
