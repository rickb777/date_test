package date_test

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rickb777/date/v2"
	"github.com/rickb777/expect"
	"os"
	"strings"
	"sync"
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
  ID   INTEGER PRIMARY KEY AUTOINCREMENT,
  D1   {D1} NOT NULL,
  D2   {D2} NOT NULL,
  D3   {D3} NOT NULL
)
`
const createMysql = `
CREATE TABLE dates (
  ID   INTEGER PRIMARY KEY AUTO_INCREMENT,
  D1   {D1} NOT NULL,
  D2   {D2} NOT NULL,
  D3   {D3} NOT NULL
) Engine=InnoDB
`

const createPostgres = `
CREATE TABLE dates (
  ID   SERIAL PRIMARY KEY,
  D1   {D1} NOT NULL,
  D2   {D2} NOT NULL,
  D3   {D3} NOT NULL
)
`

var createSql = map[string]string{
	sqlite3:  createSqlite,
	mysql:    createMysql,
	postgres: createPostgres,
}

const insertMain = `INSERT INTO dates (D1, D2, D3) VALUES `

var insertSql = map[string]string{
	sqlite3:  insertMain + "(?,?,?)",
	mysql:    insertMain + "(?,?,?)",
	postgres: insertMain + "($1,$2,$3)",
}

var examples = []date.Date{
	date.Zero,
	date.New(2000, 3, 31),
	date.New(2020, 12, 31),
}

var valuerµ = sync.Mutex{}

func TestDatesCrud_using_database_valuer_as_string(t *testing.T) {
	valuerµ.Lock()
	defer valuerµ.Unlock()
	date.Valuer = date.ValueAsString

	driver, db := connect(t)
	//defer cleanup(d.DB())

	_, err := db.Exec(`DROP TABLE IF EXISTS dates`)
	expect.Error(err).Not().ToHaveOccurred(t)

	create := strings.NewReplacer(
		"{D1}", "TEXT",
		"{D2}", "INTEGER",
		"{D3}", "DATE",
	).Replace(createSql[driver])

	//t.Log(create)

	_, err = db.Exec(create)
	expect.Error(err).Not().ToHaveOccurred(t)

	for _, e := range examples {
		n, err := db.Exec(insertSql[driver], e, int64(e), e)
		expect.Error(err).Not().ToHaveOccurred(t)
		expect.Number(n.RowsAffected()).ToBe(t, int64(1))
	}

	checkDatesValues(err, db, t)
}

func TestDatesCrud_using_database_valuer_as_integer(t *testing.T) {
	valuerµ.Lock()
	defer valuerµ.Unlock()
	date.Valuer = date.ValueAsInt

	driver, db := connect(t)
	//defer cleanup(d.DB())

	_, err := db.Exec(`DROP TABLE IF EXISTS dates`)
	expect.Error(err).Not().ToHaveOccurred(t)

	create := strings.NewReplacer(
		"{D1}", "INTEGER",
		"{D2}", "INTEGER",
		"{D3}", "DATE",
	).Replace(createSql[driver])

	//t.Log(create)

	_, err = db.Exec(create)
	expect.Error(err).Not().ToHaveOccurred(t)

	for _, e := range examples {
		n, err := db.Exec(insertSql[driver], e, int64(e), e.String())
		expect.Error(err).Not().ToHaveOccurred(t)
		expect.Number(n.RowsAffected()).ToBe(t, int64(1))
	}

	checkDatesValues(err, db, t)
}

func checkDatesValues(err error, db *sql.DB, t *testing.T) {
	rows, err := db.Query(`SELECT * FROM dates ORDER BY ID`)
	expect.Error(err).Not().ToHaveOccurred(t)
	defer rows.Close()

	j := 0
	for rows.Next() {
		var id int
		var d1, d2, d3 date.Date
		err = rows.Scan(&id, &d1, &d2, &d3)
		expect.Error(err).Not().ToHaveOccurred(t)
		expect.Any(d1).ToBe(t, examples[j])
		expect.Any(d2).ToBe(t, examples[j])
		expect.Any(d3).ToBe(t, examples[j])
		j++
	}
}
