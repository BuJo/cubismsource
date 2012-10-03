package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"jbossinfo"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)
import "code.google.com/p/gosqlite/sqlite"

var (
	sqliteDB = flag.String("db", "jvm-ram2.db", "Database for Sqlite implementation")
	site     = flag.String("site", "TST", "The site that the parsed html files belong to")

	htmlInfoRE     = regexp.MustCompile(">([A-Z#][A-Za-z ]*): </b>([^<]*)<")
	fileDateFormat = "2006-01-02 15:04:05-07:00.html"
)

const (
	SqliteTimeFormat = "2006-01-02 15:04:05-07:00"
	IsoD3Format      = "2006-01-02T15:04:05.000Z"
)

func parseJbossIntWithUnit(str string) (int, error) {
	chunks := strings.Split(str, " ")
	if len(chunks) != 2 {
		return 0, errors.New("Bad input")
	}
	nr, err := strconv.Atoi(chunks[0])
	if err != nil {
		return 0, err
	}

	switch chunks[1] {
	case "KB":
		nr *= 1024
	case "MB":
		nr *= 1024 * 1024
	case "GB":
		nr *= 1024 * 1024 * 1024
	default:
		return 0, errors.New("Bad Qualifier, must be MB/KB/GB")
	}

	return nr, nil
}

func parseHtml(filename string) (info *jbossinfo.JbossStatus, date time.Time, err error) {
	//fmt.Printf("DBG: parsing %s\n", filename)

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, date, err
	}

	matches := htmlInfoRE.FindAllSubmatch(bytes, -1 /*all*/)

	info = jbossinfo.NewStatus()
	filenameChunks := strings.Split(filename, "/")
	date, err = time.Parse(fileDateFormat, filenameChunks[len(filenameChunks)-1])
	if err != nil {
		return nil, date, err
	}

	for _, match := range matches {
		if len(match) < 3 {
			continue
		}
		key := string(match[1])

		value := string(match[2])

		switch key {
		case "Free Memory":
			info.JvmStatus.Free, err = parseJbossIntWithUnit(value)
		case "Max Memory":
			info.JvmStatus.Max, err = parseJbossIntWithUnit(value)
		case "Total Memory":
			info.JvmStatus.Total, err = parseJbossIntWithUnit(value)
		case "#Threads":
			info.Connectors[0].ThreadInfo.CurrentThreadCount, err = strconv.Atoi(value)
		}
		if err != nil {
			fmt.Printf("assignment error: %s\n", err)
		}

	}

	return info, date, nil
}

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

func main() {
	flag.Parse()

	r := bufio.NewReader(os.Stdin)

	peek, peekErr := r.Peek(5)
	if peekErr != nil {
		fmt.Printf("booh %s\n", peekErr)
	}
	fmt.Printf("startup, peek: %s\n", peek)

	queue := make(chan *InsertRequest, 100)
	ok := sqliteWriteHandler(queue)

	go func() {
		for line, err := r.ReadBytes('\n'); err == nil || err == io.EOF && len(line) > 0; line, err = r.ReadBytes('\n') {
			info, date, parseErr := parseHtml(string(line[0 : len(line)-1]))
			if parseErr != nil {
				fmt.Printf("Error parsing html: %s\n", parseErr)
				continue
			}

			queue <- &InsertRequest{*site, date, info}
		}
		queue <- nil
	}()
	<-ok
}
