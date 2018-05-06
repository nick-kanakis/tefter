package cmd

import (
	"bytes"
	"errors"
	"github.com/nicolasmanic/tefter/model"
	"github.com/nicolasmanic/tefter/repository"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

func TestAddNoteAPI(t *testing.T) {
	cases := []struct {
		checkTokenFunc   func(r *http.Request, signingKey []byte) error
		saveNoteFunc     func(*jsonNote) error
		payload          []byte
		expectedHTTPCode int
	}{
		{
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return errors.New("Unexpected Error")
			},
			expectedHTTPCode: http.StatusUnauthorized,
		}, {
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return nil
			},
			payload:          []byte(`incorrect json object`),
			expectedHTTPCode: http.StatusBadRequest,
		}, {
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return nil
			},
			saveNoteFunc: func(*jsonNote) error {
				return errors.New("Unexpected Error")
			},
			payload:          []byte(`{"title":"Shopping for weekend","memo":" Things for weekend:\n \u003e Milk\n \u003e Eggs\n \u003e Chicken breast\n","created":"2018-03-20T18:53:35.4123749+02:00","updated":"2018-03-20T18:53:35.4193801+02:00","tags":["weekend","list"],"notebook_title":"Shopping"}`),
			expectedHTTPCode: http.StatusInternalServerError,
		}, {
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return nil
			},
			saveNoteFunc: func(*jsonNote) error {
				return nil
			},
			payload:          []byte(`{"title":"Shopping for weekend","memo":" Things for weekend:\n \u003e Milk\n \u003e Eggs\n \u003e Chicken breast\n","created":"2018-03-20T18:53:35.4123749+02:00","updated":"2018-03-20T18:53:35.4193801+02:00","tags":["weekend","list"],"notebook_title":"Shopping"}`),
			expectedHTTPCode: http.StatusCreated,
		},
	}

	for _, c := range cases {
		originalSaveNote := saveNoteFunc
		originalCheckToken := checkTokenFunc
		saveNoteFunc = c.saveNoteFunc
		checkTokenFunc = c.checkTokenFunc
		defer func() {
			saveNoteFunc = originalSaveNote
			checkTokenFunc = originalCheckToken
		}()

		req, _ := http.NewRequest("POST", "/addNote", bytes.NewBuffer(c.payload))
		response := executeRequest(req)
		checkResponseCode(t, c.expectedHTTPCode, response.Code)
	}
}

func TestUpdateNoteAPI(t *testing.T) {
	cases := []struct {
		checkTokenFunc   func(r *http.Request, signingKey []byte) error
		updateNoteFunc   func(*jsonNote) error
		payload          []byte
		expectedHTTPCode int
	}{
		{
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return errors.New("Unexpected Error")
			},
			expectedHTTPCode: http.StatusUnauthorized,
		},
		{
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return nil
			},
			payload:          []byte(`incorrect json object`),
			expectedHTTPCode: http.StatusBadRequest,
		},
		{
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return nil
			},
			updateNoteFunc: func(*jsonNote) error {
				return errors.New("Unexpected Error")
			},
			payload:          []byte(`{"id":1, "title":"Shopping for weekend","memo":" Things for weekend:\n \u003e Milk\n \u003e Eggs\n \u003e Chicken breast\n","created":"2018-03-20T18:53:35.4123749+02:00","updated":"2018-03-20T18:53:35.4193801+02:00","tags":["weekend","list"],"notebook_title":"Shopping"}`),
			expectedHTTPCode: http.StatusInternalServerError,
		},
		{
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return nil
			},
			updateNoteFunc: func(*jsonNote) error {
				return nil
			},
			payload:          []byte(`{"id":1, "title":"Shopping for weekend","memo":" Things for weekend:\n \u003e Milk\n \u003e Eggs\n \u003e Chicken breast\n","created":"2018-03-20T18:53:35.4123749+02:00","updated":"2018-03-20T18:53:35.4193801+02:00","tags":["weekend","list"],"notebook_title":"Shopping"}`),
			expectedHTTPCode: http.StatusCreated,
		},
	}

	for _, c := range cases {
		originalUpdateNote := updateNoteFunc
		originalCheckToken := checkTokenFunc
		updateNoteFunc = c.updateNoteFunc
		checkTokenFunc = c.checkTokenFunc
		defer func() {
			updateNoteFunc = originalUpdateNote
			checkTokenFunc = originalCheckToken
		}()

		req, _ := http.NewRequest("PUT", "/updateNote", bytes.NewBuffer(c.payload))
		response := executeRequest(req)
		checkResponseCode(t, c.expectedHTTPCode, response.Code)
	}
}

