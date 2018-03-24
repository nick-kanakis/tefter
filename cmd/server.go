package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"strconv"
	"strings"
)

//Server add a REST API layer for manipulating notes/notebooks
type Server struct {
	Router *mux.Router
}

//NewServer returns an instance of a Server struct
func NewServer() *Server {
	return &Server{
		Router: mux.NewRouter(),
	}
}

//Initialize sets handlers to different endpoints
func (s *Server) Initialize() {
	s.Router.HandleFunc("/addNote", s.addNote).Methods("POST")
	s.Router.HandleFunc("/updateNote", s.updateNote).Methods("PUT")
	s.Router.HandleFunc("/getNotesByID/{ids}", s.getNotes).Methods("GET")
	s.Router.HandleFunc("/getNotesByNotebookTitle/{notebookTitles}", s.getNotes).Methods("GET")
	s.Router.HandleFunc("/getNotesByTags/{Tags}", s.getNotes).Methods("GET")
	s.Router.HandleFunc("/getAllNotes", s.getNotes).Methods("GET")
	s.Router.HandleFunc("/deleteNotes/{ids}", s.deleteNotes).Methods("DELETE")
	s.Router.HandleFunc("/searchBy/{keyword}", s.searchKeyword).Methods("GET")
	s.Router.HandleFunc("/updateNotebook/{oldTitle}/{newTitle}", s.updateNotebook).Methods("PUT")
	s.Router.HandleFunc("/deleteNotebooks/{notebookTitles}", s.deleteNotebooks).Methods("DELETE")
}

//RUN starts the server
func (s *Server) Run(port string) {
	fmt.Println("Server starting at port :" + port)
	log.Fatal(http.ListenAndServe(":"+port, s.Router))
}

var saveNoteFunc = addJSONNote

func (s *Server) addNote(w http.ResponseWriter, r *http.Request) {
	var jNote *jsonNote
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&jNote); err != nil {
		fmt.Printf("Error while decoding jsonNote, error msg: %v", err)
		respondWithError(w, http.StatusBadRequest, "Failed decoding note")
		return
	}
	defer r.Body.Close()

	if err := saveNoteFunc(jNote); err != nil {
		fmt.Printf("Error while saving note, error msg: %v", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJSON(w, http.StatusCreated, map[string]string{"result": "success"})
}

var updateNoteFunc = updateJSONNote

func (s *Server) updateNote(w http.ResponseWriter, r *http.Request) {
	var jNote *jsonNote
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&jNote); err != nil {
		fmt.Printf("Error while decoding jsonNote, error msg: %v", err)
		respondWithError(w, http.StatusBadRequest, "Failed decoding note")
		return
	}
	defer r.Body.Close()

	if err := updateNoteFunc(jNote); err != nil {
		fmt.Printf("Error while updating note, error msg: %v", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJSON(w, http.StatusCreated, map[string]string{"result": "success"})
}

var retrieveNotesFunc = retrieveJSONNotes

func (s *Server) getNotes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var jsonNotes []*jsonNote
	var err error
	if len(vars) == 0 {
		jsonNotes, err = retrieveNotesFunc([]int{}, []string{}, []string{}, true)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
	} else {
		//Comma separated list of ids, tags, notebookTitles
		strIDs := vars["ids"]
		strNotebookTitles := vars["notebookTitles"]
		strTags := vars["tags"]

		ids, err := parseInts(strIDs)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
		tags := parseStrings(strTags)
		notebookTitles := parseStrings(strNotebookTitles)

		jsonNotes, err = retrieveNotesFunc(ids, notebookTitles, tags, false)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
	}
	respondWithJSON(w, http.StatusOK, jsonNotes)
}

var deleteNotesFunc = delete

func (s *Server) deleteNotes(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//Comma separated list of ids
	strIDs := vars["ids"]
	ids, err := parseInts(strIDs)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}

	err = deleteNotesFunc(int64Slice(ids))
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

var deleteNotebooksFunc = deleteNotebooks

func (s *Server) deleteNotebooks(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//Comma separated  notebookTitles
	strNotebookTitles := vars["notebookTitles"]
	notebookTitles := parseStrings(strNotebookTitles)

	err := deleteNotebooksFunc(notebookTitles)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

var updateNotebookFunc = updateNotebook

func (s *Server) updateNotebook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	oldTitle := vars["oldTitle"]
	newTitle := vars["newTitle"]
	err := updateNotebookFunc(oldTitle, newTitle)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

var searchNotesFunc = search

func (s *Server) searchKeyword(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	keyword := vars["keyword"]
	notes, err := searchNotesFunc(keyword)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	jNotes, err := transformNotes2JSONNotes(notes)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
	respondWithJSON(w, http.StatusOK, jNotes)
}

//Split comma separated integers
func parseInts(str string) ([]int, error) {
	if len(str) == 0 {
		return []int{}, nil
	}
	var integers = make([]int, 0, len(str))
	strings := parseStrings(str)

	for _, s := range strings {
		if s == "" {
			continue
		}
		i, err := strconv.Atoi(s)
		if err != nil {
			return nil, err
		}
		integers = append(integers, i)
	}
	return integers, nil
}

func parseStrings(str string) []string {
	if len(str) == 0 {
		return []string{}
	}
	return strings.Split(str, ",")
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
