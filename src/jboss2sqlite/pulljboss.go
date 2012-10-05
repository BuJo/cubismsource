package main

import (
	"fmt"
	"jbossinfo"
	"net/http"
	"time"
)

var (
	JbossStatusUrls = map[string]string{
		"TST": "http://127.0.0.1:8080/status?XML=true",
	}
)

func getCurrentInfo(site string) *jbossinfo.JbossStatus {
	fmt.Printf("Trying to get xml from: %s\n", JbossStatusUrls[site])

	resp, respErr := http.Get(JbossStatusUrls[site])
	if respErr != nil {
		fmt.Printf("Can't get jboss xml: %s\n", respErr)
		return nil
	}

	info, infoErr := jbossinfo.ParseJbossInfoXML(resp.Body)
	if infoErr != nil {
		fmt.Printf("Parsing jboss xml failed: %s\n", infoErr)
		return nil
	}

	return info
}

func pullSite(site string, now time.Time, queue chan *InsertRequest) {
	fmt.Printf("pulling %s: %s\n", site, now.Format(time.RFC3339))
	info := getCurrentInfo(site)

	queue <- &InsertRequest{site, now, info}
}

func pullJboss(sites []string, interval time.Duration) {
	for _, site := range sites {
		if JbossStatusUrls[site] == "" {
			fmt.Printf("Error, %s is an unkonwn site\n", site)
			return
		}
	}

	ticker := time.Tick(interval)

	for _ = range ticker {
		queue := make(chan *InsertRequest, 10)
		ok := sqliteWriteHandler(queue)

    // remember to parallelize pullSite, but remember that the sqlite close has to come after all are done
		for _, site := range sites {
			pullSite(site, time.Now(), queue)
		}

		queue <- nil
		<-ok
	}
}
