package cmd

import (
	"errors"
	"github.com/nicolasmanic/tefter/model"
	"github.com/nicolasmanic/tefter/repository"
	"reflect"
	"testing"
)

func TestImportNotes(t *testing.T) {
	cases := []struct {
		mockNotebookDB mockNotebookDBImport
		mockNoteDB     mockNoteDBImport
		fsr            fakeFileSystemReader
		path           string
		expectedErr    error
	}{
		{
			mockNotebookDB: mockNotebookDBImport{},
			mockNoteDB:     mockNoteDBImport{},
			fsr: fakeFileSystemReader{
				rawBytes: []byte{},
				err:      errors.New("Unexpected error"),
			},
			path:        "test",
			expectedErr: errors.New("Error while reading file, error msg: Unexpected error"),
		}, {
			mockNotebookDB: mockNotebookDBImport{},
			mockNoteDB:     mockNoteDBImport{},
			fsr: fakeFileSystemReader{
				rawBytes: []byte("unmarshable text"),
				err:      nil,
			},
			path:        "test",
			expectedErr: errors.New("Could not unmarshal file at path: test, error msg: invalid character 'u' looking for beginning of value"),
		}, {
			mockNotebookDB: mockNotebookDBImport{
				err: errors.New("Unexpected error"),
			},
			mockNoteDB: mockNoteDBImport{},
			fsr: fakeFileSystemReader{
				rawBytes: []byte(`[{"id":3,"title":"title1","memo":"test\r\n","created":"2018-04-28T20:15:34.0146423+03:00","updated":"2018-04-28T20:15:34.0146423+03:00","tags":["tag1"],"notebook_title":"Default Notebook"}]`),
			},
			path:        "test",
			expectedErr: errors.New("Unexpected error"),
		},
	}

	for _, c := range cases {
		oldNotebookDB := NotebookDB
		oldNoteDB := NoteDB
		NoteDB = c.mockNoteDB
		NotebookDB = c.mockNotebookDB

		defer func() {
			NotebookDB = oldNotebookDB
			NoteDB = oldNoteDB
		}()

		err := importNotes(c.fsr, c.path)
		if !reflect.DeepEqual(c.expectedErr, err) {
			t.Errorf("Expected err to be %q but it was %q", c.expectedErr, err)
		}
	}
}

type fakeFileSystemReader struct {
	rawBytes []byte
	err      error
}

func (fsr fakeFileSystemReader) ReadFile(filepath string) ([]byte, error) {
	return fsr.rawBytes, fsr.err
}

type mockNotebookDBImport struct {
	repository.NotebookRepository
	notebook       *model.Notebook
	notebookTitles map[int64]string
	err            error
}

type mockNoteDBImport struct {
	repository.NoteRepository
	notes []*model.Note
	id    int64
	err   error
}

func (mDB mockNotebookDBImport) GetNotebookByTitle(notebooksTitle string) (*model.Notebook, error) {
	return mDB.notebook, mDB.err
}

func (mDB mockNotebookDBImport) GetAllNotebooksTitle() (map[int64]string, error) {
	return mDB.notebookTitles, mDB.err
}

func (mDB mockNoteDBImport) GetNotesByTag(tags []string) ([]*model.Note, error) {
	return mDB.notes, mDB.err
}

func (mDB mockNoteDBImport) SaveNote(note *model.Note) (int64, error) {
	return 1, nil
}

func (mDB mockNoteDBImport) GetNotes(noteIDs []int64) ([]*model.Note, error) {
	return mDB.notes, mDB.err
}
