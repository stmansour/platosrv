package db

import (
	"context"
	"database/sql"
	"plato/session"
	"time"
)

// Item defines a date and a rent amount for a property. A Item record
// is part of a group or list. The group is defined by the RSLID
//-----------------------------------------------------------------------------
type Item struct {
	IID         int64 // unique id for this record
	Title       string
	Description string
	PubDt       time.Time
	Link        string
	LastModTime time.Time // when was the record last written
	LastModBy   int64     // id of user that did the modify
	CreateTime  time.Time // when was this record created
	CreateBy    int64     // id of user that created it
}

// DeleteItem deletes the Item with the specified id from the database
//
// INPUTS
// ctx - db context
// id - IID of the record to read
//
// RETURNS
// Any errors encountered, or nil if no errors
//-----------------------------------------------------------------------------
func DeleteItem(ctx context.Context, id int64) error {
	return genericDelete(ctx, "Item", Pdb.Prepstmt.DeleteItem, id)
}

// GetItem reads and returns a Item structure
//
// INPUTS
// ctx - db context
// id - IID of the record to read
//
// RETURNS
// ErrSessionRequired if the session is invalid
// nil if the session is valid
//-----------------------------------------------------------------------------
func GetItem(ctx context.Context, id int64) (Item, error) {
	var a Item
	if !ValidateSession(ctx) {
		return a, ErrSessionRequired
	}

	fields := []interface{}{id}
	stmt, row := getRowFromDB(ctx, Pdb.Prepstmt.GetItem, fields)
	if stmt != nil {
		defer stmt.Close()
	}
	return a, ReadItem(row, &a)
}

// InsertItem writes a new Item record to the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
//-----------------------------------------------------------------------------
func InsertItem(ctx context.Context, a *Item) (int64, error) {
	sess, ok := session.GetSessionFromContext(ctx)
	if !ok {
		return a.IID, ErrSessionRequired
	}
	fields := []interface{}{
		a.Title,
		a.Description,
		a.PubDt,
		a.Link,
		sess.UID,
		sess.UID,
	}

	var err error
	a.CreateBy, a.LastModBy, a.IID, err = genericInsert(ctx, "Item", Pdb.Prepstmt.InsertItem, fields, a)
	return a.IID, err
}

// ReadItem reads a full Item structure of data from the database based
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
func ReadItem(row *sql.Row, a *Item) error {
	err := row.Scan(
		&a.IID,
		&a.Title,
		&a.Description,
		&a.PubDt,
		&a.Link,
		&a.LastModTime,
		&a.LastModBy,
		&a.CreateTime,
		&a.CreateBy,
	)
	SkipSQLNoRowsError(&err)
	return err
}

// ReadItems reads a full Item structure of data from the database based
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
func ReadItems(rows *sql.Rows, a *Item) error {
	err := rows.Scan(
		&a.IID,
		&a.Title,
		&a.Description,
		&a.PubDt,
		&a.Link,
		&a.LastModTime,
		&a.LastModBy,
		&a.CreateTime,
		&a.CreateBy,
	)
	SkipSQLNoRowsError(&err)
	return err
}

// UpdateItem updates an existing Item record in the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
//-----------------------------------------------------------------------------
func UpdateItem(ctx context.Context, a *Item) error {
	sess, ok := session.GetSessionFromContext(ctx)
	if !ok {
		return ErrSessionRequired
	}
	fields := []interface{}{
		a.Title,
		a.Description,
		a.PubDt,
		a.Link,
		sess.UID,
		a.IID, // 49
	}
	var err error
	a.LastModBy, err = genericUpdate(ctx, Pdb.Prepstmt.UpdateItem, fields)
	return updateError(err, "Item", *a)
}
