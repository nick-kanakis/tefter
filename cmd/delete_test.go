package cmd

import (
	"errors"
	"github.com/nicolasmanic/tefter/repository"
	"reflect"
	"testing"
)

func TestDeleteArgs(t *testing.T) {
	cases := []struct {
		args        []string
		expectedErr error
		mDB         mockNoteDBDelete
	}{
		{
			args:        []string{"1", "2"},
			expectedErr: nil,
			mDB: mockNoteDBDelete{
				err: nil,
			},
		}, {
			args:        []string{},
			expectedErr: nil,
			mDB: mockNoteDBDelete{
				err: nil,
			},
		}, {
			args:        []string{"a"},
			expectedErr: errors.New("Could note transform input to id for argument: a"),
			mDB:         mockNoteDBDelete{},
		}, {
			args:        []string{"1", "2"},
			expectedErr: errors.New("Error while deleting notes, error msg: Unexpected error"),
			mDB: mockNoteDBDelete{
				err: errors.New("Unexpected error"),
			},
		},
	}

	for _, c := range cases {
		oldNoteDB := NoteDB
		NoteDB = c.mDB
		//Restore interface
		defer func() {
			NoteDB = oldNoteDB
		}()

		err := deleteArgs(c.args)
		if !reflect.DeepEqual(c.expectedErr, err) {
			t.Errorf("Expected err to be %q but it was %q", c.expectedErr, err)
		}
	}
}

type mockNoteDBDelete struct {
	repository.NoteRepository
	err error
}

func (mDB mockNoteDBDelete) DeleteNotes(noteIDs []int64) error {
	return mDB.err
}
