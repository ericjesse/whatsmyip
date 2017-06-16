package main

import (
	"errors"
	"flag"
	"fmt"
	"os"

	"whatsmyip/handlers"
	"whatsmyip/logger"
	"whatsmyip/sql/schema"

	"github.com/gorilla/mux"
	"github.com/phyber/negroni-gzip/gzip"
	"github.com/urfave/negroni"

	_ "github.com/lib/pq"
)

var (
	errPortNotSet  = errors.New("The port to listen was not specified")
	errDbURLNotSet = errors.New("The URL to the database was not specified")
)

var creationCommands = []schema.CreationCommand{
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

func main() {
	log := logger.GetLogger()

	// Get the port to use from the environment variables.
	port := os.Getenv("PORT")
	dbURL := os.Getenv("DATABASE_URL")

	// Flags.
	portFlag := flag.String("port", "", "Port to listen for the services")
	dbTypeFlag := flag.String("dbType", "postgres", "Name of the database driver to use")
	dbURLFlag := flag.String("dbUrl", "", "URL to connect to the database")
	flag.Parse()

	err := validateRequiredArgs(&port, portFlag, &dbURL, dbURLFlag)
	if err != nil {
		log.Fatal(err)
	}

	log.Debugln("Connecting to the DB", dbURL)
	log.Debugln("Listening on the port", port)

	// Preparing the database.
	var sm *schema.Manager
	sm, err = schema.NewSchemaManager(*dbTypeFlag, dbURL, creationCommands)
	if sm != nil {
		defer sm.Db.Close()
	}
	if err != nil {
		log.Fatal(err)
	}

	sm.WaitForSchemaCompletion()
	log.Println("Current database version is", sm.DbVersion)

	handlers.Db = sm.Db

	// Create a Gorilla Mux router.
	router := mux.NewRouter()
	handlers.HandleIPRequest(router, "/ip")
	handlers.HandleStaticRequest(router, "/static")

	// Add the Negroni Middlewares.
	n := negroni.New()
	// Log all the requests.
	n.Use(negroni.NewLogger())
	// Manage the cases when a handler raises a panic to return 500 Internal Server Error.
	n.Use(negroni.NewRecovery())
	// Add the automatic handling of deflated request and responses.
	n.Use(gzip.Gzip(gzip.DefaultCompression))
	n.UseHandler(router)

	fmt.Fprintf(os.Stdout, "Service starting on port %s\n", port)
	n.Run(":" + port)
}

// validateRequiredArgs validates all the required args and merge values when several sources are possible (flags or env).
func validateRequiredArgs(port *string, portFlag *string, dbURL *string, dbURLFlag *string) error {
	// Validate the port to listen.
	if *port == "" && *portFlag == "" {
		return errPortNotSet
	}
	// If the port was specified as a flag, it is set into the port variable.
	if *port == "" && *portFlag != "" {
		*port = *portFlag
	}
	// Validate the URL to the DB.
	if *dbURL == "" && *dbURLFlag == "" {
		return errDbURLNotSet
	}
	// If the DB URL was specified as a flag, it is set into the dbURL variable.
	if *dbURL == "" && *dbURLFlag != "" {
		*dbURL = *dbURLFlag
	}
	return nil
}
