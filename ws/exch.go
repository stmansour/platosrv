package ws

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	db "platosrv/db/lib"
	"platosrv/session"
	util "platosrv/util/lib"
	"strconv"
	"strings"
	"time"
)

// ExchGrid contains the data from Exch that is targeted to the UI Grid that displays
// a list of Exch structs

//-------------------------------------------------------------------
//                        **** SEARCH ****
//-------------------------------------------------------------------

// ExchSearchRequestJSON is a struct suitable for describing a webservice operation.
// It is the wire format data. It will be merged into another object where JSONTime values
// are converted to time.Time
type ExchSearchRequestJSON struct {
	Cmd         string            `json:"cmd"`         // get, save, delete
	Limit       int               `json:"limit"`       // max number to return
	Offset      int               `json:"offset"`      // solution set offset
	Selected    []int             `json:"selected"`    // selected rows
	SearchLogic string            `json:"searchLogic"` // OR | AND
	Search      []GenSearch       `json:"search"`      // what fields and what values
	Sort        []ColSort         `json:"sort"`        // sort criteria
	Tickers     []string          `json:"Tickers"`     // which tickers are requested
	Dt          util.JSONDateTime `json:"Dt"`
}

// ExchSearchRequest a version of ExchSearchRequestJSON where JSONDateTime values
// are changed to time.Time values.
//---------------------------------------------------------------------------------
type ExchSearchRequest struct {
	Cmd         string
	Limit       int
	Offset      int
	Selected    []int
	SearchLogic string
	Search      []GenSearch
	Sort        []ColSort
	Tickers     []string
	Dt          time.Time
}

// ExchGrid is the structure of data for a Exch we send to the UI
type ExchGrid struct {
	Recid       int64 `json:"recid"`
	XID         int64 // unique id
	Dt          util.JSONDateTime
	Ticker      string  // the two currencies involved in this exchange rate
	Open        float64 // Opening value for this minute
	High        float64 // High value during this minute
	Low         float64 // Low value during this minute
	Close       float64 // Closing value for this minute
	CreateTime  util.JSONDateTime
	CreateBy    int64
	LastModTime util.JSONDateTime
	LastModBy   int64
	//
	// RO db.RenewOptions // contains the list of RenewOptions and context
	// RS db.RentSteps    // contains the list of RentSteps and context
}

// which fields needs to be fetched for SQL query for Exch grid
var propFieldsMap = map[string][]string{
	"XID":         {"Exch.XID"},
	"Dt":          {"Exch.Dt"},
	"Ticker":      {"Exch.Ticker"},
	"Open":        {"Exch.Open"},
	"High":        {"Exch.High"},
	"Low":         {"Exch.Low"},
	"Close":       {"Exch.Close"},
	"CreateTime":  {"Exch.CreateTime"},
	"CreateBy":    {"Exch.CreateBy"},
	"LastModTime": {"Exch.LastModTime"},
	"LastModBy":   {"Exch.LastModBy"},
}

// which fields needs to be fetched for SQL query for Exch grid
var propQuerySelectFields = []string{
	"Exch.XID",
	"Exch.Dt",
	"Exch.Ticker",
	"Exch.Open",
	"Exch.High",
	"Exch.Low",
	"Exch.Close",
	"Exch.CreateTime",
	"Exch.CreateBy",
	"Exch.LastModTime",
	"Exch.LastModBy",
}

// this is the list of fields to search for a string if the field name is blank
var propDefaultFields = []string{
	"Name",
	"City",
	"State",
	"PostalCode",
}

// SearchExchResponse is the response data for a Rental Agreement Search
type SearchExchResponse struct {
	Status  string     `json:"status"`
	Total   int64      `json:"total"`
	Records []ExchGrid `json:"records"`
}

//-------------------------------------------------------------------
//                         **** SAVE ****
//-------------------------------------------------------------------

// SaveExch is sent to save one of open time slots as a reservation
type SaveExch struct {
	Cmd    string   `json:"cmd"`
	Record ExchGrid `json:"record"`
}

//-------------------------------------------------------------------
//                         **** GET ****
//-------------------------------------------------------------------

// GetExch is the struct returned on a request for a reservation.
type GetExch struct {
	Status string   `json:"status"`
	Record ExchGrid `json:"record"`
}

//-----------------------------------------------------------------------------
//#############################################################################
//-----------------------------------------------------------------------------

// SvcHandlerExch formats a complete data record for an Exch for use
// with the w2ui Form
//
// The server command can be:
//      get
//      save
//      delete
//------------------------------------------------------------------------------
func SvcHandlerExch(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	util.Console("Entered SvcHandlerExch, d.ID = %d\n", d.ID)

	switch d.wsSearchReq.Cmd {
	case "get":
		if d.ID <= 0 && d.wsSearchReq.Limit > 0 {
			SvcSearchExch(w, r, d) // it is a query for the grid.
		} else {
			if d.ID < 0 {
				SvcErrorReturn(w, fmt.Errorf("XID is required but was not specified"))
				return
			}
			getExch(w, r, d)
		}
		break
	case "save":
		saveExch(w, r, d)
		break
	case "delete":
		deleteExch(w, r, d)
	default:
		err := fmt.Errorf("Unhandled command: %s", d.wsSearchReq.Cmd)
		SvcErrorReturn(w, err)
		return
	}
}