func TestGetNotesAPI(t *testing.T) {
	cases := []struct {
		checkTokenFunc    func(r *http.Request, signingKey []byte) error
		retrieveNotesFunc func(ids []int, notebookTitles, tags []string, getAll bool) ([]*jsonNote, error)
		url               string
		params            string
		expectedHTTPCode  int
	}{
		{
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return errors.New("Unexpected Error")
			},
			url:              "/getNotesByID/",
			params:           "1,2,3",
			expectedHTTPCode: http.StatusUnauthorized,
		}, {
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return nil
			},
			url:              "/getNotesByID/",
			params:           "abc",
			expectedHTTPCode: http.StatusInternalServerError,
		}, {
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return nil
			},
			url:    "/getNotesByID/",
			params: "1",
			retrieveNotesFunc: func(ids []int, notebookTitles, tags []string, getAll bool) ([]*jsonNote, error) {
				return nil, errors.New("Unexpected Error")
			},
			expectedHTTPCode: http.StatusInternalServerError,
		}, {
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return nil
			},
			url:    "/getNotesByID/",
			params: "1,2,3",
			retrieveNotesFunc: func(ids []int, notebookTitles, tags []string, getAll bool) ([]*jsonNote, error) {
				return nil, errors.New("Unexpected Error")
			},
			expectedHTTPCode: http.StatusInternalServerError,
		}, {
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return nil
			},
			url:    "/getNotesByID/",
			params: "1,2,3",
			retrieveNotesFunc: func(ids []int, notebookTitles, tags []string, getAll bool) ([]*jsonNote, error) {
				return mockJSONNotes(), nil
			},
			expectedHTTPCode: http.StatusOK,
		}, {
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return nil
			},
			url:    "/getNotesByNotebookTitle/",
			params: "title1,title2",
			retrieveNotesFunc: func(ids []int, notebookTitles, tags []string, getAll bool) ([]*jsonNote, error) {
				return mockJSONNotes(), nil
			},
			expectedHTTPCode: http.StatusOK,
		}, {
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return nil
			},
			url:    "/getNotesByTags/",
			params: "tag1,tag2",
			retrieveNotesFunc: func(ids []int, notebookTitles, tags []string, getAll bool) ([]*jsonNote, error) {
				return mockJSONNotes(), nil
			},
			expectedHTTPCode: http.StatusOK,
		}, {
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return nil
			},
			url:    "/getAllNotes",
			params: "",
			retrieveNotesFunc: func(ids []int, notebookTitles, tags []string, getAll bool) ([]*jsonNote, error) {
				return mockJSONNotes(), nil
			},
			expectedHTTPCode: http.StatusOK,
		},
	}

	for _, c := range cases {
		originalRetrieveNotes := retrieveNotesFunc
		originalCheckToken := checkTokenFunc
		retrieveNotesFunc = c.retrieveNotesFunc
		checkTokenFunc = c.checkTokenFunc
		defer func() {
			retrieveNotesFunc = originalRetrieveNotes
			checkTokenFunc = originalCheckToken
		}()

		req, _ := http.NewRequest("GET", c.url+c.params, nil)
		response := executeRequest(req)
		checkResponseCode(t, c.expectedHTTPCode, response.Code)

	}
}

func TestDeleteNotesAPI(t *testing.T) {
	cases := []struct {
		checkTokenFunc   func(r *http.Request, signingKey []byte) error
		deleteFunc       func(ids []int64) error
		params           string
		expectedHTTPCode int
	}{
		{
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return errors.New("Unexpected Error")
			},
			params:           "abc",
			expectedHTTPCode: http.StatusUnauthorized,
		}, {
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return nil
			},
			params:           "abc",
			expectedHTTPCode: http.StatusInternalServerError,
		}, {
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return nil
			},
			params: "1,2",
			deleteFunc: func(ids []int64) error {
				return errors.New("Unexpected error")
			},
			expectedHTTPCode: http.StatusInternalServerError,
		}, {
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return nil
			},
			params: "1,2",
			deleteFunc: func(ids []int64) error {
				return nil
			},
			expectedHTTPCode: http.StatusOK,
		},
	}

	for _, c := range cases {
		originalDeleteNotes := deleteNotesFunc
		originalCheckToken := checkTokenFunc
		deleteNotesFunc = c.deleteFunc
		checkTokenFunc = c.checkTokenFunc
		defer func() {
			deleteNotesFunc = originalDeleteNotes
			checkTokenFunc = originalCheckToken
		}()

		req, _ := http.NewRequest("DELETE", "/deleteNotes/"+c.params, nil)
		response := executeRequest(req)
		checkResponseCode(t, c.expectedHTTPCode, response.Code)

	}
}

