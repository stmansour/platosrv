package db

import (
	"context"
	"database/sql"
	"platosrv/session"
	"time"
)

// Exch defines a date and a rent amount for a property. A Exch record
// is part of a group or list. The group is defined by the RSLID
//-----------------------------------------------------------------------------
type Exch struct {
	XID         int64     // unique id for this record
	Dt          time.Time // point in time when these values are valid
	Ticker      string    // the two currencies involved in this exchange rate
	Open        float64   // Opening value for this minute
	High        float64   // High value during this minute
	Low         float64   // Low value during this minute
	Close       float64   // Closing value for this minute
	LastModTime time.Time // when was the record last written
	LastModBy   int64     // id of user that did the modify
	CreateTime  time.Time // when was this record created
	CreateBy    int64     // id of user that created it
}

// DeleteExch deletes the Exch with the specified id from the database
//
// INPUTS
// ctx - db context
// id - XID of the record to read
//
// RETURNS
// Any errors encountered, or nil if no errors
//-----------------------------------------------------------------------------
func DeleteExch(ctx context.Context, id int64) error {
	return genericDelete(ctx, "Exch", Pdb.Prepstmt.DeleteExch, id)
}

// GetExch reads and returns a Exch structure
//
// INPUTS
// ctx - db context
// id - XID of the record to read
//
// RETURNS
// ErrSessionRequired if the session is invalid
// nil if the session is valid
//-----------------------------------------------------------------------------
func GetExch(ctx context.Context, id int64) (Exch, error) {
	var a Exch
	if !ValidateSession(ctx) {
		return a, ErrSessionRequired
	}

	fields := []interface{}{id}
	stmt, row := getRowFromDB(ctx, Pdb.Prepstmt.GetExch, fields)
	if stmt != nil {
		defer stmt.Close()
	}
	return a, ReadExch(row, &a)
}

// InsertExch writes a new Exch record to the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
//-----------------------------------------------------------------------------
func InsertExch(ctx context.Context, a *Exch) (int64, error) {
	sess, ok := session.GetSessionFromContext(ctx)
	if !ok {
		return a.XID, ErrSessionRequired
	}
	fields := []interface{}{
		a.Dt,
		a.Ticker,
		a.Open,
		a.High,
		a.Low,
		a.Close,
		sess.UID,
		sess.UID,
	}

	var err error
	a.CreateBy, a.LastModBy, a.XID, err = genericInsert(ctx, "Exch", Pdb.Prepstmt.InsertExch, fields, a)
	return a.XID, err
}

// ReadExch reads a full Exch structure of data from the database based
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
func ReadExch(row *sql.Row, a *Exch) error {
	err := row.Scan(
		&a.XID,
		&a.Dt,
		&a.Ticker,
		&a.Open,
		&a.High,
		&a.Low,
		&a.Close,
		&a.LastModTime,
		&a.LastModBy,
		&a.CreateTime,
		&a.CreateBy,
	)
	SkipSQLNoRowsError(&err)
	return err
}

// ReadExchs reads a full Exch structure of data from the database based
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
func ReadExchs(rows *sql.Rows, a *Exch) error {
	err := rows.Scan(
		&a.XID,
		&a.Dt,
		&a.Ticker,
		&a.Open,
		&a.High,
		&a.Low,
		&a.Close,
		&a.LastModTime,
		&a.LastModBy,
		&a.CreateTime,
		&a.CreateBy,
	)
	SkipSQLNoRowsError(&err)
	return err
}

// UpdateExch updates an existing Exch record in the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
//-----------------------------------------------------------------------------
func UpdateExch(ctx context.Context, a *Exch) error {
	sess, ok := session.GetSessionFromContext(ctx)
	if !ok {
		return ErrSessionRequired
	}
	fields := []interface{}{
		a.Dt,
		a.Ticker,
		a.Open,
		a.High,
		a.Low,
		a.Close,
		sess.UID,
		a.XID, // 49
	}
	var err error
	a.LastModBy, err = genericUpdate(ctx, Pdb.Prepstmt.UpdateExch, fields)
	return updateError(err, "Exch", *a)
}
