package db

import "database/sql"

// PrepSQL is the structure containing all the prepared statements
type PrepSQL struct {
	GetExch    *sql.Stmt
	InsertExch *sql.Stmt
	UpdateExch *sql.Stmt
	DeleteExch *sql.Stmt

	GetItem    *sql.Stmt
	InsertItem *sql.Stmt
	UpdateItem *sql.Stmt
	DeleteItem *sql.Stmt
}
