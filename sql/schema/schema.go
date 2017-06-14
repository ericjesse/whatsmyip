// Package schema is in charge of creating and updating the database schema based upon a simple list of creation queries.
package schema

import (
	"database/sql"
	"errors"
	"log"
	"sort"
	"sync"
)

var (
	errQueryToManageSchemaVersionNotFound = errors.New("The queries to manage the schema version for the expected driver were not found")

	// DebugMode indicates if debug messages have to be logged. It has to be set by the main code.
	DebugMode bool
)

// versionQuery is a type containing the queries required to administrate the schema for a given RDBMS.
type versionQuery struct {
	tableCreationQuery string
	selectVersionQuery string
	updateVersionQuery string
}

// versionQueries is the map of versionQuery for the supported RDBMS. Only postgres is supported for now.
var versionQueries = map[string]versionQuery{
	"postgres": {"CREATE TABLE \"schemaVersion\" (version bigint NOT NULL DEFAULT 0, \"versionDeploymentDate\" timestamp without time zone NOT NULL DEFAULT now())", "SELECT MAX(version) FROM \"schemaVersion\"", "INSERT INTO \"schemaVersion\" (version) VALUES ($1)"},
}

// CreationCommand is a SQL command to create or update a database schema.
// The SchemaVersionTag specifies the order of execution.
type CreationCommand struct {
	SchemaVersionTag int32
	SQLCommand       string
}

// Manager aims at creating, updating and providing details about the current DB schema.
type Manager struct {
	DataSourceName string
	DbVersion      int32
	Db             *sql.DB
	mutex          sync.RWMutex
}

// NewSchemaManager creates a new structure of schema manager based upon the parameters.
// dbDriverName: the name of the driver, as registered.
// dataSourceName: the URL of the data source.
// creationCommands: the slice of SQL commands to create or update the schema.
func NewSchemaManager(dbDriverName string, dataSourceName string, creationCommands []CreationCommand) (*Manager, error) {
	// Open a connection to the DB.
	db, err := sql.Open(dbDriverName, dataSourceName)
	if err != nil {
		return nil, err
	}

	// Test the connection.
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	sm := &Manager{DataSourceName: dataSourceName, Db: db}
	// The schema cannot be used for the moment.
	sm.mutex.Lock()
	defer sm.mutex.Unlock()

	// Sort the commands by version tags.
	sort.Slice(creationCommands, func(i int, j int) bool {
		return creationCommands[i].SchemaVersionTag < creationCommands[j].SchemaVersionTag
	})

	sm.DbVersion, err = retrieveCurrentSchemaVersion(dbDriverName, db)
	if err != nil {
		return nil, err
	}

	var commandsToExecute []CreationCommand
	for i := 0; commandsToExecute == nil && i < len(creationCommands); i++ {
		// Find the first command more recent than the schema version.
		if creationCommands[i].SchemaVersionTag > sm.DbVersion {
			commandsToExecute = creationCommands[i:]
		}
	}

	if len(commandsToExecute) > 0 {
		sm.DbVersion, err = updateSchemaAndGetNewVersion(db, commandsToExecute, versionQueries[dbDriverName].updateVersionQuery)
		if err != nil {
			return nil, err
		}
	}

	return sm, nil
}

// Fetch the current version of the
func retrieveCurrentSchemaVersion(dbDriverName string, db *sql.DB) (int32, error) {

	// Lookup the select version query.
	if _, ok := versionQueries[dbDriverName]; !ok {
		return -1, errQueryToManageSchemaVersionNotFound
	}

	versionQuery := versionQueries[dbDriverName]
	var currentSchemaVersion int32 = -1
	var dbVersion sql.NullInt64
	if DebugMode {
		log.Println("Executing", versionQuery.selectVersionQuery)
	}
	row := db.QueryRow(versionQuery.selectVersionQuery)
	if err := row.Scan(&dbVersion); err != nil {
		// The table does not exist. Create it.
		if DebugMode {
			log.Println("Executing", versionQuery.tableCreationQuery)
		}
		_, err := db.Exec(versionQuery.tableCreationQuery)
		if err != nil {
			// The version table could not be created.
			return currentSchemaVersion, err
		}
	} else if dbVersion.Valid { // A version was found in the database.
		currentSchemaVersion = int32(dbVersion.Int64)
	}

	return currentSchemaVersion, nil
}

func updateSchemaAndGetNewVersion(db *sql.DB, commandsToExecute []CreationCommand, updateVersionQuery string) (int32, error) {
	var currentSchemaVersion int32 = -1

	// Execute all the commands for the update.
	for _, command := range commandsToExecute {
		if DebugMode {
			log.Println("Executing", command.SQLCommand)
		}
		if _, err := db.Exec(command.SQLCommand); err != nil {
			return currentSchemaVersion, err
		}

		// Update the value in the version table.
		if DebugMode {
			log.Println("Executing", updateVersionQuery, "with arguments:", command.SchemaVersionTag)
		}
		if _, err := db.Exec(updateVersionQuery, command.SchemaVersionTag); err != nil {
			return currentSchemaVersion, err
		}
		currentSchemaVersion = command.SchemaVersionTag
	}
	return currentSchemaVersion, nil
}

// WaitForSchemaCompletion blocks until the schema was successfully created or updated.
func (sm *Manager) WaitForSchemaCompletion() {
	sm.mutex.RLock()
	sm.mutex.RUnlock()
}
