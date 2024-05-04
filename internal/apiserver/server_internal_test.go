package apiserver

import (
	"bytes"
	"encoding/json"
	"firstRestAPI/internal/model"
	"firstRestAPI/internal/store/teststore"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/securecookie"
	"github.com/stretchr/testify/assert"
)

// сделай тесты, бро <3
func TestServer_AuthenticateUser(t *testing.T) {
	store := teststore.New()
	u := model.TestUser(t)
	store.User().Create(u)

	testCases := []struct {
		name         string
		cookieValue  map[interface{}]interface{}
		expectedCode int
	}{
		{
			name: "authenticated",
			cookieValue: map[interface{}]interface{}{
				"user_id": u.ID,
			},
			expectedCode: http.StatusOK,
		},
		{
			name:         "not authenticated",
			cookieValue:  nil,
			expectedCode: http.StatusUnauthorized,
		},
	}

	config := Config{
		BindAddr:      ":8080",
		LogLevel:      "debug",
		LocalHostMode: true,
	}
	secretKey := []byte("secret")
	session, err := NewRedisSessions(&config)
	if err != nil {
		return
	}
	serv := newServer(store, &config, session)
	sc := securecookie.New([]byte(secretKey), nil)
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			req, _ := http.NewRequest(http.MethodGet, "/", nil)
			cookieStr, _ := sc.Encode(sessionKeyName, tc.cookieValue)
			req.Header.Set("Cookie", fmt.Sprintf("%s=%s", sessionKeyName, cookieStr))
			serv.authenticateUser(handler).ServeHTTP(rec, req)
			assert.Equal(t, tc.expectedCode, rec.Code)
		})
	}
}

func TestServer_HandleUsersCreate(t *testing.T) {
	config := &Config{
		BindAddr:      ":8080",
		LogLevel:      "debug",
		LocalHostMode: true,
	}
	session, err := NewRedisSessions(config)
	if err != nil {
		return
	}
	s := newServer(teststore.New(), config, session)
	testCases := []struct {
		name         string
		payload      interface{}
		expeted_code int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    "user5@example.org",
				"password": "password",
				"nickname": "example",
			},
			expeted_code: http.StatusCreated,
		},
		{
			name:         "invalid_payload",
			payload:      "invalid",
			expeted_code: http.StatusBadRequest,
		},
		{
			name: "invalid params",
			payload: map[string]string{
				"email": "invalid",
			},
			expeted_code: http.StatusUnprocessableEntity,
		},
		{
			name: "no nick",
			payload: map[string]string{
				"email":    "user2@example.org",
				"password": "password",
			},
			expeted_code: http.StatusUnprocessableEntity,
		},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/users", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expeted_code, rec.Code)
		})

	}
}

func TestServer_HandleSessionsCreate(t *testing.T) {

	u := model.TestUser(t)
	store := teststore.New()
	store.User().Create(u)
	config := &Config{
		BindAddr:      ":8080",
		LogLevel:      "debug",
		LocalHostMode: true,
	}
	session, err := NewRedisSessions(config)
	if err != nil {
		return
	}
	s := newServer(store, config, session)
	testCases := []struct {
		name          string
		payload       interface{}
		expected_code int
	}{
		{
			name: "valid",
			payload: map[string]string{
				"email":    u.Email,
				"password": u.Password,
			},
			expected_code: http.StatusOK,
		},
		{
			name:          "invalid",
			payload:       "string",
			expected_code: http.StatusBadRequest,
		},
		{
			name: "invalid email",
			payload: map[string]string{
				"email":    "invalid",
				"password": u.Password,
			},
			expected_code: http.StatusUnauthorized,
		},
		{
			name: "invalid password",
			payload: map[string]string{
				"email":    u.Email,
				"password": "invalid",
			},
			expected_code: http.StatusUnauthorized,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			b := &bytes.Buffer{}
			json.NewEncoder(b).Encode(tc.payload)
			req, _ := http.NewRequest(http.MethodPost, "/sessions", b)
			s.ServeHTTP(rec, req)
			assert.Equal(t, tc.expected_code, rec.Code)
		})
	}
}
