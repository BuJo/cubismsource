package main

import (
	"bufio"
	"errors"
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

var (
	htmlInfoRE     = regexp.MustCompile(">([A-Z#][A-Za-z ]*): </b>([^<]*)<")
	fileDateFormat = "2006-01-02 15:04:05-07:00.html"
)

const (
	SqliteTimeFormat = "2006-01-02 15:04:05-07:00"
	IsoD3Format      = "2006-01-02T15:04:05.000Z"
)

func parseJbossIntWithUnit(str string) (uint, error) {
	chunks := strings.Split(str, " ")
	if len(chunks) != 2 {
		return 0, errors.New("Bad input")
	}
	nr, err := strconv.ParseUint(chunks[0], 10, 0)
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

	return uint(nr), nil
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
func importHtml(site string) {
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

			queue <- &InsertRequest{site, date, info}
		}
		queue <- nil
	}()
	<-ok
}
