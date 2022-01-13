package ws

import (
	"fmt"
	"net/http"
	db "platosrv/db/lib"
	util "platosrv/util/lib"
	"time"
)

//-------------------------------------------------------------------
//                        **** SEARCH ****
//-------------------------------------------------------------------

//-------------------------------------------------------------------
//                         **** GET ****
//-------------------------------------------------------------------

// TableSizeInfo holds the the table name and it's current size in MB
type TableSizeInfo struct {
	Name string
	Size int64 // size in MB
}

// Dashboard contains the data fields for the dashboard
type Dashboard struct {
	ExchCount int64
	ExchDtMin time.Time
	ExchDtMax time.Time
	ItemCount int64
	ItemDtMin time.Time
	ItemDtMax time.Time
	Tables    []TableSizeInfo
}

// GetDashboard is the struct used to xfer Dashboard data to a requester
type GetDashboard struct {
	Status  string    `json:"status"`
	Message string    `json:"message"`
	Record  Dashboard `json:"record"`
}

//-------------------------------------------------------------------
//                         **** SAVE ****
//-------------------------------------------------------------------

// SvcHandlerDashboard returns dashboard stats
// with the w2ui Form
//
// The server command can be:
//      get
//------------------------------------------------------------------------------
func SvcHandlerDashboard(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	util.Console("Entered SvcHandlerDashboard, d.ID = %d\n", d.ID)

	switch d.wsSearchReq.Cmd {
	case "get":
		getDashboard(w, r, d)
		break
	default:
		err := fmt.Errorf("Unhandled command: %s", d.wsSearchReq.Cmd)
		SvcErrorReturn(w, err)
		return
	}
}

// getDashboard returns the Dashboard stats
// wsdoc {
//  @Title  Get Dashboard
//	@URL /v1/Dashboard/:PRID
//  @Method  GET
//	@Synopsis Get Dashboard statistics
//  @Description  Return all fields for Dashboard :PRID
//	@Input WebGridSearchRequest
//  @Response GetDashboard
// wsdoc }
//-----------------------------------------------------------------------------
func getDashboard(w http.ResponseWriter, r *http.Request, d *ServiceData) {
	funcname := "getDashboard"
	util.Console("entered %s\n", funcname)

	var g GetDashboard

	row := db.Pdb.DB.QueryRow(`SELECT COUNT(*) as ExchCount FROM Exch`)
	if err := row.Scan(&g.Record.ExchCount); err != nil {
		SvcErrorReturn(w, err)
		return
	}
	row = db.Pdb.DB.QueryRow(`SELECT COUNT(*) as ItemCount FROM Item`)
	if err := row.Scan(&g.Record.ItemCount); err != nil {
		SvcErrorReturn(w, err)
		return
	}

	var rr TableSizeInfo
	q := `
		SELECT
			TABLE_NAME AS "Table",
			ROUND((DATA_LENGTH + INDEX_LENGTH)) AS "Size"
		FROM
			information_schema.TABLES
		WHERE
			TABLE_SCHEMA="plato"
		ORDER BY
			TABLE_NAME ASC;`
	rows, err := db.Pdb.DB.Query(q)
	
	if err != nil {
		SvcErrorReturn(w, err)
		return
	}
	for rows.Next() {
		err = rows.Scan(&rr.Name, &rr.Size)
		if err != nil {
			SvcErrorReturn(w, err)
			return
		}
		g.Record.Tables = append(g.Record.Tables, rr)
	}
	if err = rows.Err(); err != nil {
		SvcErrorReturn(w, err)
		return
	}

	row = db.Pdb.DB.QueryRow(`select MIN(Dt) AS "Min", MAX(Dt) AS "Max" FROM Exch;`)
	if err = row.Scan(&g.Record.ExchDtMin, &g.Record.ExchDtMax); err != nil {
		SvcErrorReturn(w, err)
		return
	}

	row = db.Pdb.DB.QueryRow(`select MIN(PubDt) AS "Min", MAX(PubDt) AS "Max" FROM Item;`)
	if err = row.Scan(&g.Record.ItemDtMin, &g.Record.ItemDtMax); err != nil {
		SvcErrorReturn(w, err)
		return
	}

	g.Status = "success"
	SvcWriteResponse(&g, w)
}
