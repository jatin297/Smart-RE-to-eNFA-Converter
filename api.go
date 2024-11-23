package main

import (
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"github.com/gorilla/mux"
	"github.com/jatin297/retoenfa/dto"
	. "github.com/jatin297/retoenfa/metrics"
	"github.com/jatin297/retoenfa/retoenfa"
	user2 "github.com/jatin297/retoenfa/user"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

type APIService struct {
	listenAddr string
	store      Storage
}

func NewAPIService(listenAddr string, store Storage) *APIService {
	return &APIService{
		listenAddr: listenAddr,
		store:      store,
	}
}

type funcAPI func(w http.ResponseWriter, r *http.Request) error

type errorAPI struct {
	Error string `json:"error"`
}

func writeJSON(w http.ResponseWriter, r *http.Request, status int, v any, start time.Time) error {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	RecordMetricForHttp(r.Method, r.RequestURI, status, start)
	return json.NewEncoder(w).Encode(v)
}

func withJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("Authorization")
		_, err := validateJWT(tokenString)
		if err != nil {
			writeJSON(w, r, http.StatusBadRequest, errorAPI{
				"Invalid token",
			}, time.Now())
			return
		}
		handlerFunc(w, r)
	}
}

func makeHTTPHandleFunc(f funcAPI) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			// handle error
			writeJSON(w, r, http.StatusInternalServerError, errorAPI{
				Error: err.Error(),
			}, time.Now())
		}
	}
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})

}

func createJWT(user *user2.User) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		"email":     user.Email,
		"password":  user.EncryptedPassword,
	}

	secret := os.Getenv("JWT_SECRET")
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(secret))
}

func (s *APIService) handleUserAPI(w http.ResponseWriter, r *http.Request) error {

	switch r.Method {
	case "GET":
		return s.handleGetALLUsers(w, r)
	case "POST":
		return s.handleCreateUser(w, r)
	}

	return fmt.Errorf("method not supported: %s", r.Method)
}

func (s *APIService) handleLogin(w http.ResponseWriter, r *http.Request) (err error) {

	start := time.Now()
	var login dto.Login
	if err := json.NewDecoder(r.Body).Decode(&login); err != nil {
		return err
	}

	user, err := s.store.GetUserByID(login.ID)
	if err != nil {
		return err
	}

	ok := user.ValidPassword(login.Password)
	fmt.Println(ok)
	if !ok {
		return fmt.Errorf("invalid password")
	}

	token, err := createJWT(user)
	if err != nil {
		return err
	}

	var response dto.ResponseFormat
	response.Message = fmt.Sprintf("id: %d", login.ID)
	response.Token = token

	r.RequestURI = "login"
	return writeJSON(w, r, http.StatusOK, response, start)
}

func (s *APIService) handleGetUserByID(w http.ResponseWriter, r *http.Request) error {
	start := time.Now()
	idStr := mux.Vars(r)["id"]
	id, err := strconv.Atoi(idStr)
	if err != nil {
		// If the conversion fails, return a 400 Bad Request error
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return nil
	}
	user, err := s.store.GetUserByID(id)
	if err != nil {
		return err
	}

	// if incoming api is to delete
	r.RequestURI = "/user/{id}"
	if r.Method == "DELETE" {
		err = s.handleDeleteUser(id)
		if err != nil {
			return err
		}
		message := fmt.Errorf("successfully deleted user with id %d", id)
		return writeJSON(w, r, http.StatusOK, message, start)
	}
	return writeJSON(w, r, http.StatusOK, user, start)

}

func (s *APIService) handleGetALLUsers(w http.ResponseWriter, r *http.Request) error {
	start := time.Now()
	users, err := s.store.GetAllUsers()
	if err != nil {
		return err
	}
	r.RequestURI = "/users/get_all_users"
	return writeJSON(w, r, http.StatusOK, users, start)
}

func (s *APIService) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	start := time.Now()
	var userRequest user2.User
	if err := json.NewDecoder(r.Body).Decode(&userRequest); err != nil {
		return err
	}

	user := user2.NewUser(userRequest.Name, userRequest.Email, userRequest.EncryptedPassword)
	if err := s.store.CreateUser(user); err != nil {
		return err
	}

	tokenString, err := createJWT(user)
	if err != nil {
		return err
	}

	response := dto.ResponseFormat{
		Message: "Successfully created user",
		Token:   tokenString,
	}
	r.RequestURI = "/create_user"
	return writeJSON(w, r, http.StatusOK, response, start)
}

func (s *APIService) handleDeleteUser(id int) error {
	err := s.store.DeleteUser(id)
	if err != nil {
		return err
	}
	return nil
}

func (s *APIService) convertToENFA(w http.ResponseWriter, r *http.Request) error {
	if r.Method != "GET" {
		return fmt.Errorf("invalid api method")
	}
	var re dto.RegularExpression
	var eNFA dto.ENFAResponse
	var TransitionTable dto.TransitionTable

	start := time.Now()

	if err := json.NewDecoder(r.Body).Decode(&re); err != nil {
		RecordMetrics(r.Method, r.RequestURI, http.StatusBadRequest, start, re, eNFA)
		return writeJSON(w, r, http.StatusBadRequest, fmt.Errorf("invalid request body, err: %s", err.Error), start)
	}

	trans := retoenfa.NewReToeNFA(re.RE)
	trans.StartParse()
	enfa := trans.GetEpsNFA()
	transitionTable := enfa.GenerateFormattedTransitionTable()
	eNFA.TransitionTableSize = len(transitionTable)

	RecordMetrics(r.Method, r.RequestURI, http.StatusOK, start, re, eNFA)

	TransitionTable.TransitionTable = transitionTable

	return writeJSON(w, r, http.StatusOK, TransitionTable, start)
}

func (s *APIService) Run() {
	router := mux.NewRouter()

	router.HandleFunc("/user/login", makeHTTPHandleFunc(s.handleLogin))
	router.HandleFunc("/user/{id}", withJWTAuth(makeHTTPHandleFunc(s.handleGetUserByID)))
	router.HandleFunc("/user", makeHTTPHandleFunc(s.handleUserAPI))
	router.HandleFunc("/convert", withJWTAuth(makeHTTPHandleFunc(s.convertToENFA)))
	router.Handle("/metrics", promhttp.Handler())

	log.Println("api server running on port: ", s.listenAddr)

	http.ListenAndServe(s.listenAddr, router)
}