// whereBuilder builds the where clause as needed for exchange Ticker data
// as propvided.
//
// INPUTS
//    exchSRD - data read from request
//    whr     - current WHERE clause
//
// RETURNS
//    updated whr clause string
//    any error encountered.
//------------------------------------------------------------------------------
func whereBuilder(esr *ExchSearchRequest, whr string) string {
	var qtickers []string
	var wAnd = " AND "
	var dtAnd = " AND "
	if len(esr.Tickers) > 0 {
		if len(whr) == 0 {
			whr = "WHERE "
			wAnd = ""
		}
		for i := 0; i < len(esr.Tickers); i++ {
			qtickers = append(qtickers, fmt.Sprintf("%q", esr.Tickers[i]))
		}
		whr += wAnd + "Ticker IN (" + strings.Join(qtickers, ",") + ")"
	}
	if len(whr) == 0 {
		whr = "WHERE "
		dtAnd = ""
	}
	whr += dtAnd + fmt.Sprintf("YEAR(Dt)=%d AND MONTH(Dt)=%d AND DAYOFMONTH(Dt)=%d",
		esr.Dt.Year(), int(esr.Dt.Month()), esr.Dt.Day())
	return whr
}

// SvcSearchExch generates a report of all Exch records matching the
// search criteria.
//
//	@URL /v1/Exch/
//
//-----------------------------------------------------------------------------
func SvcSearchExch(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "SvcSearchExch"
	util.Console("Entered %s\n", funcname)
	var g SearchExchResponse
	var err error
	var whereClause, orderClause string

	sess, ok := session.GetSessionFromContext(r.Context())
	if !ok {
		SvcErrorReturn(w, db.ErrSessionRequired)
		return
	}
	if sess.UID < 1 {
		SvcErrorReturn(w, db.ErrSessionRequired)
		return
	}

	var exchSRD ExchSearchRequestJSON
	err = json.Unmarshal([]byte(d.data), &exchSRD)
	if err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		SvcErrorReturn(w, e)
		return
	}

	var srd ExchSearchRequest
	if err = util.MigrateStructVals(&exchSRD, &srd); err != nil {
		e := fmt.Errorf("%s: Error with MigrateStructVals:  %s", funcname, err.Error())
		SvcErrorReturn(w, e)
		return
	}

	// Default: if date was not specified just use yesterday's exchange rates
	//--------------------------------------------------------------------------
	if srd.Dt.Year() < db.MINYEAR {
		srd.Dt = time.Now().Add(-24 * time.Hour)
	}

	whr := ""
	order := `Exch.Dt ASC` // default ORDER

	// get where clause and order clause for sql query
	// util.Console("len(d.wsSearchReq.Search) = %d\n", len(d.wsSearchReq.Search))
	HandleBlankSearchField(d, propDefaultFields)
	// util.Console("AFTER HandleBlankSearchField:  len(d.wsSearchReq.Search) = %d\n", len(d.wsSearchReq.Search))

	//------------------------------------------------------
	// use MyQueue if present, otherwise use generic...
	//------------------------------------------------------
	whereClause, orderClause = GetSearchAndSortSQL(d, propFieldsMap)
	if len(whereClause) > 0 {
		whr += "WHERE " + whereClause
	}
	whr = whereBuilder(&srd, whr)

	if len(orderClause) > 0 {
		order = orderClause
	}
	util.Console("whr = %s\n", whr)
	util.Console("order = %s\n", order)

	query := `
	SELECT {{.SelectClause}}
	FROM Exch {{.WhereClause}}
	ORDER BY {{.OrderClause}}`

	qc := db.QueryClause{
		"SelectClause": strings.Join(propQuerySelectFields, ","),
		"WhereClause":  whr,
		"OrderClause":  order,
	}

	countQuery := db.RenderSQLQuery(query, qc)

	util.Console("countQuery = %s\n", countQuery)
	g.Total, err = db.GetQueryCount(countQuery)
	if err != nil {
		SvcErrorReturn(w, err)
		return
	}
	// util.Console("g.Total = %d\n", g.Total)

	// FETCH the records WITH LIMIT AND OFFSET
	// limit the records to fetch from server, page by page
	limitAndOffsetClause := `
		LIMIT {{.LimitClause}}
		OFFSET {{.OffsetClause}};`

	// build query with limit and offset clause
	// if query ends with ';' then remove it
	queryWithLimit := query + limitAndOffsetClause

	// Add limit and offset value
	if d.wsSearchReq.Limit <= 0 || d.wsSearchReq.Limit > 500 {
		d.wsSearchReq.Limit = 500
	}
	qc["LimitClause"] = strconv.Itoa(d.wsSearchReq.Limit)
	qc["OffsetClause"] = strconv.Itoa(d.wsSearchReq.Offset)

	// get formatted query with substitution of select, where, order clause
	qry := db.RenderSQLQuery(queryWithLimit, qc)
	util.Console("SvcSearchExch: db query = %s\n", qry)

	// execute the query
	rows, err := db.Pdb.DB.Query(qry)
	if err != nil {
		SvcErrorReturn(w, err)
		return
	}
	defer rows.Close()

	i := int64(d.wsSearchReq.Offset)
	count := 0
	for rows.Next() {
		q, err := ExchRowScan(rows)
		if err != nil {
			SvcErrorReturn(w, err)
			return
		}
		q.Recid = i

		g.Records = append(g.Records, q)
		count++ // update the count only after adding the record
		if count >= d.wsSearchReq.Limit {
			break // if we've added the max number requested, then exit
		}
		i++
	}

	err = rows.Err()
	if err != nil {
		SvcErrorReturn(w, err)
		return
	}

	g.Status = "success"
	SvcWriteResponse(&g, w)
}

