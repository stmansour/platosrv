package db

import (
	"context"
	"database/sql"
)

// ItemFeed defines an item feed by its url and establishes an IFID
//-----------------------------------------------------------------------------
type ItemFeed struct {
	IFID  int64 // unique id for this record
	IID   int64 // the Item
	RSSID int64 // the feed that referenced this Item
}

// DeleteItemFeed deletes the ItemFeed with the specified id from the database
//
// INPUTS
// ctx - db context
// id - IFID of the record to read
//
// RETURNS
// Any errors encountered, or nil if no errors
//-----------------------------------------------------------------------------
func DeleteItemFeed(ctx context.Context, id int64) error {
	return genericDelete(ctx, "ItemFeed", Pdb.Prepstmt.DeleteItemFeed, id)
}

// GetItemFeed reads and returns a ItemFeed structure
//
// INPUTS
// ctx - db context
// id - IFID of the record to read
//
// RETURNS
// ErrSessionRequired if the session is invalid
// nil if the session is valid
//-----------------------------------------------------------------------------
func GetItemFeed(ctx context.Context, id int64) (ItemFeed, error) {
	var a ItemFeed
	if !ValidateSession(ctx) {
		return a, ErrSessionRequired
	}

	fields := []interface{}{id}
	stmt, row := getRowFromDB(ctx, Pdb.Prepstmt.GetItemFeed, fields)
	if stmt != nil {
		defer stmt.Close()
	}
	return a, ReadItemFeed(row, &a)
}

// InsertItemFeed writes a new ItemFeed record to the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
//-----------------------------------------------------------------------------
func InsertItemFeed(ctx context.Context, a *ItemFeed) (int64, error) {
	fields := []interface{}{
		a.IID,
		a.RSSID,
	}

	var err error
	_, _, a.IFID, err = genericInsert(ctx, "ItemFeed", Pdb.Prepstmt.InsertItemFeed, fields, a)
	return a.IFID, err
}

// ReadItemFeed reads a full ItemFeed structure of data from the database based
// on the supplied Rows pointer.
//
// INPUTS
// row - db Row pointer
// a   - pointer to struct to fill
//
// RETURNS
//
// ErrSessionRequired if the session is invalid
// nil if the session is valid
//-----------------------------------------------------------------------------
func ReadItemFeed(row *sql.Row, a *ItemFeed) error {
	err := row.Scan(
		&a.IFID,
		&a.IID,
		&a.RSSID,
	)
	SkipSQLNoRowsError(&err)
	return err
}

// ReadItemFeeds reads a full ItemFeed structure of data from the database based
// on the supplied Rows pointer.
//
// INPUTS
// row - db Row pointer
// a   - pointer to struct to fill
//
// RETURNS
//
// ErrSessionRequired if the session is invalid
// nil if the session is valid
//-----------------------------------------------------------------------------
func ReadItemFeeds(rows *sql.Rows, a *ItemFeed) error {
	err := rows.Scan(
		&a.IFID,
		&a.IID,
		&a.RSSID,
	)
	SkipSQLNoRowsError(&err)
	return err
}

// UpdateItemFeed updates an existing ItemFeed record in the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
//-----------------------------------------------------------------------------
func UpdateItemFeed(ctx context.Context, a *ItemFeed) error {
	fields := []interface{}{
		a.IID,
		a.RSSID,
		a.IFID,
	}
	var err error
	_, err = genericUpdate(ctx, Pdb.Prepstmt.UpdateItemFeed, fields)
	return updateError(err, "ItemFeed", *a)
}
