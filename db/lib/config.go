package db

import (
	"database/sql"
	"extres"
	"fmt"
	"log"
	"math/rand"
	util "platosrv/util/lib"
	"time"

	"github.com/kardianos/osext"
)

// Pdb is a struct with all variables needed by the db infrastructure
var Pdb struct {
	Prepstmt PrepSQL
	Config   extres.ExternalResources
	DB       *sql.DB
	DBFields map[string]string // map of db table fields DBFields[tablename] = field list
	Zone     *time.Location    // what timezone should the server use?
	Key      []byte            // crypto key
	Rand     *rand.Rand        // for generating Reference Numbers or other UniqueIDs
	noAuth   bool              // is authrization needed to access the db?
}

// ReadConfig will read the configuration file "config.json" if
// it exists in the current directory
func ReadConfig() error {
	folderPath, err := osext.ExecutableFolder()
	if err != nil {
		log.Fatal(err)
	}
	// fmt.Printf("Executable folder = %s\n", folderPath)
	fname := folderPath + "/config.json"
	err = extres.ReadConfig(fname, &Pdb.Config)

	Pdb.Zone, err = time.LoadLocation(Pdb.Config.Timezone)
	if err != nil {
		fmt.Printf("Error loading timezone %s : %s\n", Pdb.Config.Timezone, err.Error())
		util.Ulog("Error loading timezone %s : %s", Pdb.Config.Timezone, err.Error())
		return err
	}
	return err
}

// Init initializes the db subsystem
func Init(db *sql.DB) error {
	Pdb.DB = db
	Pdb.DBFields = map[string]string{}
	BuildPreparedStatements()
	return nil
}
