package db

import (
	"context"
	"database/sql"
	"platosrv/session"
	"time"
)

// ExchDaily defines a date and a rent amount for a property. A ExchDaily record
// is part of a group or list. The group is defined by the RSLID
// -----------------------------------------------------------------------------
type ExchDaily struct {
	XDID        int64     // unique id for this record
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

// DeleteExchDaily deletes the ExchDaily with the specified id from the database
//
// INPUTS
// ctx - db context
// id - XDID of the record to read
//
// RETURNS
// Any errors encountered, or nil if no errors
// -----------------------------------------------------------------------------
func DeleteExchDaily(ctx context.Context, id int64) error {
	return genericDelete(ctx, "ExchDaily", Pdb.Prepstmt.DeleteExchDaily, id)
}

// GetExchDaily reads and returns a ExchDaily structure
//
// INPUTS
// ctx - db context
// id - XDID of the record to read
//
// RETURNS
// ErrSessionRequired if the session is invalid
// nil if the session is valid
// -----------------------------------------------------------------------------
func GetExchDaily(ctx context.Context, id int64) (ExchDaily, error) {
	var a ExchDaily
	if !ValidateSession(ctx) {
		return a, ErrSessionRequired
	}

	fields := []interface{}{id}
	stmt, row := getRowFromDB(ctx, Pdb.Prepstmt.GetExchDaily, fields)
	if stmt != nil {
		defer stmt.Close()
	}
	return a, ReadExchDaily(row, &a)
}

// InsertExchDaily writes a new ExchDaily record to the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
// -----------------------------------------------------------------------------
func InsertExchDaily(ctx context.Context, a *ExchDaily) (int64, error) {
	sess, ok := session.GetSessionFromContext(ctx)
	if !ok {
		return a.XDID, ErrSessionRequired
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
	a.CreateBy, a.LastModBy, a.XDID, err = genericInsert(ctx, "ExchDaily", Pdb.Prepstmt.InsertExchDaily, fields, a)
	return a.XDID, err
}

// ReadExchDaily reads a full ExchDaily structure of data from the database based
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
func ReadExchDaily(row *sql.Row, a *ExchDaily) error {
	err := row.Scan(
		&a.XDID,
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

// ReadExchDailys reads a full ExchDaily structure of data from the database based
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
func ReadExchDailys(rows *sql.Rows, a *ExchDaily) error {
	err := rows.Scan(
		&a.XDID,
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

// UpdateExchDaily updates an existing ExchDaily record in the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// any error encountered or nil if no error
// -----------------------------------------------------------------------------
func UpdateExchDaily(ctx context.Context, a *ExchDaily) error {
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
		a.XDID, // 49
	}
	var err error
	a.LastModBy, err = genericUpdate(ctx, Pdb.Prepstmt.UpdateExchDaily, fields)
	return updateError(err, "ExchDaily", *a)
}
