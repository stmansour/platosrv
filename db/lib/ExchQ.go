package db

import (
	"context"
	"database/sql"
	"platosrv/session"
	"time"
)

// ExchQuarterly defines a date and a rent amount for a property. A ExchQuarterly record
// is part of a group or list. The group is defined by the RSLID
// -----------------------------------------------------------------------------
type ExchQuarterly struct {
	XQID        int64     // unique id for this record
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

// DeleteExchQuarterly deletes the ExchQuarterly with the specified id from the database
//
// INPUTS
// ctx - db context
// id - XQID of the record to read
//
// RETURNS
// Any errors encountered, or nil if no errors
// -----------------------------------------------------------------------------
func DeleteExchQuarterly(ctx context.Context, id int64) error {
	return genericDelete(ctx, "ExchQuarterly", Pdb.Prepstmt.DeleteExchQuarterly, id)
}

// GetExchQuarterly reads and returns a ExchQuarterly structure
//
// INPUTS
// ctx - db context
// id - XQID of the record to read
//
// RETURNS
// ErrSessionRequired if the session is invalid
// nil if the session is valid
// -----------------------------------------------------------------------------
func GetExchQuarterly(ctx context.Context, id int64) (ExchQuarterly, error) {
	var a ExchQuarterly
	if !ValidateSession(ctx) {
		return a, ErrSessionRequired
	}

	fields := []interface{}{id}
	stmt, row := getRowFromDB(ctx, Pdb.Prepstmt.GetExchQuarterly, fields)
	if stmt != nil {
		defer stmt.Close()
	}
	return a, ReadExchQuarterly(row, &a)
}

// InsertExchQuarterly writes a new ExchQuarterly record to the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
// -----------------------------------------------------------------------------
func InsertExchQuarterly(ctx context.Context, a *ExchQuarterly) (int64, error) {
	sess, ok := session.GetSessionFromContext(ctx)
	if !ok {
		return a.XQID, ErrSessionRequired
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
	a.CreateBy, a.LastModBy, a.XQID, err = genericInsert(ctx, "ExchQuarterly", Pdb.Prepstmt.InsertExchQuarterly, fields, a)
	return a.XQID, err
}

// ReadExchQuarterly reads a full ExchQuarterly structure of data from the database based
// on the supplied Rows pointer.
//
// INPUTS
// row - db Row pointer
// a   - pointer to struct to fill
//
// # RETURNS
//
// ErrSessionRequired if the session is invalid
// nil if the session is valid
// -----------------------------------------------------------------------------
func ReadExchQuarterly(row *sql.Row, a *ExchQuarterly) error {
	err := row.Scan(
		&a.XQID,
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

// ReadExchQuarterlys reads a full ExchQuarterly structure of data from the database based
// on the supplied Rows pointer.
//
// INPUTS
// row - db Row pointer
// a   - pointer to struct to fill
//
// # RETURNS
//
// ErrSessionRequired if the session is invalid
// nil if the session is valid
// -----------------------------------------------------------------------------
func ReadExchQuarterlys(rows *sql.Rows, a *ExchQuarterly) error {
	err := rows.Scan(
		&a.XQID,
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

// UpdateExchQuarterly updates an existing ExchQuarterly record in the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
// -----------------------------------------------------------------------------
func UpdateExchQuarterly(ctx context.Context, a *ExchQuarterly) error {
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
		a.XQID, // 49
	}
	var err error
	a.LastModBy, err = genericUpdate(ctx, Pdb.Prepstmt.UpdateExchQuarterly, fields)
	return updateError(err, "ExchQuarterly", *a)
}
