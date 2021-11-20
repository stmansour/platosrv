package ws

import (
	"fmt"
	"net/http"
	db "platosrv/db/lib"
	util "platosrv/util/lib"
)

//-------------------------------------------------------------------
//                        **** SEARCH ****
//-------------------------------------------------------------------

//-------------------------------------------------------------------
//                         **** GET ****
//-------------------------------------------------------------------

// Dashboard contains the data fields for the dashboard
type Dashboard struct {
	ExchCount int64
	ItemCount int64
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

	// var sess *session.Session
	// var ok bool
	// if sess, ok = session.GetSessionFromContext(r.Context()); !ok {
	// 	SvcErrorReturn(w, db.ErrSessionRequired)
	// 	return
	// }

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

	g.Status = "success"
	SvcWriteResponse(&g, w)
}
