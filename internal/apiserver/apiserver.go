package apiserver

import (
	"database/sql"
	"firstRestAPI/internal/store/sqlstore"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/go-redis/redis"
	"github.com/rbcervilla/redisstore"
)

func Start(config *Config) error {
	db, err := newDB(config)
	if err != nil {
		return err
	}
	defer db.Close()
	store := sqlstore.New(db)
	sessionStore, err := NewRedisSessions(config)
	if err != nil {
		return err
	}
	//sessionStore = sessions.NewCookieStore([]byte(config.SessionKey))
	s := newServer(store, config, sessionStore)
	return http.ListenAndServe(config.BindAddr, s)
}
func NewRedisSessions(config *Config) (*redisstore.RedisStore, error) {
	var RedisURL string
	if config.LocalHostMode {
		RedisURL = config.RedisURL.LocalHost
	} else {
		RedisURL = config.RedisURL.Docker
	}
	client := redis.NewClient(&redis.Options{
		Addr:     RedisURL,
		Password: "",
	})
	slog.Info("Connecting to redis")
	sessionStore, err := redisstore.NewRedisStore(client)
	if err != nil {
		slog.Error("Error connecting redis")
		return nil, err
	}
	return sessionStore, err
}

func newDB(config *Config) (*sql.DB, error) {
	slog.Info("Connecting to mysql serv\n")
	var db *sql.DB
	var err error
	if config.LocalHostMode {
		db, err = sql.Open("mysql", config.DatabaseURL.FullName)

	} else {
		db, err = sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true",
			config.DatabaseURL.User,
			config.DatabaseURL.Password,
			config.DatabaseURL.Host,
			//databaseURL.Port,
			config.DatabaseURL.DBName))
	}
	//db, err := sql.Open("mysql", databaseURL.FullName)

	if err != nil {
		slog.Error("error opening mysql\n")
		return nil, err
	}
	return db, nil
}
