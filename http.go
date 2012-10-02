package main

import "fmt"
import "log"
import "strconv"
import "io/ioutil"
import "time"
import "net/http"
import "encoding/json"
import "code.google.com/p/gosqlite/sqlite"

type appError struct {
	Error   error
	Message string
	Code    int
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

// "/1.0/metric?expression=free&start=2012-10-01T15:15:40.000Z&stop=2012-10-01T15:16:50.000Z&step=10000"
func handleMetrics(w http.ResponseWriter, r *http.Request) *appError {
	fmt.Printf("got req: %#v\n", r)

	var err appError

	var expression string
	var start time.Time
	var stop time.Time
	var step int64
	var bytes []byte

	query := r.URL.Query()

	expression = query.Get("expression")
	start, err.Error = time.Parse(IsoD3Format, query.Get("start"))
	stop, err.Error = time.Parse(IsoD3Format, query.Get("stop"))
	step, err.Error = strconv.ParseInt(query.Get("step"), 0, 0)

	if err.Error != nil {
		return &err
	}

	series := stepTimeSeries(getTimeSeries(expression, start, stop), int(step))

	bytes, err.Error = json.Marshal(series.Entries)
	if err.Error != nil {
		return &err
	}

	w.Header().Set("Content-Type", "application/json;charset=utf-8")
	w.Header().Set("Content-Length", strconv.Itoa(len(bytes)))
	w.Header().Set("Date", time.Now().Format(time.RFC1123))
	w.Header().Set("Cache-Control", "must-revalidate")
	w.Header().Set("Connection", "close")

	_, err.Error = w.Write(bytes)
	if err.Error != nil {
		return &err
	}

	//fmt.Printf("%s\n", string(bytes))

	return nil
}
func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if e := fn(w, r); e != nil { // e is *appError, not os.Error.
		http.Error(w, e.Message, e.Code)
	}
}

type TimeSeriesEntry struct {
	Name  string    `json:"symbol"`
	Time  time.Time `json:"date"`
	Value string    `json:"value"`
}
type TimeSeries struct {
	Entries []TimeSeriesEntry
}

func stepTimeSeries(series *TimeSeries, step_ms int) *TimeSeries {
	return series
}

const (
	SqliteTimeFormat = "2006-01-02 15:04:05-07:00"
	IsoD3Format      = "2006-01-02T15:04:05.000Z"
)

func getTimeSeries(field string, start, stop time.Time) *TimeSeries {
	filename := "jvm-ram.db"
	site := "TST"

	dbconn, err := sqlite.Open(filename)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	defer dbconn.Close()

	stmt, stmterr := dbconn.Prepare("select date, free from ram where site = ?1 and date between ?2 and ?3")
	if stmterr != nil {
		fmt.Print(err)
	}
	defer stmt.Finalize()

	err = stmt.Exec(site, start.Format(SqliteTimeFormat), stop.Format(SqliteTimeFormat))
	if err != nil {
		fmt.Print(err)
	}
	fmt.Printf("sql: %s\n", stmt.SQL())

	series := TimeSeries{}
	series.Entries = []TimeSeriesEntry{}
	var datestr string
	var free int

	for stmt.Next() {
		err = stmt.Scan(&datestr, &free)
		date, dateError := time.Parse(SqliteTimeFormat, datestr)
		if dateError != nil {
			fmt.Printf("bad date: %v", dateError)
		}
		entry := TimeSeriesEntry{field, date, strconv.Itoa(free)}

		//fmt.Printf("%#v\n", entry)
		series.Entries = append(series.Entries, entry)
	}

	return &series
}

func handleCubism(w http.ResponseWriter, r *http.Request) *appError {
	filename := "cubism.html"

	var err appError
	var data []byte

	data, err.Error = ioutil.ReadFile(filename)
	if err.Error != nil {
		return &err
	}

	w.Write(data)

	return nil
}

func main() {
	http.Handle("/1.0/", appHandler(handleMetrics))
	http.Handle("/cubism", appHandler(handleCubism))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
