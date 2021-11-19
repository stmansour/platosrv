package db

import (
	"strings"
)

// June 3, 2016 -- As the params change, it's easy to forget to update all the statements with the correct
// field names and the proper number of replacement characters.  I'm starting a convention where the SELECT
// fields are set into a variable and used on all the SELECT statements for that table.  The fields and
// replacement variables for INSERT and UPDATE are derived from the SELECT string.

var mySQLRpl = string("?")
var myRpl = mySQLRpl

// GenSQLInsertAndUpdateStrings generates a string suitable for SQL INSERT and UPDATE statements given the fields as used in SELECT statements.
//
//  example:
//	given this string:      "LID,BID,RAID,GLNumber,Status,Type,Name,AcctType,LastModTime,LastModBy"
//  we return these five strings:
//  1)  "BID,RAID,GLNumber,Status,Type,Name,AcctType,LastModBy"                 -- use for SELECT
//  2)  "?,?,?,?,?,?,?,?"  														-- use for INSERT
//  3)  "BID=?RAID=?,GLNumber=?,Status=?,Type=?,Name=?,AcctType=?,LastModBy=?"  -- use for UPDATE
//  4)  "LID,BID,RAID,GLNumber,Status,Type,Name,AcctType,LastModBy", 			-- use for INSERT (no PRIMARYKEY), add "WHERE LID=?"
//  5)  "?,?,?,?,?,?,?,?,?"  													-- use for INSERT (no PRIMARYKEY)
//
// Note that in this convention, we remove LastModTime from insert and update statements (the db is set up to update them by default) and
// we remove the initial ID as that number is AUTOINCREMENT on INSERTs and is not updated on UPDATE.
func GenSQLInsertAndUpdateStrings(s string) (string, string, string, string, string) {
	fields := strings.Split(s, ",")

	// mostly 0th element is ID, but it is not necessary
	s0 := fields[0]
	s2 := fields[1:] // skip the ID

	insertFields := []string{} // fields which are allowed while INSERT
	updateFields := []string{} // fields which are allowed while while UPDATE

	// remove fields which value automatically handled by database while insert and update op.
	for _, fld := range s2 {
		fld = strings.TrimSpace(fld)
		if fld == "" { // if nothing then continue
			continue
		}
		// INSERT FIELDS Inclusion
		if fld != "LastModTime" && fld != "CreateTime" { // remove these fields for INSERT
			insertFields = append(insertFields, fld)
		}
		// UPDATE FIELDS Inclusion
		if fld != "LastModTime" && fld != "CreateTime" && fld != "CreateBy" { // remove these fields for UPDATE
			updateFields = append(updateFields, fld)
		}
	}

	var s3, s4 string
	for i := range insertFields {
		if i == len(insertFields)-1 {
			s3 += myRpl
		} else {
			s3 += myRpl + ","
		}
	}

	for i, uFld := range updateFields {
		if i == len(updateFields)-1 {
			s4 += uFld + "=" + myRpl
		} else {
			s4 += uFld + "=" + myRpl + ","
		}
	}

	// list down insert fields with comma separation
	s = strings.Join(insertFields, ",")

	s5 := s0 + "," + s     // for INSERT where first val is not AUTOINCREMENT
	s6 := s3 + "," + myRpl // for INSERT where first val is not AUTOINCREMENT
	return s, s3, s4, s5, s6
}

// BuildPreparedStatements is where we build the DBFields map and create the
// prepared sql statements for queries
//
// INPUTS
//
// RETURNS
//
//------------------------------------------------------------------------------
func BuildPreparedStatements() {
	var err error
	var s1, s2, s3, flds string

	//==========================================
	// Exch
	//==========================================
	flds = "XID,Dt,Ticker,Open,High,Low,Close,CreateTime,CreateBy,LastModTime,LastModBy"
	Pdb.DBFields["Exch"] = flds
	Pdb.Prepstmt.GetExch, err = Pdb.DB.Prepare("SELECT " + flds + " FROM Exch WHERE XID=?")
	Errcheck(err)
	s1, s2, s3, _, _ = GenSQLInsertAndUpdateStrings(flds)
	Pdb.Prepstmt.InsertExch, err = Pdb.DB.Prepare("INSERT INTO Exch (" + s1 + ") VALUES(" + s2 + ")")
	Errcheck(err)
	Pdb.Prepstmt.UpdateExch, err = Pdb.DB.Prepare("UPDATE Exch SET " + s3 + " WHERE XID=?")
	Errcheck(err)
	Pdb.Prepstmt.DeleteExch, err = Pdb.DB.Prepare("DELETE from Exch WHERE XID=?")
	Errcheck(err)

	//==========================================
	// Item
	//==========================================
	flds = "IID,Title,Description,PubDt,Link,CreateTime,CreateBy,LastModTime,LastModBy"
	Pdb.DBFields["Item"] = flds
	Pdb.Prepstmt.GetItem, err = Pdb.DB.Prepare("SELECT " + flds + " FROM Item WHERE IID=?")
	Errcheck(err)
	s1, s2, s3, _, _ = GenSQLInsertAndUpdateStrings(flds)
	Pdb.Prepstmt.InsertItem, err = Pdb.DB.Prepare("INSERT INTO Item (" + s1 + ") VALUES(" + s2 + ")")
	Errcheck(err)
	Pdb.Prepstmt.UpdateItem, err = Pdb.DB.Prepare("UPDATE Item SET " + s3 + " WHERE IID=?")
	Errcheck(err)
	Pdb.Prepstmt.DeleteItem, err = Pdb.DB.Prepare("DELETE from Item WHERE IID=?")
	Errcheck(err)

}
