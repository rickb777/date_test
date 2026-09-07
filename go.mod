module github.com/rickb777/date_test

go 1.26.0

require (
	github.com/go-sql-driver/mysql v1.10.1
	github.com/jackc/pgx/v5 v5.10.0
	github.com/lib/pq v1.12.3
	github.com/mattn/go-sqlite3 v1.14.52
	github.com/rickb777/date/v2 v2.3.18
	github.com/rickb777/expect v1.3.4
)

require (
	filippo.io/edwards25519 v1.2.0 // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/govalues/decimal v0.1.36 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/puddle/v2 v2.2.2 // indirect
	github.com/rickb777/period v1.1.0 // indirect
	github.com/rickb777/plural/v2 v2.1.1 // indirect
	golang.org/x/sync v0.22.0 // indirect
	golang.org/x/text v0.41.0 // indirect
)

// for testing - check manually for consistent Git branch
//replace github.com/rickb777/date/v2 => ../date
