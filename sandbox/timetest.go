package main

import "time"
import "fmt"

func main() {
	t := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	fmt.Printf("%s\n", t.Format("2006-01-02 10:04:00"))
	fmt.Printf("%s\n", t.Format("2006-01-02 15:04:00-07:00"))
	fmt.Printf("%s\n", t.Format("2006-01-02 15:04:05.999999999 -0700 MST"))
	fmt.Printf("%s\n", t.Format(time.RFC3339))

	t2, _ := time.Parse("2006-01-02T15:04:05.999Z", "2012-10-01T15:15:40.000Z")
	fmt.Printf("%s\n", t2.Format(time.RFC3339))

	t2, _ = time.Parse("2006-01-02T15:04:05.000Z", "2012-10-01T15:15:40.000Z")
	fmt.Printf("%s\n", t2.Format(time.RFC3339))
}
