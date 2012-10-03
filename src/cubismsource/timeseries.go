package main

import "fmt"
import "time"
import "sort"

type TimeSeriesEntry struct {
	Name  string    `json:"symbol"`
	Time  time.Time `json:"date"`
	Value string    `json:"value"`
}
type TimeSeries struct {
	Entries []TimeSeriesEntry
}

func stepTimeSeries(series *TimeSeries, start, stop time.Time, step time.Duration) *TimeSeries {
	var steppedSeries TimeSeries

	if series == nil || len(series.Entries) < 1 {
		fmt.Printf("no input time series\n")
		return &steppedSeries
	}

	stepcount := stop.Sub(start).Seconds() / step.Seconds()
	fmt.Printf("Stepping, Original Source steps: %d, requested: %d\n", len(series.Entries), int(stepcount))

	for currentTime := start; currentTime.Before(stop); currentTime = currentTime.Add(step) {
		entry := TimeSeriesEntry{series.Entries[0].Name, currentTime, ""}

		i := sort.Search(len(series.Entries), func(i int) bool { return series.Entries[i].Time.After(currentTime) })
		if i <= len(series.Entries) {
			if i == 0 {
				entry.Value = series.Entries[0].Value
			} else {
				entry.Value = series.Entries[i-1].Value
			}
		}

		steppedSeries.Entries = append(steppedSeries.Entries, entry)
	}

	fmt.Printf("TimeSeries: %d more than source\n", len(steppedSeries.Entries)-len(series.Entries))

	if len(steppedSeries.Entries) != int(stepcount) {
		fmt.Printf("Have %d instead of %d steps.", len(steppedSeries.Entries), stepcount)
	}

	countNaN := 0
	for _, e := range steppedSeries.Entries {
		if e.Value == "" {
			countNaN += 1
		}
	}
	if countNaN > 0 {
		fmt.Printf("TimeSeries: have %d NaN values (%.2f%%)\n", countNaN, (float64(countNaN)/float64(len(steppedSeries.Entries)))*100)
	}

	return &steppedSeries
}
