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

// ItemGrid contains the data from Item that is targeted to the UI Grid that displays
// a list of Item structs

//-------------------------------------------------------------------
//                        **** SEARCH ****
//-------------------------------------------------------------------

// ItemSearchRequestJSON is a struct suitable for describing a webservice operation.
// It is the wire format data. It will be merged into another object where JSONTime values
// are converted to time.Time
type ItemSearchRequestJSON struct {
	Cmd         string            `json:"cmd"`         // get, save, delete
	Limit       int               `json:"limit"`       // max number to return
	Offset      int               `json:"offset"`      // solution set offset
	Selected    []int             `json:"selected"`    // selected rows
	SearchLogic string            `json:"searchLogic"` // OR | AND
	Search      []GenSearch       `json:"search"`      // what fields and what values
	Sort        []ColSort         `json:"sort"`        // sort criteria
	PubDt       util.JSONDateTime `json:"PubDt"`
}

// ItemSearchRequest a version of ItemSearchRequestJSON where JSONDateTime values
// are changed to time.Time values.
//---------------------------------------------------------------------------------
type ItemSearchRequest struct {
	Cmd         string
	Limit       int
	Offset      int
	Selected    []int
	SearchLogic string
	Search      []GenSearch
	Sort        []ColSort
	PubDt       time.Time
}

// ItemGrid is the structure of data for a Item we send to the UI
type ItemGrid struct {
	Recid       int64 `json:"recid"`
	IID         int64 // unique id
	Title       string
	Description string
	PubDt       util.JSONDateTime
	Link        string
	CreateTime  util.JSONDateTime
	CreateBy    int64
	LastModTime util.JSONDateTime
	LastModBy   int64
	//
	// RO db.RenewOptions // contains the list of RenewOptions and context
	// RS db.RentSteps    // contains the list of RentSteps and context
}

// which fields needs to be fetched for SQL query for Item grid
var itemFieldsMap = map[string][]string{
	"IID":         {"Item.IID"},
	"PubDt":       {"Item.PubDt"},
	"Title":       {"Item.Title"},
	"Description": {"Item.Description"},
	"Link":        {"Item.Link"},
	"CreateTime":  {"Item.CreateTime"},
	"CreateBy":    {"Item.CreateBy"},
	"LastModTime": {"Item.LastModTime"},
	"LastModBy":   {"Item.LastModBy"},
}

// which fields needs to be fetched for SQL query for Item grid
var itemQuerySelectFields = []string{
	"Item.IID",
	"Item.PubDt",
	"Item.Title",
	"Item.Description",
	"Item.Link",
	"Item.CreateTime",
	"Item.CreateBy",
	"Item.LastModTime",
	"Item.LastModBy",
}

// this is the list of fields to search for a string if the field name is blank
var searchDefaultFields = []string{
	"Title",
	"Description",
}

// SearchItemResponse is the response data for a Rental Agreement Search
type SearchItemResponse struct {
	Status  string     `json:"status"`
	Total   int64      `json:"total"`
	Records []ItemGrid `json:"records"`
}

//-------------------------------------------------------------------
//                         **** SAVE ****
//-------------------------------------------------------------------

// SaveItem is sent to save one of open time slots as a reservation
type SaveItem struct {
	Cmd    string   `json:"cmd"`
	Record ItemGrid `json:"record"`
}

//-------------------------------------------------------------------
//                         **** GET ****
//-------------------------------------------------------------------

// GetItem is the struct returned on a request for a reservation.
type GetItem struct {
	Status string   `json:"status"`
	Record ItemGrid `json:"record"`
}

//-----------------------------------------------------------------------------
//#############################################################################
//-----------------------------------------------------------------------------

