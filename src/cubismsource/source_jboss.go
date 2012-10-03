package main

import "fmt"
import "strconv"
import "time"
import "net/http"

var (
	JbossStatusUrls = map[string]string{
		"TST": "http://127.0.0.1:8080/status?XML=true",
	}
)

func getCurrentTimeSeries(site, field string, start, stop time.Time) *TimeSeries {
	if time.Now().Sub(start).Minutes() > 3 {
		fmt.Printf("Skipping, not supposed to look into past\n")
		return nil
	}

	fmt.Printf("Trying to get xml from: %s\n", JbossStatusUrls[site])

	resp, respErr := http.Get(JbossStatusUrls[site])
	if respErr != nil {
		fmt.Printf("Can't get jboss xml: %s\n", respErr)
		return nil
	}

	info, infoErr := GetJbossInfo(resp.Body)
	if infoErr != nil {
		fmt.Printf("Parsing jboss xml failed: %s\n", infoErr)
		return nil
	}

	series := TimeSeries{}
	series.Entries = []TimeSeriesEntry{}

	entry := TimeSeriesEntry{field, start, ""}

	switch field {
	case "free":
		entry.Value = strconv.Itoa(info.JvmStatus.Free)
	case "total":
		entry.Value = strconv.Itoa(info.JvmStatus.Total)
	case "max":
		entry.Value = strconv.Itoa(info.JvmStatus.Max)
	case "threads":
		entry.Value = strconv.Itoa(info.Connector.ThreadInfo.CurrentThreadCount)
	default:
		return nil
	}

	series.Entries = append(series.Entries, entry)

	fmt.Printf("%#v\n", entry)

	return &series
}
