package date_test

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rickb777/date"
	"github.com/rickb777/expect"
	"os"
	"testing"
)

const (
	mysql    = "mysql"
	postgres = "postgres"
	sqlite3  = "sqlite3"
)

// Environment:
// GO_DRIVER  - the driver (sqlite3, mysql, etc)
// GO_QUOTER  - the identifier quoter (ansi, mysql, none)
// GO_DSN     - the database DSN
// GO_VERBOSE - true for query logging

var verbose = false

func connect(t *testing.T) (string, *sql.DB) {
	driver, ok := os.LookupEnv("GO_DRIVER")
	if !ok {
		driver = sqlite3
	}

	dsn, ok := os.LookupEnv("GO_DSN")
	if !ok {
		dsn = "file::memory:?mode=memory&cache=shared"
	}

	db, err := sql.Open(driver, dsn)
	if err != nil {
		t.Fatalf("Warning: Unable to connect to %s (%v); test is only partially complete.\n\n", driver, err)
	}

	err = db.Ping()
	if err != nil {
		t.Fatalf("Warning: Unable to ping %s (%v); test is only partially complete.\n\n", driver, err)
	}

	fmt.Printf("Successfully connected to %s.\n", driver)
	return driver, db
}

const createSqlite = `
CREATE TABLE dates (
  ID         INTEGER PRIMARY KEY AUTOINCREMENT,
  DINTEGER   INTEGER NOT NULL,
  DSTRING1   TEXT NOT NULL,
  DSTRING2   TEXT NOT NULL,
  DDATE      DATE NOT NULL
)
`
const createMysql = `
CREATE TABLE dates (
  ID         INTEGER PRIMARY KEY AUTO_INCREMENT,
  DINTEGER   INTEGER NOT NULL,
  DSTRING1   TEXT NOT NULL,
  DSTRING2   TEXT NOT NULL,
  DDATE      DATE NOT NULL
) Engine=InnoDB
`

const createPostgres = `
CREATE TABLE dates (
  ID         SERIAL PRIMARY KEY,
  DINTEGER   INTEGER NOT NULL,
  DSTRING1   TEXT NOT NULL,
  DSTRING2   TEXT NOT NULL,
  DDATE      DATE NOT NULL
)
`

var createSql = map[string]string{
	sqlite3:  createSqlite,
	mysql:    createMysql,
	postgres: createPostgres,
}

const insertMain = `INSERT INTO dates (DINTEGER, DSTRING1, DSTRING2, DDATE) VALUES `

var insertSql = map[string]string{
	sqlite3:  insertMain + "(?,?,?,?)",
	mysql:    insertMain + "(?,?,?,?)",
	postgres: insertMain + "($1,$2,$3,$4)",
}

func TestDatesCrud_using_database(t *testing.T) {
	examples := []date.Date{
		{}, // zero date
		date.New(2000, 3, 31),
		date.New(2020, 12, 31),
	}

	driver, db := connect(t)
	//defer cleanup(d.DB())

	_, err := db.Exec(`DROP TABLE IF EXISTS dates`)
	expect.Error(err).Not().ToHaveOccurred(t)

	_, err = db.Exec(createSql[driver])
	expect.Error(err).Not().ToHaveOccurred(t)

	for _, e := range examples {
		es := date.DateString(e)
		n, err := db.Exec(insertSql[driver], e, e, es, es)
		expect.Error(err).Not().ToHaveOccurred(t)
		expect.Number(n.RowsAffected()).ToBe(t, int64(1))
	}

	rows, err := db.Query(`SELECT * FROM dates ORDER BY ID`)
	expect.Error(err).Not().ToHaveOccurred(t)
	defer rows.Close()

	j := 0
	for rows.Next() {
		var id int
		var d1, d2 date.Date
		var s1, s2 date.DateString
		err = rows.Scan(&id, &d1, &d2, &s1, &s2)
		expect.Error(err).Not().ToHaveOccurred(t)
		expect.Any(d1).ToBe(t, examples[j])
		expect.Any(d2).ToBe(t, examples[j])
		expect.Any(s1).ToBe(t, examples[j].DateString())
		expect.Any(s2).ToBe(t, examples[j].DateString())
		j++
	}
}
