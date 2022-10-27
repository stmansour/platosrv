package main

import (
	"fmt"
	db "platosrv/db/lib"
	util "platosrv/util/lib"
	"time"
)

// DBCheck scans "platodb", finds holes in the data, and does a bunch of
// consistency checking.
// ------------------------------------------------------------------------
func DBCheck() {
	var err error
	util.Console("Check the Exch Table\n")

	//----------------------------------------------
	// Iterate through the mappings we store...
	//----------------------------------------------
	for k, v := range Tickers {
		if v > 0 {
			util.Console("\nProcessing %s\n", k)
			if err = scanExch(k); err != nil {
				util.Console("Error in scanExch: %s\n", err.Error)
			}
		}
	}
}

func scanExch(t string) error {
	var a db.Exch
	var ldt time.Time
	var ldiff time.Duration
	var count int64
	var errors int64
	var warnings int64

	twoDays := 48*time.Hour + 4*time.Minute // essentially 2 days and change
	fiveMinutes := 5 * time.Minute

	qry := fmt.Sprintf("SELECT %s FROM Exch WHERE Ticker = \"%s\" ORDER BY Dt ASC", db.Pdb.DBFields["Exch"], t)
	// util.Console("qry = %s\n", qry)
	rows, err := db.Pdb.DB.Query(qry)
	if err != nil {
		return err
	}
	defer rows.Close()

	count = 0
	for rows.Next() {
		if err := db.ReadExchs(rows, &a); err != nil {
			return err
		}
		if count == 0 {
			util.Console("data begins on %s\n", a.Dt.Format("Mon 2006-01-02 15:04:05"))
			ldt = a.Dt
			count++
			continue
		}

		//---------------------------------------------------------------------------
		// The exchange rate is published for every minute of the day. If we see
		// that the last value checked was 1 min after its predecessor and this
		// value is > 1 min, then we have a gap
		//---------------------------------------------------------------------------
		diff := a.Dt.Sub(ldt)
		if diff.Minutes() > float64(1.0) && ldiff.Minutes() <= float64(1.0) {
			if a.Dt.Weekday() == time.Sunday && (twoDays-diff).Minutes() < fiveMinutes.Minutes() {
				util.Console("** WARNING ***  weekend gap at %s\n", ldt.Format("Mon 2006-01-02 15:04:05"))
				warnings++
			} else {
				util.Console("*** ERROR ***  gap = %s  begining at %s\n", diff, ldt.Format("Mon 2006-01-02 15:04:05"))
				errors++
			}
		}

		ldt = a.Dt   // this was the time of the last record checked
		ldiff = diff // this was the time delta from the last record checked
		count++
	}

	if err := rows.Err(); err != nil {
		return err
	}

	util.Console("   records processed: %d\n", count)
	util.Console("   errors: %d\n", errors)
	util.Console("   warnings: %d\n", warnings)

	return nil
}
