package main

import (
	"context"
	"fmt"
	"os"
	db "platosrv/db/lib"
	util "platosrv/util/lib"
	"strings"
	"time"
)

// TestItem checks the basic db functions for the Exch struct
// -----------------------------------------------------------------------------
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
// -----------------------------------------------------------------------------
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

// TestExchDaily checks the basic db functions for the ExchDaily struct
// -----------------------------------------------------------------------------
func TestExchDaily(ctx context.Context) {
	var err error
	dt := time.Date(2020, time.March, 23, 0, 0, 0, 0, time.UTC)
	util.Console("Entered TestExchDaily\n")
	rs := db.ExchDaily{
		XDID:   0,
		Dt:     dt,
		Ticker: "SMTEST",
		Open:   float64(175.45),
		High:   float64(187.12),
		Low:    float64(171.34),
		Close:  float64(185.62),
	}
	var delid, id int64
	if id, err = db.InsertExchDaily(ctx, &rs); err != nil {
		fmt.Printf("Error inserting first ExchDaily: %s\n", err)
		os.Exit(1)
	}

	// Insert another for delete...
	rs.Dt = time.Date(2020, time.March, 24, 0, 0, 0, 0, time.UTC)
	if delid, err = db.InsertExchDaily(ctx, &rs); err != nil {
		fmt.Printf("Error inserting second ExchDaily: %s\n", err)
		os.Exit(1)
	}
	if err = db.DeleteExchDaily(ctx, delid); err != nil {
		fmt.Printf("Error deleting ExchDaily: %s\n", err)
		os.Exit(1)
	}

	var rs1 db.ExchDaily
	if rs1, err = db.GetExchDaily(ctx, id); err != nil {
		fmt.Printf("error in GetExchDaily: %s\n", err.Error())
		os.Exit(1)
	}
	rs1.Low = float64(172.34)
	if err = db.UpdateExchDaily(ctx, &rs1); err != nil {
		fmt.Printf("Error updating ExchDaily: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Success! Delete, Get, Insert, and Update ExchDaily\n")
}

// TestExchMonthly checks the basic db functions for the ExchMonthly struct
// -----------------------------------------------------------------------------
func TestExchMonthly(ctx context.Context) {
	var err error
	dt := time.Date(2020, time.March, 1, 0, 0, 0, 0, time.UTC)
	util.Console("Entered TestExchMonthly\n")
	rs := db.ExchMonthly{
		XMID:   0,
		Dt:     dt,
		Ticker: "SMTEST",
		Open:   float64(175.45),
		High:   float64(187.12),
		Low:    float64(171.34),
		Close:  float64(185.62),
	}
	var delid, id int64
	if id, err = db.InsertExchMonthly(ctx, &rs); err != nil {
		fmt.Printf("Error inserting ExchMonthly: %s\n", err)
		os.Exit(1)
	}

	// Insert another for delete...
	rs.Dt = time.Date(2020, time.April, 1, 0, 0, 0, 0, time.UTC)
	if delid, err = db.InsertExchMonthly(ctx, &rs); err != nil {
		fmt.Printf("Error inserting ExchMonthly: %s\n", err)
		os.Exit(1)
	}
	if err = db.DeleteExchMonthly(ctx, delid); err != nil {
		fmt.Printf("Error deleting ExchMonthly: %s\n", err)
		os.Exit(1)
	}

	var rs1 db.ExchMonthly
	if rs1, err = db.GetExchMonthly(ctx, id); err != nil {
		fmt.Printf("error in GetExchMonthly: %s\n", err.Error())
		os.Exit(1)
	}
	rs1.Low = float64(172.34)
	if err = db.UpdateExchMonthly(ctx, &rs1); err != nil {
		fmt.Printf("Error updating ExchMonthly: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Success! Delete, Get, Insert, and Update ExchMonthly\n")
}

// TestExchQuarterly checks the basic db functions for the ExchQuarterly struct
// -----------------------------------------------------------------------------
func TestExchQuarterly(ctx context.Context) {
	var err error
	dt := time.Date(2020, time.March, 1, 0, 0, 0, 0, time.UTC)
	util.Console("Entered TestExchQuarterly\n")
	rs := db.ExchQuarterly{
		XQID:   0,
		Dt:     dt,
		Ticker: "SMTEST",
		Open:   float64(175.45),
		High:   float64(187.12),
		Low:    float64(171.34),
		Close:  float64(185.62),
	}
	var delid, id int64
	if id, err = db.InsertExchQuarterly(ctx, &rs); err != nil {
		fmt.Printf("Error inserting ExchQuarterly: %s\n", err)
		os.Exit(1)
	}

	// Insert another for delete...
	rs.Dt = time.Date(2020, time.June, 1, 0, 0, 0, 0, time.UTC)
	if delid, err = db.InsertExchQuarterly(ctx, &rs); err != nil {
		fmt.Printf("Error inserting ExchQuarterly: %s\n", err)
		os.Exit(1)
	}
	if err = db.DeleteExchQuarterly(ctx, delid); err != nil {
		fmt.Printf("Error deleting ExchQuarterly: %s\n", err)
		os.Exit(1)
	}

	var rs1 db.ExchQuarterly
	if rs1, err = db.GetExchQuarterly(ctx, id); err != nil {
		fmt.Printf("error in GetExchQuarterly: %s\n", err.Error())
		os.Exit(1)
	}
	rs1.Low = float64(172.34)
	if err = db.UpdateExchQuarterly(ctx, &rs1); err != nil {
		fmt.Printf("Error updating ExchQuarterly: %s\n", err)
		os.Exit(1)
	}

	fmt.Printf("Success! Delete, Get, Insert, and Update ExchQuarterly\n")
}

// TestRSSFeed checks the basic db functions for the RSSFeed struct
// -----------------------------------------------------------------------------
func TestRSSFeed(ctx context.Context) {
	var err error
	var r db.RSSFeed
	var id, id2 int64
	fmt.Printf("\n------------------------------------------------\n")
	r.URL = "https://rss.nytimes.com/services/xml/rss/nyt/World.xml"
	r.FLAGS = int64(0)
	if id, err = db.InsertRSSFeed(ctx, &r); err != nil {
		fmt.Printf("Error inserting RSSFeed: %s\n", err)
		os.Exit(1)
	}
	if id != r.RSSID {
		fmt.Printf("Error id does not match RSSID\n")
		os.Exit(1)
	}

	var r2 db.RSSFeed
	r2.URL = "https://rss.nytimes.com/services/xml/rss/nyt/World.xml"
	if _, err = db.InsertRSSFeed(ctx, &r2); err != nil { // this should fail because we already have that URL
		// make sure this is error 1062
		if !strings.Contains(err.Error(), "Error 1062: Duplicate") {
			fmt.Printf("Duplicate key was added to RSSFeed table!!!\n")
			os.Exit(1)
		}
	}
	r2.URL = "https://rss.nytimes.com/services/xml/rss/nyt/Education.xml"
	if id2, err = db.InsertRSSFeed(ctx, &r2); err != nil {
		fmt.Printf("Error inserting RSSFeed: %s\n", err)
		os.Exit(1)
	}
	if id2 != r2.RSSID {
		fmt.Printf("Error id2 does not match RSSID\n")
		os.Exit(1)
	}
	// now update it
	r2.URL = "https://rss.nytimes.com/services/xml/rss/nyt/Politics.xml"
	if err = db.UpdateRSSFeed(ctx, &r2); err != nil {
		fmt.Printf("Error updating RSSFeed: %s\n", err)
		os.Exit(1)
	}
	if err = db.DeleteRSSFeed(ctx, r2.RSSID); err != nil {
		fmt.Printf("Error deleting RSSFeed: %s\n", err)
		os.Exit(1)
	}

	// add a few more so we have more than one in the db...
	var feeds = []string{"Space", "Computers", "Business", "NYReligion"}
	for i := 0; i < len(feeds); i++ {
		r2.URL = fmt.Sprintf("https://rss.nytimes.com/services/xml/rss/nyt/%s.xml", feeds[i])
		if id2, err = db.InsertRSSFeed(ctx, &r2); err != nil {
			fmt.Printf("Error inserting RSSFeed: %s\n", err)
			os.Exit(1)
		}
		if id2 != r2.RSSID {
			fmt.Printf("Error id2 does not match RSSID\n")
			os.Exit(1)
		}
	}
	fmt.Printf("------------------------------------------------\n")
	fmt.Printf("Any errors above were part of the test.\n")
	fmt.Printf("Hitting this point means everything worked.\n\n")
}

// TestItemFeed checks the basic db functions for the ItemFeed struct
// -----------------------------------------------------------------------------
func TestItemFeed(ctx context.Context) {
	var err error
	var r db.ItemFeed
	var id, id2 int64
	r.RSSID = 1
	r.IID = 1
	if id, err = db.InsertItemFeed(ctx, &r); err != nil {
		fmt.Printf("Error inserting ItemFeed: %s\n", err)
		os.Exit(1)
	}
	if id != r.IFID {
		fmt.Printf("Error id does not match IFID\n")
		os.Exit(1)
	}

	var r2 db.ItemFeed
	r2.IID = 1
	r2.RSSID = 1
	if _, err = db.InsertItemFeed(ctx, &r2); err != nil { // this should fail because we already have that URL
		// make sure this is error 1062
		if !strings.Contains(err.Error(), "Error 1062: Duplicate") {
			fmt.Printf("Duplicate key was added to ItemFeed table!!!\n")
			os.Exit(1)
		}
	}
	r2.IID = 2
	r2.RSSID = 1
	if id2, err = db.InsertItemFeed(ctx, &r2); err != nil {
		fmt.Printf("Error inserting ItemFeed: %s\n", err)
		os.Exit(1)
	}
	if id2 != r2.IFID {
		fmt.Printf("Error id2 does not match IFID\n")
		os.Exit(1)
	}
	// now update it
	r2.IID = 3
	if err = db.UpdateItemFeed(ctx, &r2); err != nil {
		fmt.Printf("Error updating ItemFeed: %s\n", err)
		os.Exit(1)
	}
	// Insert an itemFeed and delete it
	r2.IFID = 0
	r2.IID = 4
	r2.RSSID = 1
	if id2, err = db.InsertItemFeed(ctx, &r2); err != nil {
		fmt.Printf("Error inserting ItemFeed: %s\n", err)
		os.Exit(1)
	}
	if id2 != r2.IFID {
		fmt.Printf("Error id2 does not match IFID\n")
		os.Exit(1)
	}
	if err = db.DeleteItemFeed(ctx, r2.IFID); err != nil {
		fmt.Printf("Error deleting ItemFeed: %s\n", err)
		os.Exit(1)
	}
	fmt.Printf("Success! Delete, Get, Insert, and Update ItemFeed\n")
}
