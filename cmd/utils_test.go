package cmd

import (
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
	oldNotebookDB := NotebookDB
	oldNoteDB := NoteDB
	NoteDB = mockNoteDBUtils{}
	NotebookDB = mockNotebookDBUtils{}
	//Restore interface
	defer func() {
		NotebookDB = oldNotebookDB
		NoteDB = oldNoteDB
	}()

	noteMap := collectNotesFromDB([]int{1}, []string{"notebookTitle"}, []string{"testTag"}, false)
	if len(noteMap) != 4 {
		t.Error("Error while collecting notes per id, notebook title and tags")
	}
	allNotes := collectNotesFromDB([]int{}, []string{}, []string{}, true)
	if len(allNotes) != 1 {
		t.Error("Error while collectingall notes")
	}
}

type mockNotebookDBUtils struct {
	repository.NotebookRepository
}

type mockNoteDBUtils struct {
	repository.NoteRepository
}

func (mDB mockNoteDBUtils) GetNotesByTag(tags []string) ([]*model.Note, error) {
	note1 := model.NewNote("testTitle", "testMemo", repository.DEFAULT_NOTEBOOK_ID, []string{})
	note1.ID = 1
	note2 := model.NewNote("testTitle2", "testMemo2", repository.DEFAULT_NOTEBOOK_ID, []string{})
	note2.ID = 2
	return []*model.Note{note1, note2}, nil
}

func (mDB mockNotebookDBUtils) GetNotebookByTitle(notebooksTitle string) (*model.Notebook, error) {
	notebook := model.NewNotebook(notebooksTitle)
	notebook.ID = 1
	note := model.NewNote("testTitle", "testMemo", notebook.ID, []string{})
	note.ID = 3
	notebook.AddNote(note)
	return notebook, nil
}

func (mDB mockNoteDBUtils) GetNotes(noteIDs []int64) ([]*model.Note, error) {
	note := model.NewNote("testTitle4", "testMemo", repository.DEFAULT_NOTEBOOK_ID, []string{})
	note.ID = 4
	return []*model.Note{note}, nil
}