func TestDeleteNotebooksAPI(t *testing.T) {
	cases := []struct {
		checkTokenFunc      func(r *http.Request, signingKey []byte) error
		deleteNotebooksFunc func(titles []string) error
		params              string
		expectedHTTPCode    int
	}{
		{
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return errors.New("Unexpected Error")
			},
			params:           "title1,title2",
			expectedHTTPCode: http.StatusUnauthorized,
		}, {
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return nil
			},
			deleteNotebooksFunc: func(titles []string) error {
				return errors.New("Unexpected Error")
			},
			params:           "title1",
			expectedHTTPCode: http.StatusInternalServerError,
		}, {
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return nil
			},
			deleteNotebooksFunc: func(titles []string) error {
				return nil
			},
			params:           "title1",
			expectedHTTPCode: http.StatusOK,
		},
	}

	for _, c := range cases {
		originalDeleteNotebooks := deleteNotebooksFunc
		originalCheckToken := checkTokenFunc
		deleteNotebooksFunc = c.deleteNotebooksFunc
		checkTokenFunc = c.checkTokenFunc
		defer func() {
			deleteNotebooksFunc = originalDeleteNotebooks
			checkTokenFunc = originalCheckToken
		}()

		req, _ := http.NewRequest("DELETE", "/deleteNotebooks/"+c.params, nil)
		response := executeRequest(req)
		checkResponseCode(t, c.expectedHTTPCode, response.Code)

	}
}

func TestUpdateNotebooksAPI(t *testing.T) {
	cases := []struct {
		checkTokenFunc     func(r *http.Request, signingKey []byte) error
		updateNotebookFunc func(oldTitle, newTitle string) error
		expectedHTTPCode   int
	}{
		{
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return errors.New("Unexpected Error")
			},
			expectedHTTPCode: http.StatusUnauthorized,
		}, {
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return nil
			},
			updateNotebookFunc: func(oldTitle, newTitle string) error {
				return errors.New("Unexpected Error")
			},
			expectedHTTPCode: http.StatusInternalServerError,
		}, {
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return nil
			},
			updateNotebookFunc: func(oldTitle, newTitle string) error {
				return nil
			},
			expectedHTTPCode: http.StatusOK,
		},
	}

	for _, c := range cases {
		originalUpdateNotebook := updateNotebookFunc
		originalCheckToken := checkTokenFunc
		updateNotebookFunc = c.updateNotebookFunc
		checkTokenFunc = c.checkTokenFunc
		defer func() {
			updateNotebookFunc = originalUpdateNotebook
			checkTokenFunc = originalCheckToken
		}()
		req, _ := http.NewRequest("PUT", "/updateNotebook/oldTitle_1/newTitle_1", nil)
		response := executeRequest(req)
		checkResponseCode(t, c.expectedHTTPCode, response.Code)
	}

}

func TestSearchNotesAPI(t *testing.T) {
	cases := []struct {
		checkTokenFunc   func(r *http.Request, signingKey []byte) error
		searchNotesFunc  func(keyword string) ([]*model.Note, error)
		notebookDB       mockNotebookDBAPI
		expectedHTTPCode int
	}{
		{
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return errors.New("Unexpected Error")
			},
			expectedHTTPCode: http.StatusUnauthorized,
		}, {
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return nil
			},
			searchNotesFunc: func(keyword string) ([]*model.Note, error) {
				return nil, errors.New("Unexpected Error")
			},
			expectedHTTPCode: http.StatusInternalServerError,
		}, {
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return nil
			},
			searchNotesFunc: func(keyword string) ([]*model.Note, error) {
				return nil, nil
			},
			notebookDB: mockNotebookDBAPI{
				notebookTitles: map[int64]string{1: "testTitle", 2: "testTitle2"},
				err:            nil,
			},
			expectedHTTPCode: http.StatusOK,
		}, {
			checkTokenFunc: func(r *http.Request, signingKey []byte) error {
				return nil
			},
			searchNotesFunc: func(keyword string) ([]*model.Note, error) {
				note1 := model.NewNote("testTitle", "testMemo", 1, []string{})
				return []*model.Note{note1}, nil
			},
			notebookDB: mockNotebookDBAPI{
				notebookTitles: map[int64]string{1: "testTitle", 2: "testTitle2"},
				err:            errors.New("Unexpected Error"),
			},
			expectedHTTPCode: http.StatusInternalServerError,
		},
	}

	for _, c := range cases {
		originalSearchNotes := searchNotesFunc
		oldNotebookDB := NotebookDB
		originalCheckToken := checkTokenFunc
		NotebookDB = c.notebookDB
		searchNotesFunc = c.searchNotesFunc
		checkTokenFunc = c.checkTokenFunc
		defer func() {
			searchNotesFunc = originalSearchNotes
			NotebookDB = oldNotebookDB
			checkTokenFunc = originalCheckToken
		}()

		req, _ := http.NewRequest("GET", "/searchBy/Title", nil)
		response := executeRequest(req)
		checkResponseCode(t, c.expectedHTTPCode, response.Code)
	}
}

