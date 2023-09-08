package ingest

import (
	"bufio"
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

const (
	bearSQLLiteLocation = "./data/database.sqllite"
	completedIndicator  = "- [x]"
)

type BearTaskItem struct {
	Date      time.Time
	Name      string
	Completed bool
}

func Ingest() []BearTaskItem {
	fmt.Println("ingest")

	//return util.GenerateDummyData(5)

	return extractBearData()
}

func extractBearData() []BearTaskItem {

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
func retrieveDataFromBear(db *sql.DB) []BearTaskItem {
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

	entries := make([]BearTaskItem, 0)

	for rows.Next() {
		var ZTITLE, ZTEXT string
		if err = rows.Scan(&ZTITLE, &ZTEXT); err != nil {
			panic(err)
		}

		dailyFlag := false
		scanner := bufio.NewScanner(strings.NewReader(ZTEXT))
		for scanner.Scan() {
			item := scanner.Text()
			// All daily checklist blocks begin with this
			if strings.Contains(item, "## Daily Routine") {
				dailyFlag = true
			}
			// only ingest daily goals, ignore all other content in note
			if dailyFlag {
				entries = addLineEntryList(entries, item, ZTITLE)
			}
			if strings.TrimSpace(item) == "" {
				dailyFlag = false
			}
		}
	}

	return entries
}

func addLineEntryList(entryList []BearTaskItem, item string, date string) []BearTaskItem {
	item = strings.TrimSpace(item)

	if strings.Contains(item, completedIndicator) {
		eventTime, err := time.Parse("2006-01-02", date)
		if err != nil {
			panic(err)
		}
		goalEntry := BearTaskItem{Name: cleanName(item, completedIndicator), Date: eventTime, Completed: true}
		entryList = append(entryList, goalEntry)
	}
	return entryList
}

func cleanName(item, statusIndicator string) string {
	item = strings.Replace(item, statusIndicator, "", -1)
	item = strings.TrimSpace(item)
	return item
}
