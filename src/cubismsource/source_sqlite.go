package main

import "fmt"
import "time"
import "code.google.com/p/gosqlite/sqlite"

const (
	SqliteTimeFormat = "2006-01-02 15:04:05-07:00"
	IsoD3Format      = "2006-01-02T15:04:05.000Z"
)

func getSqliteTimeSeries(site, field string, start, stop time.Time) *TimeSeries {
	filename := "jvm-ram.db"

	dbconn, err := sqlite.Open(filename)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer dbconn.Close()

	switch field {
	case "site", "date", "free", "max", "total", "threads": // all good
	default:
		fmt.Printf("Bad field: %s\n", field)
		return nil
	}

	stmt, stmterr := dbconn.Prepare("select date, " + field + " from ram where site = ?1 and date between ?2 and ?3 order by date")
	if stmterr != nil {
		fmt.Print(err)
	}
	defer stmt.Finalize()

	err = stmt.Exec(site, start.Format(SqliteTimeFormat), stop.Format(SqliteTimeFormat))
	if err != nil {
		fmt.Print(err)
	}
	//fmt.Printf("sql: %s\n", stmt.SQL())

	series := TimeSeries{}
	series.Entries = []TimeSeriesEntry{}
	var datestr string
	var value string

	for stmt.Next() {
		err = stmt.Scan(&datestr, &value)
		date, dateError := time.Parse(SqliteTimeFormat, datestr)
		if dateError != nil {
			fmt.Printf("bad date: %v", dateError)
		}
		entry := TimeSeriesEntry{field, date, value}

		//fmt.Printf("%#v\n", entry)
		series.Entries = append(series.Entries, entry)
	}

	if len(series.Entries) == 0 {
		fmt.Printf("No Values for: %s\n", stmt.SQL())
	}

	return &series
}
