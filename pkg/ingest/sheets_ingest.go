package ingest

import (
	"gopkg.in/Iwark/spreadsheet.v2"
	"os"
	"strconv"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
)

const sheet = "1e-vL6bAHhuUZU2xCVDSSvgCAQMKCFMnPocn169ULpSg"

func Sheets() []Timeseries {

	// Columns:
	//	Date, Start Page, End Page, Minutes, Total Pages, Pages / Min, Book

	data, err := os.ReadFile("secret.json")
	if err != nil {
		panic(err)
	}
	conf, err := google.JWTConfigFromJSON(data, spreadsheet.Scope)
	if err != nil {
		panic(err)
	}
	client := conf.Client(context.TODO())

	service := spreadsheet.NewServiceWithClient(client)
	spreadsheet, err := service.FetchSpreadsheet(sheet)
	if err != nil {
		panic(err)
	}

	sheets := []string{"Reading 2021", "Reading 2022", "Reading 2023"}

	entries := make([]Timeseries, 0)

	for _, sheetName := range sheets {
		sheet, err := spreadsheet.SheetByTitle(sheetName)
		if err != nil {
			panic(err)
		}

		for _, row := range sheet.Rows {

			// skip first
			if row[0].Value == "Date" || row[0].Value == "" {
				continue
			}

			date, err := time.Parse("1/2/06", row[0].Value)
			if err != nil {
				panic(err)
			}
			minutes, err := strconv.Atoi(row[3].Value)
			if err != nil {
				minutes = 0
			}

			entries = append(entries, Timeseries{
				Date:  date,
				Name:  row[6].Value,
				Value: int64(minutes),
			})
		}
	}

	return entries
}
