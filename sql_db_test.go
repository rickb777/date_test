package date_test

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	. "github.com/onsi/gomega"
	"github.com/rickb777/date/v2"
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
	g := NewGomegaWithT(t)

	examples := []date.Date{
		date.Zero,
		date.New(2000, 3, 31),
		date.New(2020, 12, 31),
	}

	driver, db := connect(t)
	//defer cleanup(d.DB())

	_, err := db.Exec(`DROP TABLE IF EXISTS dates`)
	g.Expect(err).NotTo(HaveOccurred())

	_, err = db.Exec(createSql[driver])
	g.Expect(err).NotTo(HaveOccurred())

	for _, e := range examples {
		n, err := db.Exec(insertSql[driver], e, e, e, e)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(n).NotTo(BeEquivalentTo(1))
	}

	rows, err := db.Query(`SELECT * FROM dates ORDER BY ID`)
	g.Expect(err).NotTo(HaveOccurred())
	defer rows.Close()

	j := 0
	for rows.Next() {
		var id int
		var d1, d2 date.Date
		var s1, s2 date.Date
		err = rows.Scan(&id, &d1, &d2, &s1, &s2)
		g.Expect(err).NotTo(HaveOccurred())
		g.Expect(d1).To(Equal(examples[j]))
		g.Expect(d2).To(Equal(examples[j]))
		g.Expect(s1).To(Equal(examples[j]))
		g.Expect(s2).To(Equal(examples[j]))
		j++
	}
}
