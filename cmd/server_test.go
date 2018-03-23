package cmd

import (
	"bytes"
	"encoding/json"
	"github.com/nicolasmanic/tefter/model"
	"github.com/nicolasmanic/tefter/repository"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"
)

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

func TestAddNoteAPI(t *testing.T) {
	originalSaveNote := saveNoteFunc
	defer func() {
		saveNoteFunc = originalSaveNote
	}()
	saveNoteFunc = func(*jsonNote) error {
		return nil
	}

	payload := []byte(`{"title":"Shopping for weekend","memo":" Things for weekend:\n \u003e Milk\n \u003e Eggs\n \u003e Chicken breast\n","created":"2018-03-20T18:53:35.4123749+02:00","updated":"2018-03-20T18:53:35.4193801+02:00","tags":["weekend","list"],"notebook_title":"Shopping"}`)

	req, _ := http.NewRequest("POST", "/addNote", bytes.NewBuffer(payload))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)
}

func TestUpdateNoteAPI(t *testing.T) {
	originalUpdateNote := updateNoteFunc
	defer func() {
		updateNoteFunc = originalUpdateNote
	}()
	updateNoteFunc = func(*jsonNote) error {
		return nil
	}

	payload := []byte(`{"id":1, "title":"Shopping for weekend","memo":" Things for weekend:\n \u003e Milk\n \u003e Eggs\n \u003e Chicken breast\n","created":"2018-03-20T18:53:35.4123749+02:00","updated":"2018-03-20T18:53:35.4193801+02:00","tags":["weekend","list"],"notebook_title":"Shopping"}`)

	req, _ := http.NewRequest("PUT", "/updateNote", bytes.NewBuffer(payload))
	response := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, response.Code)
}

func TestGetNotesByIDAPI(t *testing.T) {
	originalRetrieveNotes := retrieveNotesFunc
	defer func() {
		retrieveNotesFunc = originalRetrieveNotes
	}()

	retrieveNotesFunc = func(ids []int, notebookTitles, tags []string, getAll bool) ([]*jsonNote, error) {
		if len(ids) != 2 || len(notebookTitles) != 0 || len(tags) != 0 && getAll {
			t.Error("Wrong arguments at retrieveNotesFunc for searching by id")
		}
		return mockJSONNotes(), nil
	}

	req, _ := http.NewRequest("GET", "/getNotesByID/1,2", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var jNotes []*jsonNote
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&jNotes); err != nil {
		t.Errorf("Failed decoding returned json notes (requested by id), error msg:%v", err)
	}
	if len(jNotes) != 2 {
		t.Error("Failed decoding returned json notes (requested by id)")
	}
	for _, jNote := range jNotes {
		if len(jNote.Tags) != 1 && len(jNote.Tags) != 2 {
			t.Error("Failed decoding returned json notes (requested by id)")
		}
		if jNote.ID != 1 && jNote.ID != 2 {
			t.Error("Failed decoding returned json notes (requested by id)")
		}
	}
}

func TestGetNotesByNotebookTitleAPI(t *testing.T) {
	originalRetrieveNotes := retrieveNotesFunc
	defer func() {
		retrieveNotesFunc = originalRetrieveNotes
	}()

	retrieveNotesFunc = func(ids []int, notebookTitles, tags []string, getAll bool) ([]*jsonNote, error) {
		if len(ids) != 0 || len(notebookTitles) != 3 || len(tags) != 0 && getAll {
			t.Error("Wrong arguments at retrieveNotesFunc for searching by notebookTitle")
		}
		return mockJSONNotes(), nil
	}

	req, _ := http.NewRequest("GET", "/getNotesByNotebookTitle/title1,title2,title3", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var jNotes []*jsonNote
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&jNotes); err != nil {
		t.Errorf("Failed decoding returned json notes (requested by notebookTitle), error msg:%v", err)
	}
	if len(jNotes) != 2 {
		t.Error("Failed decoding returned json notes (requested by notebookTitle)")
	}
	for _, jNote := range jNotes {
		if len(jNote.Tags) != 1 && len(jNote.Tags) != 2 {
			t.Error("Failed decoding returned json notes (requested by notebookTitle)")
		}
		if jNote.ID != 1 && jNote.ID != 2 {
			t.Error("Failed decoding returned json notes (requested by notebookTitle)")
		}
	}
}

