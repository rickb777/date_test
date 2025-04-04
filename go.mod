module github.com/rickb777/date_test

go 1.24.1

require (
	github.com/go-sql-driver/mysql v1.9.1
	github.com/jackc/pgx/v4 v4.18.3
	github.com/lib/pq v1.10.9
	github.com/mattn/go-sqlite3 v1.14.24
	github.com/rickb777/date/v2 v2.1.8
	github.com/rickb777/expect v0.9.0
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/google/go-cmp v0.7.0 // indirect
	github.com/govalues/decimal v0.1.36 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20240606120523-5a60cdf6a761 // indirect
	github.com/jackc/pgtype v1.14.4 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rickb777/period v1.0.9 // indirect
	github.com/rickb777/plural v1.4.2 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	golang.org/x/crypto v0.36.0 // indirect
	golang.org/x/text v0.23.0 // indirect
)

// for testing - check manually for consistent Git branch
//replace github.com/rickb777/date/v2 => ../date
