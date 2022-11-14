package main

import (
	"context"
	"database/sql"
	"extres"
	"flag"
	"fmt"
	"os"
	db "platosrv/db/lib"
	"platosrv/session"
	util "platosrv/util/lib"
	"platosrv/ws"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// App is the global application structure
var App struct {
	db        *sql.DB
	ctx       context.Context // context to use for this app
	DBName    string
	DBUser    string
	Port      int      // port on which platosrv listens
	LogFile   *os.File // where to log messages
	fname     string
	startName string
	Warnings  bool // true if we want to show warnings
}

func readCommandLineArgs() {
	portPtr := flag.Int("p", 8277, "port on which platosrv server listens")
	vptr := flag.Bool("v", false, "Show version, then exit")
	wptr := flag.Bool("w", false, "Don't show warnings")
	flag.Parse()
	if *vptr {
		fmt.Printf("Version:   %s\n", ws.GetVersionNo())
		os.Exit(0)
	}
	App.Port = *portPtr
	App.Warnings = !*wptr
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
	util.Console("Initiating DB\n")
	db.Init(App.db) // initializes database
	util.Console("Building session table\n")
	session.Init(10, db.Pdb.Config) // we must have login sessions
	util.Console("Building prepared statements\n")
	db.BuildPreparedStatements() // the prepared statement for db access

	//------------------------------------------------------------------------
	// Create a session that this process can use for accessing the database
	//------------------------------------------------------------------------
	util.Console("Create db context\n")
	now := time.Now()
	App.ctx = context.Background()
	expire := now.Add(4 * time.Hour) // 4 hours
	sess := session.New(
		"dbtest-app"+fmt.Sprintf("%010x", expire.Unix()), // token
		"dbtest",      // username
		"dbtest-app",  // name string
		int64(-99998), // uid
		"",            // image url
		-1,            // security role id
		&expire)       // expiredt
	App.ctx = session.SetSessionContextKey(App.ctx, sess)

	DBCheck()
	createExchDaily(App.ctx)
}
