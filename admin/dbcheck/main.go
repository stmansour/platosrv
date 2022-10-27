package main

import (
	"database/sql"
	"extres"
	"flag"
	"fmt"
	"os"
	db "platosrv/db/lib"
	util "platosrv/util/lib"
	"platosrv/ws"

	_ "github.com/go-sql-driver/mysql"
)

// App is the global application structure
var App struct {
	db        *sql.DB
	DBName    string
	DBUser    string
	Port      int      // port on which platosrv listens
	LogFile   *os.File // where to log messages
	fname     string
	startName string
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

	util.Console("*** PLATODB CHECKER ***\n")
	util.Console("Using database: %s , host = %s, port = %d\n", db.Pdb.Config.PlatoDbname, db.Pdb.Config.PlatoDbhost, db.Pdb.Config.PlatoDbport)

	// Get the database...
	// s := "<awsdbusername>:<password>@tcp(<rdsinstancename>:3306)/accord"
	s := extres.GetSQLOpenString(db.Pdb.Config.PlatoDbname, &db.Pdb.Config)
	if App.db, err = sql.Open("mysql", s); err != nil {
		util.Console("sql.Open for database=%s, dbuser=%s: Error = %v\n", db.Pdb.Config.PlatoDbname, db.Pdb.Config.PlatoDbuser, err)
		os.Exit(1)
	}
	util.Ulog("successfully opened database %q as user %q on %s\n", db.Pdb.Config.PlatoDbname, db.Pdb.Config.PlatoDbuser, db.Pdb.Config.PlatoDbhost)
	defer App.db.Close()

	if err = App.db.Ping(); nil != err {
		util.Console("could not ping database %q as user %q on %s\n", db.Pdb.Config.PlatoDbname, db.Pdb.Config.PlatoDbuser, db.Pdb.Config.PlatoDbhost)
		util.Console("error: %s\n", err.Error())
		util.Console("App.db.Ping for database=%s, dbuser=%s: Error = %v\n", db.Pdb.Config.PlatoDbname, db.Pdb.Config.PlatoDbuser, err)
		os.Exit(1)
	}
	db.Init(App.db)              // initializes database
	db.BuildPreparedStatements() // the prepared statement for db access

	DBCheck()
}