func TestGetNotesByNotebookTagsAPI(t *testing.T) {
	originalRetrieveNotes := retrieveNotesFunc
	defer func() {
		retrieveNotesFunc = originalRetrieveNotes
	}()

	retrieveNotesFunc = func(ids []int, notebookTitles, tags []string, getAll bool) ([]*jsonNote, error) {
		if len(ids) != 0 || len(notebookTitles) != 0 || len(tags) != 1 && getAll {
			t.Error("Wrong arguments at retrieveNotesFunc for searching by tags")
		}
		return mockJSONNotes(), nil
	}

	req, _ := http.NewRequest("GET", "/getNotesByTags/tag1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var jNotes []*jsonNote
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&jNotes); err != nil {
		t.Errorf("Failed decoding returned json notes (requested by tag), error msg:%v", err)
	}
	if len(jNotes) != 2 {
		t.Error("Failed decoding returned json notes (requested by tag)")
	}
	for _, jNote := range jNotes {
		if len(jNote.Tags) != 1 && len(jNote.Tags) != 2 {
			t.Error("Failed decoding returned json notes (requested by tag)")
		}
		if jNote.ID != 1 && jNote.ID != 2 {
			t.Error("Failed decoding returned json notes (requested by tag)")
		}
	}
}

func TestGetAllNotesAPI(t *testing.T) {
	originalRetrieveNotes := retrieveNotesFunc
	defer func() {
		retrieveNotesFunc = originalRetrieveNotes
	}()

	retrieveNotesFunc = func(ids []int, notebookTitles, tags []string, getAll bool) ([]*jsonNote, error) {
		if len(ids) != 0 || len(notebookTitles) != 0 || len(tags) != 0 && !getAll {
			t.Error("Wrong arguments at retrieveNotesFunc for all notes")
		}
		return mockJSONNotes(), nil
	}

	req, _ := http.NewRequest("GET", "/getAllNotes", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var jNotes []*jsonNote
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&jNotes); err != nil {
		t.Errorf("Failed decoding returned json notes, error msg:%v", err)
	}
	if len(jNotes) != 2 {
		t.Error("Failed decoding returned json notes")
	}
	for _, jNote := range jNotes {
		if len(jNote.Tags) != 1 && len(jNote.Tags) != 2 {
			t.Error("Failed decoding returned json notes")
		}
		if jNote.ID != 1 && jNote.ID != 2 {
			t.Error("Failed decoding returned json notes")
		}
	}
}

func TestDeleteNotesAPI(t *testing.T) {
	originalDeleteNotes := deleteNotesFunc
	defer func() {
		deleteNotesFunc = originalDeleteNotes
	}()
	deleteNotesFunc = func(ids []int64) error {
		if len(ids) != 2 {
			t.Errorf("Incorrect number of ids passed for notes deletion, ids received: %v", len(ids))
		}
		return nil
	}

	req, _ := http.NewRequest("DELETE", "/deleteNotes/1,2", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestDeleteNotebooksAPI(t *testing.T) {
	originalDeleteNotebooks := deleteNotebooksFunc
	defer func() {
		deleteNotebooksFunc = originalDeleteNotebooks
	}()
	deleteNotebooksFunc = func(titles []string) error {
		if len(titles) != 1 {
			t.Errorf("Incorrect number of titles passed for notebooks deletion, titles received: %v", len(titles))
		}
		return nil
	}

	req, _ := http.NewRequest("DELETE", "/deleteNotebooks/title1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestUpdateNotebooksAPI(t *testing.T) {
	originalUpdateNotebook := updateNotebookFunc
	defer func() {
		updateNotebookFunc = originalUpdateNotebook
	}()
	updateNotebookFunc = func(oldTitle, newTitle string) error {
		if oldTitle != "oldTitle_1" || newTitle != "newTitle_1" {
			t.Errorf("Missmatch titles expected newTitle: newTitle_1, oldTitle: oldTitle_1, got %v, %v", newTitle, oldTitle)
		}
		return nil
	}

	req, _ := http.NewRequest("PUT", "/updateNotebook/oldTitle_1/newTitle_1", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)
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

func TestSearchNotesAPI(t *testing.T) {
	originalSearchNotes := searchNotesFunc
	oldNotebookDB := NotebookDB
	NotebookDB = mockNotebookDBAPI{}
	defer func() {
		searchNotesFunc = originalSearchNotes
		NotebookDB = oldNotebookDB
	}()

	searchNotesFunc = func(keyword string) ([]*model.Note, error) {
		note1 := model.NewNote("testTitle", "testMemo", 1, []string{})
		note1.ID = 1
		note2 := model.NewNote("testTitle2", "testMemo2", 2, []string{})
		note2.ID = 2
		return []*model.Note{note1, note2}, nil
	}

	req, _ := http.NewRequest("GET", "/searchBy/Title", nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	var jNotes []*jsonNote
	decoder := json.NewDecoder(response.Body)
	if err := decoder.Decode(&jNotes); err != nil {
		t.Errorf("Failed decoding returned json notes, error msg:%v", err)
	}
	if len(jNotes) != 2 {
		t.Error("Failed decoding returned json notes")
	}
}

type mockNotebookDBAPI struct {
	repository.NotebookRepository
}

func (mDB mockNotebookDBAPI) GetAllNotebooksTitle() (map[int64]string, error) {
	return map[int64]string{1: "testTitle", 2: "testTitle2"}, nil
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
