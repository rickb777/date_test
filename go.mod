module github.com/rickb777/date_test

go 1.21

toolchain go1.21.5

require (
	github.com/go-sql-driver/mysql v1.8.1
	github.com/jackc/pgx/v4 v4.18.2
	github.com/lib/pq v1.10.9
	github.com/mattn/go-sqlite3 v1.14.22
	github.com/onsi/gomega v1.33.0
	github.com/rickb777/date/v2 v2.0.12-beta
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/gofrs/uuid v4.4.0+incompatible // indirect
	github.com/google/go-cmp v0.6.0 // indirect
	github.com/govalues/decimal v0.1.23 // indirect
	github.com/jackc/chunkreader/v2 v2.0.1 // indirect
	github.com/jackc/pgconn v1.14.3 // indirect
	github.com/jackc/pgio v1.0.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgproto3/v2 v2.3.3 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgtype v1.14.0 // indirect
	github.com/kr/pretty v0.3.1 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/rickb777/period v1.0.4-beta // indirect
	github.com/rickb777/plural v1.4.2 // indirect
	github.com/rogpeppe/go-internal v1.12.0 // indirect
	github.com/shopspring/decimal v1.3.1 // indirect
	golang.org/x/crypto v0.22.0 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

// for testing - check manually for consistent Git branch
//replace github.com/rickb777/date/v2 => ../date
