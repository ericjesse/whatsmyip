// Package maintenance aims at running maintenance tasks on the database.
package maintenance

import (
	"database/sql"
	"time"
	"whatsmyip/logger"
)

const (
	purgeIPCheckQuery string = `DELETE FROM "ipCheck" WHERE instant < $1`
	purgeGeoLocQuery  string = `DELETE FROM "geoLoc" WHERE instant < $1`
)

var (
	// The default data retention duration is 6 weeks.
	dataRetentionDuration time.Duration = (time.Hour * 24 * 7 * 6)
	log                                 = logger.GetLogger()
)

// Start schedules all the tasks for maintenance of the database.
func Start(db *sql.DB, dataRetentionDurationFlag string) {
	drd, err := time.ParseDuration(dataRetentionDurationFlag)
	if err == nil {
		// If the duration passed as parameter could be read successfully.
		dataRetentionDuration = drd
	}

	log.Debugln("Starting the DB maintenance jobs.")
	go func() {
		ticker := time.NewTicker(time.Minute * 2)
		for {
			hasError := false
			obsolescenceInstant := time.Now().Add(-dataRetentionDuration)
			log.Debugf("Purging data older than %v\n", obsolescenceInstant)
			_, err := db.Exec(purgeGeoLocQuery, obsolescenceInstant)
			if err != nil {
				log.Println(err)
				hasError = true
			}
			_, err = db.Exec(purgeIPCheckQuery, obsolescenceInstant)
			if err != nil {
				log.Println(err)
				hasError = true
			}
			if hasError {
				log.Debugln("Data purged with errors (see above)")
			} else {
				log.Debugln("Data purged successfully")
			}
			<-ticker.C
		}
	}()
}
