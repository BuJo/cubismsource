package main

import (
	"code.google.com/p/gosqlite/sqlite"
	"fmt"
	"jbossinfo"
	"time"
)

type InsertRequest struct {
	site      string
	date      time.Time
	jbossinfo *jbossinfo.JbossStatus
}

func sqliteWriteHandler(queue chan *InsertRequest) chan bool {
	ok := make(chan bool)
	filename := *sqliteDB

	dbconn, err := sqlite.Open(filename)
	if err != nil {
		fmt.Println(err)
		return nil
	}
	stmt, stmterr := dbconn.Prepare("INSERT INTO jvmmetrics (site, date, free, max, total, threads, xml) VALUES (?1, ?2, ?3, ?4, ?5, ?6, ?7)")
	if stmterr != nil {
		fmt.Print(err)
		return nil
	}

	go func() {
		fmt.Printf("Waiting for parses\n")
		for {
			request := <-queue
			//fmt.Printf("got one request: %#v\n", request)

			if request == nil || request.site == "" {
				stmt.Finalize()
				dbconn.Close()
				close(queue)
				ok <- true
				return
			}

			xml, xmlErr := jbossinfo.InfoXML(request.jbossinfo)
			if xmlErr != nil {
				fmt.Printf("error generating xml: %s\n", xmlErr)
				continue
			}

			err = stmt.Exec(request.site, request.date.Format(SqliteTimeFormat), request.jbossinfo.JvmStatus.Free, request.jbossinfo.JvmStatus.Max, request.jbossinfo.JvmStatus.Total, request.jbossinfo.Connectors[0].ThreadInfo.CurrentThreadCount, xml)
			if err != nil {
				fmt.Print(err)
			}
			//fmt.Printf("SQL: %s\n", stmt.SQL())
			stmt.Next()

		}

	}()

	return ok
}
