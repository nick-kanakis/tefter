package cmd

import (
	"crypto/rand"
	"encoding/json"
	"errors"
	jwt "github.com/dgrijalva/jwt-go"
	"github.com/dgrijalva/jwt-go/request"
	"github.com/gorilla/mux"
	"github.com/nicolasmanic/tefter/model"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

//Server add a REST API layer for manipulating notes/notebooks
type Server struct {
	signingKey []byte
	Router     *mux.Router
}

//NewServer returns an instance of a Server struct
func NewServer() *Server {
	return &Server{
		Router: mux.NewRouter(),
	}
}

//Initialize initialize signing key & sets handlers to different endpoints
func (s *Server) Initialize() {
	//generate random signing key
	s.signingKey = make([]byte, 32)
	_, err := rand.Read(s.signingKey)
	if err != nil {
		log.Fatalln("Failed to generate signing key")
	}

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
	s.Router.HandleFunc("/login", s.login).Methods("POST")
}

//Run starts the server
func (s *Server) Run(port string) {
	log.Println("Server starting at port: " + port)
	log.Fatal(http.ListenAndServe(":"+port, s.Router))
}

var saveNoteFunc = addJSONNote
var checkTokenFunc = checkToken

func (s *Server) addNote(w http.ResponseWriter, r *http.Request) {
	if err := checkTokenFunc(r, s.signingKey); err != nil {
		log.Printf("Invalid token, failed with message: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Authorization failed")
		return
	}

	var jNote *jsonNote
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&jNote); err != nil {
		log.Printf("Error while decoding jsonNote, error msg: %v", err)
		respondWithError(w, http.StatusBadRequest, "Failed decoding note")
		return
	}
	defer r.Body.Close()

	if err := saveNoteFunc(jNote); err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, map[string]string{"result": "success"})
}

var updateNoteFunc = updateJSONNote

func (s *Server) updateNote(w http.ResponseWriter, r *http.Request) {
	if err := checkTokenFunc(r, s.signingKey); err != nil {
		log.Printf("Invalid token, failed with message: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Authorization failed")
		return
	}

	var jNote *jsonNote
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&jNote); err != nil {
		log.Printf("Error while decoding jsonNote, error msg: %v", err)
		respondWithError(w, http.StatusBadRequest, "Failed decoding note")
		return
	}
	defer r.Body.Close()

	if err := updateNoteFunc(jNote); err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusCreated, map[string]string{"result": "success"})
}

var retrieveNotesFunc = retrieveJSONNotes

func (s *Server) getNotes(w http.ResponseWriter, r *http.Request) {
	if err := checkTokenFunc(r, s.signingKey); err != nil {
		log.Printf("Invalid token, failed with message: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Authorization failed")
		return
	}

	vars := mux.Vars(r)
	var jsonNotes []*jsonNote
	var err error

	//Comma separated list of ids, tags, notebookTitles
	strIDs := vars["ids"]
	strNotebookTitles := vars["notebookTitles"]
	strTags := vars["tags"]

	ids, err := parseInts(strIDs)
	if err != nil {
		log.Printf("Error while parsing ids, error msg: %v", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	tags := parseStrings(strTags)
	notebookTitles := parseStrings(strNotebookTitles)

	jsonNotes, err = retrieveNotesFunc(ids, notebookTitles, tags, false)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, jsonNotes)
}

var deleteNotesFunc = delete

func (s *Server) deleteNotes(w http.ResponseWriter, r *http.Request) {
	if err := checkTokenFunc(r, s.signingKey); err != nil {
		log.Printf("Invalid token, failed with message: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Authorization failed")
		return
	}

	vars := mux.Vars(r)
	//Comma separated list of ids
	strIDs := vars["ids"]
	ids, err := parseInts(strIDs)
	if err != nil {
		log.Printf("Error while parsing ids, error msg: %v", err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	err = deleteNotesFunc(int64Slice(ids))
	if err != nil {
		log.Print(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

var deleteNotebooksFunc = deleteNotebooks

func (s *Server) deleteNotebooks(w http.ResponseWriter, r *http.Request) {
	if err := checkTokenFunc(r, s.signingKey); err != nil {
		log.Printf("Invalid token, failed with message: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Authorization failed")
		return
	}

	vars := mux.Vars(r)
	//Comma separated  notebookTitles
	strNotebookTitles := vars["notebookTitles"]
	notebookTitles := parseStrings(strNotebookTitles)

	err := deleteNotebooksFunc(notebookTitles)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

var updateNotebookFunc = updateNotebook

func (s *Server) updateNotebook(w http.ResponseWriter, r *http.Request) {
	if err := checkTokenFunc(r, s.signingKey); err != nil {
		log.Printf("Invalid token, failed with message: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Authorization failed")
		return
	}

	vars := mux.Vars(r)
	oldTitle := vars["oldTitle"]
	newTitle := vars["newTitle"]
	err := updateNotebookFunc(oldTitle, newTitle)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

var searchNotesFunc = search

func (s *Server) searchKeyword(w http.ResponseWriter, r *http.Request) {
	if err := checkTokenFunc(r, s.signingKey); err != nil {
		log.Printf("Invalid token, failed with message: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Authorization failed")
		return
	}

	vars := mux.Vars(r)
	keyword := vars["keyword"]
	notes, err := searchNotesFunc(keyword)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	jNotes, err := transformNotes2JSONNotes(notes)
	if err != nil {
		log.Println(err)
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	respondWithJSON(w, http.StatusOK, jNotes)
}

func (s *Server) login(w http.ResponseWriter, r *http.Request) {
	var accountRequest *model.Account
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&accountRequest); err != nil {
		log.Printf("Error while decoding account, error msg: %v", err)
		respondWithError(w, http.StatusBadRequest, "Failed decoding account")
		return
	}
	account, err := AccountDB.GetAccount(accountRequest.Username)
	if err != nil {
		log.Printf("Error retrieving username, error msg: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Error retrieving username")
		return
	}
	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(accountRequest.Password)); err != nil {
		log.Printf("Username and password don't match, error msg: %v", err)
		respondWithError(w, http.StatusUnauthorized, "Username and password don't match")
		return
	}
	//token will be valid for 24 hours
	exp := time.Now().Add(24 * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"exp": exp.Unix(),
		"sub": account.Username,
	})
	signedToken, err := token.SignedString(s.signingKey)
	if err != nil {
		log.Printf("Could note sign token, error msg: %v", err)
		respondWithError(w, http.StatusInternalServerError, "Could note sign token")
		return
	}
	respondWithJSON(w, http.StatusOK, map[string]string{"token": signedToken})
}

func checkToken(r *http.Request, signingKey []byte) error {
	token, err := request.ParseFromRequest(r, request.AuthorizationHeaderExtractor, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("Unexpected signing method")
		}
		return signingKey, nil
	})
	if err != nil {
		return err
	}
	claims := token.Claims.(jwt.MapClaims)
	return claims.Valid()
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
