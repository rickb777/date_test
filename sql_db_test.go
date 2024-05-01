package date_test

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	. "github.com/onsi/gomega"
	"github.com/rickb777/date/v2"
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
	g := NewGomegaWithT(t)

	valuerµ.Lock()
	defer valuerµ.Unlock()
	date.Valuer = date.ValueAsString

	driver, db := connect(t)
	//defer cleanup(d.DB())

	_, err := db.Exec(`DROP TABLE IF EXISTS dates`)
	g.Expect(err).NotTo(HaveOccurred())

	create := strings.NewReplacer(
		"{D1}", "TEXT",
		"{D2}", "INTEGER",
		"{D3}", "DATE",
	).Replace(createSql[driver])

	//t.Log(create)

	_, err = db.Exec(create)
	g.Expect(err).NotTo(HaveOccurred())

	for _, e := range examples {
		n, err := db.Exec(insertSql[driver], e, int64(e), e)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(n).NotTo(BeEquivalentTo(1))
	}

	checkDatesValues(err, db, g)
}

func TestDatesCrud_using_database_valuer_as_integer(t *testing.T) {
	g := NewGomegaWithT(t)

	valuerµ.Lock()
	defer valuerµ.Unlock()
	date.Valuer = date.ValueAsInt

	driver, db := connect(t)
	//defer cleanup(d.DB())

	_, err := db.Exec(`DROP TABLE IF EXISTS dates`)
	g.Expect(err).NotTo(HaveOccurred())

	create := strings.NewReplacer(
		"{D1}", "INTEGER",
		"{D2}", "INTEGER",
		"{D3}", "DATE",
	).Replace(createSql[driver])

	//t.Log(create)

	_, err = db.Exec(create)
	g.Expect(err).NotTo(HaveOccurred())

	for _, e := range examples {
		n, err := db.Exec(insertSql[driver], e, int64(e), e.String())
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(n).NotTo(BeEquivalentTo(1))
	}

	checkDatesValues(err, db, g)
}

func checkDatesValues(err error, db *sql.DB, g *WithT) {
	rows, err := db.Query(`SELECT * FROM dates ORDER BY ID`)
	g.Expect(err).NotTo(HaveOccurred())
	defer rows.Close()

	j := 0
	for rows.Next() {
		var id int
		var d1, d2, d3 date.Date
		err = rows.Scan(&id, &d1, &d2, &d3)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(d1).To(Equal(examples[j]))
		g.Expect(d2).To(Equal(examples[j]))
		g.Expect(d3).To(Equal(examples[j]))
		j++
	}
}
