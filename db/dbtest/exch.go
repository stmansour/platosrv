package main

import (
	"context"
	"fmt"
	"mojo/util"
	"os"
	db "platosrv/db/lib"
	"time"
)

// TestItem checks the basic db functions for the Exch struct
//-----------------------------------------------------------------------------
func TestItem(ctx context.Context) {
	var err error
	util.Console("Entered TestItem\n")
	dt := time.Date(2020, time.March, 23, 0, 0, 0, 0, time.UTC)
	rs := db.Item{
		IID:         0,
		Title:       "Big Time Article With Deep Meaning",
		Description: "Big time article is really important. Read it.",
		PubDt:       dt,
		Link:        "http://example.com/bigTimeArticle.html",
	}
	var delid, id int64
	if id, err = db.InsertItem(ctx, &rs); err != nil {
		fmt.Printf("Error inserting Item: %s\n", err)
		os.Exit(1)
	}

	// Insert another for delete...
	rs.Link = "http://example.com/bigTimeArticle2.html" // ensure unique
	if delid, err = db.InsertItem(ctx, &rs); err != nil {
		fmt.Printf("Error inserting Item: %s\n", err)
		os.Exit(1)
	}
	if err = db.DeleteItem(ctx, delid); err != nil {
		fmt.Printf("Error deleting Item: %s\n", err)
		os.Exit(1)
	}

	var rs1 db.Item
	if rs1, err = db.GetItem(ctx, id); err != nil {
		fmt.Printf("error in GetItem: %s\n", err.Error())
		os.Exit(1)
	}
	rs.Link = "http://example.com/bigTimeArticle3.html" // ensure unique
	if err = db.UpdateItem(ctx, &rs1); err != nil {
		fmt.Printf("Error updating Item: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Success! Delete, Get, Insert, and Update Item\n")
}

// TestExch checks the basic db functions for the Exch struct
//-----------------------------------------------------------------------------
func TestExch(ctx context.Context) {
	var err error
	dt := time.Date(2020, time.March, 23, 0, 0, 0, 0, time.UTC)
	util.Console("Entered TestExch\n")
	rs := db.Exch{
		XID:    0,
		Dt:     dt,
		Ticker: "SMTEST",
		Open:   float64(175.45),
		High:   float64(187.12),
		Low:    float64(171.34),
		Close:  float64(185.62),
	}
	var delid, id int64
	if id, err = db.InsertExch(ctx, &rs); err != nil {
		fmt.Printf("Error inserting Exch: %s\n", err)
		os.Exit(1)
	}

	// Insert another for delete...
	rs.Dt = dt.Add(1 * time.Second) // need to ensure no duplicate entries
	if delid, err = db.InsertExch(ctx, &rs); err != nil {
		fmt.Printf("Error inserting Exch: %s\n", err)
		os.Exit(1)
	}
	if err = db.DeleteExch(ctx, delid); err != nil {
		fmt.Printf("Error deleting Exch: %s\n", err)
		os.Exit(1)
	}

	var rs1 db.Exch
	if rs1, err = db.GetExch(ctx, id); err != nil {
		fmt.Printf("error in GetExch: %s\n", err.Error())
		os.Exit(1)
	}
	rs1.Low = float64(172.34)
	if err = db.UpdateExch(ctx, &rs1); err != nil {
		fmt.Printf("Error updating Exch: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Success! Delete, Get, Insert, and Update Exch\n")
}
