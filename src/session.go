package fox

import (
	"errors"
	"math/rand"
	"os"
	"time"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/robin-andreasson/fox/parser"
)

type SessionOptions struct {
	Secret  string // String used in session id hashing
	TimeOut int    // Milliseconds until session is cleared on the server. If not set, defaults to 24 hours
	Path    string // Path to the db file
	Cookie  CookieAttributes
}

var sessionOpt SessionOptions

var memoryStorage = map[string]string{}

const initdb string = `
    CREATE TABLE IF NOT EXISTS sessions (
    	sessID TEXT PRIMARY KEY,
    	payload TEXT,
    	timeout INT
    );
    
	CREATE INDEX IF NOT EXISTS "timeout_i" ON sessions ("timeout" ASC);
`

/*
Initialize Sessions

NOTE: sqlite3 is used as session store meaning that a gcc compiler is needed
*/
func Session(options SessionOptions) error {

	if options.Secret == "" {
		return errors.New("empty secret will lead to unsafe id hashing")
	}

	if options.TimeOut == 0 {
		options.TimeOut = 1000 * 60 * 60 * 24
	}

	if _, err := os.Stat(options.Path); os.IsNotExist(err) {
		return errors.New("sql file does not exist")
	}

	if ext, err := Ext(options.Path); err != nil || (ext != "db" && ext != "sql") {
		return errors.New("'db' and 'sql' are the only valid file extensions")
	}

	db, err := sql.Open("sqlite3", options.Path)

	if err != nil {
		return err
	}

	defer db.Close()

	db.Exec(initdb)

	sessionOpt = options

	return nil
}

func handleSession(sessID string, c *Context) error {
	if sessID == "" {
		return errors.New("no session id")
	}

	db, err := sql.Open("sqlite3", sessionOpt.Path)

	if err != nil {
		return err
	}

	rand.Seed(time.Now().Unix())
	if rand.Float64() <= 0.1 {
		if _, err := db.Exec("DELETE FROM sessions WHERE timeout<=?", time.Now().UnixMilli()); err != nil {
			return err
		}
	}

	stmt, err := db.Prepare("SELECT payload FROM sessions WHERE sessID=?")

	if err != nil {
		return err
	}

	row := stmt.QueryRow(sessID)

	payload := ""

	if err := row.Scan(&payload); err != nil {
		return err
	}

	if payload == "" {
		return errors.New("session id is not in storage")
	}

	return parser.JSONUnmarshal(payload, &c.Session)
}
