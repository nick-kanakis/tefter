package cmd

import (
	"errors"
	"github.com/nicolasmanic/tefter/model"
	"github.com/nicolasmanic/tefter/repository"
	"reflect"
	"testing"
)

func TestSearch(t *testing.T) {
	cases := []struct {
		noteDB      mockNoteDBSearch
		keyword     string
		expectedErr error
	}{
		{
			noteDB: mockNoteDBSearch{
				notes: []*model.Note{
					model.NewNote("testTitle1", "testMemo1", repository.DEFAULT_NOTEBOOK_ID, []string{}),
				},
				err: nil,
			},
			keyword:     "",
			expectedErr: nil,
		}, {
			noteDB: mockNoteDBSearch{
				err: errors.New("Unexpected Errors"),
			},
			keyword:     "ab",
			expectedErr: errors.New("Error retrieving Notes from DB, error msg: Unexpected Errors"),
		},
	}

	for _, c := range cases {
		oldNoteDB := NoteDB
		NoteDB = c.noteDB
		//Restore interface
		defer func() {
			NoteDB = oldNoteDB
		}()

		_, err := search(c.keyword)
		if !reflect.DeepEqual(c.expectedErr, err) {
			t.Errorf("Expected err to be %q but it was %q", c.expectedErr, err)
		}
	}
}

type mockNoteDBSearch struct {
	repository.NoteRepository
	notes []*model.Note
	err   error
}

func (mDB mockNoteDBSearch) GetNotes(noteIDs []int64) ([]*model.Note, error) {
	return mDB.notes, mDB.err
}

func (mDB mockNoteDBSearch) SearchNotesByKeyword(keyword string) ([]*model.Note, error) {
	return mDB.notes, mDB.err
}
