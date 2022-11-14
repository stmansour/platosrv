package db

import (
	"context"
	"database/sql"
	"platosrv/session"
	"time"
)

// ExchMonthly defines a date and a rent amount for a property. A ExchMonthly record
// is part of a group or list. The group is defined by the RSLID
// -----------------------------------------------------------------------------
type ExchMonthly struct {
	XMID        int64     // unique id for this record
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

// DeleteExchMonthly deletes the ExchMonthly with the specified id from the database
//
// INPUTS
// ctx - db context
// id - XMID of the record to read
//
// RETURNS
// Any errors encountered, or nil if no errors
// -----------------------------------------------------------------------------
func DeleteExchMonthly(ctx context.Context, id int64) error {
	return genericDelete(ctx, "ExchMonthly", Pdb.Prepstmt.DeleteExchMonthly, id)
}

// GetExchMonthly reads and returns a ExchMonthly structure
//
// INPUTS
// ctx - db context
// id - XMID of the record to read
//
// RETURNS
// ErrSessionRequired if the session is invalid
// nil if the session is valid
// -----------------------------------------------------------------------------
func GetExchMonthly(ctx context.Context, id int64) (ExchMonthly, error) {
	var a ExchMonthly
	if !ValidateSession(ctx) {
		return a, ErrSessionRequired
	}

	fields := []interface{}{id}
	stmt, row := getRowFromDB(ctx, Pdb.Prepstmt.GetExchMonthly, fields)
	if stmt != nil {
		defer stmt.Close()
	}
	return a, ReadExchMonthly(row, &a)
}

// InsertExchMonthly writes a new ExchMonthly record to the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
// -----------------------------------------------------------------------------
func InsertExchMonthly(ctx context.Context, a *ExchMonthly) (int64, error) {
	sess, ok := session.GetSessionFromContext(ctx)
	if !ok {
		return a.XMID, ErrSessionRequired
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
	a.CreateBy, a.LastModBy, a.XMID, err = genericInsert(ctx, "ExchMonthly", Pdb.Prepstmt.InsertExchMonthly, fields, a)
	return a.XMID, err
}

// ReadExchMonthly reads a full ExchMonthly structure of data from the database based
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
func ReadExchMonthly(row *sql.Row, a *ExchMonthly) error {
	err := row.Scan(
		&a.XMID,
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

// ReadExchMonthlys reads a full ExchMonthly structure of data from the database based
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
func ReadExchMonthlys(rows *sql.Rows, a *ExchMonthly) error {
	err := rows.Scan(
		&a.XMID,
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

// UpdateExchMonthly updates an existing ExchMonthly record in the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
// -----------------------------------------------------------------------------
func UpdateExchMonthly(ctx context.Context, a *ExchMonthly) error {
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
		a.XMID, // 49
	}
	var err error
	a.LastModBy, err = genericUpdate(ctx, Pdb.Prepstmt.UpdateExchMonthly, fields)
	return updateError(err, "ExchMonthly", *a)
}
