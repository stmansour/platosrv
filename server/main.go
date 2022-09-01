// faa  a program to scrape the FAA directory site.
package main

import (
	"database/sql"
	"extres"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	db "platosrv/db/lib"
	"platosrv/session"
	util "platosrv/util/lib"
	"platosrv/ws"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

// App is the global data structure for this app
var App struct {
	db        *sql.DB
	DBName    string
	DBUser    string
	Port      int      // port on which platosrv listens
	LogFile   *os.File // where to log messages
	fname     string
	startName string
}

// HomeHandler serves static http content such as the .css files
func HomeHandler(w http.ResponseWriter, r *http.Request) {
	u := r.URL.Path;
	util.Console("***\n***\n*** URL Request: %s\n***\n***\n", u)
	if strings.Contains(u, ".") {
		util.Console("*** CONTAINS .\n")
		Chttp.ServeHTTP(w, r)
	} else if strings.Contains(u,"sim" ) {
		util.Console("*** CONTAINS sim\n")
		Chttp.ServeHTTP(w, r)
	} else {
		util.Console("*** REDIRECT to home\n")
		http.Redirect(w, r, "/home/", http.StatusFound)
	}
}

// Chttp is a server mux for handling unprocessed html page requests.
// For example, a .css file or an image file.
var Chttp = http.NewServeMux()

func initHTTP() {
	Chttp.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static/"))))
	http.HandleFunc("/", HomeHandler)
	http.HandleFunc("/home/", ws.HomeUIHandler)
	http.HandleFunc("/v1/", ws.V1ServiceHandler)
}

func readCommandLineArgs() {
	portPtr := flag.Int("p", 8277, "port on which platosrv server listens")
	vptr := flag.Bool("v", false, "Show version, then exit")
	flag.Parse()
	if *vptr {
		fmt.Printf("Version:   %s\n", ws.GetVersionNo())
		os.Exit(0)
	}
	App.Port = *portPtr
}

func main() {
	var err error
	readCommandLineArgs()
	err = db.ReadConfig()
	if err != nil {
		fmt.Printf("Error in db.ReadConfig: %s\n", err.Error())
		os.Exit(1)
	}

	//==============================================
	// Open the logfile and begin logging...
	//==============================================
	App.LogFile, err = os.OpenFile("platosrv.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	util.ErrCheck(err)
	defer App.LogFile.Close()
	log.SetOutput(App.LogFile)
	util.Ulog("*** platosrv CURRENCY EXCHANGE ANALYSIS ***\n")
	util.Ulog("Using database: %s , host = %s, port = %d\n", db.Pdb.Config.PlatoDbname, db.Pdb.Config.PlatoDbhost, db.Pdb.Config.PlatoDbport)

	// Get the database...
	// s := "<awsdbusername>:<password>@tcp(<rdsinstancename>:3306)/accord"
	s := extres.GetSQLOpenString(db.Pdb.Config.PlatoDbname, &db.Pdb.Config)
	fmt.Printf("s = %q\n",s)
	if App.db, err = sql.Open("mysql", s); err != nil {
		fmt.Printf("sql.Open for database=%s, dbuser=%s: Error = %v\n", db.Pdb.Config.PlatoDbname, db.Pdb.Config.PlatoDbuser, err)
		os.Exit(1)
	}
	util.Ulog("successfully opened database %q as user %q on %s\n", db.Pdb.Config.PlatoDbname, db.Pdb.Config.PlatoDbuser, db.Pdb.Config.PlatoDbhost)
	defer App.db.Close()

	if err = App.db.Ping(); nil != err {
		util.Ulog("could not ping database %q as user %q on %s\n", db.Pdb.Config.PlatoDbname, db.Pdb.Config.PlatoDbuser, db.Pdb.Config.PlatoDbhost)
		util.Ulog("error: %s\n", err.Error())
		fmt.Printf("App.db.Ping for database=%s, dbuser=%s: Error = %v\n", db.Pdb.Config.PlatoDbname, db.Pdb.Config.PlatoDbuser, err)
		os.Exit(1)
	}
	db.Init(App.db)                 // initializes database
	session.Init(10, db.Pdb.Config) // we must have login sessions
	db.BuildPreparedStatements()    // the prepared statement for db access
	initHTTP()
	util.Ulog("platosrv initiating HTTP service on port %d\n", App.Port)
	fmt.Printf("Using database: %s , host = %s, port = %d\n", db.Pdb.Config.PlatoDbname, db.Pdb.Config.PlatoDbhost, db.Pdb.Config.PlatoDbport)

	//go http.ListenAndServeTLS(fmt.Sprintf(":%d", App.Port+1), App.CertFile, App.KeyFile, nil)
	err = http.ListenAndServe(fmt.Sprintf(":%d", App.Port), nil)
	if nil != err {
		fmt.Printf("*** Error on http.ListenAndServe: %v\n", err)
		util.Ulog("*** Error on http.ListenAndServe: %v\n", err)
		os.Exit(1)
	}
}
