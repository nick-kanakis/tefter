package cmd

import (
	"errors"
	"github.com/nicolasmanic/tefter/model"
	"github.com/nicolasmanic/tefter/repository"
	"reflect"
	"testing"
)

func TestInt64Slice(t *testing.T) {
	tests := []struct {
		input    []int
		expected []int64
	}{
		{
			input:    []int{0, 1, 3},
			expected: []int64{0, 1, 3},
		},
		{
			input:    []int{-2, -3},
			expected: []int64{-2, -3},
		},
		{
			input:    []int{},
			expected: []int64{},
		},
	}

	for _, test := range tests {
		result := int64Slice(test.input)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("Expected: %v, got: %v", test.expected, result)
		}
	}
}

func TestNoteMap2slice(t *testing.T) {
	tests := []struct {
		input    map[int64]*model.Note
		expected int
	}{
		{
			input:    make(map[int64]*model.Note),
			expected: 0,
		}, {
			input: map[int64]*model.Note{
				1: model.NewNote("test1", "test1", 1, []string{"tag1"}),
				2: model.NewNote("test2", "test2", 1, []string{"tag2"}),
			},
			expected: 2,
		},
	}

	for _, test := range tests {
		result := noteMap2Slice(test.input)
		if test.expected != len(result) {
			t.Errorf("Expected slice size of: %v , got: %v", test.expected, len(result))
		}
	}
}

func TestTagMap2Slice(t *testing.T) {
	tests := []struct {
		input    map[string]bool
		expected int
	}{
		{
			input:    make(map[string]bool),
			expected: 0,
		},
		{
			input: map[string]bool{
				"test1": true,
				"test2": false,
			},
			expected: 2,
		},
	}

	for _, test := range tests {
		result := tagMap2Slice(test.input)
		if test.expected != len(result) {
			t.Errorf("Expected slice size of: %v , got: %v", test.expected, len(result))
		}
	}
}

func TestCollectNotesFromDB(t *testing.T) {
	cases := []struct {
		noteDB         mockNoteDBUtils
		notebookDB     mockNotebookDBUtils
		ids            []int
		notebookTitles []string
		tags           []string
		getAll         bool
		expectedErr    error
	}{
		{
			noteDB: mockNoteDBUtils{
				notes: []*model.Note{model.NewNote("testTitle", "testMemo", repository.DEFAULT_NOTEBOOK_ID, []string{})},
				err:   nil,
			},
			notebookDB: mockNotebookDBUtils{
				notebook: model.NewNotebook("title"),
				err:      nil,
			},
			ids:            []int{1},
			notebookTitles: []string{"title"},
			tags:           []string{"tag"},
			getAll:         false,
			expectedErr:    nil,
		}, {
			noteDB: mockNoteDBUtils{
				notes: nil,
				err:   errors.New("Unexpected error"),
			},
			notebookDB:     mockNotebookDBUtils{},
			ids:            []int{},
			notebookTitles: []string{},
			tags:           []string{},
			getAll:         true,
			expectedErr:    errors.New("Error while retrieving all notes, error msg: Unexpected error"),
		}, {
			noteDB: mockNoteDBUtils{
				notes: nil,
				err:   errors.New("Unexpected error"),
			},
			notebookDB:     mockNotebookDBUtils{},
			ids:            []int{1},
			notebookTitles: []string{},
			tags:           []string{},
			getAll:         false,
			expectedErr:    errors.New("Error while retrieving notes by id, error msg: Unexpected error"),
		}, {
			noteDB: mockNoteDBUtils{},
			notebookDB: mockNotebookDBUtils{
				notebook: nil,
				err:      errors.New("Unexpected error"),
			},
			ids:            []int{},
			notebookTitles: []string{"title"},
			tags:           []string{},
			getAll:         false,
			expectedErr:    errors.New("Error while retrieving notebook by title, error msg: Unexpected error"),
		}, {
			noteDB: mockNoteDBUtils{
				notes: nil,
				err:   errors.New("Unexpected error"),
			},
			notebookDB:     mockNotebookDBUtils{},
			ids:            []int{},
			notebookTitles: []string{},
			tags:           []string{"tags"},
			getAll:         false,
			expectedErr:    errors.New("Error while retrieving notes by tag, error msg: Unexpected error"),
		},
	}

	for _, c := range cases {
		oldNotebookDB := NotebookDB
		oldNoteDB := NoteDB
		NoteDB = c.noteDB
		NotebookDB = c.notebookDB
		//Restore interface
		defer func() {
			NotebookDB = oldNotebookDB
			NoteDB = oldNoteDB
		}()

		_, err := collectNotesFromDB(c.ids, c.notebookTitles, c.tags, c.getAll)
		if !reflect.DeepEqual(err, c.expectedErr) {
			t.Errorf("Expected err to be %q but it was %q", c.expectedErr, err)
		}
	}

}

type mockNotebookDBUtils struct {
	repository.NotebookRepository
	notebook *model.Notebook
	err      error
}

type mockNoteDBUtils struct {
	repository.NoteRepository
	notes []*model.Note
	err   error
}

func (mDB mockNoteDBUtils) GetNotesByTag(tags []string) ([]*model.Note, error) {
	return mDB.notes, mDB.err
}

func (mDB mockNoteDBUtils) GetNotes(noteIDs []int64) ([]*model.Note, error) {
	return mDB.notes, mDB.err
}

func (mDB mockNotebookDBUtils) GetNotebookByTitle(notebooksTitle string) (*model.Notebook, error) {
	return mDB.notebook, mDB.err
}