// SvcHandlerItem formats a complete data record for an Item for use
// with the w2ui Form
//
// The server command can be:
//      get
//      save
//      delete
//------------------------------------------------------------------------------
func SvcHandlerItem(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	util.Console("Entered SvcHandlerItem, d.ID = %d\n", d.ID)

	switch d.wsSearchReq.Cmd {
	case "get":
		if d.ID <= 0 && d.wsSearchReq.Limit > 0 {
			SvcSearchItem(w, r, d) // it is a query for the grid.
		} else {
			if d.ID < 0 {
				SvcErrorReturn(w, fmt.Errorf("field IID is required but was not specified"))
				return
			}
			getItem(w, r, d)
		}
	case "save":
		saveItem(w, r, d)
	case "delete":
		deleteItem(w, r, d)
	default:
		err := fmt.Errorf("unhandled command: %s", d.wsSearchReq.Cmd)
		SvcErrorReturn(w, err)
		return
	}
}

// SvcSearchItem generates a report of all Item records matching the
// search criteria.
//
//	@URL /v1/Item/
//
//-----------------------------------------------------------------------------
func SvcSearchItem(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "SvcSearchItem"
	util.Console("Entered %s\n", funcname)
	var g SearchItemResponse
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

	var ItemSRD ItemSearchRequestJSON
	err = json.Unmarshal([]byte(d.data), &ItemSRD)
	if err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		SvcErrorReturn(w, e)
		return
	}

	var srd ItemSearchRequest
	if err = util.MigrateStructVals(&ItemSRD, &srd); err != nil {
		e := fmt.Errorf("%s: Error with MigrateStructVals:  %s", funcname, err.Error())
		SvcErrorReturn(w, e)
		return
	}

	whr := ""

	order := `Item.PubDt,Item.Title ASC` // default ORDER

	// get where clause and order clause for sql query
	// util.Console("len(d.wsSearchReq.Search) = %d\n", len(d.wsSearchReq.Search))
	HandleBlankSearchField(d, searchDefaultFields)
	// util.Console("AFTER HandleBlankSearchField:  len(d.wsSearchReq.Search) = %d\n", len(d.wsSearchReq.Search))

	//------------------------------------------------------
	// use MyQueue if present, otherwise use generic...
	//------------------------------------------------------
	whereClause, orderClause = GetSearchAndSortSQL(d, itemFieldsMap)
	if len(whereClause) > 0 {
		whr += "WHERE " + whereClause
	}

	util.Console("len(d.wsSearchReq.Search) = %d\n", len(d.wsSearchReq.Search))

	//------------------------------------------------------------------------
	// GetSearchAndSortSQL is busted for the search string we want here...
	// It needs to be fixed. Until then, here is some work-around code.
	//------------------------------------------------------------------------
	if len(d.wsSearchReq.Search) > 0 {
		s := d.wsSearchReq.Search[0].Value
		whr = "WHERE (Title LIKE \"%" + s + "%\" OR Description LIKE \"%" + s + "%\")"
	}
	if srd.PubDt.Year() >= db.MINYEAR {
		sTmp := "WHERE "
		if len(whr) > 0 {
			sTmp = " AND "
		}
		whr += sTmp + fmt.Sprintf("YEAR(PubDt)=%d AND MONTH(PubDt)=%d AND DAYOFMONTH(PubDt)=%d",
			srd.PubDt.Year(), int(srd.PubDt.Month()), srd.PubDt.Day())
	}

	if len(orderClause) > 0 {
		order = orderClause
	}
	util.Console("whr = %s\n", whr)
	util.Console("order = %s\n", order)

	query := `
	SELECT {{.SelectClause}}
	FROM Item {{.WhereClause}}
	ORDER BY {{.OrderClause}}`

	qc := db.QueryClause{
		"SelectClause": strings.Join(itemQuerySelectFields, ","),
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
	qc["LimitClause"] = strconv.Itoa(d.wsSearchReq.Limit)
	qc["OffsetClause"] = strconv.Itoa(d.wsSearchReq.Offset)

	// get formatted query with substitution of select, where, order clause
	qry := db.RenderSQLQuery(queryWithLimit, qc)
	util.Console("SvcSearchItem: db query = %s\n", qry)

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
		q, err := ItemRowScan(rows)
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

// ItemRowScan scans a result from sql row and dump it in a
// ItemGrid struct
//
// RETURNS
//  Item
//-----------------------------------------------------------------------------
func ItemRowScan(rows *sql.Rows) (ItemGrid, error) {
	var q ItemGrid
	err := rows.Scan(
		&q.IID,
		&q.PubDt,
		&q.Title,
		&q.Description,
		&q.Link,
		&q.CreateTime,
		&q.CreateBy,
		&q.LastModTime,
		&q.LastModBy,
	)
	return q, err
}

// deleteItem deletes a payment type from the database
// wsdoc {
//  @Title  Delete Item
//	@URL /v1/Item/IID
//  @Method  POST
//	@Synopsis Delete a Payment Type
//  @Desc  This service deletes a Item.
//	@Input WebGridDelete
//  @Response SvcStatusResponse
// wsdoc }
//-----------------------------------------------------------------------------
func deleteItem(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "deleteItem"
	util.Console("Entered %s\n", funcname)
	util.Console("record data = %s\n", d.data)
	SvcWriteSuccessResponse(w)
}

// SaveItem returns the requested Item
// wsdoc {
//  @Title  Save Item
//	@URL /v1/Item/IID
//  @Method  GET
//	@Synopsis Update the information on a Item with the supplied data, create if necessary.
//  @Description  This service creates a Item if IID == 0 or updates a Item if IID > 0 with
//  @Description  the information supplied. All fields must be supplied.
//	@Input SaveItem
//  @Response SvcStatusResponse
// wsdoc }
//-----------------------------------------------------------------------------
func saveItem(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "saveItem"
	util.Console("Entered %s\n", funcname)
	util.Console("record data = %s\n", d.data)
	util.Console("IID = %d\n", d.ID)

	var foo SaveItem

	data := []byte(d.data)
	err := json.Unmarshal(data, &foo)

	if err != nil {
		e := fmt.Errorf("%s: Error with json.Unmarshal:  %s", funcname, err.Error())
		SvcErrorReturn(w, e)
		return
	}

	// util.Console("read foo.  foo.Record.IID = %d, foo.Record.Name = %s\n", foo.Record.IID, foo.Record.Name)
	var p db.Item
	if err = util.MigrateStructVals(&foo.Record, &p); err != nil {
		e := fmt.Errorf("%s: Error with MigrateStructVals:  %s", funcname, err.Error())
		SvcErrorReturn(w, e)
		return
	}
	if p.IID < 1 {
		if _, err = db.InsertItem(r.Context(), &p); err != nil {
			e := fmt.Errorf("%s: Error with db.CreateItem:  %s", funcname, err.Error())
			SvcErrorReturn(w, e)
			return
		}
	} else {
		if err = db.UpdateItem(r.Context(), &p); err != nil {
			e := fmt.Errorf("%s: Error with db.UpdateItem:  %s", funcname, err.Error())
			SvcErrorReturn(w, e)
			return
		}
	}
	// util.Console("UpdateItem completed successfully\n")
	SvcWriteSuccessResponseWithID(w, p.IID)
}

// ItemUpdate updates the supplied Item in the database with the supplied
// info. It only allows certain fields to be updated.
//-----------------------------------------------------------------------------
func ItemUpdate(p *ItemGrid, d *ServiceData) error {
	util.Console("entered ItemUpdate\n")
	return nil
}

// GetItem returns the requested Item
// wsdoc {
//  @Title  Get Item
//	@URL /v1/Item/:IID
//  @Method  GET
//	@Synopsis Get information on a Item
//  @Description  Return all fields for Item :IID
//	@Input WebGridSearchRequest
//  @Response GetItem
// wsdoc }
//-----------------------------------------------------------------------------
func getItem(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "getItem"
	util.Console("entered %s\n", funcname)
	var g GetItem
	a, err := db.GetItem(r.Context(), d.ID)
	if err != nil {
		SvcErrorReturn(w, err)
		return
	}
	if a.IID == d.ID {
		var gg ItemGrid
		util.MigrateStructVals(&a, &gg)
		gg.Recid = gg.IID
		g.Record = gg
	} else {
		err = fmt.Errorf("could not find Item with IID = %d", d.ID)
		SvcErrorReturn(w, err)
		return
	}
	g.Status = "success"
	SvcWriteResponse(&g, w)
}
