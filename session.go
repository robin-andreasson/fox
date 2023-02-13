package fox

import (
	"log"
	"math/rand"
	"os"
	"time"

	"database/sql"

	_ "github.com/mattn/go-sqlite3"

	"github.com/robin-andreasson/fox/parser"
)

type SessionOptions struct {
	Secret           string  // String used in session id hashing
	TimeOut          int     // Milliseconds until session is expired in the session store, defaults to 24 hours
	ClearProbability float64 // value between 0 - 100 that represents the chance of fox clearing expired sessions
	Path             string  // path to the session store
	Cookie           CookieAttributes

	init bool
}

var sessionOpt SessionOptions

const dbinit string = `
    CREATE TABLE IF NOT EXISTS sessions (
    	sessID TEXT PRIMARY KEY,
    	payload TEXT,
    	timeout INT
    );
    
	CREATE INDEX IF NOT EXISTS "timeout_i" ON sessions ("timeout" ASC);
`

/*
Initialize Sessions

returns Session middleware

NOTE: sqlite3 is used as session store meaning that a gcc compiler is needed
*/
func Session(options SessionOptions) {

	if options.Secret == "" {
		log.Panic("zero value secret will lead to unsafe id hashing")
	}

	if options.TimeOut == 0 {
		options.TimeOut = 1000 * 60 * 60 * 24
	}

	if options.ClearProbability < 0 || options.ClearProbability > 100 {
		log.Panic("invalid value for ClearProbability, acceptable values are between 0 and 100")
	}

	if _, err := os.Stat(options.Path); os.IsNotExist(err) {
		log.Panic("could not find target session store")
	}

	if extension, err := Ext(options.Path); err != nil || (extension != "db" && extension != "sql") {
		log.Panic("invalid session store extension, sql or db is required")
	}

	if err := os.Truncate(options.Path, 0); err != nil {
		log.Panic("could not clear session store before initialization")
	}

	db, err := sql.Open("sqlite3", options.Path)

	if err != nil {
		log.Panic("error opening session store")
	}

	defer db.Close()

	if _, err := db.Exec(dbinit); err != nil {
		log.Panic("error initializing table and timeout index")
	}

	sessionOpt = options
	sessionOpt.ClearProbability = sessionOpt.ClearProbability / 100
	sessionOpt.init = true
}

func handleSession(sessID string, c *Context) {
	if !sessionOpt.init {
		return
	}

	if sessID == "" {
		return
	}

	db, err := sql.Open("sqlite3", sessionOpt.Path)

	if err != nil {
		return
	}

	rand.Seed(time.Now().Unix())
	if rand.Float64() <= sessionOpt.ClearProbability {
		if _, err := db.Exec("DELETE FROM sessions WHERE timeout<=?", time.Now().UnixMilli()); err != nil {
			return
		}
	}

	stmt, err := db.Prepare("SELECT payload FROM sessions WHERE sessID=?")

	if err != nil {
		return
	}

	row := stmt.QueryRow(sessID)

	payload := ""

	if err := row.Scan(&payload); err != nil {
		return
	}

	if payload == "" {
		return
	}

	parser.JSONUnmarshal(payload, &c.Session)
}