// ExchRowScan scans a result from sql row and dump it in a
// ExchGrid struct
//
// RETURNS
//  Exch
//-----------------------------------------------------------------------------
func ExchRowScan(rows *sql.Rows) (ExchGrid, error) {
	var q ExchGrid
	err := rows.Scan(
		&q.XID,
		&q.Dt,
		&q.Ticker,
		&q.Open,
		&q.High,
		&q.Low,
		&q.Close,
		&q.CreateTime,
		&q.CreateBy,
		&q.LastModTime,
		&q.LastModBy,
	)
	return q, err
}

// deleteExch deletes a payment type from the database
// wsdoc {
//  @Title  Delete Exch
//	@URL /v1/Exch/XID
//  @Method  POST
//	@Synopsis Delete a Payment Type
//  @Desc  This service deletes a Exch.
//	@Input WebGridDelete
//  @Response SvcStatusResponse
// wsdoc }
//-----------------------------------------------------------------------------
func deleteExch(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "deleteExch"
	util.Console("Entered %s\n", funcname)
	util.Console("record data = %s\n", d.data)
	SvcWriteSuccessResponse(w)
}

// SaveExch returns the requested Exch
// wsdoc {
//  @Title  Save Exch
//	@URL /v1/Exch/XID
//  @Method  GET
//	@Synopsis Update the information on a Exch with the supplied data, create if necessary.
//  @Description  This service creates a Exch if XID == 0 or updates a Exch if XID > 0 with
//  @Description  the information supplied. All fields must be supplied.
//	@Input SaveExch
//  @Response SvcStatusResponse
// wsdoc }
//-----------------------------------------------------------------------------
func saveExch(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "saveExch"
	util.Console("Entered %s\n", funcname)
	util.Console("record data = %s\n", d.data)
	util.Console("XID = %d\n", d.ID)

	var foo SaveExch

	data := []byte(d.data)
	err := json.Unmarshal(data, &foo)

	if err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		SvcErrorReturn(w, e)
		return
	}

	// util.Console("read foo.  foo.Record.XID = %d, foo.Record.Name = %s\n", foo.Record.XID, foo.Record.Name)
	var p db.Exch
	if err = util.MigrateStructVals(&foo.Record, &p); err != nil {
		e := fmt.Errorf("%s: Error with MigrateStructVals:  %s", funcname, err.Error())
		SvcErrorReturn(w, e)
		return
	}
	if p.XID < 1 {
		if _, err = db.InsertExch(r.Context(), &p); err != nil {
			e := fmt.Errorf("%s: Error with db.CreateExch:  %s", funcname, err.Error())
			SvcErrorReturn(w, e)
			return
		}
	} else {
		if err = db.UpdateExch(r.Context(), &p); err != nil {
			e := fmt.Errorf("%s: Error with db.UpdateExch:  %s", funcname, err.Error())
			SvcErrorReturn(w, e)
			return
		}
	}
	// util.Console("UpdateExch completed successfully\n")
	SvcWriteSuccessResponseWithID(w, p.XID)
}

// ExchUpdate updates the supplied Exch in the database with the supplied
// info. It only allows certain fields to be updated.
//-----------------------------------------------------------------------------
func ExchUpdate(p *ExchGrid, d *ServiceData) error {
	util.Console("entered ExchUpdate\n")
	return nil
}

// GetExch returns the requested Exch
// wsdoc {
//  @Title  Get Exch
//	@URL /v1/Exch/:XID
//  @Method  GET
//	@Synopsis Get information on a Exch
//  @Description  Return all fields for Exch :XID
//	@Input WebGridSearchRequest
//  @Response GetExch
// wsdoc }
//-----------------------------------------------------------------------------
func getExch(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "getExch"
	util.Console("entered %s\n", funcname)
	var g GetExch
	a, err := db.GetExch(r.Context(), d.ID)
	if err != nil {
		SvcErrorReturn(w, err)
		return
	}
	if a.XID == d.ID {
		var gg ExchGrid
		util.MigrateStructVals(&a, &gg)
		gg.Recid = gg.XID
		g.Record = gg
	} else {
		err = fmt.Errorf("Could not find Exch with XID = %d", d.ID)
		SvcErrorReturn(w, err)
		return
	}
	g.Status = "success"
	SvcWriteResponse(&g, w)
}
