package main

import "fmt"
import "log"
import "strconv"
import "io/ioutil"
import "time"
import "net/http"
import "encoding/json"

type appError struct {
	Error   error
	Message string
	Code    int
}

type appHandler func(http.ResponseWriter, *http.Request) *appError

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

	if site == "" {
		site = "TST"
	}

	//series := stepTimeSeries(getSqliteTimeSeries(site, expression, start, stop), start, stop, step)
	series := stepTimeSeries(getCurrentTimeSeries(site, expression, start, stop), start, stop, step)
	if series == nil {
		err.Message = "Something bad happened with steps"
		return &err
	}

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
		http.Error(w, e.Message+":"+e.Error.Error(), e.Code)
	}
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
