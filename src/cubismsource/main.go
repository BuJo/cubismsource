package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	bindServer = flag.String("bind", ":8080", "Address/Port to bind to, default: *:8080")
	httpRoot   = flag.String("root", "assets", "HTTP Root")
	intImpl    = flag.String("impl", "jboss", "INTERNAL specifying server behaviour (sqlite, jboss)")

	Usage = func() {
		fmt.Printf("%s Usage:\n", os.Args[0])
		flag.PrintDefaults()
	}
)

type appError struct {
	Error   error
	Message string
	Code    int
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

func getJsonMetrics(site, expression string, start, stop time.Time, step time.Duration) ([]byte, error) {
	var series *TimeSeries

	if *intImpl == "sqlite" {
		series = getSqliteTimeSeries(site, expression, start, stop)
	} else if *intImpl == "jboss" {
		series = getCurrentTimeSeries(site, expression, start, stop)
	}

	series = stepTimeSeries(series, start, stop, step)

	if series == nil {
		return nil, errors.New("Something bad happened with steps")
	}

	bytes, err := json.Marshal(series.Entries)
	if err != nil {
		return nil, errors.New("marshalling time series failed")
	}

	return bytes, nil
}

// "/1.0/metric?expression=free&start=2012-10-01T15:15:40.000Z&stop=2012-10-01T15:16:50.000Z&step=10000"
func handleMetrics(w http.ResponseWriter, r *http.Request) *appError {
	fmt.Printf("Metrics Request: %#v\n", r.URL.Query())

	err := appError{nil, "Something went Wrong", 500}

	var site string
	var expression string
	var start time.Time
	var stop time.Time
	var step time.Duration

	var bytes []byte

	query := r.URL.Query()

	site = query.Get("site")
	expression = query.Get("expression")
	start, err.Error = time.Parse(IsoD3Format, query.Get("start"))
	stop, err.Error = time.Parse(IsoD3Format, query.Get("stop"))
	step, err.Error = time.ParseDuration(query.Get("step") + "ms")

	if err.Error != nil {
		return &err
	}

	// XXX: js does not send (parseable) the time zone, temporary hackfix here :(
	//start = time.Date(start.Year(), start.Month(), start.Day(), start.Hour(), start.Minute(), start.Second(), start.Nanosecond(), time.UTC)
	//stop = time.Date(stop.Year(), stop.Month(), stop.Day(), stop.Hour(), stop.Minute(), stop.Second(), stop.Nanosecond(), time.UTC)
	start = start.Add(time.Hour * 1)
	stop = stop.Add(time.Hour * 1)

	if site == "" {
		site = "TST"
	}

	bytes, err.Error = getJsonMetrics(site, expression, start, stop, step)
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
		http.Error(w, e.Message+":"+e.Error.Error(), e.Code)
	}
}

func handleCubism(w http.ResponseWriter, r *http.Request) *appError {
	filename := *httpRoot + "/cubism.html"

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
	flag.Parse()

	http.Handle("/1.0/", appHandler(handleMetrics))
	http.Handle("/cubism", appHandler(handleCubism))

	log.Fatal(http.ListenAndServe(*bindServer, nil))
}
