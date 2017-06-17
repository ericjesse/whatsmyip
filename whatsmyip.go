package main

import (
	"errors"
	"flag"
	"os"

	"whatsmyip/handlers"
	"whatsmyip/logger"
	"whatsmyip/sql/maintenance"
	"whatsmyip/sql/schema"

	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/urfave/negroni"

	_ "github.com/lib/pq"
)

var (
	// Errors.
	errPortNotSet  = errors.New("The port to listen was not specified")
	errDbURLNotSet = errors.New("The URL to the database was not specified")

	// Flags.
	port, dbType, dbURL, dataRetentionDuration string

	// Schema creation commands.
	creationCommands = []schema.CreationCommand{
		{SchemaVersionTag: 1, SQLCommand: `CREATE SEQUENCE "ipCheckSeq" INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 START 1 CACHE 1`},
		{SchemaVersionTag: 2, SQLCommand: `CREATE TABLE "ipCheck" (
				id BIGINT NOT NULL DEFAULT nextval('"ipCheckSeq"'::regclass),
				instant timestamp without time zone NOT NULL DEFAULT now(),
				"ipAddressV4" character varying(15) NOT NULL,
				source character varying(30) NOT NULL,
				CONSTRAINT "ipCheckPkey" PRIMARY KEY (id))`},
		{SchemaVersionTag: 3, SQLCommand: `CREATE SEQUENCE "geoLocSeq" INCREMENT 1 MINVALUE 1 MAXVALUE 9223372036854775807 START 1 CACHE 1`},
		{SchemaVersionTag: 4, SQLCommand: `CREATE TABLE "geoLoc" (
				id BIGINT NOT NULL DEFAULT nextval('"ipCheckSeq"'::regclass),
				instant timestamp without time zone NOT NULL DEFAULT now(),
				"ipAddressV4" character varying(15) NOT NULL,
				provider character varying(200) ,
				city character varying(200),
				country character varying(100),
				"countryCode" character varying(5),
				region character varying(200),
				timezone character varying(200),
				"zipCode" character varying(15),
				latitude double precision,
				longitude double precision,
				CONSTRAINT "geoLocPkey" PRIMARY KEY (id))`},
	}
)

func init() {

	// Flags.
	flag.StringVar(&port, "port", os.Getenv("PORT"), "Port to listen for the services")
	flag.StringVar(&dbType, "dbType", "postgres", "Name of the database driver to use")
	flag.StringVar(&dbURL, "dbUrl", os.Getenv("DATABASE_URL"), "URL to connect to the database")
	flag.StringVar(&dataRetentionDuration, "dataRetentionDuration", os.Getenv("DATA_RETENTION"), `Duration for the data retention. Default is 6 weeks. A duration string is a possibly signed sequence of
		decimal numbers, each with optional fraction and a unit suffix,	such as "300ms", "-1.5h" or "2h45m".
		Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h""`)
}

func main() {
	flag.Parse()
	log := logger.GetLogger()

	err := validateRequiredArgs(port, dbURL)
	if err != nil {
		log.Fatal(err)
	}

	// ***************************************
	// DB Part
	// ***************************************
	// Preparing the database.
	log.Debugln("Connecting to the DB", dbURL)
	var sm *schema.Manager
	sm, err = schema.NewSchemaManager(dbType, dbURL, creationCommands)
	if sm != nil {
		// Close the DB connection when the program ends.
		defer sm.Db.Close()
	}
	if err != nil {
		log.Fatal(err)
	}

	sm.WaitForSchemaCompletion()
	log.Println("Current database version is", sm.DbVersion)

	// Set the DB connection to the handlers.
	handlers.Db = sm.Db
	maintenance.Start(sm.Db, dataRetentionDuration)

	// ***************************************
	// HTTP Part
	// ***************************************

	// Create a Gorilla Mux router.
	router := mux.NewRouter()
	// Add the Negroni Middlewares.
	n := negroni.New()
	// Log all the requests.
	n.Use(negroni.NewLogger())
	// Manage the cases when a handler raises a panic to return 500 Internal Server Error.
	n.Use(negroni.NewRecovery())
	// Add the automatic handling of deflated request and responses.
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.UseHandler(router)

	// Adding the handlers.
	handlers.HandleIPRequest(router, "/ip")
	handlers.HandleStaticRequest(router, "/static")

	log.Printf("Service starting on port %s\n", port)
	n.Run(":" + port)
}

// validateRequiredArgs validates all the required args and merge values when several sources are possible (flags or env).
func validateRequiredArgs(portFlag string, dbURLFlag string) error {
	// Validate the port to listen.
	if portFlag == "" {
		return errPortNotSet
	}
	// Validate the URL to the DB.
	if dbURLFlag == "" {
		return errDbURLNotSet
	}
	return nil
}
