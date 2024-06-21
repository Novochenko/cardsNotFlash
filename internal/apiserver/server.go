package apiserver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"firstRestAPI/internal/model"
	"firstRestAPI/internal/store"
	"fmt"
	"io"
	"log/slog"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/fatih/structs"
	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	// "github.com/mitchellh/mapstructure"
)

const (
	sessionKeyName        = "session"
	MaxImageSize          = 31457280 // 30mb in bytes
	MaxImagesCount        = 10
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
	private.HandleFunc("/showusingtime", s.HandleCardsShowUsingTime()).Methods(http.MethodGet, http.MethodOptions)
	private.HandleFunc("/updatecardflag", s.HandleCardFlagUp()).Methods(http.MethodPost, http.MethodOptions)
	private.HandleFunc("/sessionquit", s.SessionsQuit()).Methods(http.MethodGet, http.MethodOptions)
	private.HandleFunc("/groupcreate", s.HandleGroupCreate()).Methods(http.MethodPost, http.MethodOptions)
	private.HandleFunc("/groupedit", s.HandleGroupEdit()).Methods(http.MethodPost, http.MethodOptions)
	private.HandleFunc("/groupdelete", s.HandleGroupDelete()).Methods(http.MethodPost, http.MethodOptions)
	private.HandleFunc("/groupshow", s.HandleGroupShow()).Methods(http.MethodPost, http.MethodOptions)
	private.HandleFunc("/showallgroups", s.HandleShowAllGroups()).Methods(http.MethodGet, http.MethodOptions)
	private.HandleFunc("/showgroupusingtime", s.HandleGroupShowUsingTime()).Methods(http.MethodPost, http.MethodOptions)
	private.HandleFunc("/lkshow", s.HandleLKShow()).Methods(http.MethodGet, http.MethodOptions)
	private.HandleFunc("/lkdescriptionedit", s.HandleLKDescriptionEdit()).Methods(http.MethodPost, http.MethodOptions)
	private.HandleFunc("/pfpupload", s.HandlePFPUpload()).Methods(http.MethodPost, http.MethodOptions)
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
func (s *server) HandleLKDescriptionEdit() http.HandlerFunc {
	type request struct {
		UserDescription string `json:"user_description"`
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
		lk := &model.UserLK{
			UserID:          id.(int64),
			UserDescription: req.UserDescription,
		}
		if err = s.store.UserLK().LKDescriptionEdit(lk); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		s.respond(w, r, http.StatusOK, nil)
	}
}

//	func download(values map[string]interface{}) (b *bytes.Buffer, dataType string, err error) {
//		// Prepare a form that you will submit to that URL.
//		b = &bytes.Buffer{}
//		w := multipart.NewWriter(b)
//		defer w.Close()
//		imagesCount := 0
//		for key, r := range values {
//			if x, ok := r.(io.Closer); ok {
//				defer x.Close()
//			}
//			// Add an image file
//			if x, ok := r.(*os.File); ok {
//				if x == nil {
//					continue
//				}
//				xStat, err := x.Stat()
//				if err != nil {
//					return nil, "", err
//				}
//				if size := xStat.Size(); size > MaxImageSize {
//					return nil, "", store.ErrMaxSizeAttained
//				}
//				if imagesCount >= MaxImagesCount {
//					slog.Error("imagesCount > maxImagesCount")
//					continue
//				}
//				part, err := w.CreateFormFile(key, x.Name())
//				if err != nil {
//					return nil, "", err
//				}
//				if _, err := io.Copy(part, x); err != nil {
//					return nil, "", err
//				}
//				imagesCount++
//			}
//			if x, ok := r.(string); ok {
//				// Add other fields
//				if err = w.WriteField(key, x); err != nil {
//					return nil, "", err
//				}
//			}
//			if x, ok := r.(int64); ok {
//				// Add other fields
//				if err = w.WriteField(key, strconv.Itoa(int(x))); err != nil {
//					return nil, "", err
//				}
//			}
//			dataType = w.FormDataContentType()
//		}
//		// Don't forget to close the multipart writer.
//		// If you don't close it, your request will be missing the terminating boundary.
//		return
//	}

/*
Server -> Client send of multipart/form-data, dataType must be used first then sequentaly rw in a go func() and error from done channel handling
usage template:

w.Header().Set("Content-Type", dataType)
w.WriteHeader(http.StatusOK)

	go func() {
		_, err := io.Copy(w, pr)
		if err != nil {
			error handling
	}

}()

	if err = <-done; err != nil {
		error handling
	}

close(done)

	if err = closer(); err != nil {
		error handling
	}
*/
func download(values map[string]interface{}) (*io.PipeReader, func() error, string, chan error) {
	imagesCount := 0
	var dataType string
	done := make(chan error)
	pipeReader, pipeWriter := io.Pipe()
	w := multipart.NewWriter(pipeWriter)
	go func(pw *io.PipeWriter) {
		defer w.Close()
		for key, r := range values {
			slog.Info(key)
			switch value := r.(type) {
			case *os.File:
				if value == nil {
					done <- store.ErrNilFile
					return
				}
				xStat, err := value.Stat()
				if err != nil {
					done <- err
					return
				}
				if size := xStat.Size(); size > MaxImageSize {
					done <- store.ErrMaxSizeAttained
					return
				}
				if imagesCount >= MaxImagesCount {
					done <- store.ErrMaxImagesAttained
					return
				}
				part, err := w.CreateFormFile(key, value.Name())
				if err != nil {
					done <- err
					return
				}
				if _, err := io.Copy(part, value); err != nil {
					done <- err
					return
				}
				imagesCount++
			case string:
				if err := w.WriteField(key, value); err != nil {
					done <- err
				}
			case int64:
				if err := w.WriteField(key, strconv.Itoa(int(value))); err != nil {
					done <- err
				}
			}
		}
		dataType = w.FormDataContentType()
		done <- nil
	}(pipeWriter)
	closer := func() error {
		if err := pipeReader.Close(); err != nil {
			return err
		}
		if err := pipeWriter.Close(); err != nil {
			return err
		}
		return nil
	}
	return pipeReader, closer, dataType, done
}

// Close() обязательно
func mustOpen(f string) (*os.File, error) {
	r, err := os.Open(f)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (s *server) HandleLKShow() http.HandlerFunc {
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
		lk := &model.UserLK{
			UserID: id.(int64),
		}
		if err = s.store.UserLK().LKShow(lk); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		strID := strconv.Itoa(int(id.(int64)))
		file, err := mustOpen(s.store.Images() + "/" + "pfpimages" + "/" + strID + ".png")
		if err != nil {
			slog.Info("no pfp image")
		}
		if !os.IsNotExist(err) {
			defer file.Close()
			lk.Image = file
		}
		m := structs.Map(lk)
		pr, closer, dataType, done := download(m)
		w.Header().Set("Content-Type", dataType)
		w.WriteHeader(http.StatusOK)
		go func() {
			n, err := io.Copy(w, pr)
			slog.Info(strconv.Itoa(int(n)))
			if err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
		}()
		if err = <-done; err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		close(done)
		if err = closer(); err != nil {
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
		//w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(cards); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}

func (s *server) HandleGroupShowUsingTime() http.HandlerFunc {
	type request struct {
		GroupID int64 `json:"group_id"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		// stroka, err := io.ReadAll(r.Body)
		// if err != nil {
		// 	slog.Info("Pizdets")
		// 	return
		// }
		// slog.Info(string(stroka))
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
		g := &model.Group{
			UserID:  id.(int64),
			GroupID: req.GroupID,
		}
		cards, err := s.store.Group().ShowUsingTime(g)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
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
		session.Options.MaxAge = -1
		session.Save(r, w)
		w.WriteHeader(http.StatusOK)
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
		w.Header().Set("Content-Type", "application/json")
		if err = json.NewEncoder(w).Encode(group.GroupID); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
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
		GroupID string `json:"group_id"`
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
		groupID, err := strconv.ParseInt(req.GroupID, 10, 64)
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}
		group := &model.Group{
			UserID:  id.(int64),
			GroupID: groupID,
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

func (s *server) uploadCard(r *http.Request, values map[string]interface{}, userID int64) (uuidSlice uuid.UUIDs, err error) {
	reader, err := r.MultipartReader()
	if err != nil {
		return nil, err
	}
	b := bytes.NewBuffer(nil)
	uuidSlice = make(uuid.UUIDs, 0)
	for {
		part, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		switch part.FormName() {
		case "front_side_images":
			uuid := uuid.New()
			cardImage := &model.CardImages{
				UserID:  userID,
				ImageID: uuid,
			}
			file, err := s.store.CardImages().Add(cardImage, true)
			if err != nil {
				s.store.CardImages().DeleteOnError(uuidSlice)
				return uuidSlice, err
			}
			if _, err := io.Copy(file, part); err != nil {
				s.store.CardImages().DeleteOnError(uuidSlice)
				return uuidSlice, err
			}
			file.Close()
			uuidSlice = append(uuidSlice, uuid)
		case "back_side_images":
			uuid := uuid.New()
			cardImage := &model.CardImages{
				UserID:  userID,
				ImageID: uuid,
			}
			file, err := s.store.CardImages().Add(cardImage, false)
			if err != nil {
				s.store.CardImages().DeleteOnError(uuidSlice)
				return uuidSlice, err
			}
			if _, err := io.Copy(file, part); err != nil {
				s.store.CardImages().DeleteOnError(uuidSlice)
				return uuidSlice, err
			}
			file.Close()
			uuidSlice = append(uuidSlice, uuid)
		case "group_id":
			b.ReadFrom(part)
			tmp, err := strconv.Atoi(b.String())
			err = errors.New("blabla")
			if err != nil {
				s.store.CardImages().DeleteOnError(uuidSlice)
				return uuidSlice, err
			}
			values["group_id"] = int64(tmp)
			b.Reset()
		case "front_side":
			b.ReadFrom(part)
			values["front_side"] = b.String()
			b.Reset()
		case "back_side":
			b.ReadFrom(part)
			values["back_side"] = b.String()
			b.Reset()
		}
	}
	return uuidSlice, err
}

//	func (s *server) HandleCardCreate() http.HandlerFunc {
//		type request struct {
//			FrontSide string `json:"front_side"`
//			BackSide  string `json:"back_side"`
//			GroupID   int64  `json:"group_id"`
//		}
//		type reqImg struct {
//			FrontSide       string     `json:"front_side"`
//			BackSide        string     `json:"back_side"`
//			GroupID         int64      `json:"group_id"`
//			FrontsideImages []*bytes.Buffer `json:"front_side_images"`
//			BacksideImages  []*bytes.Buffer `json:"back_side_images"`
//		}
//		return func(w http.ResponseWriter, r *http.Request) {
//			req := &request{}
//			if err := json.NewDecoder(r.Body).Decode(req); err != nil {
//				s.error(w, r, http.StatusBadRequest, err)
//				return
//			}
//			session, err := s.sessionStore.Get(r, sessionKeyName)
//			if err != nil {
//				s.error(w, r, http.StatusInternalServerError, err)
//				return
//			}
//			id, ok := session.Values["user_id"]
//			if !ok {
//				s.error(w, r, http.StatusUnauthorized, err)
//				return
//			}
//			card := &model.Card{
//				UserID:    id.(int64),
//				FrontSide: req.FrontSide,
//				BackSide:  req.BackSide,
//				GroupID:   req.GroupID,
//			}
//			if err := s.store.Card().Create(card); err != nil {
//				s.error(w, r, http.StatusUnprocessableEntity, err)
//				return
//			}
//		}
//	}
func (s *server) HandleCardCreate() http.HandlerFunc {

	type reqImg struct {
		FrontSide       string           `structs:"front_side" mapstructure:"front_side"`
		BackSide        string           `structs:"back_side" mapstructure:"back_side"`
		GroupID         int64            `structs:"group_id" mapstructure:"group_id"`
		FrontsideImages []*io.PipeReader `structs:"front_side_images" mapstructure:"front_side_images"`
		BacksideImages  []*io.PipeReader `structs:"back_side_images" mapstructure:"back_side_images"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		req := &reqImg{}
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
		m := structs.Map(req)
		uuIDs, err := s.uploadCard(r, m, id.(int64))
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		card := &model.Card{
			UserID:    id.(int64),
			FrontSide: m["front_side"].(string),
			BackSide:  m["back_side"].(string),
			GroupID:   m["group_id"].(int64),
		}
		if err := s.store.Card().Create(card); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
		}
		for _, v := range uuIDs {
			if err := s.store.CardImages().CardIDUpdate(card.ID, v); err != nil {
				s.error(w, r, http.StatusInternalServerError, err)
				return
			}
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
	}
}

func (s *server) HandlePFPUpload() http.HandlerFunc {
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
		strID := strconv.Itoa(int(id.(int64)))
		file, err := os.Create(s.store.Images() + "/" + "pfpimages" + "/" + strID + ".png")
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		defer file.Close()
		if _, err := io.Copy(file, r.Body); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
	}
}
func (s *server) HandleShowAllGroups() http.HandlerFunc {
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
		u := &model.User{
			ID: id.(int64),
		}
		groups, err := s.store.User().ShowALLGroups(u)
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(groups); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}
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
		// cardImages, err := s.store.CardImages().Show(cards)
		// if err != nil {
		// 	s.error(w, r, http.StatusInternalServerError, err)
		// 	return
		// }
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
