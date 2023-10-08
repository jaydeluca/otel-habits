package ingest

import (
	"bufio"
	"database/sql"
	"fmt"
	"github.com/jaydeluca/otel-habits/pkg/models"
	"github.com/jaydeluca/otel-habits/pkg/util"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	bearSQLLiteLocation = "./data/database.sqllite"
	completedIndicator  = "- [x]"
)

func BearData() []models.Timeseries {
	dayCount := os.Getenv("GENERATE_DATA_DAY_COUNT")
	if dayCount != "" {
		dayInt, err := strconv.Atoi(dayCount)
		if err != nil {
			panic(err)
		}
		fmt.Sprintf("Env var set to generate %d days worth of models", dayInt)

		return util.GenerateDummyData(dayInt)
	}
	return extractBearData()
}

func extractBearData() []models.Timeseries {

	fmt.Printf("-- Beginning Ingestion from Bear --\n")

	if _, err := os.Stat(bearSQLLiteLocation); os.IsNotExist(err) {
		panic(fmt.Sprintf("No file found at location %s\\", bearSQLLiteLocation))
	}

	fmt.Printf("Found bear database input file at location: %s\n", bearSQLLiteLocation)

	bearDB, err := sql.Open("sqlite3", bearSQLLiteLocation)
	if err != nil {
		panic(err)
	}
	defer func(bearDB *sql.DB) {
		err := bearDB.Close()
		if err != nil {
			fmt.Printf("Error closing connection: %v", err)
		}
	}(bearDB)

	return retrieveDataFromBear(bearDB)
}

/*
Query SQLlite database for records
where ZTITLE is the date of the entry, and the body is the markdown checklist / notes
*/
func retrieveDataFromBear(db *sql.DB) []models.Timeseries {
	rows, err := db.Query("SELECT ZTITLE, ZTEXT FROM ZSFNOTE z where ZTEXT like '%#daily%' order by ZTITLE DESC")
	if err != nil {
		panic(fmt.Sprintf("failed selecting notes: %v", err))
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			fmt.Printf("error closing rows: %v", err)
		}
	}(rows)

	entries := make([]models.Timeseries, 0)

	for rows.Next() {
		var ZTITLE, ZTEXT string
		if err = rows.Scan(&ZTITLE, &ZTEXT); err != nil {
			panic(err)
		}

		dailyFlag := false
		scanner := bufio.NewScanner(strings.NewReader(ZTEXT))
		for scanner.Scan() {
			item := scanner.Text()
			item = strings.TrimSpace(item)

			// All daily checklist blocks begin with this
			if strings.Contains(item, "## Daily Routine") {
				dailyFlag = true
			}
			// only ingest daily goals, ignore all other content in note
			if dailyFlag {
				if strings.Contains(item, completedIndicator) {
					eventTime, err := time.Parse("2006-01-02", ZTITLE)
					if err != nil {
						panic(err)
					}
					goalEntry := models.Timeseries{Name: cleanName(item, completedIndicator), Date: eventTime, Value: 1}
					entries = append(entries, goalEntry)
				}
			}
			if strings.TrimSpace(item) == "" {
				dailyFlag = false
			}
		}
	}

	return entries
}

func cleanName(item, statusIndicator string) string {
	item = strings.Replace(item, statusIndicator, "", -1)
	item = strings.TrimSpace(item)
	return item
}
