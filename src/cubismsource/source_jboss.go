package main

import (
	"fmt"
	"io"
	"jbossinfo"
	"net/http"
	"strconv"
	"time"
)

var (
	JbossStatusUrls = map[string]string{
		"TST": "http://127.0.0.1:8080/status?XML=true",
	}
)

func getFieldValueFromXml(xml io.Reader, field string) string {
	var value string

	info, infoErr := jbossinfo.ParseJbossInfoXML(xml)
	if infoErr != nil {
		fmt.Printf("Parsing jboss xml failed: %s\n", infoErr)
		return ""
	}

	switch field {
	case "free":
		value = strconv.FormatUint(uint64(info.JvmStatus.Free), 10)
	case "total":
		value = strconv.FormatUint(uint64(info.JvmStatus.Total), 10)
	case "max":
		value = strconv.FormatUint(uint64(info.JvmStatus.Max), 10)
	case "used":
		value = strconv.FormatUint(uint64(info.JvmStatus.Total-info.JvmStatus.Free), 10)
	case "threads":
		threadCount := 0
		for _, connector := range info.Connectors {
			threadCount += connector.ThreadInfo.CurrentThreadsBusy
		}
		value = strconv.Itoa(threadCount)
	case "maxRequestTime":
		maxRequestTime := 0
		for _, connector := range info.Connectors {
			for _, worker := range connector.Workers {
				if worker.RequestProcessingTime > maxRequestTime {
					maxRequestTime = worker.RequestProcessingTime
				}
			}
		}
		value = strconv.Itoa(maxRequestTime)
	default:
		return ""
	}

	return value
}

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

	series := TimeSeries{}
	series.Entries = []TimeSeriesEntry{}

	entry := TimeSeriesEntry{field, start, getFieldValueFromXml(resp.Body, field)}

	series.Entries = append(series.Entries, entry)

	fmt.Printf("%#v\n", entry)

	return &series
}
