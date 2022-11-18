// exchupdate is used to add the daily, weekly, Weekly, and quarterly
// tables for the Exch data.
// -----------------------------------------------------------------------
package main

import (
	"context"
	"fmt"
	db "platosrv/db/lib"
	util "platosrv/util/lib"
	"sort"
	"strings"
	"time"
)

// createExchWeekly builds  the table with the daily average
func createExchWeekly(ctx context.Context) {
	var err error
	var errors, totErrors int64
	var warnings, totWarnings int64
	var aTickers []string

	util.Console("Create/update ExchWeekly\n")

	//----------------------------------------------
	// Iterate through the mappings we store...
	//----------------------------------------------
	totWarnings = 0
	totErrors = 0
	for k, v := range Tickers {
		if v > 0 {
			aTickers = append(aTickers, k)
		}
	}
	util.Console("aTickers before sort: %v\n", aTickers)
	sort.Strings(aTickers)
	util.Console("aTickers after sort: %v\n", aTickers)
	for i := 0; i < len(aTickers); i++ {
		k := aTickers[i]
		util.Console("\nProcessing %s\n", k)
		if errors, warnings, err = WeeklyExch(k); err != nil {
			util.Console("Error in scanExch: %s\n", err)
		}
		totErrors += errors
		totWarnings += warnings
	}
	util.Console("\nFinished\nTotal Errors: %d\n", totErrors)
	if App.Warnings {
		util.Console("Total Warnings: %d\n", totWarnings)
	}
}

// WeeklyExch reads the minute based data and saves a daily average in the ExchWeekly table.
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
func WeeklyExch(t string) (int64, int64, error) {
	var a db.Exch       // record to read from db
	var x db.ExchWeekly // where we keep updated totals
	var count, n int64
	var errors int64
	var warnings int64

	util.Console("WeeklyExch:  processing t = %s\n", t)
	var qry = fmt.Sprintf(
		`SELECT %s FROM Exch WHERE Ticker = "%s" AND YEAR(Dt)>2010 ORDER BY Dt ASC`,
		db.Pdb.DBFields["Exch"], t)
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
		if err = db.ReadExchs(rows, &a); err != nil {
			errors++
			return errors, warnings, err
		}

		//----------------------------------------------------------------
		// If the day has changed then create the average for the day...
		//----------------------------------------------------------------
		_, w := a.Dt.ISOWeek()

		if int64(w) != x.ISOWeek {

			// util.Console("\nWeek changed\n")
			// util.Console("a: %d\n", w)
			// util.Console("x: %d\n", x.ISOWeek)
			if n > 0 { // skip if we haven't collected anything yet
				x.Open /= float64(n)
				x.Close /= float64(n)
				x.High /= float64(n)
				x.Low /= float64(n)

				// write or update this record
				if err = writeUpdateExchWeekly(&x, t); err != nil {
					errors++
					return errors, warnings, err
				}

				// initialize for next record
				x.Open = 0.0
				x.Close = 0.0
				x.High = 0.0
				x.Low = 0.0
				x.ISOWeek = int64(w)
				n = 0
			}
			x.Dt = time.Date(a.Dt.Year(), a.Dt.Month(), a.Dt.Day(), 0, 0, 0, 0, time.UTC) // set this correctly whether we've collected anything or not
			x.ISOWeek = int64(w)
			// util.Console("x.Dt set to %s  ISOWeek = %d\n", x.Dt.Format(util.RRDATEREPORTFMT), x.ISOWeek)
		}

		//---------------------------------------------------
		// Still on the same week, add this to the totals...
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

	//----------------------------------------
	// Save anything that we've collected...
	//----------------------------------------
	if n > 0 {
		x.Open /= float64(n)
		x.Close /= float64(n)
		x.High /= float64(n)
		x.Low /= float64(n)

		// write or update this record
		if err = writeUpdateExchWeekly(&x, t); err != nil {
			errors++
			return errors, warnings, nil
		}
	}
	return errors, warnings, nil
}

// Write the specified record. If it exists, update it with this information.
//
// INPUTS
// x = pointer to struct to write
// t = ticker
// ------------------------------------------------------------------------------
func writeUpdateExchWeekly(x *db.ExchWeekly, t string) error {
	var err error
	// Try to insert...
	if _, err = db.InsertExchWeekly(App.ctx, x); err == nil {
		return nil // if that worked, we're done
	}

	util.Console("\nInsertExchWeekly failed\n")

	//-----------------------------------------------------------
	// If the error was Duplicate Entry, then we just update...
	//-----------------------------------------------------------
	errs := err.Error()
	if !strings.Contains(errs, "Duplicate entry") {
		return err // if the error was something else, then we're done, just return the error
	}

	qry := fmt.Sprintf(
		`SELECT %s FROM ExchWeekly WHERE Ticker = "%s" AND YEAR(Dt)=%d AND ISOWeek=%d`,
		db.Pdb.DBFields["ExchWeekly"], t, x.Dt.Year(), x.ISOWeek)
	row := db.Pdb.DB.QueryRow(qry)
	var x1 db.ExchWeekly
	if err = db.ReadExchWeekly(row, &x1); err != nil {
		return err
	}

	x1.Open = x.Open
	x1.Close = x.Close
	x1.High = x.High
	x1.Low = x.Low
	return db.UpdateExchWeekly(App.ctx, &x1)
}
