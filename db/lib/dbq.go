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

	GetRSSFeed    *sql.Stmt
	InsertRSSFeed *sql.Stmt
	UpdateRSSFeed *sql.Stmt
	DeleteRSSFeed *sql.Stmt

	GetItemFeed    *sql.Stmt
	InsertItemFeed *sql.Stmt
	UpdateItemFeed *sql.Stmt
	DeleteItemFeed *sql.Stmt
}
