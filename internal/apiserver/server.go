package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"firstRestAPI/internal/model"
	"firstRestAPI/internal/store"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

const (
	sessionKeyName        = "session"
	ctxKeyUser     ctxKey = iota
	ctxKeyRequestID
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type ctxKey int8

type server struct {
	router       *mux.Router
	logger       *slog.Logger
	store        store.Store
	sessionStore sessions.Store
}

func newServer(store store.Store, config *Config, sessionStore sessions.Store) *server {
	s := &server{
		router: mux.NewRouter(),
		logger: func() *slog.Logger {
			var log *slog.Logger
			switch config.LogLevel {
			case "error":
				log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
			case "info":
				log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
			case "warn":
				log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelWarn}))
			case "debug":
				log = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
			}
			return log
		}(),
		store:        store,
		sessionStore: sessionStore,
	}
	s.configureRouter()
	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logReqeust)
	// s.router.Use(cors.New(
	// 	cors.Options{
	// 		AllowedOrigins: []string{"*"},
	// 		AllowedHeaders: []string{"*"},
	// 		Debug:          false,
	// 		AllowedMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodOptions, http.MethodDelete}},
	// ).Handler)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE"}),
		handlers.AllowedHeaders([]string{"*"}),
		handlers.AllowCredentials()))
	s.router.HandleFunc("/users", s.HandleUsersCreate()).Methods("POST")
	s.router.HandleFunc("/sessions", s.HandleSessionsCreate()).Methods("POST")
	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/show", s.HandleShow()).Methods(http.MethodGet)
	private.HandleFunc("/createcard", s.HandleCardCreate()).Methods("POST")
	private.HandleFunc("/deletecard", s.HandleDeleteCard()).Methods(http.MethodPost)
	private.HandleFunc("/editcard", s.HandleCardEdit()).Methods(http.MethodPost)
	private.HandleFunc("/whoami", s.handleWhoami()).Methods("GET")
	private.HandleFunc("/showusingtime", s.HandleCardsShowUsingTime()).Methods(http.MethodGet)
	private.HandleFunc("/updatecardflag", s.HandleCardFlagUp()).Methods(http.MethodPost)
}

func (s *server) setRequestID(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestID, id)))
	})
}

func (s *server) logReqeust(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.With(
			slog.String("remote_addr", r.RemoteAddr),
			slog.String("request_id", r.Context().Value(ctxKeyRequestID).(string)),
		)
		logger.Info(fmt.Sprintf("started %s %s", r.Method, r.RequestURI))
		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		logger.Info(
			fmt.Sprintf("completed with %d %s in %v",
				rw.code,
				http.StatusText(rw.code),
				time.Since(start).String(),
			),
		)
	})
}

func (s *server) HandleCardFlagUp() http.HandlerFunc {
	type request struct {
		ID int64 `json:"card_id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		session, err := s.sessionStore.Get(r, sessionKeyName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}
		card := &model.Card{
			UserID: id.(int64),
			ID:     req.ID,
		}
		if err := s.store.Card().CardFlagUp(card); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (s *server) HandleCardsShowUsingTime() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionKeyName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}
		card := &model.Card{
			UserID: id.(int64),
		}
		cards, err := s.store.Card().ShowUsingTime(card)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(cards); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (s *server) HandleCardEdit() http.HandlerFunc {
	type request struct {
		CardID    int64  `json:"card_id"`
		FrontSide string `json:"front_side"`
		BackSide  string `json:"back_side"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		session, err := s.sessionStore.Get(r, sessionKeyName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}
		card := &model.Card{
			UserID:    id.(int64),
			ID:        req.CardID,
			FrontSide: req.FrontSide,
			BackSide:  req.BackSide,
		}
		if err := s.store.Card().Edit(card); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
	}
}

func (s *server) HandleSessionsCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil || !u.ComparePassword(req.Password) {
			s.error(w, r, http.StatusUnauthorized, errIncorrectEmailOrPassword)
			return
		}
		//sesKey := "someSession"
		// sessionName должен будет генерироваться через uuid и сохраняться
		session, err := s.sessionStore.Get(r, sessionKeyName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		session.Values["user_id"] = u.ID
		if err := s.sessionStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)

	}
}

func (s *server) HandleCardCreate() http.HandlerFunc {
	type request struct {
		FrontSide string `json:"front_side"`
		BackSide  string `json:"back_side"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		session, err := s.sessionStore.Get(r, sessionKeyName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}

		card := &model.Card{
			UserID:    id.(int64),
			FrontSide: req.FrontSide,
			BackSide:  req.BackSide,
		}
		if err := s.store.Card().Create(card); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) HandleDeleteCard() http.HandlerFunc {
	type request struct {
		ID int64 `json:"card_id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		session, err := s.sessionStore.Get(r, sessionKeyName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, err)
			return
		}
		card := &model.Card{
			ID:     req.ID,
			UserID: id.(int64),
		}
		if err := s.store.Card().Delete(card); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) HandleShow() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionKeyName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		card := &model.Card{
			UserID: id.(int64),
		}
		cards, err := s.store.Card().Show(card)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(cards); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (s *server) handleWhoami() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.respond(w, r, http.StatusOK, r.Context().Value(ctxKeyUser).(*model.User))
	}
}

func (s *server) authenticateUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionKeyName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		id, ok := session.Values["user_id"]
		if !ok {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		u, err := s.store.User().Find(id.(int64))
		if err != nil {
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyUser, u)))
	})
}

func (s *server) HandleUsersCreate() http.HandlerFunc {
	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Nickname string `json:"nickname"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		if err := s.store.UserLK().FindByNickname(req.Nickname); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		u := &model.User{
			Email:    req.Email,
			Password: req.Password,
		}
		ulk := &model.UserLK{
			Nickname: req.Nickname,
		}
		if err := ulk.Validate(); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		ulk.UserID = u.ID
		if err := s.store.UserLK().Create(ulk); err != nil {
			if err = s.store.User().Delete(u.ID); err != nil {
				s.error(w, r, http.StatusBadRequest, err)
				return
			}
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		u.Sanitize()
		s.respond(w, r, http.StatusCreated, u, ulk)
		//s.respond(w, r, http.StatusCreated, u)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}
func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data ...interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
