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

func pullJboss(sites []string, interval time.Duration) {

	queue := make(chan *InsertRequest, 100)
	ok := sqliteWriteHandler(queue)

	ticker := time.Tick(interval)

	go func() {
		for now := range ticker {
			for _, site := range sites {
				info := getCurrentInfo(site)

				queue <- &InsertRequest{site, now, info}
			}
		}
		queue <- nil
	}()
	<-ok
}
