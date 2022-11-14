// exchupdate is used to add the daily, weekly, monthly, and quarterly
// tables for the Exch data.
// -----------------------------------------------------------------------
package main

import (
	"context"
	"fmt"
	db "platosrv/db/lib"
	util "platosrv/util/lib"
	"strings"
	"time"
)

// createExchDaily builds  the table with the daily average
func createExchDaily(ctx context.Context) {
	var err error
	var errors, totErrors int64
	var warnings, totWarnings int64
	util.Console("Create/update ExchDaily\n")

	//----------------------------------------------
	// Iterate through the mappings we store...
	//----------------------------------------------
	totWarnings = 0
	totErrors = 0
	for k, v := range Tickers {
		if v > 0 {
			util.Console("\nProcessing %s\n", k)
			if errors, warnings, err = dailyExch(k); err != nil {
				util.Console("Error in scanExch: %s\n", err)
			}
			totErrors += errors
			totWarnings += warnings
		}
	}
	util.Console("\nFinished\nTotal Errors: %d\n", totErrors)
	if App.Warnings {
		util.Console("Total Warnings: %d\n", totWarnings)
	}
}

// dailyExch reads the minute based data and saves a daily average in the ExchDaily table.
//
// INPUTS
//
//	t = ticker
//
// RETURNS
//
//	error count
//	warning count
//	nil - if no errors, otherwise the error that stopped the processing...
//
// -----------------------------------------------------------------------------------------
func dailyExch(t string) (int64, int64, error) {
	var a db.Exch      // record to read from db
	var x db.ExchDaily // where we keep updated totals
	var count, n int64
	var errors int64
	var warnings int64
	var qry string

	util.Console("dailyExch:  processing t = %s\n", t)
	qry = fmt.Sprintf("SELECT %s FROM Exch WHERE Ticker = \"%s\" AND Dt > \"2010-01-01\" ORDER BY Dt ASC", db.Pdb.DBFields["Exch"], t)
	rows, err := db.Pdb.DB.Query(qry)
	if err != nil {
		util.Console("error win db.Pdb.DB.Query: %s", err)
		errors++
		return errors, warnings, err
	}
	defer rows.Close()

	count = 0
	n = 0
	x.Ticker = t
	x.Dt = time.Date(2010, time.January, 1, 0, 0, 0, 0, time.UTC) // a valid date, prior to the first db record date. (Exch data begins in 2011)
	for rows.Next() {
		if err := db.ReadExchs(rows, &a); err != nil {
			errors++
			return errors, warnings, err
		}

		//----------------------------------------------------------------
		// If the day has changed then create the average for the day...
		//----------------------------------------------------------------
		if a.Dt.Day() != x.Dt.Day() {
			if n > 0 { // skip if we haven't collected anything yet
				x.Open /= float64(n)
				x.Close /= float64(n)
				x.High /= float64(n)
				x.Low /= float64(n)

				// write or update this record
				if err = writeUpdateExchDaily(&x, t); err != nil {
					errors++
					return errors, warnings, err
				}

				// initialize for next record
				x.Open = 0.0
				x.Close = 0.0
				x.High = 0.0
				x.Low = 0.0
				n = 0
			}
			x.Dt = time.Date(a.Dt.Year(), a.Dt.Month(), a.Dt.Day(), 0, 0, 0, 0, time.UTC) // set this correctly whether we've collected anything or not
		}

		//---------------------------------------------------
		// Still on the same day, add this to the totals...
		//---------------------------------------------------
		n++
		x.Open += a.Open
		x.Close += a.Close
		x.High += a.High
		x.Low += a.Low

		count++
	}

	if err := rows.Err(); err != nil {
		errors++
		return errors, warnings, err
	}

	return errors, warnings, nil
}

// Write the specified record. If it exists, update it with this information.
//
// INPUTS
// x = pointer to struct to write
// t = ticker
// ------------------------------------------------------------------------------
func writeUpdateExchDaily(x *db.ExchDaily, t string) error {
	var err error
	// Try to insert...
	if _, err = db.InsertExchDaily(App.ctx, x); err == nil {
		return nil // if that worked, we're done
	}

	//-----------------------------------------------------------
	// If the error was Duplicate Entry, then we just update...
	//-----------------------------------------------------------
	errs := err.Error()
	if !strings.Contains(errs, "Duplicate entry") {
		return err // if the error was something else, then we're done, just return the error
	}

	qry := fmt.Sprintf(
		`SELECT %s FROM ExchDaily WHERE Ticker = "%s" AND YEAR(Dt)=%d AND MONTH(Dt)=%d AND DAYOFMONTH(Dt)=%d`,
		db.Pdb.DBFields["ExchDaily"], t, x.Dt.Year(), x.Dt.Month(), x.Dt.Day())
	row := db.Pdb.DB.QueryRow(qry)
	var x1 db.ExchDaily
	if err = db.ReadExchDaily(row, &x1); err != nil {
		return err
	}

	x1.Open = x.Open
	x1.Close = x.Close
	x1.High = x.High
	x1.Low = x.Low
	return db.UpdateExchDaily(App.ctx, &x1)
}
