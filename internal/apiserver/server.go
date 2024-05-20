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
	slog.Info("s.router.ServeHTTP(w, r) прошел")
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"https://localhost:10443"}), // тут воровская звезда
		handlers.AllowedMethods([]string{"GET", "POST", "OPTIONS", "HEAD", "PUT"}),
		handlers.AllowedHeaders([]string{"Content-Type", "Accept", "Origin", "X-Request-ID", "Allow", "Set-Cookie", "Cookie"}),
		handlers.AllowCredentials(),
	))
	s.router.StrictSlash(true)
	s.router.Use(s.setRequestID)
	s.router.Use(s.logReqeust)
	s.router.HandleFunc("/users", s.HandleUsersCreate()).Methods(http.MethodPost, http.MethodOptions, http.MethodHead, http.MethodGet)
	s.router.HandleFunc("/sessions", s.HandleSessionsCreate()).Methods(http.MethodPost, http.MethodOptions, http.MethodHead, http.MethodGet)
	private := s.router.PathPrefix("/private").Subrouter()
	private.Use(s.authenticateUser)
	private.HandleFunc("/show", s.HandleShow()).Methods(http.MethodGet, http.MethodOptions)
	private.HandleFunc("/createcard", s.HandleCardCreate()).Methods(http.MethodPost, http.MethodOptions)
	private.HandleFunc("/deletecard", s.HandleDeleteCard()).Methods(http.MethodPost, http.MethodOptions)
	private.HandleFunc("/editcard", s.HandleCardEdit()).Methods(http.MethodPost, http.MethodOptions)
	private.HandleFunc("/whoami", s.handleWhoami()).Methods(http.MethodGet, http.MethodOptions)
	private.HandleFunc("/showusingtime", s.HandleCardsShowUsingTime()).Methods(http.MethodPost, http.MethodOptions)
	private.HandleFunc("/updatecardflag", s.HandleCardFlagUp()).Methods(http.MethodPost, http.MethodOptions)
	private.HandleFunc("/sessionquit", s.SessionsQuit()).Methods(http.MethodGet, http.MethodOptions)
	private.HandleFunc("/groupcreate", s.HandleGroupCreate()).Methods(http.MethodPost, http.MethodOptions)
	private.HandleFunc("/groupedit", s.HandleGroupEdit()).Methods(http.MethodPost, http.MethodOptions)
	private.HandleFunc("/groupdelete", s.HandleGroupDelete()).Methods(http.MethodPost, http.MethodOptions)
	private.HandleFunc("/groupshow", s.HandleGroupShow()).Methods(http.MethodPost, http.MethodOptions)
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
	type request struct {
		GroupID int64 `json:"group_id"`
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
			UserID:  id.(int64),
			GroupID: req.GroupID,
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
func (s *server) SessionsQuit() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, err := s.sessionStore.Get(r, sessionKeyName)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		session.Options.MaxAge = 0
		session.Save(r, w)
		w.WriteHeader(200)
	}
}
func (s *server) HandleGroupCreate() http.HandlerFunc {
	type request struct {
		GroupName string `json:"group_name"`
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
		group := &model.Group{
			UserID:    id.(int64),
			GroupName: req.GroupName,
		}
		if err := s.store.Group().Create(group); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		s.respond(w, r, 200, group.GroupID)
	}
}
func (s *server) HandleGroupDelete() http.HandlerFunc {
	type request struct {
		GroupID int64 `json:"group_id"`
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
		group := &model.Group{
			UserID:  id.(int64),
			GroupID: req.GroupID,
		}
		if err := s.store.Group().Delete(group); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
	}
}
func (s *server) HandleGroupEdit() http.HandlerFunc {
	type request struct {
		GroupID   int64  `json:"group_id"`
		GroupName string `json:"group_name"`
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
		group := &model.Group{
			UserID:    id.(int64),
			GroupID:   req.GroupID,
			GroupName: req.GroupName,
		}
		if err := s.store.Group().Edit(group); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
	}
}
func (s *server) HandleGroupShow() http.HandlerFunc {
	type request struct {
		GroupID int64 `json:"group_id"`
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
			s.error(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}
		group := &model.Group{
			UserID:  id.(int64),
			GroupID: req.GroupID,
		}
		cards, err := s.store.Group().Show(group)
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
		session.Options.Path = "/private"
		session.Options.SameSite = http.SameSiteStrictMode
		session.Options.Secure = true
		session.Options.HttpOnly = true
		session.Values["user_id"] = u.ID
		if err := s.sessionStore.Save(r, w, session); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		//s.respond(w, r, http.StatusOK, nil)

	}
}

func (s *server) HandleCardCreate() http.HandlerFunc {
	type request struct {
		FrontSide string `json:"front_side"`
		BackSide  string `json:"back_side"`
		GroupID   int64  `json:"group_id"`
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
			GroupID:   req.GroupID,
		}
		if err := s.store.Card().Create(card); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
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
		slog.Info(fmt.Sprintf("Request in HUC: %s, %s, %s", req.Nickname, req.Email, req.Password))
		slog.Info(fmt.Sprintf("*http.Request in HUC: %s, %s, %s", r.Host, r.Method, r.RequestURI))
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
		slog.Info(fmt.Sprintf("model.User: %s, %s, %d", u.Email, u.Password, u.ID))
		if err := s.store.User().Create(u); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}
		ulk.UserID = u.ID
		if err := s.store.UserLK().Create(ulk, u); err != nil {
			if err = s.store.User().Delete(u.ID); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
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
	slog.Error(fmt.Sprintf("Error: %s", err.Error()))
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}
func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data ...interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