func TestLoginAPI(t *testing.T) {
	cases := []struct {
		payload          []byte
		accountDB        mockAccountDBAPI
		expectedHTTPCode int
	}{
		{
			payload:          []byte("Unmarshable input"),
			expectedHTTPCode: http.StatusBadRequest,
		}, {
			payload: []byte(`{"username":"mockedUser", "password": "mockedPassword"}`),
			accountDB: mockAccountDBAPI{
				err: errors.New("Unexpected Error"),
			},
			expectedHTTPCode: http.StatusInternalServerError,
		}, {
			payload: []byte(`{"username":"mockedUser2", "password": "mockedPassword2"}`),
			accountDB: mockAccountDBAPI{
				username: "mockedUser",
				password: "mockedPassword",
				err:      nil,
			},
			expectedHTTPCode: http.StatusUnauthorized,
		}, {
			payload: []byte(`{"username":"mockedUser", "password": "mockedPassword"}`),
			accountDB: mockAccountDBAPI{
				username: "mockedUser",
				password: "mockedPassword",
				err:      nil,
			},
			expectedHTTPCode: http.StatusOK,
		},
	}

	for _, c := range cases {
		oldAccountDB := AccountDB
		AccountDB = c.accountDB
		defer func() {
			AccountDB = oldAccountDB
		}()

		req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(c.payload))
		response := executeRequest(req)
		checkResponseCode(t, c.expectedHTTPCode, response.Code)
	}
}

func TestParseInts(t *testing.T) {
	tests := []struct {
		input    string
		expected []int
	}{
		{
			input:    "1,2,23,123",
			expected: []int{1, 2, 23, 123},
		},
		{
			input:    "",
			expected: []int{},
		},
		{
			input:    "1,2,3,",
			expected: []int{1, 2, 3},
		},
	}

	for _, test := range tests {
		output, err := parseInts(test.input)
		if err != nil {
			t.Errorf("Could not parse %v, error msg: %v", test.input, err)
		}

		if !reflect.DeepEqual(test.expected, output) {
			t.Errorf("Could not parse integer expected %v, got %v", test.expected, output)
		}
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	mockServer := NewServer()
	mockServer.Initialize()
	mockServer.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

type mockAccountDBAPI struct {
	repository.AccountRepository
	username string
	password string
	err      error
}

func (mDB mockAccountDBAPI) GetAccount(string) (*model.Account, error) {
	hashedPass, _ := bcrypt.GenerateFromPassword([]byte(mDB.password), 10)
	return &model.Account{
		Username: mDB.username,
		Password: string(hashedPass)}, mDB.err
}

type mockNotebookDBAPI struct {
	repository.NotebookRepository
	notebookTitles map[int64]string
	err            error
}

func (mDB mockNotebookDBAPI) GetAllNotebooksTitle() (map[int64]string, error) {
	return mDB.notebookTitles, mDB.err
}

func mockJSONNotes() []*jsonNote {
	jNote1 := &jsonNote{
		ID:            1,
		Title:         "My Title",
		Memo:          "My Memo",
		Created:       time.Now(),
		LastUpdated:   time.Now(),
		Tags:          []string{"tag1"},
		NotebookTitle: "Notebook",
	}

	jNote2 := &jsonNote{
		ID:            2,
		Title:         "My Title 2",
		Memo:          "My Memo 2",
		Created:       time.Now(),
		LastUpdated:   time.Now(),
		Tags:          []string{"tag2", "tag2_1"},
		NotebookTitle: "Notebook 2",
	}

	return []*jsonNote{jNote1, jNote2}
}
