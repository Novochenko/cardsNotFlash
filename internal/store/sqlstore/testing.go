package sqlstore

import (
	"database/sql"
	"fmt"
	"strings"
	"testing"
)

// создает коннект и функцию удаления всех таблиц из бд
func TestDB(t *testing.T, database string) (*sql.DB, func(...string)) {
	t.Helper()

	db, err := sql.Open("mysql", database)
	if err != nil {
		t.Fatal(err)
	}

	if err := db.Ping(); err != nil {
		t.Fatal(err)
	}
	return db, func(tables ...string) {
		if len(tables) > 0 {
			db.Exec(fmt.Sprintf("TRUNCATE %s;", strings.Join(tables, ", ")))
		}
		db.Close()
	}
}
