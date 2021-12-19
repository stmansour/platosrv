package db

import (
	"context"
	"database/sql"
	"platosrv/session"
	"time"
)

// RSSFeed defines a date and a rent amount for a property. A RSSFeed record
// is part of a group or list. The group is defined by the RSLID
//-----------------------------------------------------------------------------
type RSSFeed struct {
	RSSID       int64 // unique id for this record
	URL         string
	FLAGS       int64
	LastModTime time.Time // when was the record last written
	LastModBy   int64     // id of user that did the modify
	CreateTime  time.Time // when was this record created
	CreateBy    int64     // id of user that created it
}

// DeleteRSSFeed deletes the RSSFeed with the specified id from the database
//
// INPUTS
// ctx - db context
// id - RSSID of the record to read
//
// RETURNS
// Any errors encountered, or nil if no errors
//-----------------------------------------------------------------------------
func DeleteRSSFeed(ctx context.Context, id int64) error {
	return genericDelete(ctx, "RSSFeed", Pdb.Prepstmt.DeleteRSSFeed, id)
}

// GetRSSFeed reads and returns a RSSFeed structure
//
// INPUTS
// ctx - db context
// id - RSSID of the record to read
//
// RETURNS
// ErrSessionRequired if the session is invalid
// nil if the session is valid
//-----------------------------------------------------------------------------
func GetRSSFeed(ctx context.Context, id int64) (RSSFeed, error) {
	var a RSSFeed
	if !ValidateSession(ctx) {
		return a, ErrSessionRequired
	}

	fields := []interface{}{id}
	stmt, row := getRowFromDB(ctx, Pdb.Prepstmt.GetRSSFeed, fields)
	if stmt != nil {
		defer stmt.Close()
	}
	return a, ReadRSSFeed(row, &a)
}

// InsertRSSFeed writes a new RSSFeed record to the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
//-----------------------------------------------------------------------------
func InsertRSSFeed(ctx context.Context, a *RSSFeed) (int64, error) {
	sess, ok := session.GetSessionFromContext(ctx)
	if !ok {
		return a.RSSID, ErrSessionRequired
	}
	fields := []interface{}{
		a.URL,
		a.FLAGS,
		sess.UID,
		sess.UID,
	}

	var err error
	a.CreateBy, a.LastModBy, a.RSSID, err = genericInsert(ctx, "RSSFeed", Pdb.Prepstmt.InsertRSSFeed, fields, a)
	return a.RSSID, err
}

// ReadRSSFeed reads a full RSSFeed structure of data from the database based
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
func ReadRSSFeed(row *sql.Row, a *RSSFeed) error {
	err := row.Scan(
		&a.RSSID,
		&a.URL,
		&a.FLAGS,
		&a.LastModTime,
		&a.LastModBy,
		&a.CreateTime,
		&a.CreateBy,
	)
	SkipSQLNoRowsError(&err)
	return err
}

// ReadRSSFeeds reads a full RSSFeed structure of data from the database based
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
func ReadRSSFeeds(rows *sql.Rows, a *RSSFeed) error {
	err := rows.Scan(
		&a.RSSID,
		&a.URL,
		&a.FLAGS,
		&a.LastModTime,
		&a.LastModBy,
		&a.CreateTime,
		&a.CreateBy,
	)
	SkipSQLNoRowsError(&err)
	return err
}

// UpdateRSSFeed updates an existing RSSFeed record in the database
//
// INPUTS
// ctx - db context
// a   - pointer to struct to fill
//
// RETURNS
// id of the record just inserted
// any error encountered or nil if no error
//-----------------------------------------------------------------------------
func UpdateRSSFeed(ctx context.Context, a *RSSFeed) error {
	sess, ok := session.GetSessionFromContext(ctx)
	if !ok {
		return ErrSessionRequired
	}
	fields := []interface{}{
		a.URL,
		a.FLAGS,
		sess.UID,
		a.RSSID, // 49
	}
	var err error
	a.LastModBy, err = genericUpdate(ctx, Pdb.Prepstmt.UpdateRSSFeed, fields)
	return updateError(err, "RSSFeed", *a)
}
