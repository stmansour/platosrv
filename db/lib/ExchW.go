package db

import (
	"context"
	"database/sql"
	"platosrv/session"
	"time"
)

// ExchWeekly defines a date and a rent amount for a property. A ExchWeekly record
// is part of a group or list. The group is defined by the RSLID
// -----------------------------------------------------------------------------
type ExchWeekly struct {
	XWID        int64     // unique id for this record
	Dt          time.Time // point in time when these values are valid
	ISOWeek     int64     // Week number for the year in Dt.  Dt day and month are the first day of this isoweek.
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

// DeleteExchWeekly deletes the ExchWeekly with the specified id from the database
//
// INPUTS
// ctx - db context
// id - XWID of the record to read
//
// RETURNS
// Any errors encountered, or nil if no errors
// -----------------------------------------------------------------------------
func DeleteExchWeekly(ctx context.Context, id int64) error {
	return genericDelete(ctx, "ExchWeekly", Pdb.Prepstmt.DeleteExchWeekly, id)
}

// GetExchWeekly reads and returns a ExchWeekly structure
//
// INPUTS
// ctx - db context
// id - XWID of the record to read
//
// RETURNS
// ErrSessionRequired if the session is invalid
// nil if the session is valid
// -----------------------------------------------------------------------------
func GetExchWeekly(ctx context.Context, id int64) (ExchWeekly, error) {
	var a ExchWeekly
	if !ValidateSession(ctx) {
		return a, ErrSessionRequired
	}

	fields := []interface{}{id}
	stmt, row := getRowFromDB(ctx, Pdb.Prepstmt.GetExchWeekly, fields)
	if stmt != nil {
		defer stmt.Close()
	}
	return a, ReadExchWeekly(row, &a)
}

// InsertExchWeekly writes a new ExchWeekly record to the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
// -----------------------------------------------------------------------------
func InsertExchWeekly(ctx context.Context, a *ExchWeekly) (int64, error) {
	sess, ok := session.GetSessionFromContext(ctx)
	_, w := a.Dt.ISOWeek()
	a.ISOWeek = int64(w)
	if !ok {
		return a.XWID, ErrSessionRequired
	}
	fields := []interface{}{
		a.Dt,
		a.ISOWeek,
		a.Ticker,
		a.Open,
		a.High,
		a.Low,
		a.Close,
		sess.UID,
		sess.UID,
	}

	var err error
	a.CreateBy, a.LastModBy, a.XWID, err = genericInsert(ctx, "ExchWeekly", Pdb.Prepstmt.InsertExchWeekly, fields, a)
	return a.XWID, err
}

// ReadExchWeekly reads a full ExchWeekly structure of data from the database based
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
func ReadExchWeekly(row *sql.Row, a *ExchWeekly) error {
	err := row.Scan(
		&a.XWID,
		&a.Dt,
		&a.ISOWeek,
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

// ReadExchWeeklys reads a full ExchWeekly structure of data from the database based
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
func ReadExchWeeklys(rows *sql.Rows, a *ExchWeekly) error {
	err := rows.Scan(
		&a.XWID,
		&a.Dt,
		&a.ISOWeek,
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

// UpdateExchWeekly updates an existing ExchWeekly record in the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
// -----------------------------------------------------------------------------
func UpdateExchWeekly(ctx context.Context, a *ExchWeekly) error {
	sess, ok := session.GetSessionFromContext(ctx)
	_, w := a.Dt.ISOWeek()
	a.ISOWeek = int64(w)

	if !ok {
		return ErrSessionRequired
	}
	fields := []interface{}{
		a.Dt,
		a.ISOWeek,
		a.Ticker,
		a.Open,
		a.High,
		a.Low,
		a.Close,
		sess.UID,
		a.XWID, // 49
	}
	var err error
	a.LastModBy, err = genericUpdate(ctx, Pdb.Prepstmt.UpdateExchWeekly, fields)
	return updateError(err, "ExchWeekly", *a)
}
